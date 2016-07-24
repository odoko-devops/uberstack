#!/bin/bash

###########################################
# Define parameters
AGENT_VER="v1.0.2"
USAGE="$0 <rancher host> <rancher access key> <rancher secret key> [<labels>]"

RANCHER_HOST=${1?$USAGE}
AUTH="${2?$USAGE}:${3?$USAGE}"
LABELS=${4-''}

###########################################
# Install dependencies
sudo apt-get install -y jq

###########################################
# Introspect required information
INTERNAL_IP=$(ip add show eth0 | awk '/inet / {print $2}' | cut -d/ -f1)
RANCHER_ENV=$(curl -su "${AUTH}" http://${RANCHER_HOST}/v1/accounts | jq -r .data[0].id)
RANCHER_URL=$(curl -su "${AUTH}" http://${RANCHER_HOST}/v1/registrationtokens?projectId=${RANCHER_ENV} | jq -r .data[0].registrationUrl)

###########################################
# Install rancher agent
sudo docker run \
    -e CATTLE_AGENT_IP="$INTERNAL_IP" \
    -e CATTLE_HOST_LABELS="$LABELS" \
    -d --privileged --name rancher-bootstrap \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/lib/rancher:/var/lib/rancher \
      rancher/agent:$AGENT_VER $RANCHER_URL
