package org.blab.v2k.vcas.consumer;

import java.time.Duration;
import java.util.Properties;
import java.util.Set;

public class ConsumerTest {
  public static void main(String[] args) throws Exception {
    System.setProperty("org.slf4j.simpleLogger.defaultLogLevel", "debug");

    var properties = new Properties();

    properties.put(ConsumerConfig.BROKER_ADDR, "localhost");
    properties.put(ConsumerConfig.BROKER_PORT, 20041);
    properties.put(ConsumerConfig.DB_ADDR, "jdbc:postgresql://localhost:5432/postgres");
    properties.put(ConsumerConfig.DB_USER, "postgres");
    properties.put(ConsumerConfig.DB_PASS, "postgres");

    var consumer = new VcasConsumer(properties);

    consumer.subscribe(Set.of("test.*"));

    System.out.println(consumer.getTopics());

    while (true) {
      var res = consumer.poll(Duration.ofSeconds(1));

      for (var p : res.values()) {
        for (var r : p) {
          System.out.printf("%s: %s\n", r.getTopic(), r.getValue());
        }
      }
    }
  }
}
