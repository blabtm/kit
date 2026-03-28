package org.blab.v2k.em_es;

import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

public class Application {
  public static void main(String[] args) throws Exception {
    try (var environment = StreamExecutionEnvironment.getExecutionEnvironment()) {
      Estimator.create(environment);
      environment.execute();
    }
  }
}
