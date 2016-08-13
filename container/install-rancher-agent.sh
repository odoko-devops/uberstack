#!/bin/bash

###########################################
# Define parameters
AGENT_VER="v1.0.2"
USAGE="$0 <rancher host> <rancher access key> <rancher secret key> [<labels>]"

RANCHER_HOSTNAME=${1?$USAGE}
AUTH="${2?$USAGE}:${3?$USAGE}"
INTERFACE="${4-eth0}"
LABELS=${5-''}

###########################################
# Install dependencies
sudo wget -O /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
sudo chmod +x /usr/local/bin/jq

###########################################
# Introspect required information
INTERNAL_IP=$(ip add show ${INTERFACE} | awk '/inet / {print $2}' | cut -d/ -f1)
echo "Internal IP=${INTERNAL_IP}"
RANCHER_ENV=$(curl -su "${AUTH}" http://${RANCHER_HOSTNAME}/v1/accounts | jq -r .data[0].id)
echo "Rancher environment=${RANCHER_ENV}"
RANCHER_URL=$(curl -su "${AUTH}" "http://${RANCHER_HOSTNAME}/v1/projects/${RANCHER_ENV}/registrationtokens?state=active&limit=-1" | jq -r .data[0].registrationUrl)
echo "Rancher Registration URL=${RANCHER_URL}"

###########################################
# Install rancher agent
sudo docker run \
    -e CATTLE_AGENT_IP="$INTERNAL_IP" \
    -e CATTLE_HOST_LABELS="$LABELS" \
    -d --privileged --name rancher-bootstrap \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/lib/rancher:/var/lib/rancher \
      rancher/agent:$AGENT_VER $RANCHER_URL
