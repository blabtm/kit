package org.blab.v2k.em_es;

import lombok.experimental.SuperBuilder;
import org.apache.commons.math3.linear.RealVector;

@SuperBuilder
public class Problem extends org.blab.v2k.math.Problem {
  public Problem(org.blab.v2k.math.Problem.ProblemBuilder<?, ?> builder) {
    super(builder);
  }

  @Override
  public int getM() {
    return 2;
  }

  @Override
  public void populate(RealVector src, RealVector dst) {
    dst.setEntry(1, src.getEntry(0));
    dst.setEntry(4, src.getEntry(1));
  }

  @Override
  public void extract(RealVector src, RealVector dst) {
    dst.setEntry(0, src.getEntry(1));
    dst.setEntry(1, src.getEntry(4));
  }
}
