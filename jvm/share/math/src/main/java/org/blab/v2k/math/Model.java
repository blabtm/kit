package org.blab.v2k.math;

import lombok.Getter;
import org.apache.commons.math3.linear.RealMatrix;
import org.apache.commons.math3.linear.RealVector;
import org.apache.commons.math3.util.Pair;

import java.io.Serializable;

@Getter
public abstract class Model implements Serializable {
  /**
   * @return Number of parameters.
   */
  public abstract int getN();

  /**
   * @return Number of equations.
   */
  public abstract int getM();

  public abstract Pair<RealVector, RealMatrix> evaluate(RealVector p);

}
