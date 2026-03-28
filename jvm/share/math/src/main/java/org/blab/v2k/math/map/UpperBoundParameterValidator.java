package org.blab.v2k.math.map;

import org.apache.commons.math3.util.Pair;

public class UpperBoundParameterValidator implements ParameterValidator {
  private final double bound;

  public UpperBoundParameterValidator(double bound) {
    this.bound = bound;
  }

  @Override
  public Pair<Double, Double> validate(double parameter) {
    return new Pair<>(
            bound + 1 - Math.sqrt(parameter * parameter + 1),
            -parameter / Math.sqrt(parameter * parameter + 1)
    );
  }
}
