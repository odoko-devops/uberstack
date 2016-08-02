#!/bin/sh

export DOCKER_HOST=tcp://172.17.0.1:4243

echo "Copying plugins..."
cp /plugins/* /var/jenkins_home/plugins/

echo "About to attempt Docker login..."
while true; do
  docker login -u ${USERNAME} -p ${PASSWORD} ${DOCKER_HOSTNAME}
  if [ $? = 0 ]; then
    sleep 10
    docker login -u ${USERNAME} -p ${PASSWORD} ${DOCKER_HOSTNAME}
    echo "Logged into Docker Registry"
    break
  fi

  echo "Waiting for Registry to start..."
  sleep 10
done &

/bin/tini -- /usr/local/bin/jenkins.sh
