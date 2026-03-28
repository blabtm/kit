package org.blab.v2k.math.tex.ast;

import org.blab.v2k.math.tex.lex.Token;

public class GroupExpression extends Expression {
  public Expression x;

  public GroupExpression(Token t, Expression x) {
    super(t);
    this.x = x;
  }
}
