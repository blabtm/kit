#!/bin/sh

set -e

case $(uname -m) in
  "aarch64") target="arm64" ;;
  *) target="amd64" ;;
esac

uri="https://github.com/blabtm/v2k/releases/download/platform/sm/v0.0.0/sm_0.0.0_${target}"
home="$HOME/.sm"
bin="$home/bin"
exe="$bin/sm"

mkdir -p "$home"
mkdir -p "$bin"

curl --fail --location --progress-bar --output "$exe" "$uri"
cat <<EOF > "$home/config.yaml"
remote:
  host: localhost
  port: 8080
EOF

