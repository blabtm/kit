package org.blab.v2k.io.aggregation;

import org.apache.flink.api.common.functions.AggregateFunction;
import org.apache.flink.api.java.tuple.Tuple2;
import org.blab.v2k.math.Parameter;

public class AverageAggregator implements AggregateFunction<Parameter, Tuple2<Parameter, Integer>, Parameter> {
  @Override
  public Tuple2<Parameter, Integer> createAccumulator() {
    return new Tuple2<>(null, 0);
  }

  @Override
  public Tuple2<Parameter, Integer> add(Parameter parameter, Tuple2<Parameter, Integer> a) {
    if (a.f0 == null) {
      a.f0 = parameter;
      a.f1 = 1;
    } else {
      a.f0.add(parameter.getValue());
      a.f1 += 1;
    }

    return a;
  }

  @Override
  public Parameter getResult(Tuple2<Parameter, Integer> a) {
    a.f0.setValue(a.f0.getValue() / a.f1);
    return a.f0;
  }

  @Override
  public Tuple2<Parameter, Integer> merge(Tuple2<Parameter, Integer> a, Tuple2<Parameter, Integer> b) {
    a.f0.add(b.f0.getValue());
    a.f1 += b.f1;

    return a;
  }
}
