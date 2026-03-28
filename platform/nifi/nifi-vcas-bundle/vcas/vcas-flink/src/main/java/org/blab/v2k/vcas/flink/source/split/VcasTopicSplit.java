package org.blab.v2k.vcas.flink.source.split;

import lombok.Getter;
import org.apache.flink.api.connector.source.SourceSplit;

@Getter
public class VcasTopicSplit implements SourceSplit {
    private final String name;

    public VcasTopicSplit(String name) {
        this.name = name;
    }

    @Override
    public String splitId() {
        return this.name;
    }
}
