package org.blab.v2k.vcas.flink.source.reader;

import org.apache.flink.api.connector.source.SourceReaderContext;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.connector.base.source.reader.RecordEmitter;
import org.apache.flink.connector.base.source.reader.SingleThreadMultiplexSourceReaderBase;
import org.apache.flink.connector.base.source.reader.splitreader.SplitReader;
import org.blab.v2k.vcas.consumer.ConsumerRecord;
import org.blab.v2k.vcas.flink.source.split.VcasTopicSplit;

import java.util.Map;
import java.util.function.Supplier;

public class VcasSourceReader<T>
        extends SingleThreadMultiplexSourceReaderBase<ConsumerRecord, T, VcasTopicSplit, VcasTopicSplit> {
  public VcasSourceReader(
          Supplier<SplitReader<ConsumerRecord, VcasTopicSplit>> splitReaderSupplier,
          RecordEmitter<ConsumerRecord, T, VcasTopicSplit> recordEmitter,
          Configuration config,
          SourceReaderContext context) {
    super(splitReaderSupplier, recordEmitter, config, context);
  }

  @Override
  protected void onSplitFinished(Map<String, VcasTopicSplit> map) {
  }

  @Override
  protected VcasTopicSplit initializedState(VcasTopicSplit split) {
    return split;
  }

  @Override
  protected VcasTopicSplit toSplitType(String s, VcasTopicSplit split) {
    return split;
  }
}
