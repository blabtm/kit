#!/bin/bash

echo "Validating configurations..."

mkdir -p .config && rm -rf .config/*

find config -name config.json | while read file; do
  svc="$(sed 's|/config.json$||g' <<< "$file")"
  svc="$(sed 's|^config/||g' <<< "$svc")"

  cue vet -c "github.com/blabtm/v2k/model/$svc" -d "#Config" "$file"

  if [[ "$?" != "0" ]]; then
      echo "Malformed configuration. Terminating."
      exit 1
  fi
done

# find config -name config.cue | while read file; do
#   src="$(sed 's/\/config\.cue//g' <<< "$file")"
#   dst="$(sed 's/config\///g' <<< "$src")"
# 
#   echo "$file"
#   mkdir -p ".config/$dst"
# 
#   cue vet -c "config/$dst/config.json"
# 
#   cue export "github.com/blabtm/v2k/$src" > ".config/$dst/config.json"
# 
#   if [[ "$?" != "0" ]]; then
#       echo "Malformed configuration. Terminating."
#       rm -rf .config/*
#       exit 1
#   fi
# done
