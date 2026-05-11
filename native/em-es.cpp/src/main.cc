#include "schemaregistry/rest/model/Schema.h"
#include <any>
#include <atomic>
#include <cassert>
#include <csignal>
#include <format>
#include <kafka/Error.h>
#include <kafka/KafkaConsumer.h>
#include <kafka/KafkaProducer.h>
#include <kafka/ProducerCommon.h>
#include <optional>
#include <spdlog/spdlog.h>

#include <schemaregistry/rest/SchemaRegistryClient.h>
#include <schemaregistry/serdes/protobuf/ProtobufDeserializer.h>
#include <schemaregistry/serdes/protobuf/ProtobufSerializer.h>

#include <v2k/beam/AxisSize.pb.h>
#include <v2k/beam/Current.pb.h>
#include <v2k/beam/Emittance.pb.h>
#include <v2k/beam/EnergySpread.pb.h>

#include <v2k/ccd.h>
#include <v2k/em_es.h>

using namespace schemaregistry::rest;
using namespace schemaregistry::serdes;
using namespace schemaregistry::serdes::protobuf;
using namespace kafka::clients::consumer;
using namespace kafka::clients::producer;

namespace v2k {

struct TumblingWindow {
  std::size_t begin;
  std::size_t end;

  bool operator==(TumblingWindow const &other) const {
    return begin == other.begin && end == other.end;
  }

  struct Hasher {
    size_t operator()(v2k::TumblingWindow const &window) const {
      std::size_t h1 = std::hash<size_t>{}(window.begin);
      std::size_t h2 = std::hash<size_t>{}(window.end);

      return h1 ^ (h2 << 1);
    }
  };
};

class TumblingWindowAssigner {
private:
  std::size_t window_size_;
  std::list<v2k::TumblingWindow> windows_;

public:
  TumblingWindowAssigner(std::size_t window_size)
      : window_size_(window_size) {};

  using Ptr = std::shared_ptr<TumblingWindowAssigner>;

  static Ptr make(std::size_t window_size) {
    return std::make_shared<TumblingWindowAssigner>(window_size);
  }

  v2k::TumblingWindow Assign(std::size_t time) {
    if (windows_.empty()) {
      return *windows_.insert(windows_.end(), {
                                                  .begin = time,
                                                  .end = time + window_size_,
                                              });
    }

    std::list<v2k::TumblingWindow>::iterator it = windows_.begin();

    while (it != windows_.end() && it->end <= time) {
      it++;
    }

    if (it != windows_.end() && it->begin <= time) {
      return *it;
    }

    std::size_t last = (--it)->end;
    std::size_t begin = last + ((time - last) / window_size_) * window_size_;

    return *windows_.insert(it, {
                                    .begin = begin,
                                    .end = begin + window_size_,
                                });
  }

  std::list<v2k::TumblingWindow> Completed(std::size_t watermark) {
    std::list<v2k::TumblingWindow> r;
    std::list<v2k::TumblingWindow>::iterator it = windows_.begin();

    while (it != windows_.end() && watermark >= it->end) {
      r.push_back(*it);
      it = windows_.erase(it);
    }

    return r;
  }
};

class Stream {
  struct State {
    std::size_t count;
    double sum;
  };

private:
  std::any data_;
  std::size_t watermark_;
  v2k::TumblingWindowAssigner::Ptr assigner_;
  std::unordered_map<std::size_t, v2k::Stream::State> state_;

public:
  Stream(std::any data) : data_(data) {}

  using Ptr = std::shared_ptr<Stream>;

  static Ptr make(std::any data) { return std::make_shared<Stream>(data); }

  std::any const &Data() { return data_; }
  std::size_t Watermark() const { return watermark_; }

  void SetWindowAssigner(v2k::TumblingWindowAssigner::Ptr assigner) {
    assigner_ = assigner;
  }

  void Push(std::size_t time, double value) {
    auto const window = assigner_->Assign(time);

    if (!state_.contains(window.begin)) {
      state_.insert({window.begin, State{.count = 0, .sum = 0}});
    }

    auto &state = state_.at(window.begin);

    state.count += 1;
    state.sum += value;
    watermark_ = time;
  }

  std::optional<double> Trigger(v2k::TumblingWindow const window) {
    if (!state_.contains(window.begin)) {
      return std::nullopt;
    }

    auto const &state = state_.at(window.begin);
    double const avg = state.sum / state.count;

    state_.erase(window.begin);

    return avg;
  }
};

class StreamGroup {
private:
  std::vector<v2k::Stream::Ptr> streams_;
  v2k::TumblingWindowAssigner::Ptr assigner_;

public:
  StreamGroup(v2k::TumblingWindowAssigner::Ptr assigner)
      : assigner_(assigner) {}

  void Add(v2k::Stream::Ptr const stream) {
    stream->SetWindowAssigner(assigner_);
    streams_.push_back(stream);
  }

  std::size_t Watermark() {
    std::size_t w = UINT64_MAX;

    for (auto const &s : streams_) {
      auto const sw = s->Watermark();

      if (sw < w) {
        w = sw;
      }
    }

    return w;
  }

  std::list<v2k::TumblingWindow> Completed() {
    return assigner_->Completed(Watermark());
  }
};

}; // namespace v2k

std::atomic_bool running = true;

void stop(int sig) {
  if (running) {
    running = false;
  }
}

int main(int argc, char **argv) {
  (void)argc;
  (void)argv;

  spdlog::set_level(spdlog::level::debug);
  signal(SIGINT, stop);

  char const *const dir = getenv("CONFIG_DIR");
  char const *const rp_uri = getenv("RP_URI");
  char const *const sr_uri = getenv("SR_URI");

  v2k::ccd::Config const ccd_config{std::format("{}/ccd/config.yaml", dir)};
  v2k::ems::Config const ems_config{std::format("{}/em-es/config.yaml", dir)};

  kafka::Properties const props({{"bootstrap.servers", {{rp_uri}}}});
  KafkaConsumer consumer(props);
  KafkaProducer producer(props);

  std::set<std::string> topics;
  v2k::StreamGroup group(
      v2k::TumblingWindowAssigner::make(ems_config.window * 1000));
  std::vector<v2k::ems::AxisRealization> realizations;

  auto const sr_config =
      std::make_shared<ClientConfiguration>(std::vector<std::string>{sr_uri});
  auto const sr_client = SchemaRegistryClient::newClient(sr_config);

  std::unordered_map<std::string, std::string> rule_config{};
  auto deser_config = DeserializerConfig(std::nullopt, false, rule_config);

  auto const cur_deser =
      std::make_unique<ProtobufDeserializer<v2k::beam::Current>>(
          sr_client, nullptr, deser_config);
  auto const axis_deser =
      std::make_unique<ProtobufDeserializer<v2k::beam::AxisSize>>(
          sr_client, nullptr, deser_config);

  auto cur_stream = v2k::Stream::make(nullptr);
  auto ccd_stream = std::map<std::string, v2k::Stream::Ptr>();

  for (auto it = ems_config.cams.begin(); it != ems_config.cams.end(); ++it) {
    auto const &setup = it->second;
    auto const &config = ccd_config.cams.at(it->first);

    if (setup.axes.x.weight != 0) {
      realizations.push_back({
          .config = config.axes.x,
          .setup = setup.axes.x,
          .measurement = 0,
      });

      auto const topic = std::format("vepp.ccd.{}.sigma_x", it->first);
      auto const stream = v2k::Stream::make(realizations.size() - 1);

      topics.insert(topic);
      ccd_stream.insert({topic, stream});
      group.Add(stream);
    }

    if (setup.axes.z.weight != 0) {
      realizations.push_back({
          .config = config.axes.z,
          .setup = setup.axes.z,
          .measurement = 0,
      });

      auto const topic = std::format("vepp.ccd.{}.sigma_z", it->first);
      auto const stream = v2k::Stream::make(realizations.size() - 1);

      topics.insert(topic);
      ccd_stream.insert({topic, stream});
      group.Add(stream);
    }
  }

  topics.insert("vepp.currents.fz");
  group.Add(cur_stream);
  consumer.subscribe(topics);

  v2k::beam::Emittance emittance;
  v2k::beam::EnergySpread energy_spread;

  while (running) {
    auto records = consumer.poll(std::chrono::milliseconds(500));

    for (const auto &record : records) {
      if (record.error()) {
        spdlog::error(record.error().message());
        continue;
      }

      // Remove Confluent wire header and message index
      char const *const payload = &((char const *)record.value().data())[6];
      std::size_t const size = record.value().size() - 6;

      if ("vepp.currents.fz" == record.topic()) {
        SerializationContext ser_ctx;
        ser_ctx.topic = record.topic();
        ser_ctx.serde_type = SerdeType::Value;
        ser_ctx.serde_format = SerdeFormat::Protobuf;
        ser_ctx.headers = std::nullopt;

        std::vector<uint8_t> payload(
            static_cast<const uint8_t *>(record.value().data()),
            static_cast<const uint8_t *>(record.value().data()) +
                record.value().size());

        auto const cur = cur_deser->deserialize(ser_ctx, payload);
        cur_stream->Push(cur->time(), cur->value());
      } else {
        SerializationContext ser_ctx;
        ser_ctx.topic = record.topic();
        ser_ctx.serde_type = SerdeType::Value;
        ser_ctx.serde_format = SerdeFormat::Protobuf;
        ser_ctx.headers = std::nullopt;

        std::vector<uint8_t> payload(
            static_cast<const uint8_t *>(record.value().data()),
            static_cast<const uint8_t *>(record.value().data()) +
                record.value().size());

        auto const axis_size = axis_deser->deserialize(ser_ctx, payload);
        ccd_stream.at(record.topic())
            ->Push(axis_size->time(), axis_size->value());
      }

      auto const completed = group.Completed();

      for (auto const window : completed) {
        spdlog::debug("window {} - {} completed", window.begin, window.end);

        double estimation[] = {0.1, 0.1};

        auto const cur_opt = cur_stream->Trigger(window);
        if (!cur_opt.has_value()) {
          spdlog::error("there are no records for the current");
          continue;
        }

        bool ok = true;

        for (auto const &[name, stream] : ccd_stream) {
          auto const realization_index =
              std::any_cast<std::size_t>(stream->Data());
          auto &realization = realizations.at(realization_index);

          auto const size_opt = stream->Trigger(window);
          if (!size_opt.has_value()) {
            spdlog::error("there are no record for the {}", name);
            ok = false;
            break;
          }

          realization.measurement = size_opt.value();
        }

        if (!ok) {
          continue;
        }

        if (!v2k::ems::Estimate(realizations, estimation)) {
          spdlog::error("estimation: failed");
        }

        auto deliveryCb = [](const RecordMetadata &metadata,
                             const kafka::Error &error) {
          if (!error) {
            std::cout << "Message delivered: " << metadata.toString()
                      << std::endl;
          } else {
            std::cerr << "Message failed to be delivered: " << error.message()
                      << std::endl;
          }
        };

        emittance.set_time(window.end);
        emittance.set_emittance(estimation[0]);
        emittance.set_current(cur_opt.value());

        auto const em_data = new char[emittance.ByteSizeLong()];
        emittance.SerializeToArray(em_data, emittance.ByteSizeLong());

        producer.send(kafka::clients::producer::ProducerRecord(
                          "vepp.emittance", kafka::NullKey,
                          kafka::Value(em_data, emittance.ByteSizeLong())),
                      deliveryCb);

        energy_spread.set_time(window.end);
        energy_spread.set_energy_spread(estimation[1]);
        energy_spread.set_current(cur_opt.value());

        auto const es_data = new char[energy_spread.ByteSizeLong()];
        emittance.SerializeToArray(es_data, energy_spread.ByteSizeLong());

        producer.send(kafka::clients::producer::ProducerRecord(
                          "vepp.energy_spread", kafka::NullKey,
                          kafka::Value(em_data, energy_spread.ByteSizeLong())),
                      deliveryCb);
      }
    }
  }

  consumer.close();

  return 0;
}
