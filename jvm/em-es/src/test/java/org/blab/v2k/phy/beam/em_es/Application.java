package org.blab.v2k.phy.beam.em_es;

import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

public class Application {
  public static void main(String[] args) throws Exception {
    try (var environment = StreamExecutionEnvironment.createLocalEnvironment()) {
      environment.setParallelism(1);
      Estimator.create(environment);
      environment.execute();
    }
  }
}
