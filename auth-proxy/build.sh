#!/bin/sh

VERSION=1.0.1
docker build -t odoko/auth-proxy:$VERSION .
docker push odoko/auth-proxy:$VERSION
