#!/bin/bash

echo "Synchronizing deployments."

mkdir -p .deploy && rm -rf .deploy/*

while read -r src; do
    path=$(sed 's/\./\//g' <<< "$src")
    path=$(sed -E 's/(native|py|golang|jvm|platform)\///g' <<< "$path")
    path="./.deploy/$path"

    echo "Sync: $src > $path"
    mkdir -p $path
    cp -r ./$src/deploy/* $path
done < <(find . | grep deploy/service.yaml | sed 's/\/deploy\/service.yaml//g' | sed 's/\.\///g')

find .deploy -type f -print0 | while IFS= read -r -d '' file; do
    if [[ ! -d "$file" && ! "$file" == *".md" ]]; then
        echo "Inject: $file"
        envsubst < "$file" > "$file.tmp" && mv "$file.tmp" "$file"
    fi
done

