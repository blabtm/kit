from confluent_kafka import Producer
from confluent_kafka.schema_registry import SchemaRegistryClient
from confluent_kafka.schema_registry.protobuf import ProtobufSerializer
from confluent_kafka.schema_registry import record_subject_name_strategy
from confluent_kafka.serialization import SerializationContext, MessageField

from v2k.beam.AxisSize_pb2 import AxisSize
from v2k.beam.Current_pb2 import Current

import yaml
import time
import math
import os
import random

root = os.getenv('CONFIG_DIR')
regc = SchemaRegistryClient({'url': os.getenv('SR_URI')})
prod = Producer({'bootstrap.servers': os.getenv('RP_URI')})

with open(f'{root}/ccd/config.yaml', 'r') as file:
    ccdConf = yaml.safe_load(file)

with open(f'{root}/em-es/config.yaml', 'r') as file:
    svcConf = yaml.safe_load(file)

with open(f'{root}/em-es.sim/config.yaml', 'r') as file:
    simConf = yaml.safe_load(file)

i0 = float(simConf['service']['initialCurrentMilliAmp'])
lt = simConf['service']['beamLifeTimeSec']
np = float(simConf['service']['noisePercent']) / 100

cuSer = ProtobufSerializer(Current, regc, conf = {
    'use.deprecated.format': False,
    'subject.name.strategy': record_subject_name_strategy
})

szSer = ProtobufSerializer(AxisSize, regc, conf = {
    'use.deprecated.format': False,
    'subject.name.strategy': record_subject_name_strategy
})

for t in range(0, lt, 1):
    rt = int((time.time() * 1000))

    cu = i0 * math.exp(-float(t) / float(lt))
    em = math.sqrt(cu) * 1.640487e-5
    es = math.cbrt(cu) * 0.000185

    prod.produce(topic = 'vepp.beam.cur', value = cuSer(
        Current(time=rt, value=cu),
        SerializationContext('vepp.beam.cur', MessageField.VALUE)
    ))

    for cam in svcConf['service']['cams']:
        bx = ccdConf['cams'][cam]['axes']['x']['beta']
        dx = ccdConf['cams'][cam]['axes']['x']['dispersion']
        bz = ccdConf['cams'][cam]['axes']['z']['beta']
        dz = ccdConf['cams'][cam]['axes']['z']['dispersion']

        if svcConf['service']['cams'][cam]['axes']['x']['weight'] != 0:
            x = math.sqrt(em * bx + (dx * es) ** 2)
            n = x * np * random.random() * ((-1) ** random.randint(1, 2))
            t = f'vepp.ccd.{cam}.sigma_x'
            prod.produce(topic = t, value = szSer(
                AxisSize(time=rt, value=x+n),
                SerializationContext(t, MessageField.VALUE)
            ))
        
        if svcConf['service']['cams'][cam]['axes']['z']['weight'] != 0:
            z = math.sqrt(em * bz + (dz * es) ** 2)
            n = z * np * random.random() * ((-1) ** random.randint(1, 2))
            t = f'vepp.ccd.{cam}.sigma_z'
            prod.produce(topic = t, value = szSer(
                AxisSize(time=rt, value=z+n),
                SerializationContext(t, MessageField.VALUE)
            ))

    time.sleep(1)
