#!/bin/sh

set -e

case $(uname -m) in
  "aarch64") target="arm64" ;;
  *) target="amd64" ;;
esac

vers="0.0.0"
uri="https://github.com/blabtm/v2k/releases/download/platform/sm/v${vers}/sm_${vers}_${target}"
home="$HOME/.sm"
bin="$home/bin"
exe="$bin/sm"

mkdir -p "$home"
mkdir -p "$bin"

curl --fail --location --progress-bar --output "$exe" "$uri"
chmod +x "$exe"
cat <<EOF > "$home/config.yaml"
remote:
  host: localhost
  port: 8080
EOF
