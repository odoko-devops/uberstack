#!/bin/bash

ACTION_MESSAGE="Available actions: cert, passwd, run, cron"
ACTION=${1?$ACTION_MESSAGE}

if [ "$ACTION" = "passwd" ]; then
  htpasswd -Bbn ${USERNAME} ${PASSWORD} >> /auth/htpasswd

elif [ "$ACTION" = "run" ]; then
 
  while true; do
    MY_IP=$(curl -s https://api.ipify.org)
    DOCKER_IP=$(getent hosts ${DOCKER_HOSTNAME} | cut -d" " -f1)
    if [ "$MY_IP" = "$DOCKER_IP" ]; then
      break
    fi
    echo "Waiting for $DOCKER_HOSTNAME to point to $MY_IP (currently $DOCKER_IP)"
    sleep 5
  done
  echo "$DOCKER_HOSTNAME now set correctly to $MY_IP"
  echo

  if [ ! -e /etc/letsencrypt/live/${DOCKER_HOSTNAME}/domain.crt ]; then
    while true; do
      /certbot/letsencrypt-auto certonly --keep-until-expiring --standalone --non-interactive --agree-tos -d ${DOCKER_HOSTNAME} --email $EMAIL
      if [ $? = 0 ]; then break; fi 
      sleep 60
    done
   (cd /etc/letsencrypt/live/${DOCKER_HOSTNAME}; cp privkey.pem domain.key; cat cert.pem chain.pem > domain.crt; chmod 755 domain.*)
  fi

  if [ ! -e /auth/htpasswd ]; then
    htpasswd -Bbn ${USERNAME} ${PASSWORD} > /auth/htpasswd
  fi
  /bin/registry serve /etc/docker/registry/config.yml

elif [ "$ACTION" = "cron" ]; then
  touch /var/log/cron.log
  cron && tail -f /var/log/cron.log

else
  echo $ACTION_MESSAGE
  exit
fi
