package org.blab.v2k.vcas.consumer.schema;

import java.io.IOException;
import java.io.Serializable;

public interface RecordDeserializationSchema<T> extends Serializable {
  T deserialize(String v) throws IOException;
}
