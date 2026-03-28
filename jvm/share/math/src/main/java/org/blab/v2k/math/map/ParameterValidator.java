package org.blab.v2k.math.map;

import org.apache.commons.math3.util.Pair;

import java.io.Serializable;

public interface ParameterValidator extends Serializable {
  Pair<Double, Double> validate(double parameter);
}
