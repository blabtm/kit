package org.blab.v2k.math.root;

import org.apache.commons.math3.exception.DimensionMismatchException;
import org.apache.commons.math3.fitting.leastsquares.MultivariateJacobianFunction;
import org.apache.commons.math3.linear.Array2DRowRealMatrix;
import org.apache.commons.math3.linear.LUDecomposition;
import org.apache.commons.math3.linear.RealVector;
import org.apache.commons.math3.util.IntegerSequence;
import org.blab.v2k.math.map.MultivariateParameterValidator;

import java.util.List;

public class NewtonSolver {
  private final MultivariateJacobianFunction function;

  private double tolerance;
  private IntegerSequence.Incrementor iterationCounter;
  private MultivariateParameterValidator validator;

  public NewtonSolver(MultivariateJacobianFunction function) {
    this.function = function;
    this.tolerance = 1e-10;
    this.iterationCounter = IntegerSequence.Incrementor.create()
            .withMaximalCount(1000);
  }

  public NewtonSolver withMaxIterations(int max) {
    this.iterationCounter = this.iterationCounter.withMaximalCount(max);
    return this;
  }

  public NewtonSolver withTolerance(double tolerance) {
    this.tolerance = tolerance;
    return this;
  }

  public NewtonSolver withParameterValidator(MultivariateParameterValidator validator) {
    this.validator = validator;
    return this;
  }

  public RealVector solve(List<Integer> p, RealVector s) {
    var nj = new Array2DRowRealMatrix(p.size(), p.size());

    RealVector mp = s.copy(); // mapping parameter
    RealVector mx = null; // mapped vector
    RealVector mj = null; // mapped jacobian

    while (true) {
      iterationCounter.increment();

      if (validator != null) {
        var t = validator.validate(mp);

        mx = t.getFirst();
        mj = t.getSecond();
      } else {
        mx = mp;
      }

      var f = function.value(mx);
      var v = f.getFirst();

      if (p.size() != v.getDimension()) {
        throw new DimensionMismatchException(p.size(), v.getDimension());
      }

      for (int i = 0; i < p.size(); ++i) {
        for (int j = 0; j < p.size(); ++j) {
          var der = f.getSecond().getEntry(i, p.get(j));

          if (mj != null) {
            der *= mj.getEntry(p.get(j));
          }

          nj.setEntry(i, j, der);
        }
      }

      var nd = new LUDecomposition(nj)
              .getSolver()
              .solve(f.getFirst());

      if (nd.getNorm() < tolerance) {
        break;
      }

      for (int i = 0; i < p.size(); ++i) {
        mp.addToEntry(p.get(i), -nd.getEntry(i));
      }
    }

    if (validator != null) {
      mx = validator.validate(mp).getFirst();
    } else {
      mx = mp;
    }

    return mx;
  }
}
