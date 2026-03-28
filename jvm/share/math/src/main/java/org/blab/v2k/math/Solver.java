package org.blab.v2k.math;

import lombok.experimental.SuperBuilder;
import org.apache.commons.math3.util.IntegerSequence;
import org.blab.v2k.math.check.ConvergenceChecker;
import org.blab.v2k.math.check.CostConvergenceChecker;
import org.blab.v2k.math.check.OneOfConvergenceChecker;
import org.blab.v2k.math.check.RelativeCostConvergenceChecker;

import javax.naming.ConfigurationException;
import java.io.Serializable;

@SuperBuilder(toBuilder = true)
public abstract class Solver {
  protected IntegerSequence.Incrementor evaluationCounter;
  protected IntegerSequence.Incrementor iterationCounter;
  protected ConvergenceChecker convergenceChecker;

  public abstract Problem.Evaluation solve(Problem problem);

  public static abstract class SolverBuilder<C extends Solver, B extends SolverBuilder<C, B>> {
    public B fromConfig(Config config) throws ConfigurationException {
      this.evaluationCounter = IntegerSequence.Incrementor.create()
              .withMaximalCount(config.getMaxEvaluations());

      this.iterationCounter = IntegerSequence.Incrementor.create()
              .withMaximalCount(config.getMaxIterations());

      var checker = new OneOfConvergenceChecker();

      if (config.absoluteCost != null) {
        checker = checker.withChecker(new CostConvergenceChecker(config.absoluteCost));
      }

      if (config.relativeCost != null) {
        checker = checker.withChecker(new RelativeCostConvergenceChecker(config.relativeCost));
      }

      if (checker.size() == 0) {
        throw new ConfigurationException("at least one checker required, but zero given");
      }

      this.convergenceChecker = checker;

      return self();
    }
  }

  @lombok.Getter
  @lombok.Setter
  public static class Config implements Serializable {
    protected Integer maxEvaluations;
    protected Integer maxIterations;
    protected Double absoluteCost;
    protected Double relativeCost;
  }
}
