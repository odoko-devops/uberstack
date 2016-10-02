#!/bin/bash -e

####################################################################
# UBERSTACK build script
# 
# This build runs in two phases - it prepares a Docker image
# containing the Go executable, and then calls itself inside a 
# container (built from that image) to build the Go libraries.
#
# Options to this script are:
#  * build: rebuild the build container, and dependencies
#  * remote: build the 'uberstack-remote-agent' tool
#  * verbose: when building the container, explain what you are doing

ARGS=$*

QUIET=-q

while [ $# -gt 0 ]; do
  case $1 in
    IN_CONTAINER)
      IN_CONTAINER=true
      ;;
    remote)
      BUILD_REMOTE=true
      ;;
    verbose)
      QUIET=
      ;;
    test)
      TEST=true
      ;;
  esac
  shift
done

if [ -z $IN_CONTAINER ]; then
  HERE=$(cd `dirname $0`;pwd)
  docker build $QUIET -t odoko/uberstack-build .

  docker run -v $HERE/build.sh:/go/build.sh -v $HERE/bin:/go/bin/darwin_amd64 -v $HERE/tests:/go/tests -v $HERE/src:/go/src odoko/uberstack-build IN_CONTAINER $ARGS
    
else
  echo "Building local resources..."
  GOOS=darwin GOARCH=amd64 go install github.com/odoko-devops/uberstack/uberstack

  if [ "$BUILD_REMOTE" = "true" ]; then
    echo "Building remote resources..."
    GOOS=linux GOARCH=amd64 go build -o /build/uberstack-remote-agent uberstack-remote-agent
  fi
  if [ "$TEST" = "true" ]; then
    echo "Executing tests..."
    go test github.com/odoko-devops/uberstack/...
  fi
fi

  
