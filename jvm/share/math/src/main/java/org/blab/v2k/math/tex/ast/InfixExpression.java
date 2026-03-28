package org.blab.v2k.math.tex.ast;

import org.blab.v2k.math.tex.lex.Token;

public class InfixExpression extends Expression {
  public Expression x;
  public Expression y;

  public InfixExpression(Token t, Expression x, Expression y) {
    super(t);
    this.x = x;
    this.y = y;
  }
}
