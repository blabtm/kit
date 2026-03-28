package org.blab.v2k.math.check;

import org.blab.v2k.math.Problem;

import java.util.LinkedList;
import java.util.List;

public class OneOfConvergenceChecker implements ConvergenceChecker {
  private final List<ConvergenceChecker> checkers = new LinkedList<>();

  public OneOfConvergenceChecker withChecker(ConvergenceChecker checker) {
    checkers.add(checker);
    return this;
  }

  public int size() {
    return checkers.size();
  }

  @Override
  public boolean check(Problem.Evaluation previous, Problem.Evaluation current) {
    for (var chk : checkers) {
      if (chk.check(previous, current)) {
        return true;
      }
    }

    return false;
  }
}
