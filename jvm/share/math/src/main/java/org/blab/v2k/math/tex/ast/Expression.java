package org.blab.v2k.math.tex.ast;

import org.blab.v2k.math.tex.lex.Token;

public abstract class Expression {
  public Token t;

  public Expression(Token t) {
    this.t = t;
  }

  public void visit(Visitor v) {
    v.visit(this);
  }
}
