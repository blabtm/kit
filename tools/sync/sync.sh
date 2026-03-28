#!/bin/bash

rm -r ./.docs
rm -r ./.deploy

docs_base="http://$SM_HOST:$SM_PORT/v1/docs"

while read -r src; do
    path=$(sed 's/\./\//g' <<< "$src")
    path=$(sed -E 's/(native|py|go|jvm|platform)\///g' <<< "$path")

    deploy=".deploy/$path" 
    docs=".docs/$path"

    echo "Sync: $src > $deploy, $docs"

    mkdir -p $deploy
    mkdir -p $docs

    cp -r ./$src/deploy/* ./$deploy

    if [[ -f "./$src/README.md" ]]; then
        assets="$docs_base/$path/assets"
        cp ./$src/README.md ./$docs
        sed -i "s|src=\"\./assets|src=\"$assets|g" ./$docs/README.md 
    fi

    if [[ -d "./$src/assets" ]]; then
        cp -r ./$src/assets ./$docs
    fi
done < <(find . | grep deploy/service.yaml | sed 's/\/deploy\/service.yaml//g' | sed 's/\.\///g')

find .deploy -type f -print0 | while IFS= read -r -d '' file; do
    if [[ ! -d "$file" && ! "$file" == *".md" ]]; then
        echo "Inject: $file"
        envsubst < "$file" > "$file.tmp" && mv "$file.tmp" "$file"
    fi
done
