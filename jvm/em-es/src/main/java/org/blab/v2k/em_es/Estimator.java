package org.blab.v2k.em_es;

import io.confluent.kafka.serializers.protobuf.KafkaProtobufDeserializer;
import io.confluent.kafka.serializers.protobuf.KafkaProtobufSerializer;
import io.confluent.kafka.serializers.subject.RecordNameStrategy;
import org.apache.commons.math3.linear.Array2DRowRealMatrix;
import org.apache.commons.math3.linear.ArrayRealVector;
import org.apache.commons.math3.linear.RealMatrix;
import org.apache.commons.math3.linear.RealVector;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.typeinfo.TypeInformation;
import org.apache.flink.api.java.tuple.Tuple2;
import org.apache.flink.connector.kafka.sink.KafkaPartitioner;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.connector.kafka.source.reader.deserializer.KafkaRecordDeserializationSchema;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.api.functions.ProcessFunction;
import org.apache.flink.streaming.api.windowing.assigners.TumblingEventTimeWindows;
import org.apache.flink.util.Collector;
import org.apache.flink.util.ConfigurationException;
import org.apache.flink.util.OutputTag;
import org.apache.kafka.common.TopicPartition;
import org.blab.v2k.*;
import org.blab.v2k.io.Record;
import org.blab.v2k.io.aggregation.AverageAggregator;
import org.blab.v2k.io.aggregation.ObservationAggregator;
import org.blab.v2k.io.index.IndexAssigner;
import org.blab.v2k.io.watermark.PunctuatedWatermarkGeneratorSupplier;
import org.blab.v2k.math.Parameter;
import org.blab.v2k.math.map.CompoundParameterValidator;
import org.blab.v2k.math.map.ExponentialParameterValidator;
import org.blab.v2k.math.optimization.GaussNewtonOptimizer;
import tools.jackson.databind.ObjectMapper;
import tools.jackson.dataformat.yaml.YAMLFactory;

import java.io.File;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.stream.Stream;

public class Estimator {
  public static <T extends Record> DataStream<T> stream(
          TopicPartition partition,
          StreamExecutionEnvironment environment) {
    var builder = KafkaSource.<T>builder()
            .setBootstrapServers(System.getenv("RP_URI"))
            .setPartitions(Set.of(partition))
            .setStartingOffsets(OffsetsInitializer.latest())
            .setDeserializer(KafkaRecordDeserializationSchema.valueOnly(
                    KafkaProtobufDeserializer.class,
                    Map.of("schema.registry.url", System.getenv("SR_URI"))
            ));

    return environment.fromSource(builder.build(), WatermarkStrategy
            .<T>forGenerator(new PunctuatedWatermarkGeneratorSupplier<>())
            .withTimestampAssigner((r, t) -> r.getTime()), partition.topic());
  }

  public static void create(StreamExecutionEnvironment environment) throws Exception {
    var config = resolveConfig();
    var window = config.getWindow();
    var partitions = new HashMap<String, TopicPartition>();
    var currentPartition = new TopicPartition(Config.TOPIC_CURRENT, 0);

    partitions.put(Config.TOPIC_CURRENT, currentPartition);

    var activeAxes = config.getCams().entrySet().stream()
            .flatMap(cs -> Stream.of(cs.getValue().getAxes().getX(), cs.getValue().getAxes().getZ()))
            .filter(r -> r.getWeight() > 0)
            .toList();

    for (var axis : activeAxes) {
      partitions.put(axis.getChannel(), new TopicPartition(axis.getChannel(), 0));
    }

    DataStream<CurrentOuterClass.Current> currentStream = stream(partitions.get(Config.TOPIC_CURRENT), environment);
    DataStream<Parameter> stream = currentStream
            .map(CurrentOuterClass.Current::getValue)
            .map(new IndexAssigner(-1, 0))
            .windowAll(TumblingEventTimeWindows.of(Duration.ofSeconds(window)))
            .aggregate(new AverageAggregator());

    int i = 0;

    for (var axis : activeAxes) {
      DataStream<AxisSizeOuterClass.AxisSize> axisSizeStream = stream(partitions.get(axis.getChannel()), environment);

      stream = stream.union(axisSizeStream
              .map(AxisSizeOuterClass.AxisSize::getValue)
              .map(new IndexAssigner(i++, 0))
              .windowAll(TumblingEventTimeWindows.of(Duration.ofSeconds(window)))
              .aggregate(new AverageAggregator())
      );
    }

    var model = new Model();
    var observation = new Array2DRowRealMatrix(i, model.getN());
    var weight = new ArrayRealVector(activeAxes.size());

    i = 0;

    for (var axis : activeAxes) {
      observation.setEntry(i, 2, axis.getConfig().getBeta());
      observation.setEntry(i, 3, axis.getConfig().getDispersion());
      weight.setEntry(i, axis.getWeight());

      i += 1;
    }

    var problem = Problem.builder()
            .model(model)
            .observation(observation)
            .observationWeight(weight)
            .initialGuess(new ArrayRealVector(new double[]{0, 0}))
            .parameterValidator(new CompoundParameterValidator()
                    .withValidator(0, new ExponentialParameterValidator())
                    .withValidator(1, new ExponentialParameterValidator()))
            .build();

    var em = new OutputTag<>("em", TypeInformation.of(EmittanceOuterClass.Emittance.class));
    var es = new OutputTag<>("es", TypeInformation.of(EnergySpreadOuterClass.EnergySpread.class));
    var log = new OutputTag<>("log", TypeInformation.of(LogOuterClass.Log.class));

    var solverConfig = config.getSolver();

    var result = stream
            .keyBy(e -> 0)
            .window(TumblingEventTimeWindows.of(Duration.ofSeconds(window)))
            .aggregate(new ObservationAggregator(observation, 1))
            .process(new ProcessFunction<Tuple2<RealVector, RealMatrix>, Void>() {
              @Override
              public void processElement(
                      Tuple2<RealVector, RealMatrix> observation,
                      ProcessFunction<Tuple2<RealVector, RealMatrix>, Void>.Context context,
                      Collector<Void> collector) {
                try {
                  var optimum = GaussNewtonOptimizer.builder()
                          .fromConfig(solverConfig)
                          .build()
                          .solve(problem.withObservation(observation.f1));

                  var c = observation.f0.getEntry(0);
                  var e = optimum.getPoint().getEntry(0);
                  var s = optimum.getPoint().getEntry(1);

                  context.output(em, EmittanceOuterClass.Emittance.newBuilder()
                          .setTime(context.timestamp())
                          .setCurrent(c)
                          .setEmittance(e)
                          .build());

                  context.output(es, EnergySpreadOuterClass.EnergySpread.newBuilder()
                          .setTime(context.timestamp())
                          .setCurrent(c)
                          .setEnergySpread(s)
                          .build());
                } catch (Exception e) {
                  context.output(log, LogOuterClass.Log.newBuilder()
                          .setTime(context.timestamp())
                          .setLevel(LogOuterClass.Log.Level.ERROR)
                          .setMessage(e.getMessage())
                          .build());
                }
              }
            });

    result.getSideOutput(em).sinkTo(KafkaSink.<EmittanceOuterClass.Emittance>builder()
            .setBootstrapServers(System.getenv("RP_URI"))
            .setRecordSerializer(KafkaRecordSerializationSchema.<EmittanceOuterClass.Emittance>builder()
                    .setTopic(Config.TOPIC_EMITTANCE)
                    .setPartitioner((KafkaPartitioner<Object>) (o, bytes, bytes1, s, ints) -> 0)
                    .setKafkaValueSerializer(KafkaProtobufSerializer.class, Map.of(
                            "schema.registry.url", System.getenv("SR_URI"),
                            "value.subject.name.strategy", RecordNameStrategy.class.getName()
                    ))
                    .build())
            .build());

    result.getSideOutput(es).sinkTo(KafkaSink.<EnergySpreadOuterClass.EnergySpread>builder()
            .setBootstrapServers(System.getenv("RP_URI"))
            .setRecordSerializer(KafkaRecordSerializationSchema.<EnergySpreadOuterClass.EnergySpread>builder()
                    .setTopic(Config.TOPIC_ENERGY_SPREAD)
                    .setPartitioner((KafkaPartitioner<Object>) (o, bytes, bytes1, s, ints) -> 0)
                    .setKafkaValueSerializer(KafkaProtobufSerializer.class, Map.of(
                            "schema.registry.url", System.getenv("SR_URI"),
                            "value.subject.name.strategy", RecordNameStrategy.class.getName()
                    ))
                    .build())
            .build());

    result.getSideOutput(log).sinkTo(KafkaSink.<LogOuterClass.Log>builder()
            .setBootstrapServers(System.getenv("RP_URI"))
            .setRecordSerializer(KafkaRecordSerializationSchema.<LogOuterClass.Log>builder()
                    .setTopic(Config.TOPIC_LOG)
                    .setPartitioner((KafkaPartitioner<Object>) (o, bytes, bytes1, s, ints) -> 0)
                    .setKafkaValueSerializer(KafkaProtobufSerializer.class, Map.of(
                            "schema.registry.url", System.getenv("SR_URI"),
                            "value.subject.name.strategy", RecordNameStrategy.class.getName()
                    ))
                    .build())
            .build());
  }

  private static Config resolveConfig() throws ConfigurationException {
    var yml = new ObjectMapper(new YAMLFactory());
    var dir = System.getenv("CONFIG_DIR");

    var appConfig = yml.readValue(new File(String.format("%s/em-es/config.yaml", dir)), Config.class);
    var camConfig = yml.readValue(new File(String.format("%s/ccd/config.yaml", dir)), org.blab.v2k.ccd.Config.class)
            .getCams();

    for (var cam : appConfig.getCams().entrySet()) {
      var conf = camConfig.get(cam.getKey());

      if (conf == null) {
        throw new ConfigurationException(String.format("could not find config for %s", cam.getKey()));
      }

      var axes = cam.getValue().getAxes();

      cam.getValue().setConfig(conf);
      axes.getX().setConfig(conf.getAxes().getX());
      axes.getZ().setConfig(conf.getAxes().getZ());
      axes.getX().setChannel(String.format("vepp.ccd.%s.sigma_x", cam.getKey()));
      axes.getZ().setChannel(String.format("vepp.ccd.%s.sigma_z", cam.getKey()));
    }

    return appConfig;
  }
}
