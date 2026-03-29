#include <chrono>
#include <kafka/KafkaConsumer.h>

using namespace kafka::clients::consumer;

int main(int argc, char **argv) {
  (void)argc;
  (void)argv;

  KafkaConsumer consumer({{"bootstrap.servers", "localhost:9091"}});

  consumer.subscribe({kafka::Topic{"test"}});

  while (true) {
    auto records = consumer.poll(std::chrono::milliseconds(100));

    for (const auto &record : records) {
      if (!record.error()) {
        std::cout << record.topic() << std::endl;
        std::cout << record.value().toString() << std::endl;
      }
    }
  }

  consumer.close();

  return 0;
}
