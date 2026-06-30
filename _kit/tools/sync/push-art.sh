#!/bin/bash

echo "Synchronizing deployments..."

mkdir -p deploy && rm -rf deploy/*

find . | grep deploy/service.cue | while read file; do
  src="$(sed 's|/deploy/service.cue$||g' <<< "$file")"
  svc="$(sed -E 's/^\.\/(native|py|golang|jvm|platform)\///g' <<< "$src")"

  echo "$svc"
  mkdir -p "deploy/$svc"

  cue export "$file" --out yaml > "deploy/$svc/service.yaml"

  if [[ "$?" != "0" ]]; then
      echo "Malformed deployment file. Terminating."
      rm -rf .deploy/*
      exit 1
  fi

  cp -rf "$src/deploy/art" "deploy/$svc/art"

  find "deploy/$svc/art" -type f | while read art; do
    envsubst < "$art" > "$art.tmp"
    mv "$art.tmp" "$art"
    sed -i -E 's/\%\{([^}]*)\}/${\1}/g' "$art"
  done
done
