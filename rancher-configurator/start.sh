#!/bin/sh

function wait_for_rancher() {
  while true; do
    curl -s http://${RANCHER}/v1/ > /dev/null
    if [ $? = 0 ]; then 
      sleep 5
      break
    fi
  
    echo "Waiting for Rancher to start..."
    sleep 10
  done
  echo "Rancher detected."
}

function set_api_host() {
  API_URL=$(curl -s ${RANCHER}/v1/settings/api.host | jq -r .links.self)
  cat > /tmp/rancher-data <<EOF
{
  "id": "1as!api.host",
  "type": "activeSetting",
  "name": "api.host",
  "activeValue": null,
  "inDb": false,
  "source": null,
  "value": "http://${RANCHER}"
}
EOF
  curl -s $API_URL -d @/tmp/rancher-data -H "Content-type: application/json" > /dev/null

  echo "API Host set."
}

function register_docker_registry() {
  RANCHER_ENV=$(curl -s http://${RANCHER}/v1/accounts | jq -r .data[0].id)
  cat > /tmp/rancher-data <<EOF
{
  "type": "registry",
  "serverAddress": "${DOCKER}",
  "blockDevicePath": null,
  "created": null,
  "description": null,
  "driverName": null,
  "externalId": null,
  "kind": null,
  "name": null,
  "removed": null,
  "uuid": null,
  "volumeAccessMode": null
}
EOF
  REGISTRY_ID=$(curl -s http://${RANCHER}/v1/projects/${RANCHER_ENV}/registry -d @/tmp/rancher-data -H "Content-type: application/json)| jq .id)
  cat > /tmp/rancher-data <<EOF
{
  "type": "registryCredential",
  "registryId": "${REGISTRY_ID}",
  "email": "${EMAIL}",
  "publicValue": "${USERNAME}",
  "secretValue": "${PASSWORD}",
  "created": null,
  "description": null,
  "kind": null,
  "name": null,
  "removed": null,
  "uuid": null
}
EOF
  curl -s http://${RANCHER}/v1/projects/${RANCHER_ENV}/registries/${REGISTRY_ID}/credential -d @/tmp/rancher-data -H "Content-type: application/json" > /dev/null
  
  echo "Docker registry registered"
}
  
function enable_auth() {
  cat > /tmp/rancher-data <<EOF
{
  "accessMode":"unrestricted",
  "name":"${USERNAME}",
  "id":null,
  "type":"localAuthConfig",
  "enabled":true,
  "password":"${PASSWORD}",
  "username":"${USERNAME}"}
EOF
  curl -s http://${RANCHER}/v1/localauthconfig -d @/tmp/rancher-data -H "Content-type: application/json" > /dev/null
  echo "Rancher auth enabled"
}

wait_for_rancher
set_api_host
register_docker_registry
enable_auth

