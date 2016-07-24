#!/bin/sh

VERSION=1.0.0
docker build -t odoko/jenkins:${VERSION} .
docker push odoko/jenkins:${VERSION}
