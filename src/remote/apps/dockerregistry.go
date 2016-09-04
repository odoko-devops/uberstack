package apps

var DockerCompose = `
registry:
  restart: always
  image: odoko/registry:2
  ports:
    - 443:443
  environment:
    DOCKER_HOSTNAME: ${DOCKER_HOSTNAME}
    EMAIL: ${EMAIL}
    USERNAME: ${USERNAME}
    PASSWORD: ${PASSWORD}
    REGISTRY_HTTP_TLS_CERTIFICATE: /etc/letsencrypt/live/${DOCKER_HOSTNAME}/domain.crt
    REGISTRY_HTTP_TLS_KEY: /etc/letsencrypt/live/${DOCKER_HOSTNAME}/domain.key
    REGISTRY_AUTH: htpasswd
    REGISTRY_AUTH_HTPASSWD_PATH: /auth/htpasswd
    REGISTRY_AUTH_HTPASSWD_REALM: Registry Realm
  volumes:
    - /opt/docker-registry:/var/lib/registry
    - /etc/letsencrypt/:/etc/letsencrypt
    - /etc/docker/auth:/auth

letsencryt-cron:
  restart: always
  image: odoko/registry:2
  volumes:
    - /etc/letsencrypt:/etc/letsencrypt
  command:
    - cron
`