package org.blab.v2k.vcas.flink.source.enumerator;

import lombok.SneakyThrows;
import org.apache.flink.core.io.SimpleVersionedSerializer;

import java.io.*;

public class VcasSourceEnumStateSerializer implements SimpleVersionedSerializer<VcasSourceEnumState> {
  @Override
  public int getVersion() {
    return 0;
  }

  @SneakyThrows
  public byte[] serialize(VcasSourceEnumState state) throws IOException {
    var out = new ByteArrayOutputStream();
    var buf = new DataOutputStream(out);

    buf.writeInt(state.assigned.size());

    for (String split: state.assigned) {
      buf.writeUTF(split);
    }

    buf.writeInt(state.unassigned.size());

    for (String split: state.unassigned) {
      buf.writeUTF(split);
    }

    buf.flush();

    return out.toByteArray();
  }

  @Override
  public VcasSourceEnumState deserialize(int v, byte[] bin) throws IOException {
    var buf = new DataInputStream(new ByteArrayInputStream(bin));
    var out = new VcasSourceEnumState();

    var assigned = buf.readInt();

    for (int i = 0; i < assigned; ++i) {
      out.assigned.add(buf.readUTF());
    }

    var unassigned = buf.readInt();

    for (int i = 0; i < unassigned; ++i) {
      out.unassigned.offer(buf.readUTF());
    }

    return out;
  }
}
