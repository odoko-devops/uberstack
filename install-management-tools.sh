#!/bin/sh

cd /tmp
curl -Ls https://github.com/docker/compose/releases/download/1.6.2/docker-compose-`uname -s`-`uname -m` > docker-compose
chmod +x docker-compose

export DOMAIN=$1
export EMAIL=$2
export USERNAME=$3
export PASSWORD=$4
export RANCHER_HOST=$5
export JENKINS_HOST=$6
export RANCHER_PORT=8080
export JENKINS_PORT=8081

sudo mkdir /jenkins
sudo chown 1000 /jenkins

./docker-compose up -d
