package org.blab.v2k.em_es;

import org.apache.commons.math3.linear.Array2DRowRealMatrix;
import org.apache.commons.math3.linear.ArrayRealVector;
import org.apache.commons.math3.linear.RealMatrix;
import org.apache.commons.math3.linear.RealVector;
import org.apache.commons.math3.util.Pair;

public class Model extends org.blab.v2k.math.Model {
  @Override
  public int getN() {
    return 5;
  }

  @Override
  public int getM() {
    return 1;
  }

  @Override
  public Pair<RealVector, RealMatrix> evaluate(RealVector p) {
    var r = new ArrayRealVector(getM());
    var j = new Array2DRowRealMatrix(getM(), getN());

    r.setEntry(0, Math.pow(p.getEntry(0), 2) - p.getEntry(1) * p.getEntry(2) - Math.pow(p.getEntry(3) * p.getEntry(4), 2));

    j.setEntry(0, 0, 2 * p.getEntry(0));
    j.setEntry(0, 1, -p.getEntry(2));
    j.setEntry(0, 2, -p.getEntry(1));
    j.setEntry(0, 3, -2 * (p.getEntry(3) * p.getEntry(4)) * p.getEntry(4));
    j.setEntry(0, 4, -2 * (p.getEntry(3) * p.getEntry(4)) * p.getEntry(3));

    return new Pair<>(r, j);
  }
}
