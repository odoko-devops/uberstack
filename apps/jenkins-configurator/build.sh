#!/bin/sh

docker build -t odoko/jenkins-configurator .
docker push odoko/jenkins-configurator

