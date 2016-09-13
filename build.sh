#!/bin/bash -e

ARGS=$*

QUIET=-q

while [ $# -gt 0 ]; do
  case $1 in
    IN_CONTAINER)
      IN_CONTAINER=true
      ;;
    build)
      BUILD_CONTAINER=true
      ;;
    remote)
      BUILD_REMOTE=true
      ;;
    push)
      PUSH_IMAGE=true
      ;;
    verbose)
      QUIET=
      ;;
  esac
  shift
done

if [ -z $IN_CONTAINER ]; then
  HERE=$(cd `dirname $0`;pwd)
  if [ "$BUILD_CONTAINER" != "" ]; then
    docker build $QUIET -t odoko/docker-stack-build .
  fi

  docker run -v $HERE/bin:/build -v $HERE/src:/odoko/golibs/src odoko/docker-stack-build IN_CONTAINER $ARGS
else
  echo "Building local resources..."
  GOOS=darwin GOARCH=amd64 go build -o /build/uberstack installer
  GOOS=darwin GOARCH=amd64 go build -o /build/foo foo 

  if [ "$BUILD_REMOTE" = "true" ]; then
    echo "Building remote resources..."
    GOOS=linux GOARCH=amd64 go build -o /build/remote remote
  fi
fi
