package org.blab.v2k.math.tex.ast;

import org.blab.v2k.math.tex.lex.Token;

public class PrefixExpression extends Expression {
  public Expression x;

  public PrefixExpression(Token t, Expression x) {
    super(t);
    this.x = x;
  }
}
