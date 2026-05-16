#!/bin/bash

echo "Synchronizing configurations."

mkdir -p .config && rm -rf .config/*

find config -name config.cue | while read file; do
  src="$(sed 's/\/config\.cue//g' <<< "$file")"
  dst="$(sed 's/config\///g' <<< "$src")"

  mkdir -p ".config/$dst"
  cue export "v2k.org/$src" -e config > ".config/$dst/config.json"
done
