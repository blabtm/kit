package org.blab.v2k.vcas.flink.source.reader;

import org.apache.flink.api.connector.source.SourceOutput;
import org.apache.flink.connector.base.source.reader.RecordEmitter;
import org.blab.v2k.vcas.consumer.ConsumerRecord;
import org.blab.v2k.vcas.consumer.schema.RecordDeserializationSchema;
import org.blab.v2k.vcas.flink.source.split.VcasTopicSplit;

import java.io.IOException;

public class VcasRecordEmitter<T> implements RecordEmitter<ConsumerRecord, T, VcasTopicSplit> {
  private final RecordDeserializationSchema<T> schema;

  public VcasRecordEmitter(RecordDeserializationSchema<T> schema) {
    this.schema = schema;
  }

  @Override
  public void emitRecord(ConsumerRecord record, SourceOutput<T> output, VcasTopicSplit split) throws IOException {
    output.collect(schema.deserialize(record.getValue()), record.getTimestamp());
  }
}
