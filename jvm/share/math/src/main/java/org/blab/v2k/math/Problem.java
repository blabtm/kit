package org.blab.v2k.math;

import lombok.Getter;
import lombok.Setter;
import lombok.ToString;
import lombok.experimental.SuperBuilder;
import org.apache.commons.math3.linear.Array2DRowRealMatrix;
import org.apache.commons.math3.linear.ArrayRealVector;
import org.apache.commons.math3.linear.RealMatrix;
import org.apache.commons.math3.linear.RealVector;
import org.blab.v2k.math.auto.AutoProblemBuilder;
import org.blab.v2k.math.map.MultivariateParameterValidator;

import javax.tools.*;
import java.io.Serializable;

@Getter
@SuperBuilder
public abstract class Problem implements Serializable {
  protected final Model model;
  protected final RealVector initialGuess;
  protected RealMatrix observation;

  protected MultivariateParameterValidator parameterValidator;
  protected RealVector observationWeight;

  public Problem(ProblemBuilder<?, ?> builder) {
    this.model = builder.model;
    this.observation = builder.observation;
    this.parameterValidator = builder.parameterValidator;
    this.observationWeight = builder.observationWeight;
    this.initialGuess = builder.initialGuess;

    if (this.observationWeight == null) {
      this.observationWeight = new ArrayRealVector(getObservationSize(), 1.0);
    }
  }

  public Problem withObservation(RealMatrix observation) {
    this.observation = observation;
    return this;
  }

  public static AutoProblemBuilder auto() {
    return new AutoProblemBuilder();
  }

  /**
   * @return Number of unknowns.
   */
  public abstract int getM();

  public abstract void populate(RealVector src, RealVector dst);

  public abstract void extract(RealVector src, RealVector dst);

  public Evaluation evaluate(RealVector parameters) {
    var ev = new Evaluation();

    ev.residual = new ArrayRealVector(getObservationSize());
    ev.jacobian = new Array2DRowRealMatrix(getObservationSize(), getM());

    for (int i = 0; i < observation.getRowDimension(); ++i) {
      var v = observation.getRowVector(i);

      populate(parameters, v);

      var m = model.evaluate(v);
      var r = m.getFirst();
      var j = m.getSecond();

      ev.residual.setSubVector(i * model.getM(), r);

      for (int n = 0; n < model.getM(); ++n) {
        var t = new ArrayRealVector(getM());

        extract(j.getRowVector(n), t);
        ev.jacobian.setRowVector(i * model.getM() + n, t);
      }
    }

    ev.point = parameters;

    return ev;
  }

  public int getObservationSize() {
    return observation.getRowDimension() * model.getM();
  }

  @Getter
  @Setter
  @ToString
  public static class Evaluation {
    private RealVector point;
    private RealVector residual;
    private RealMatrix jacobian;

    private int iterations;
    private int evaluations;

    public double getCost() {
      return residual.getNorm();
    }
  }
}
