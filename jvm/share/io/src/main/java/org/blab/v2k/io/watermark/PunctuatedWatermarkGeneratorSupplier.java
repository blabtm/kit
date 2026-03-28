package org.blab.v2k.io.watermark;

import org.apache.flink.api.common.eventtime.Watermark;
import org.apache.flink.api.common.eventtime.WatermarkGenerator;
import org.apache.flink.api.common.eventtime.WatermarkGeneratorSupplier;
import org.apache.flink.api.common.eventtime.WatermarkOutput;

public class PunctuatedWatermarkGeneratorSupplier<T> implements WatermarkGeneratorSupplier<T> {
  @Override
  public WatermarkGenerator<T> createWatermarkGenerator(Context context) {
    return new WatermarkGenerator<T>() {
      @Override
      public void onEvent(T event, long eventTimestamp, WatermarkOutput output) {
        output.emitWatermark(new Watermark(eventTimestamp));
      }

      @Override
      public void onPeriodicEmit(WatermarkOutput output) {
        // nothing to do
      }
    };
  }
}
