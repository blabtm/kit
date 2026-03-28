package org.blab.v2k.math.tex.lex;

import java.util.Map;

public record Token(Type type, String literal) {
  public enum Type {
    EndOfFile,
    Symbol,
    IntegerNumber,
    RealNumber,
    Equality,
    Addition,
    Subtraction,
    Multiplication,
    Division,
    Power,
    LeftParenthesis,
    RightParenthesis,
    LeftBracket,
    RightBracket;

    public static final Map<String, Type> KEYWORDS = Map.ofEntries(
            Map.entry("\\frac", Division)
    );
  }
}
