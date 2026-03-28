package org.blab.v2k.vcas.consumer.schema;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class DoubleDeserializationSchema implements RecordDeserializationSchema<Double> {
  @Override
  public Double deserialize(String v) {
    try {
      return Double.parseDouble(v);
    } catch (Exception e) {
      log.error("could not parse value as double", e);
      return 0.0;
    }
  }
}
