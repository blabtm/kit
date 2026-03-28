package org.blab.v2k.em_es;

import org.blab.v2k.math.Solver;

import java.util.Map;

@lombok.Getter
@lombok.Setter
public class Config {
  public static final String TOPIC_CURRENT = "vepp.currents.fz";
  public static final String TOPIC_EMITTANCE = "vepp.em";
  public static final String TOPIC_ENERGY_SPREAD = "vepp.es";
  public static final String TOPIC_LOG = "vepp.svc.em-es.log";

  private Integer window = 3;
  private Solver.Config solver;
  private Map<String, CamSetup> cams;

  @lombok.Getter
  @lombok.Setter
  public static class CamSetup {
    private transient org.blab.v2k.ccd.Config.CamConfig config;

    private AxesSetup axes;

    @lombok.Getter
    @lombok.Setter
    public static class AxesSetup {
      private AxisSetup x;
      private AxisSetup z;
    }

    @lombok.Getter
    @lombok.Setter
    public static class AxisSetup {
      private transient org.blab.v2k.ccd.Config.CamConfig.AxisConfig config;
      private transient String channel;

      private Double weight = 1.0;
    }
  }
}
