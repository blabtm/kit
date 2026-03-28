package org.blab.v2k.math.optimization;

import lombok.experimental.SuperBuilder;
import org.apache.commons.math3.analysis.function.Power;
import org.apache.commons.math3.exception.NullArgumentException;
import org.apache.commons.math3.linear.*;
import org.blab.v2k.math.Problem;
import org.blab.v2k.math.Solver;

@SuperBuilder
public class GaussNewtonOptimizer extends Solver {
  private static final double REGULARIZATION = 1e-10;

  public Problem.Evaluation solve(Problem problem) {
    var ic = iterationCounter;
    var cc = convergenceChecker;
    var pv = problem.getParameterValidator();
    var ow = problem.getObservationWeight();

    if (cc == null) {
      throw new NullArgumentException();
    }

    var weight = ow.map(new Power(2));

    var m = new Array2DRowRealMatrix(problem.getM(), problem.getM());
    var b = new ArrayRealVector(problem.getM());

    RealVector mp = problem.getInitialGuess().copy(); // mapping parameters
    RealVector mx = null; // mapped vector
    RealVector mj = null; // mapped jacobian

    Problem.Evaluation prv = null;
    Problem.Evaluation cur = null;

    while (true) {
      ic.increment();

      if (pv != null) {
        var v = pv.validate(mp);

        mx = v.getFirst();
        mj = v.getSecond();
      } else {
        mx = mp;
      }

      prv = cur;
      cur = problem.evaluate(mx);

      if (cc.check(prv, cur)) {
        cur.setIterations(ic.getCount());
        return cur;
      }

      var res = cur.getResidual();
      var jac = cur.getJacobian();

      if (mj != null) {
        for (int i = 0; i < jac.getRowDimension(); ++i) {
          for (int j = 0; j < jac.getColumnDimension(); ++j) {
            if (Double.isNaN(jac.getEntry(i, j)))
              jac.setEntry(i, j, 0.0);

            jac.multiplyEntry(i, j, mj.getEntry(j));
          }
        }
      }

      for (int q = 0; q < problem.getM(); ++q) {
        for (int s = 0; s < problem.getM(); ++s) {
          m.setEntry(q, s, 0);

          for (int i = 0; i < problem.getObservationSize(); ++i) {
            m.addToEntry(q, s, weight.getEntry(i) * jac.getEntry(i, q) 
              * jac.getEntry(i, s));
          }
        }

        b.setEntry(q, 0);

        for (int i = 0; i < problem.getObservationSize(); ++i) {
          b.addToEntry(q, -weight.getEntry(i) * res.getEntry(i) 
            * jac.getEntry(i, q));
        }
      }

      RealVector d = null;

      try {
        d = new LUDecomposition(m)
                .getSolver()
                .solve(b);
      } catch (SingularMatrixException ignored) {
        for (int i = 0; i < m.getRowDimension(); ++i) {
          m.addToEntry(i, i, REGULARIZATION);
        }

        d = new LUDecomposition(m)
                .getSolver()
                .solve(b);
      }

      mp = mp.add(d);
    }
  }
}
