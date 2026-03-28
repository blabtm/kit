package org.blab.v2k.math.check;

import org.blab.v2k.math.Problem;

import java.io.Serializable;

public interface ConvergenceChecker extends Serializable {
  boolean check(Problem.Evaluation previous, Problem.Evaluation current);
}
