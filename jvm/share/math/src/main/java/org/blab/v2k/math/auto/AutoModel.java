package org.blab.v2k.math.auto;

import org.blab.v2k.math.tex.ast.Parser;
import org.blab.v2k.math.tex.gen.AutomaticDerivativeEquationGenerator;
import org.blab.v2k.math.tex.lex.Scanner;
import org.blab.v2k.math.Model;

import javax.tools.*;
import java.net.URI;
import java.util.HashMap;
import java.util.List;

public class AutoModel {
  private static int counter = 0;

  public static Model compile(List<String> e) throws Exception {
    var m = e.size();
    var n = new HashMap<String, Integer>();
    var c = new HashMap<String, Integer>();

    var pfx = new StringBuilder();
    var exp = new StringBuilder();
    var col = new StringBuilder();

    for (int i = 0; i < m; ++i) {
      var eqn = e.get(i);
      var gen = new AutomaticDerivativeEquationGenerator();

      new Parser(new Scanner(eqn), n, c)
              .parse()
              .visit(gen);

      exp.append(String.format("var f%d = %s;\n", i, gen.collect()));
      col.append(String.format("v.setEntry(%d, f%d.getValue());\n", i, i));
    }

    for (var ce : c.entrySet()) {
      pfx.append(String.format("var c%d = new DerivativeStructure(%d, %d, %s);\n",
              ce.getValue(), n.size(), 1, ce.getKey()));
    }

    for (var ne : n.entrySet()) {
      pfx.append(String.format("var p%d = new DerivativeStructure(%d, %d, %d, p.getEntry(%d));\n",
              ne.getValue(), n.size(), 1, ne.getValue(), ne.getValue()));
    }

    for (int i = 0; i < m; ++i) {
      for (int j = 0; j < n.size(); ++j) {
        col.append(String.format("j.setEntry(%d, %d, f%d.getPartialDerivative(", i, j, i));

        for (int p = 0; p < j; ++p) {
          col.append('0').append(',');
        }

        col.append('1').append(',');

        for (int p = j + 1; p < n.size(); ++p) {
          col.append('0').append(',');
        }

        col.deleteCharAt(col.length() - 1);
        col.append("));\n");
      }
    }

    var src = String.format("""
                      package org.blab.math;
            
                      import org.apache.commons.math3.analysis.differentiation.DerivativeStructure;
                      import org.apache.commons.math3.linear.Array2DRowRealMatrix;
                      import org.apache.commons.math3.linear.ArrayRealVector;
                      import org.apache.commons.math3.linear.RealMatrix;
                      import org.apache.commons.math3.linear.RealVector;
                      import org.apache.commons.math3.util.Pair;
                      import java.util.List;
                      import java.util.ArrayList;
            
                      public class Model%d extends Model {
                        @Override
                        public int getN() {
                          return %d;
                        }
            
                        @Override
                        public int getM() {
                          return %d;
                        }
            
                        @Override
                        public Pair<RealVector, RealMatrix> evaluate(RealVector p) {
                          var v = new ArrayRealVector(getM());
                          var j = new Array2DRowRealMatrix(getM(), getN());
            
                          %s
            
                          %s
            
                          %s
            
                          return new Pair<>(v, j);
                        }
                      }
            """, counter, n.size(), m, pfx, exp, col);

    JavaCompiler compiler = ToolProvider.getSystemJavaCompiler();
    SimpleJavaFileObject file = new SimpleJavaFileObject(
            URI.create(String.format("string:///org/blab/math/Model%d.java", counter)),
            JavaFileObject.Kind.SOURCE
    ) {
      @Override
      public CharSequence getCharContent(boolean ignoreEncodingErrors) {
        return src;
      }
    };

    DiagnosticCollector<JavaFileObject> diagnostics = new DiagnosticCollector<>();
    StandardJavaFileManager manager = compiler.getStandardFileManager(diagnostics, null, null);

    var units = List.of(file);

    JavaCompiler.CompilationTask task = compiler.getTask(
            null,
            manager,
            diagnostics,
            List.of("-d", "build/classes/java/main"),
            null,
            units);

    if (!task.call()) {
      diagnostics.getDiagnostics().forEach(System.err::println);
      throw new RuntimeException();
    }

    return (Model) Class.forName(String.format("org.blab.math.Model%d", counter++))
            .getDeclaredConstructor()
            .newInstance();
  }
}
