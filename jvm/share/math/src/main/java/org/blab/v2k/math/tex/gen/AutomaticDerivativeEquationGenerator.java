package org.blab.v2k.math.tex.gen;


import org.blab.v2k.math.tex.ast.*;

public class AutomaticDerivativeEquationGenerator implements Visitor {
  private final StringBuilder builder;

  public AutomaticDerivativeEquationGenerator() {
    this.builder = new StringBuilder();
  }

  @Override
  public void visit(Expression e) {
    switch (e) {
      case VariableExpression v -> visitVariableExpression(v);
      case NumericExpression n -> visitNumericExpression(n);
      case InfixExpression i -> visitInfixExpression(i);
      case PrefixExpression p -> visitPrefixExpression(p);
      case GroupExpression g -> visitGroupExpression(g);
      default -> throw new RuntimeException();
    }
  }

  private void visitVariableExpression(VariableExpression e) {
    builder.append('(').append('p').append(e.index).append(')');
  }

  private void visitNumericExpression(NumericExpression e) {
    builder.append('(').append("c").append(e.index).append(')');
  }

  private void visitInfixExpression(InfixExpression e) {
    builder.append('(');

    e.x.visit(this);

    switch (e.t.type()) {
      case Addition -> builder.append(".add").append('(');
      case Subtraction -> builder.append(".subtract").append('(');
      case Multiplication -> builder.append(".multiply").append('(');
      case Division -> builder.append(".divide").append('(');
      case Power -> builder.append(".pow").append('(');
      default -> throw new RuntimeException();
    }

    e.y.visit(this);

    builder.append(')').append(')');
  }

  private void visitPrefixExpression(PrefixExpression e) {
    builder.append('(');

    e.x.visit(this);

    switch (e.t.type()) {
      case Subtraction -> builder.append(".negate()");
      default -> throw new RuntimeException();
    }

    builder.append(')');
  }

  private void visitGroupExpression(GroupExpression e) {
    builder.append('(');
    e.x.visit(this);
    builder.append(')');
  }

  public String collect() {
    return builder.toString();
  }
}
