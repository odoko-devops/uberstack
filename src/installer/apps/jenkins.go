package apps

import (
	"log"
	"utils"
	"fmt"
	"installer/model"
)

func create_jenkins_mount_point(hostConfig model.HostConfig) {
	log.Println("Create Mount Point for Jenkins")
	command := fmt.Sprintf("docker-machine ssh %s \"sudo mkdir /jenkins ; sudo chown 1000 /jenkins\"", hostConfig.Name)
	utils.Execute(command, nil, "")
}

var dockerCompose = `
jenkins:
  image: odoko/jenkins:2.7.1-odoko01
  ports:
   - 8081:8080
  environment:
    JENKINS_OPTS:
    DOCKER_HOSTNAME: ${DOCKER_HOSTNAME}
    USERNAME: ${USERNAME}
    PASSWORD: ${PASSWORD}
    PLUGINS: git
  volumes:
   - /jenkins:/var/jenkins_home
   - /var/run/docker.sock:/var/run/docker.sock
`