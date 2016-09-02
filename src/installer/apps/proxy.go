package apps

var dockercompose = `
auth-proxy:
  image: odoko/auth-proxy:1.0.1
  ports:
    - 80:80
  environment:
    JENKINS_HOSTNAME: ${JENKINS_HOSTNAME}
    JENKINS_PORT: 8081
    RANCHER_HOSTNAME: ${RANCHER_HOSTNAME}
    RANCHER_PORT: 8080
  `