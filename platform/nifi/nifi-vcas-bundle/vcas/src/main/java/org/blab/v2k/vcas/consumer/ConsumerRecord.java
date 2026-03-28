package org.blab.v2k.vcas.consumer;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.Arrays;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class ConsumerRecord {
  public static DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern("dd.MM.yyyy HH_mm_ss.SSS");

  private String topic;
  private Long timestamp;
  private String value;

  public static ConsumerRecord parse(String buf) {
    var rec = new ConsumerRecord();

    rec.setTimestamp(System.currentTimeMillis());

    Arrays.stream(buf.split("\\|"))
            .map(Field::parse)
            .forEach(f -> {
              switch (f.name) {
                case "n", "name":
                  rec.setTopic(f.value);
                  break;
                case "v", "val", "value":
                  rec.setValue(f.value);
                  break;
                case "t", "time":
                  rec.setTimestamp(LocalDateTime.parse(f.value, dateTimeFormatter)
                          .toInstant(ZoneOffset.ofHours(7))
                          .toEpochMilli());
                  break;
              }
            });

    return rec;
  }

  record Field(String name, String value) {
    static Field parse(String f) {
      var t = f.split(":");

      if (t.length != 2) {
        throw new IllegalArgumentException("malformed field: \"" + f + "\"");
      }

      return new Field(t[0], t[1]);
    }
  }
}
