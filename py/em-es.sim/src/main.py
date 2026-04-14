import yaml
import time
import math
import os
import random
import socket

root = os.getenv('CONFIG_DIR')
vcas_host = os.getenv('MQ_HOST')
vcas_port = os.getenv('MQ_VCAS_PORT')

print(f'Running with {vcas_host}:{vcas_port}')

with open(f'{root}/ccd/config.yaml', 'r') as file:
    ccd_conf = yaml.safe_load(file)

with open(f'{root}/em-es/config.yaml', 'r') as file:
    svc_conf = yaml.safe_load(file)

with open(f'{root}/em-es.sim/config.yaml', 'r') as file:
    sim_conf = yaml.safe_load(file)

i0 = float(sim_conf['service']['initialCurrentMilliAmp'])
lifetime = sim_conf['service']['beamLifeTimeSec']
noise = float(sim_conf['service']['noisePercent']) / 100

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect((vcas_host, int(vcas_port)))

for t in range(0, lifetime, 1):
    real_time = int((time.time() * 1000))
    i = i0 * math.exp(-float(t) / float(lifetime))
    em = math.sqrt(i) * 1.640487e-5
    es = math.cbrt(i) * 0.000185

    print(f'True: {em:e}, {es:e}')

    sock.sendall(f'name:VEPP/CURRENT|method:set|value:{i}\n'.encode())

    for cam in svc_conf['service']['cams']:
        cam_upper = cam.upper()

        bx = ccd_conf['cams'][cam]['axes']['x']['beta']
        dx = ccd_conf['cams'][cam]['axes']['x']['dispersion']
        bz = ccd_conf['cams'][cam]['axes']['z']['beta']
        dz = ccd_conf['cams'][cam]['axes']['z']['dispersion']

        if svc_conf['service']['cams'][cam]['axes']['x']['weight'] != 0:
            sx = math.sqrt(em * bx + (dx * es) ** 2)
            nx = sx * noise * random.random() * ((-1) ** random.randint(1, 2))
            vx = sx + nx
            topic = f'VEPP/CCD/{cam_upper}/sigma_x'
            sock.sendall(f'name:{topic}|method:set|value:{vx}\n'.encode())

        if svc_conf['service']['cams'][cam]['axes']['z']['weight'] != 0:
            sz = math.sqrt(em * bz + (dz * es) ** 2)
            nz = sz * noise * random.random() * ((-1) ** random.randint(1, 2))
            vz = sz + nz
            topic = f'VEPP/CCD/{cam_upper}/sigma_z'
            sock.sendall(f'name:{topic}|method:set|value:{vz}\n'.encode())

    time.sleep(1)

sock.close()
