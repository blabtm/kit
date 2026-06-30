#!/bin/bash

echo "Setting up Grafana instance."

BASE="$GRAFANA_HOST:$GRAFANA_PORT"
CRED="$GRAFANA_ADMIN_USER:$GRAFANA_ADMIN_PASS"

curl -Is "http://$BASE" > /dev/null

if [[ "$?" != "0" ]]; then
  echo "Grafana is down. Skipping."
  exit 0
fi

if [[ -n "$GRAFANA_TOKEN" ]]; then
  echo "Looks like you're trying to setup existing Grafana instance."
  echo "If you want to setup again, remove token from Grafana and from platform variables."
  echo "Updating gcx settings to match platform variables."

  gcx login v2k-adm             \
    --server "$GRAFANA_SERVER"  \
    --token "$GRAFANA_TOKEN"    \
    --yes

  exit 0
fi

result=$(curl -s                                           \
           -X POST                                         \
           -H 'Content-Type: application/json'             \
           -d '{"name":"v2k-adm", "role": "Admin"}'        \
           "http://${CRED}@${BASE}/api/serviceaccounts")

code=$(echo "$result" | jq --raw-output ".statusCode")

if [[ "$code" == "400" ]]; then
  echo "Service account 'v2k-adm' already exists."

  result=$(curl -s "http://${CRED}@${BASE}/api/serviceaccounts/search?query=v2k-adm")
  id=$(echo "$result" | jq --raw-output ".serviceAccounts[0].id")

  echo "Service account 'v2k-adm' found. ID: $id."
else
  id=$(echo "$result" | jq --raw-output ".id")

  echo "Service account 'v2k-adm' created. ID: $id."
fi

echo "Using service account 'v2k-adm' with ID: $id."

result=$(curl -s                                                        \
           -X POST                                                      \
           -H 'Content-Type: application/json'                          \
           -d '{"name": "v2k-adm"}'                                     \
           "http://${CRED}@${BASE}/api/serviceaccounts/${id}/tokens")

code=$(echo "$result" | jq --raw-output ".statusCode")

if [[ "$code" == "400" ]]; then
  echo -e "Token 'v2k-adm' already exists."
  echo -e "You must specifiy token in the \"GRAFANA_TOKEN\" environment variable."
  echo -e "If you lost the token, delete it and run this script again."
  exit 1;
fi

token=$(echo "$result" | jq --raw-output ".key")

sed -i -e "s/GRAFANA_TOKEN=\"\"/GRAFANA_TOKEN=\""$token"\"/g" \
  $ROOT_DIR/platform/env/secret.env

echo "Token created and stored in platform/env/secret.env."
echo "Updating gcx settings to match platform variables."

gcx login v2k-adm             \
  --server "$GRAFANA_SERVER"  \
  --token "$GRAFANA_TOKEN"    \
  --yes

exit 0
