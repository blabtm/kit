#!/bin/bash

echo "Setup Appsmith..."

ADDR="http://$APPSMITH_HOST:$APPSMITH_PORT/api/v1"
WS="69aebe6f21fc7b25434cffda"

git config --file .gitmodules --get-regexp url  \
  | awk '{ print $2 }'                          \
  | grep "\.app"                                \
  | while IFS='' read -r url; do
  echo "  $url"

  curl -X POST "$ADDR/git/import/$WS" \
    --header "Content-Type: application/json" \
    --data-raw "{
      \"remoteUrl\": \"$url\"
    }"
done

