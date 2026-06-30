#!/bin/bash

echo "Setting up Redpanda instance."

if [[ -n "$(rpk profile list | grep v2k)" ]]; then
  echo "Looks like you're trying to setup existing Redpanda instance."
  echo "Sync."

  rpk profile set brokers=$RP_HOST:$RP_LAN_PORT
  rpk profile set schema_registry.addresses=$RP_HOST:$RP_SR_PORT

  exit 0
fi

rpk profile create v2k
rpk profile set brokers=$RP_HOST:$RP_LAN_PORT
rpk profile set schema_registry.addresses=$RP_HOST:$RP_SR_PORT
