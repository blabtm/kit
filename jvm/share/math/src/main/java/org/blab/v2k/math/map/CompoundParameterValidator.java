package org.blab.v2k.math.map;

import lombok.AllArgsConstructor;
import lombok.Getter;
import org.apache.commons.math3.linear.ArrayRealVector;
import org.apache.commons.math3.linear.RealVector;
import org.apache.commons.math3.util.Pair;

import java.io.Serializable;
import java.util.LinkedList;
import java.util.List;

@AllArgsConstructor
@Getter
class Validator implements Serializable {
  private final Integer index;
  private final ParameterValidator validator;
}

public class CompoundParameterValidator implements MultivariateParameterValidator {
  private final List<Validator> validators = new LinkedList<>();

  public CompoundParameterValidator withValidator(int parameter, ParameterValidator validator) {
    validators.add(new Validator(parameter, validator));
    return this;
  }

  @Override
  public Pair<RealVector, RealVector> validate(RealVector parameters) {
    var mp = parameters.copy();
    var mj = new ArrayRealVector(parameters.getDimension(), 1);

    for (var p : validators) {
      var i = p.getIndex();
      var v = p.getValidator().validate(parameters.getEntry(p.getIndex()));

      mp.setEntry(i, v.getFirst());
      mj.setEntry(i, v.getSecond());
    }

    return new Pair<>(mp, mj);
  }
}
