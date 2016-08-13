#!/bin/sh

VERSION=2.7.1-odoko07
docker build -t odoko/jenkins:${VERSION} .
docker push odoko/jenkins:${VERSION}
