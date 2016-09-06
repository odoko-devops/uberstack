#!/bin/sh

if [ "$1" = "local" ]; then
  LOCAL=true
fi

if [ -n $LOCAL ]; then
  GOOS=darwin GOARCH=amd64 go build -o /build/local local 
fi

GOBIN=/build go install uberstack installer remote 

