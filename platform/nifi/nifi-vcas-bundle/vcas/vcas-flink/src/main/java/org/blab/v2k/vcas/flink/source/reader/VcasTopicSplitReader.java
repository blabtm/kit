package org.blab.v2k.vcas.flink.source.reader;

import org.apache.flink.connector.base.source.reader.RecordsWithSplitIds;
import org.apache.flink.connector.base.source.reader.splitreader.SplitReader;
import org.apache.flink.connector.base.source.reader.splitreader.SplitsChange;
import org.blab.v2k.vcas.consumer.ConsumerRecord;
import org.blab.v2k.vcas.consumer.VcasConsumer;
import org.blab.v2k.vcas.flink.source.split.VcasTopicSplit;

import javax.annotation.Nullable;
import java.io.IOException;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.ExecutionException;

public class VcasTopicSplitReader implements SplitReader<ConsumerRecord, VcasTopicSplit> {
  private final VcasConsumer consumer;

  public VcasTopicSplitReader(Properties properties) throws Exception {
    try {
      this.consumer = new VcasConsumer(properties);
    } catch (IOException | ExecutionException | InterruptedException e) {
      throw new RuntimeException(e);
    }
  }

  @Override
  public RecordsWithSplitIds<ConsumerRecord> fetch() {
    try {
      return new ConsumerRecordsWithSplits(consumer.poll(Duration.ofSeconds(0)));
    } catch (InterruptedException e) {
      throw new RuntimeException(e);
    }
  }

  @Override
  public void handleSplitsChanges(SplitsChange<VcasTopicSplit> splitsChange) {
    try {
      consumer.subscribe(splitsChange.splits().stream().map(VcasTopicSplit::getName).toList());
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }

  @Override
  public void wakeUp() {
  }

  @Override
  public void close() throws Exception {
    consumer.close();
  }

  private static class ConsumerRecordsWithSplits implements RecordsWithSplitIds<ConsumerRecord> {
    private final Iterator<Map.Entry<String, LinkedList<ConsumerRecord>>> splitIterator;
    private Iterator<ConsumerRecord> recordIterator;

    private ConsumerRecordsWithSplits(Map<String, LinkedList<ConsumerRecord>> records) {
      this.splitIterator = records.entrySet().iterator();
    }

    @Nullable
    @Override
    public String nextSplit() {
      if (splitIterator.hasNext()) {
        var current = splitIterator.next();
        recordIterator = current.getValue().iterator();
        return current.getKey();
      }

      recordIterator = null;

      return null;
    }

    @Nullable
    @Override
    public ConsumerRecord nextRecordFromSplit() {
      if (recordIterator.hasNext()) {
        return recordIterator.next();
      }

      return null;
    }

    @Override
    public Set<String> finishedSplits() {
      return Set.of();
    }
  }
}
