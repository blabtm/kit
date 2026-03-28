package org.blab.v2k.ccd;

import java.util.Map;

@lombok.Getter
@lombok.Setter
public class Config {
  private Map<String, CamConfig> cams;

  @lombok.Getter
  @lombok.Setter
  public static class CamConfig {
    private AxesConfig axes;

    @lombok.Getter
    @lombok.Setter
    public static class AxesConfig {
      private AxisConfig x;
      private AxisConfig z;
    }

    @lombok.Getter
    @lombok.Setter
    public static class AxisConfig {
      private Double beta;
      private Double dispersion;
    }
  }
}
