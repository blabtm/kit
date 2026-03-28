package org.blab.v2k.math;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;
import lombok.ToString;

@Getter
@Setter
@AllArgsConstructor
@ToString
public class Parameter {
  private double value;
  private Index index;

  @Getter
  @Setter
  @AllArgsConstructor
  public static class Index {
    private int row;
    private int col;
  }

  public void add(double value) {
    this.value += value;
  }
}
