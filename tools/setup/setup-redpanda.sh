#!/bin/bash

echo "Setting up Redpanda instance."

if [[ -n "$(rpk profile list | grep v2k)" ]]; then
  echo -e "Looks like you're trying to setup existing Redpanda instance."
  echo -e "Try to sync instead."
  exit 0
fi

rpk profile set brokers=$RP_HOST:$RP_LAN_PORT
rpk profile set schema_registry.addresses=$RP_HOST:$RP_SR_PORT
