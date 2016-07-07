#!/bin/sh

sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates
sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" > /tmp/docker.list
sudo mv /tmp/docker.list /etc/apt/sources.list.d/
sudo apt-get update
sudo apt-get install -y docker-engine
sudo gpasswd -a ubuntu docker
