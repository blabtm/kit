package org.blab.v2k.math.map;

import org.apache.commons.math3.util.Pair;

public class SoftplusParameterValidator implements ParameterValidator {
  @Override
  public Pair<Double, Double> validate(double parameter) {
    var exp = Math.exp(parameter);

    return new Pair<>(
            Math.log(1 + exp),
            exp / (1 + exp)
    );
  }
}
