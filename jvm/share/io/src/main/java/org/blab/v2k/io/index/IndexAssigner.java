package org.blab.v2k.io.index;

import lombok.RequiredArgsConstructor;
import org.apache.flink.api.common.functions.MapFunction;
import org.blab.v2k.math.Parameter;

@RequiredArgsConstructor
public class IndexAssigner implements MapFunction<Double, Parameter> {
  private final int r;
  private final int c;

  @Override
  public Parameter map(Double v) {
    return new Parameter(v, new Parameter.Index(r, c));
  }
}
