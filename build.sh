#!/bin/sh

if [ "$1" = "local" ]; then
  LOCAL=true
fi

go get gopkg.in/yaml.v2
if [ -n $LOCAL ]; then
  GOOS=darwin GOARCH=amd64 go build -o /build/local local 
fi

GOBIN=/build go install uberstack installer remote 

