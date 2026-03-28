package org.blab.v2k.math.tex.ast;

import org.blab.v2k.math.tex.lex.Scanner;
import org.blab.v2k.math.tex.lex.Token;

import java.util.Map;

public class Parser {
  private final Scanner scanner;
  private final Map<String, Integer> symbols;
  private final Map<String, Integer> constants;

  private Token cur = null;
  private Token prv = null;
  private boolean buf = false;
  private int idx = 0;

  public Parser(Scanner scanner, Map<String, Integer> symbols, Map<String, Integer> constants) {
    this.scanner = scanner;
    this.symbols = symbols;
    this.constants = constants;
  }

  public Expression parse() {
    var x = parseExpression();

    if (next().type() != Token.Type.Equality) {
      throw new RuntimeException();
    }

    var y = parseExpression();

    return new InfixExpression(new Token(Token.Type.Subtraction, "-"), x, y);
  }

  private Token peek() {
    if (!buf) {
      buf = true;
      prv = cur;
      cur = scanner.next();
    }

    return cur;
  }

  private Token next() {
    if (!buf) {
      prv = cur;
      cur = scanner.next();
    }

    buf = false;

    return cur;
  }

  private Expression parseExpression() {
    return parseExpression0();
  }

  private Expression parseExpression0() {
    var cur = parseExpression1();

    while (true) {
      var tok = peek();

      switch (tok.type()) {
        case Addition, Subtraction:
          tok = next();
          cur = new InfixExpression(tok, cur, parseExpression1());
          continue;
      }

      break;
    }

    return cur;
  }

  private Expression parseExpression1() {
    var cur = parseExpression2();

    while (true) {
      var tok = peek();

      switch (tok.type()) {
        case Multiplication, Power:
          tok = next();

          var x = cur;
          var y = parseExpression1();

          cur = new InfixExpression(tok, x, y);

          continue;
      }

      break;
    }

    return cur;
  }

  private Expression parseExpression2() {
    var tok = peek();

    if (tok.type() == Token.Type.Subtraction) {
      return new PrefixExpression(next(), parseExpression2());
    }

    return parseExpression3();
  }

  private Expression parseExpression3() {
    var tok = peek();

    switch (tok.type()) {
      case LeftParenthesis:
        var res = new GroupExpression(next(), parseExpression0());

        if (next().type() != Token.Type.RightParenthesis) {
          throw new RuntimeException();
        }

        return res;
      case Division:
        return new InfixExpression(
                next(),
                parseExpression4(),
                parseExpression4()
        );
      case Symbol:
        var s = next();

        if (!symbols.containsKey(s.literal())) {
          symbols.put(s.literal(), symbols.size());
        }

        return new VariableExpression(s, symbols.get(s.literal()));
      case IntegerNumber, RealNumber:
        var c = next();

        if (!constants.containsKey(c.literal())) {
          constants.put(c.literal(), constants.size());
        }

        return new NumericExpression(c, constants.get(c.literal()));
      default:
        throw new RuntimeException();
    }
  }

  private Expression parseExpression4() {
    if (next().type() != Token.Type.LeftBracket) {
      throw new RuntimeException();
    }

    var res = new GroupExpression(cur, parseExpression0());

    if (next().type() != Token.Type.RightBracket) {
      throw new RuntimeException();
    }

    return res;
  }
}
