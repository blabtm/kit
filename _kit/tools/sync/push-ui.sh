#!/bin/bash

ui_dir="platform/sm/ui/schema"
index="$ui_dir/index.json"

mkdir -p "$ui_dir" && rm -rf "$ui_dir/*"

cat <<EOF > "$index"
{
  "services": [
EOF

first="1"

find . | grep deploy/service.cue | while read file; do
  src="$(sed 's|/deploy/service.cue$||g' <<< "$file")"
  svc="$(sed -E 's/^\.\/(native|py|golang|jvm|platform)\///g' <<< "$src")"
  svc_dir="$(sed 's|/|\.|g' <<< "$svc")"

  if [[ ! -f "$src/config.cue" ]]; then
    continue
  fi

  echo "$svc"
  mkdir -p "$ui_dir/$svc_dir"

  cue export --out jsonschema -e "#Config" "$src/config.cue" \
    > "$ui_dir/$svc_dir/schema.json"

  if [[ "$?" != "0" ]]; then
      echo "Malformed configuration file. Terminating."
      exit 1
  fi

  sed -i -e 's/%23/#/g' "$ui_dir/$svc_dir/schema.json"
  sed -i '/$schema/d' "$ui_dir/$svc_dir/schema.json"

  if [[ "$first" == "1" ]]; then
    first="0"
  else
    echo "," >> "$index"
  fi

  echo -n "    \"$svc_dir\"" >> "$index"
done

echo "" >> "$index"
echo "  ]" >> "$index"
echo "}" >> "$index"
