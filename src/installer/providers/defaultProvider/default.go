package defaultProvider

import (
	"installer/model"
	"log"
	"utils"
	"fmt"
)

type DefaultProvider struct {
	Config model.Config
}

func (p DefaultProvider) Configure(config model.Config) {
	p.Config = config
}

func (p DefaultProvider) HostDestroy(host model.HostConfig) (bool, error) {
	log.Printf("Destroy host: %s\n", host.Name)
	command := fmt.Sprintf("docker-machine rm -f %s", host.Name)
	utils.Execute(command, nil, "")
	return true, nil
}

func (p DefaultProvider) AddUbuntuToDockerGroup(host model.HostConfig) {
	log.Printf("Add ubuntu user to docker unix group on host %s\n", host.Name)
	command := fmt.Sprintf("docker-machine ssh %s \"sudo gpasswd -a ubuntu docker\"", host.Name)
	utils.Execute(command, nil, "")
}

func (p DefaultProvider) RegenerateCerts(host model.HostConfig) {
	command := fmt.Sprintf("docker-machine regenerate-certs -f %s", host.Name)
	utils.Execute(command, nil, "")
}

func (p DefaultProvider) StartApps(config model.Config, state model.State, host model.HostConfig) {
/*
	log.Println("Deploy Applications")
	env := utils.Environment {
		"DOCKER_HOSTNAME": config["apps"]["docker"]["name"],
	"RANCHER_HOSTNAME": config["apps"]["rancher"]["name"],
	"JENKINS_HOSTNAME": config["apps"]["jenkins"]["name"],
	"EMAIL": config["auth"]["email"],
	"USERNAME": config["auth"]["username"],
	"PASSWORD": config["auth"]["password"],
	"DOCKER_TLS_VERIFY": "1",
	"DOCKER_HOST": "tcp://%s:2376" % config["aws"]["management-host"]["elastic-ip"],
	"DOCKER_CERT_PATH": "/odoko/.docker/machine/machines/management",
	"DOCKER_MACHINE_NAME": "management"
	}
	utils.Execute("docker-compose up -d", env, cwd)
 */
}

func (p DefaultProvider) StartRancherAgent(config model.Config, state model.State, host model.HostConfig) {
	/*
	if config["hosts"].has_key(hostname):
		labels = ",".join(["%s=%s" % (k, v) for k, v in	config["hosts"][hostname]["labels"].items()])
	else :
		labels = ""

	execute("docker-machine scp golibs/bin/rancheragent docker-host%s:" % count)
	execute("docker-machine ssh %s ./rancheragent -interface eth0 -rancher=%s -access_key=%s -secret_key=%s -labels=%s" %
		(hostname,
		config["apps"]["rancher"]["name"],
		rancher["api-access-key"],
		rancher["api-secret-key"],
		labels))

def enable_rancher(config, host):
  rancher_host=config["local"]["rancher"]["ip"]
  rancher_server.wait_for_rancher(rancher_host)
  rancher_server.set_api_host(rancher_host)
  access_key, secret_key = rancher_server.get_keys(rancher_host)

  with open("install-rancher-agent.sh") as f:
      script = f.read()
  script = script.replace("${1?$USAGE}", rancher_host)
  script = script.replace("${2?$USAGE}", access_key)
  script = script.replace("${3?$USAGE}", secret_key)
  script = script.replace("${4-eth0}", "eth1")
  script = "cat <<'EOF' | docker-machine ssh %s\n%s\nEOF" % (host, script)
  write_script("/state/run", script)
  ask("state/run")
  return access_key, secret_key
  */
}

func (p DefaultProvider) getDockerEnvironment(host model.HostConfig) {
	/*
	  RE=re.compile(r"export (.*)=\"(.*)\"")
	  execute("docker-machine regenerate-certs -f management")
	  result=execute("docker-machine env --shell management")

	  env={}
	  for line in result.split("\n"):
	    m=RE.match(line)
	    if m:
	      env[m.group(1)] = m.group[2]
	  return env
	*/
}

func (p DefaultProvider) getRancherEnvironment() {
	/*
	  if hostname == "management":
	    print "export DOCKER_TLS_VERIFY=1"
	    print "export DOCKER_HOST=tcp://%s:2376" % config["aws"]["management-host"]["elastic-ip"]
	    print "export DOCKER_CERT_PATH=~/.docker/machine/machines/management"
	    print "export DOCKER_MACHINE_NAME=management"
	  if rancher == "remote":
	    print "export RANCHER_URL=http://%s" % config["apps"]["rancher"]["name"]
	    print "export RANCHER_ACCESS_KEY=%s" % config["rancher"]["api-access-key"]
	    print "export RANCHER_SECRET_KEY=%s" % config["rancher"]["api-secret-key"]
	  elif rancher == "local-rancher":
	    print "export RANCHER_URL=http://%s" % config["local"]["rancher"]["ip"]
	    print "export RANCHER_ACCESS_KEY=%s" % config["local-rancher"]["api-access-key"]
	    print "export RANCHER_SECRET_KEY=%s" % config["local-rancher"]["api-secret-key"]
	*/

}