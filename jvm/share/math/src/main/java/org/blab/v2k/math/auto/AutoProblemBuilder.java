package org.blab.v2k.math.auto;

import org.blab.v2k.math.Problem;

import javax.tools.*;
import java.net.URI;
import java.util.ArrayList;
import java.util.List;

public class AutoProblemBuilder extends Problem.ProblemBuilder<Problem, AutoProblemBuilder> {
  private static int counter = 0;
  protected List<Integer> parameters = new ArrayList<>();

  public AutoProblemBuilder respect(int p) {
    parameters.add(p);
    return self();
  }

  public AutoProblemBuilder respect(List<Integer> p) {
    parameters.addAll(p);
    return self();
  }

  @Override
  protected AutoProblemBuilder self() {
    return this;
  }

  @Override
  public Problem build() {
    try {
      return compile();
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }

  public Problem compile() throws Exception {
    var n = parameters.size();
    var p = new StringBuilder();
    var e = new StringBuilder();

    for (int i = 0; i < n; ++i) {
      p.append(String.format("dst.setEntry(%d, src.getEntry(%d));\n", parameters.get(i), i));
      e.append(String.format("dst.setEntry(%d, src.getEntry(%d));\n", i, parameters.get(i)));
    }

    var src = String.format("""
                      package org.blab.math;
            
                      import org.apache.commons.math3.analysis.differentiation.DerivativeStructure;
                      import org.apache.commons.math3.linear.Array2DRowRealMatrix;
                      import org.apache.commons.math3.linear.ArrayRealVector;
                      import org.apache.commons.math3.linear.RealMatrix;
                      import org.apache.commons.math3.linear.RealVector;
                      import java.util.List;
                      import java.util.ArrayList;
            
                      public class Problem%d extends Problem {
                        protected Problem%d(ProblemBuilder<?, ?> b) {
                          super(b);
                        }
            
                        @Override
                        public int getM() {
                          return %d;
                        }
            
                        @Override
                        public void populate(RealVector src, RealVector dst) {
                          %s
                        }
            
                        @Override
                        public void extract(RealVector src, RealVector dst) {
                          %s
                        }
                      }
            """, counter, counter, n, p, e);

    JavaCompiler compiler = ToolProvider.getSystemJavaCompiler();
    SimpleJavaFileObject file = new SimpleJavaFileObject(
            URI.create(String.format("string:///org/blab/math/Problem%d.java", counter)),
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

    return (Problem) Class.forName(String.format("org.blab.math.Problem%d", counter++))
            .getDeclaredConstructor(Problem.ProblemBuilder.class)
            .newInstance(this);
  }
}
