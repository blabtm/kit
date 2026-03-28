package org.blab.v2k.math.check;

import org.blab.v2k.math.Problem;

public class RelativeCostConvergenceChecker implements ConvergenceChecker {
  private final double tolerance;

  public RelativeCostConvergenceChecker(double tolerance) {
    this.tolerance = tolerance;
  }

  @Override
  public boolean check(Problem.Evaluation previous, Problem.Evaluation current) {
    if (previous == null) {
      return false;
    }

    var pc = previous.getCost();
    var cc = current.getCost();

    return Math.abs(pc - cc) < tolerance * cc;
  }
}
