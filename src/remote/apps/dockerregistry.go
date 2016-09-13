package apps

import (
	"fmt"
	"utils"
	"model"
)

func Registry_Install(config model.Config, host model.HostConfig, app model.AppConfig) {
	dockerHost := app.Config["host"]
	authRealm := model.GetAuthRealm(config, app.Config["auth-realm"])
	email := authRealm.Users[0].Email
	username := authRealm.Users[0].Username
	password := authRealm.Users[0].Password
	command := fmt.Sprintf(`docker run -d -p 443:443 \
		--restart=always \
		-e DOCKER_HOSTNAME=%s \
		-e EMAIL=%s \
		-e USERNAME=%s \
		-e PASSWORD=%s \
		-e REGISTRY_HTTP_TLS_CERTIFICATE=/etc/letsencrypt/live/%s/domain.crt \
		-e REGISTRY_HTTP_TLS_KEY=/etc/letsencrypt/live/%s/domain.key \
		-e REGISTRY_AUTH=htpasswd \
		-e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
		-e REGISTRY_AUTH_HTPASSWD_REALM=Registry-Realm \
		-v /opt/docker-registry:/var/lib/registry \
		-v /etc/letsencrypt/:/etc/letsencrypt \
		-v /etc/docker/auth:/auth \
		odoko/registry:2`,
		dockerHost,
		email,
		username,
		password,
		dockerHost,
		dockerHost)
	utils.ExecuteRemote(host.Name, command, nil, "")

	command = `docker run -d --restart=always \
			-v /etc/letsencrypt:/etc/letsencrypt \
			odoko/registry:2 cron`
	utils.ExecuteRemote(host.Name, command, nil, "")
}