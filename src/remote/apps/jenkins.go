package apps

import (
	"utils"
	"fmt"
	"model"
)

func Jenkins_Install(config model.Config, hostConfig model.HostConfig, app model.AppConfig) error {
	dockerHost := app.Config["docker-host"]
	authRealm := model.GetAuthRealm(config, app.Config["auth-realm"])
	username := authRealm.Users[0].Username
	password := authRealm.Users[0].Password

	command := "sudo mkdir /jenkins"
	utils.ExecuteRemote(hostConfig.Name, command, nil, "")

	command = "sudo chown 1000 /jenkins"
	utils.ExecuteRemote(hostConfig.Name, command, nil, "")

	command = fmt.Sprintf(`docker run -d -p 8081:8080 \
					-e JENKINS_OPTS= \
					-e DOCKER_HOSTNAME=%s \
					-e USERNAME=%s \
					-e PASSWORD=%s \
					-e PLUGINS=git \
					-v /jenkins:/var/jenkins_home \
					-v /var/run/docker.sock:/var/run/docker.sock \
					odoko/jenkins:2.7.1-odoko01`,
					dockerHost,
					username,
					password)
	utils.ExecuteRemote(hostConfig.Name, command, nil, "")
	return nil
}