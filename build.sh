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
IN_COUNTAINER=true

while [ $# -gt 0 ]; do
  case $1 in
    IN_CONTAINER)
      IN_CONTAINER=true
      ;;
    notest)
      SKIP_TESTS=true
      ;;
    remote)
      BUILD_REMOTE=true
      ;;
    verbose)
      QUIET=
      ;;
  esac
  shift
done

if [ -z $SKIP_TESTS ] ; then
  export UBER_HOME=`pwd`/tests 
  echo "Executing tests..."
  go test github.com/odoko-devops/uberstack/...

  echo "Execute integration tests..."
  (cd src/github.com/odoko-devops/uberstack/integration; godog)
fi

echo "Building local resources..."
GOOS=darwin GOARCH=amd64 go install github.com/odoko-devops/uberstack/cmd/uberstack

if [ "$BUILD_REMOTE" = "true" ]; then
  echo "Building remote resources..."
  GOOS=linux GOARCH=amd64 go build github.com/odoko-devops/uberstack/cmd/uberstack-remote-agent
fi
