package org.blab.v2k.math.tex.ast;

import org.blab.v2k.math.tex.lex.Token;

public class VariableExpression extends Expression {
  public int index;

  public VariableExpression(Token t, int index) {
    super(t);
    this.index = index;
  }
}
