package org.blab.v2k.math.map;

import org.apache.commons.math3.util.Pair;

public class SquareParameterValidator implements ParameterValidator {
  @Override
  public Pair<Double, Double> validate(double parameter) {
    return new Pair<>(
            parameter * parameter,
            2 * parameter
    );
  }
}
