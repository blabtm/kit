package org.blab.v2k.math;

import lombok.experimental.SuperBuilder;
import org.apache.commons.math3.linear.Array2DRowRealMatrix;
import org.apache.commons.math3.linear.ArrayRealVector;
import org.apache.commons.math3.linear.RealMatrix;
import org.apache.commons.math3.linear.RealVector;
import org.apache.commons.math3.util.IntegerSequence;
import org.apache.commons.math3.util.Pair;
import org.blab.v2k.math.check.CostConvergenceChecker;
import org.blab.v2k.math.check.OneOfConvergenceChecker;
import org.blab.v2k.math.check.RelativeCostConvergenceChecker;
import org.blab.v2k.math.map.CompoundParameterValidator;
import org.blab.v2k.math.map.ExponentialParameterValidator;
import org.blab.v2k.math.optimization.GaussNewtonOptimizer;

class TestModel extends org.blab.v2k.math.Model {
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


@SuperBuilder
class TestProblem extends Problem {
  public TestProblem(Problem.ProblemBuilder<?, ?> builder) {
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

public class GaussNewtonOptimizerTest {
  public static void main(String[] args) {
    var obs = new Array2DRowRealMatrix(new double[][]{
            {0.0221291412, 0.0, 4.3573654, 0.0008886465, 0.0},
            {0.0247131933, 0.0, 5.4344124, -0.0, 0.0},
            {0.0140476865, 0.0, 1.7553643, 0.37415184, 0.0},
            {0.0170902865, 0.0, 2.5989238, 0.0, 0.0}
    });

    var problem = TestProblem.builder()
            .model(new TestModel())
            .observation(obs)
            .initialGuess(new ArrayRealVector(new double[]{1, 1}))
            .parameterValidator(new CompoundParameterValidator()
                    .withValidator(0, new ExponentialParameterValidator())
                    .withValidator(1, new ExponentialParameterValidator()))
            .build();

    var opt = GaussNewtonOptimizer.builder()
            .evaluationCounter(IntegerSequence.Incrementor.create().withMaximalCount(1000))
            .iterationCounter(IntegerSequence.Incrementor.create().withMaximalCount(1000))
            .convergenceChecker(new OneOfConvergenceChecker()
                    .withChecker(new CostConvergenceChecker(1e-7))
                    .withChecker(new RelativeCostConvergenceChecker(1e-5)))
            .build().solve(problem);

    System.out.println(opt);
  }
}
