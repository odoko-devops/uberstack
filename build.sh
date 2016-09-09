#!/bin/bash -e

ARGS=$*
while [ $# -gt 0 ]; do
  case $1 in
    IN_CONTAINER)
      IN_CONTAINER=true
      ;;
    build)
      BUILD_CONTAINER=true
      ;;
    local)
      BUILD_LOCAL=true
      ;;
    push)
      PUSH_IMAGE=true
      ;;
  esac
  shift
done

if [ -z $IN_CONTAINER ]; then
  HERE=$(cd `dirname $0`;pwd)
  if [ "$BUILD_CONTAINER" != "" ]; then
    docker build -t odoko/docker-stack-build -f Dockerfile.build .
  fi

  docker run -v $HERE/bin:/build -v $HERE/src:/odoko/golibs/src odoko/docker-stack-build IN_CONTAINER $ARGS
  echo "Building container..."
  docker build -t odoko/docker-stack .
  if [ "$PUSH_IMAGE" != "" ]; then
    docker push odoko/docker-stack
  fi
else
  if [ "$BUILD_LOCAL" = "true" ]; then
    echo "Building local resources..."
    GOOS=darwin GOARCH=amd64 go build -o /build/uberstack local 
  fi

  echo "Building container resources..."
  GOBIN=/build go install installer remote rancheragent 
fi
