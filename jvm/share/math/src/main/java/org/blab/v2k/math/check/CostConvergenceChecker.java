package org.blab.v2k.math.check;

import org.blab.v2k.math.Problem;

public class CostConvergenceChecker implements ConvergenceChecker {
  private final double tolerance;

  public CostConvergenceChecker(double tolerance) {
    this.tolerance = tolerance;
  }

  @Override
  public boolean check(Problem.Evaluation previous, Problem.Evaluation current) {
    return Math.abs(current.getCost()) < tolerance;
  }
}
