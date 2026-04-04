#!/bin/bash

SR_BASE="http://$RP_HOST:$RP_SR_PORT/subjects"

find schema -type f -print0 | while IFS= read -r -d '' path; do
  name=$(sed -E 's/(schema\/|\.proto)//g' <<< "$path")
  name=$(sed 's|/|.|g' <<< "$name")

  echo "Installing schema: $name ($path)"

  rpk registry schema create $name --schema $path
done
