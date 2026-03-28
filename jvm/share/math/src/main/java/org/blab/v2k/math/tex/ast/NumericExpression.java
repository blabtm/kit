package org.blab.v2k.math.tex.ast;

import org.blab.v2k.math.tex.lex.Token;

public class NumericExpression extends Expression {
  public int index;

  public NumericExpression(Token t, int index) {
    super(t);
    this.index = index;
  }
}
