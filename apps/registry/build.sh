#!/bin/sh

docker build -t odoko/registry:2 .
docker push odoko/registry:2

