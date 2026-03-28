package org.blab.v2k.math.map;

import org.apache.commons.math3.linear.RealVector;
import org.apache.commons.math3.util.Pair;

import java.io.Serializable;

public interface MultivariateParameterValidator extends Serializable {
  Pair<RealVector, RealVector> validate(RealVector parameters);
}
