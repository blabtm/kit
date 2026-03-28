package org.blab.v2k.vcas.flink.source.enumerator;

import org.apache.flink.api.connector.source.SplitEnumerator;
import org.apache.flink.api.connector.source.SplitEnumeratorContext;
import org.blab.v2k.vcas.flink.source.split.VcasTopicSplit;

import javax.annotation.Nullable;
import java.io.IOException;
import java.util.*;

public class VcasSourceEnumerator implements SplitEnumerator<VcasTopicSplit, VcasSourceEnumState> {
  private final SplitEnumeratorContext<VcasTopicSplit> context;
  private final VcasSourceEnumState state;

  public VcasSourceEnumerator(SplitEnumeratorContext<VcasTopicSplit> context, Set<String> topics) {
    this.context = context;
    this.state = new VcasSourceEnumState(topics);
  }

  public VcasSourceEnumerator(
          SplitEnumeratorContext<VcasTopicSplit> context, VcasSourceEnumState state) {
    this.context = context;
    this.state = state;
  }

  @Override
  public void start() {
  }

  @Override
  public void handleSplitRequest(int i, @Nullable String s) {
  }

  @Override
  public void addSplitsBack(List<VcasTopicSplit> list, int i) {
    for (VcasTopicSplit split : list) {
      state.assigned.remove(split.getName());
      state.unassigned.offer(split.getName());
    }
  }

  @Override
  public void addReader(int i) {
    for (var t : state.unassigned) {
      context.assignSplit(new VcasTopicSplit(t), i);
    }
  }

  @Override
  public VcasSourceEnumState snapshotState(long l) throws Exception {
    return state;
  }

  @Override
  public void close() throws IOException {
  }
}
