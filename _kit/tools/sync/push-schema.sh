#!/bin/bash

echo "Synchronizing schemas."

curl -Is "http://$RP_HOST:$RP_SR_PORT" > /dev/null

if [[ "$?" != "0" ]]; then
  echo "Schema registry is down. Skipping."
  exit 0
fi

SR_BASE="http://$RP_HOST:$RP_SR_PORT/subjects"

find schema -type f -print0 | while IFS= read -r -d '' path; do
  name=$(sed -E 's/(schema\/|\.proto)//g' <<< "$path")
  name=$(sed 's|/|.|g' <<< "$name")

  rpk registry schema create $name --schema $path
done
