package org.blab.v2k.vcas.flink.source;

import org.apache.flink.api.connector.source.*;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.core.io.SimpleVersionedSerializer;
import org.blab.v2k.vcas.consumer.schema.RecordDeserializationSchema;
import org.blab.v2k.vcas.flink.source.enumerator.VcasSourceEnumState;
import org.blab.v2k.vcas.flink.source.enumerator.VcasSourceEnumStateSerializer;
import org.blab.v2k.vcas.flink.source.enumerator.VcasSourceEnumerator;
import org.blab.v2k.vcas.flink.source.reader.VcasRecordEmitter;
import org.blab.v2k.vcas.flink.source.reader.VcasSourceReader;
import org.blab.v2k.vcas.flink.source.reader.VcasTopicSplitReader;
import org.blab.v2k.vcas.flink.source.split.VcasTopicSplit;
import org.blab.v2k.vcas.flink.source.split.VcasTopicSplitSerializer;

import java.util.Properties;
import java.util.Set;

public class VcasSource<T> implements Source<T, VcasTopicSplit, VcasSourceEnumState> {
  private final Properties properties;
  private final Set<String> topics;
  private final RecordDeserializationSchema<T> schema;

  public VcasSource(
          Properties properties, Set<String> topics, RecordDeserializationSchema<T> schema) {
    this.properties = properties;
    this.topics = topics;
    this.schema = schema;
  }

  @Override
  public Boundedness getBoundedness() {
    return Boundedness.CONTINUOUS_UNBOUNDED;
  }

  @Override
  public SplitEnumerator<VcasTopicSplit, VcasSourceEnumState> createEnumerator(
          SplitEnumeratorContext<VcasTopicSplit> context) throws Exception {
    return new VcasSourceEnumerator(context, topics);
  }

  @Override
  public SplitEnumerator<VcasTopicSplit, VcasSourceEnumState> restoreEnumerator(
          SplitEnumeratorContext<VcasTopicSplit> context, VcasSourceEnumState state) throws Exception {
    return new VcasSourceEnumerator(context, state);
  }

  @Override
  public SimpleVersionedSerializer<VcasSourceEnumState> getEnumeratorCheckpointSerializer() {
    return new VcasSourceEnumStateSerializer();
  }

  @Override
  public SimpleVersionedSerializer<VcasTopicSplit> getSplitSerializer() {
    return new VcasTopicSplitSerializer();
  }

  @Override
  public SourceReader<T, VcasTopicSplit> createReader(SourceReaderContext context) {
    return new VcasSourceReader<T>(
            () -> {
              try {
                return new VcasTopicSplitReader(properties);
              } catch (Exception e) {
                throw new RuntimeException(e);
              }
            },
            new VcasRecordEmitter<>(schema),
            new Configuration(),
            context);
  }
}
