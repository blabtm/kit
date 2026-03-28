package org.blab.v2k.math.map;

import org.apache.commons.math3.util.Pair;

public class ExponentialParameterValidator implements ParameterValidator {
  @Override
  public Pair<Double, Double> validate(double parameter) {
    return new Pair<>(
            Math.exp(parameter),
            Math.exp(parameter)
    );
  }
}
