#!/bin/sh

VERSION=latest
docker build -t odoko/jenkins:${VERSION} .
docker push odoko/jenkins:${VERSION}
