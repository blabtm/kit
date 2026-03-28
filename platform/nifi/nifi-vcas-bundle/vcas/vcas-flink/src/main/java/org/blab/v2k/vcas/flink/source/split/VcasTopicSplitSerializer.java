package org.blab.v2k.vcas.flink.source.split;

import org.apache.flink.core.io.SimpleVersionedSerializer;
import java.io.*;

public class VcasTopicSplitSerializer implements SimpleVersionedSerializer<VcasTopicSplit> {
  @Override
  public int getVersion() {
    return 0;
  }

  @Override
  public byte[] serialize(VcasTopicSplit split) throws IOException {
    var out = new ByteArrayOutputStream();
    var buf = new DataOutputStream(out);

    buf.writeUTF(split.getName());
    buf.flush();

    return out.toByteArray();
  }

  @Override
  public VcasTopicSplit deserialize(int v, byte[] bin) throws IOException {
    return new VcasTopicSplit(new DataInputStream(new ByteArrayInputStream(bin)).readUTF());
  }
}
