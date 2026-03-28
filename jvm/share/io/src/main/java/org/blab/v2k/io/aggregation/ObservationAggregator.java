package org.blab.v2k.io.aggregation;

import lombok.RequiredArgsConstructor;
import org.apache.commons.math3.linear.ArrayRealVector;
import org.apache.commons.math3.linear.RealMatrix;
import org.apache.commons.math3.linear.RealVector;
import org.apache.flink.api.common.functions.AggregateFunction;
import org.apache.flink.api.java.tuple.Tuple2;
import org.blab.v2k.math.Parameter;

@RequiredArgsConstructor
public class ObservationAggregator implements AggregateFunction<Parameter, Tuple2<RealVector, RealMatrix>, Tuple2<RealVector, RealMatrix>> {
  private final RealMatrix base;
  private final Integer attachmentSize;

  @Override
  public Tuple2<RealVector, RealMatrix> createAccumulator() {
    return new Tuple2<>(new ArrayRealVector(attachmentSize), base.copy());
  }

  @Override
  public Tuple2<RealVector, RealMatrix> add(Parameter parameter, Tuple2<RealVector, RealMatrix> accumulator) {
    int r = parameter.getIndex().getRow();
    int c = parameter.getIndex().getCol();

    if (r == -1) {
      accumulator.f0.setEntry(c, parameter.getValue());
    } else {
      accumulator.f1.setEntry(r, c, parameter.getValue());
    }

    return accumulator;
  }

  @Override
  public Tuple2<RealVector, RealMatrix> getResult(Tuple2<RealVector, RealMatrix> accumulator) {
    return accumulator;
  }

  @Override
  public Tuple2<RealVector, RealMatrix> merge(Tuple2<RealVector, RealMatrix> a, Tuple2<RealVector, RealMatrix> b) {
    return a;
  }
}
