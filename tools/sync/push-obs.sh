#!/bin/bash

if [[ -z "$GRAFANA_TOKEN" ]]; then
  echo -e "You must setup Grafana first."
  exit 1
fi

BASE="$GRAFANA_HOST:$GRAFANA_PORT"
CRED="$GRAFANA_ADMIN_USER:$GRAFANA_ADMIN_PASS"

curl -Is "http://${BASE}" > /dev/null

if [[ "$?" != "0" ]]; then
  echo "Grafana is down. Skipping."
  exit 0
fi

mkdir -p .obs.dash && rm -rf .obs.dash/*
mkdir -p .obs.dash/folders.v1.folder.grafana.app
mkdir -p .obs.dash/dashboards.v2.dashboard.grafana.app
mkdir -p .obs.prov && rm -rf .obs.prov/*

find . | grep -E "/obs/.*\.yaml" | while read file; do
  kind=$(cat "$file" | yq --raw-output ".kind")

  if [[ "$kind" == "Folder" ]]; then
    cp "$file" ".obs.dash/folders.v1.folder.grafana.app/"
  fi

  if [[ "$kind" == "Dashboard" ]]; then
    cp "$file" ".obs.dash/dashboards.v2.dashboard.grafana.app/"
  fi
done

find ./platform/obs/provisioning -type f | while read file; do
  name="$(echo "$file" | rev | cut -d'/' -f1,2 | rev)"
  mkdir -p ".obs.prov/$(echo "$name" | cut -d '/' -f1)"
  envsubst < "$file" > ".obs.prov/$name"
done

curl -s -X POST \
  "http://${CRED}@${BASE}/api/admin/provisioning/datasources/reload" > /dev/null
gcx resources push -p ./.obs.dash
