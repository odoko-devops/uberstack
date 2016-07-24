#!/bin/sh

export DOCKER_HOST=tcp://172.17.0.1:4243

while true; do
  docker login -u ${USERNAME} -p ${PASSWORD} ${DOCKER_DOMAIN}
  if [ $? = 0 ]; then
    sleep 10
    docker login -u ${USERNAME} -p ${PASSWORD} ${DOCKER_DOMAIN}
    break
  fi

  echo "Waiting for Registry to start..."
  sleep 10
done
echo "Logged into Docker Registry"

/bin/tini -- /usr/local/bin/jenkins.sh
