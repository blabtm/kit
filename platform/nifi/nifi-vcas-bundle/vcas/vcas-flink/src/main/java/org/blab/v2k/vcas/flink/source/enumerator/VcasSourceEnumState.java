package org.blab.v2k.vcas.flink.source.enumerator;

import java.util.HashSet;
import java.util.LinkedList;
import java.util.Queue;
import java.util.Set;

public class VcasSourceEnumState {
  Queue<String> unassigned;
  Set<String> assigned;

  public VcasSourceEnumState() {
    unassigned = new LinkedList<>();
    assigned = new HashSet<>();
  }

  public VcasSourceEnumState(Set<String> topics) {
    unassigned = new LinkedList<>(topics);
    assigned = new HashSet<>();
  }
}
