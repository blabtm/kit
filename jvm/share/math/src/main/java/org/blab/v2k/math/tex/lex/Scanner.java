package org.blab.v2k.math.tex.lex;

import java.util.Iterator;

public class Scanner implements Iterator<Token> {
  private final char[] expression;

  private int begin;
  private int current;

  public Scanner(String expression) {
    this.expression = expression.toCharArray();
    this.begin = 0;
    this.current = 0;
  }

  @Override
  public boolean hasNext() {
    return current < expression.length;
  }

  @Override
  public Token next() {
    skip();

    if (!hasNext()) {
      return new Token(Token.Type.EndOfFile, null);
    }

    if (Character.isDigit(expression[current])) {
      return nextNumber();
    }

    if (expression[current] == '\\' || Character.isLetter(expression[current])) {
      return nextSymbol();
    }

    var literal = String.copyValueOf(expression, current, 1);

    return switch (expression[current++]) {
      case '=' -> new Token(Token.Type.Equality, literal);
      case '+' -> new Token(Token.Type.Addition, literal);
      case '-' -> new Token(Token.Type.Subtraction, literal);
      case '*' -> new Token(Token.Type.Multiplication, literal);
      case '^' -> new Token(Token.Type.Power, literal);
      case '(' -> new Token(Token.Type.LeftParenthesis, literal);
      case ')' -> new Token(Token.Type.RightParenthesis, literal);
      case '{' -> new Token(Token.Type.LeftBracket, literal);
      case '}' -> new Token(Token.Type.RightBracket, literal);
      default -> throw new RuntimeException("malformed input");
    };
  }

  private void skip() {
    while (hasNext() && Character.isWhitespace(expression[current])) {
      current += 1;
    }
  }

  private Token nextNumber() {
    begin = current;

    while (hasNext() && Character.isDigit(expression[current])) {
      current += 1;
    }

    var type = Token.Type.IntegerNumber;

    if (hasNext() && expression[current] == '.') {
      type = Token.Type.RealNumber;
      current += 1;

      if (!hasNext() || !Character.isDigit(expression[current])) {
        throw new RuntimeException("malformed input");
      }

      while (hasNext() && Character.isDigit(expression[current])) {
        current += 1;
      }
    }

    return new Token(type, String.copyValueOf(expression, begin, (current - begin)));
  }

  private Token nextSymbol() {
    begin = current++;

    while (hasNext() && (Character.isLetter(expression[current]) || expression[current] == '_')) {
      current += 1;
    }

    if (current - begin == 0) {
      throw new RuntimeException("malformed input");
    }

    var literal = String.copyValueOf(expression, begin, (current - begin));
    var type = Token.Type.Symbol;

    if (Token.Type.KEYWORDS.containsKey(literal)) {
      type = Token.Type.KEYWORDS.get(literal);
    }

    return new Token(type, literal);
  }
}
