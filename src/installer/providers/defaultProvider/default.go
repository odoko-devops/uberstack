package defaultProvider

import (
	"installer/model"
	"log"
	"utils"
	"fmt"
	"strings"
	"remote/apps"
)

type DefaultProvider struct {
	Config model.Config
}

func (p DefaultProvider) Configure(config model.Config) {
	p.Config = config
}

func (p DefaultProvider) HostDestroy(host model.HostConfig) (bool, error) {
	log.Printf("Destroy host: %s\n", host.Name)

	command := fmt.Sprintf("docker-machine -s /state/machine rm -f %s", host.Name)
	utils.Execute(command, nil, "")
	return true, nil
}

func (p DefaultProvider) AddUbuntuToDockerGroup(host model.HostConfig) {
	log.Printf("Add ubuntu user to docker unix group on host %s\n", host.Name)
	command := fmt.Sprintf("docker-machine -s /state/machine ssh %s \"sudo gpasswd -a ubuntu docker\"", host.Name)
	utils.Execute(command, nil, "")
}

func (p DefaultProvider) RegenerateCerts(host model.HostConfig) {
	log.Printf("Regenerating certificates for %s\n", host.Name)
	command := fmt.Sprintf("docker-machine -s /state/machine regenerate-certs -f %s", host.Name)
	utils.Execute(command, nil, "")
}

func (p DefaultProvider) UploadSelf(host model.HostConfig) {
	log.Printf("Upload configuration utility to %s\n", host.Name)
	command := fmt.Sprintf("docker-machine -s /state/machine scp /usr/local/bin/remote %s:", host.Name)
	utils.Execute(command, nil, "")
}

func (p DefaultProvider) StartApps(config model.Config, state *model.State, host model.HostConfig, skip *model.SkipList) {
	log.Printf("Installing apps on %s\n", host.Name)
	for _, app := range host.Apps {
		log.Printf("Installing app: %s...\n", app.Type)
		switch app.Type {
		case "registry":
			if !skip.Avoid(model.SkipDockerRegistry) {
				apps.Registry_Install(config, host, app)
			}
		case "rancher-server":
			if !skip.Avoid(model.SkipRancherServer) {
				apps.Rancher_InstallServer(config, state, host, app)
			}
		case "jenkins":
			if !skip.Avoid(model.SkipJenkins) {
				apps.Jenkins_Install(config, host, app)
			}
		case "http-proxy":
			if !skip.Avoid(model.SkipProxy) {
				apps.Proxy_Install(config, host, app)
			}
		case "vpn":
			if !skip.Avoid(model.SkipVpn) {
				//Not yet sorted:
				//apps.Vpn_Install(config, state, host, app)
			}
		default:
			log.Panic("Unknown app: %v", app.Type)
		}
	}
	log.Println("Apps installed")
}

func (p DefaultProvider) StartRancherAgent(config model.Config, state *model.State, provider model.ProviderConfig, host model.HostConfig) {

	labels := make([]string, len(host.Labels))
	i:=0
	for k, v := range host.Labels {
		labels[i] = k + "=" + v
		i++
	}

	providerState := state.Provider[provider.Name]
	command := fmt.Sprintf(`./remote \
	                   -interface=%s \
	                   -rancher=%s \
	                   -access_key=%s \
	                   -secret_key=%s \
	                   -labels=%s \
	                   rancher-agent`,
		host.RancherInterface,
		providerState.RancherUrl,
		providerState.AccessKey,
		providerState.SecretKey,
		strings.Join(labels, ","))
	utils.ExecuteRemote(host.Name, command, nil, "")
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

	  if hostname == "management":
	    print "export DOCKER_TLS_VERIFY=1"
	    print "export DOCKER_HOST=tcp://%s:2376" % config["aws"]["management-host"]["elastic-ip"]
	    print "export DOCKER_CERT_PATH=~/.docker/machine/machines/management"
	    print "export DOCKER_MACHINE_NAME=management"
	*/
}

func (p DefaultProvider) GetRancherEnvironment(state *model.State, provider model.ProviderConfig) {
	providerState := state.Provider[provider.Name]
	fmt.Printf(`
export RANCHER_URL=http://%s
export RANCHER_ACCESS_KEY=%s
export RANCHER_SECRET_KEY=%s\n`,
		providerState.RancherUrl,
		providerState.AccessKey,
		providerState.SecretKey)
}

func (p DefaultProvider) List() {
	utils.Execute("docker-machine -s /state/machine ls", nil, "")
}