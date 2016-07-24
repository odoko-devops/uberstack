#!/bin/sh

cd /tmp
curl -Ls https://github.com/docker/compose/releases/download/1.6.2/docker-compose-`uname -s`-`uname -m` > docker-compose
chmod +x docker-compose

sudo bash -c "echo 'DOCKER_OPTS=\"-H tcp://172.17.0.1:4243 -H unix:///var/run/docker.sock\"' >> /etc/default/docker"
sleep 5
sudo service docker restart
sleep 15 # wait so that Docker has properly restarted before Jenkins attempts to log-in

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
