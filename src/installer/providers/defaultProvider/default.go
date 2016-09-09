package defaultProvider

import (
	"installer/model"
	"log"
	"utils"
	"fmt"
	"strings"
	"remote/apps"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
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
			log.Panic("Unknown app: " + app.Type)
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

func (p DefaultProvider) ListHosts() {
	utils.Execute("docker-machine -s /state/machine ls", nil, "")
}


/***********************************************************************
 * Identify stacks within $UBER_HOME directory
 */
func getUberstacks(uberHome string) []string {
	files, _ := ioutil.ReadDir(uberHome + "/uberstacks")
	uberstacks := make([]string, 0, len(files))
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yml") {
			s := strings.Split(f.Name(), ".")
			if len(s) == 2 && s[1] == "yml" {
				uberstacks = uberstacks[0: len(uberstacks) + 1]
				uberstacks[len(uberstacks) - 1] = s[0]
			}
		}
	}
	return uberstacks
}

/***********************************************************************
 * Identify stacks within $UBER_HOME directory
 */
func getStacks(uberHome string) []string {
	files, _ := ioutil.ReadDir(uberHome + "/stacks")
	stacks := make([]string, 0, len(files))
	for _, f := range files {
		s := strings.Split(f.Name(), ".")
		if len(s) == 1 {
			stacks = stacks[0: len(stacks) + 1]
			stacks[len(stacks) - 1] = s[0]
		}
	}
	return stacks
}

/***********************************************************************
 * Read a single uberstack from its config yaml file
 */
func (p DefaultProvider) GetUberstack(uberHome string, name string) model.Uberstack {
	bytes, err := ioutil.ReadFile(uberHome + "/uberstacks/" + name + ".yml")
	utils.Check(err)
	uberstack := model.Uberstack{}
	err = yaml.Unmarshal(bytes, &uberstack)
	utils.Check(err)
	return uberstack
}

/***********************************************************************
 * Execute ls command
 */
func ListUberstacks(uberHome string) {
	fmt.Println("\nStacks:")
	fmt.Println("-------")
	stacks := getStacks(uberHome)
	for stack_id := range stacks {
		fmt.Println(stacks[stack_id])
	}
	fmt.Println("\nUberstacks:")
	fmt.Println("-----------")
	uberstacks := getUberstacks(uberHome)
	for uber_id := range uberstacks {
		fmt.Println(uberstacks[uber_id])
	}
}

/***********************************************************************
 * Build a suitable environment for execution
 */
func getParametersFor(uberstack model.Uberstack, env string, state_file string) utils.Environment {
	params := getParametersFromEnvironmentAndUberstack(uberstack, env)
	addParametersFromState(env, state_file, &params)
	checkRequiredUberstackVariables(uberstack, params)
	return params
}

func getParametersFromEnvironmentAndUberstack(uberstack model.Uberstack, env string) utils.Environment {
	environ := os.Environ()
	params := utils.Environment{}

	for _, v := range environ {
		s := strings.SplitN(v, "=", 2)
		name := s[0]
		value := s[1]
		params[name] = value
	}

	for k, v := range uberstack.Environments[env] {
		params[k] = v
	}

	return params
}

func addParametersFromState(env string, state_file string, params *utils.Environment) {
	state := model.LoadState("/state/state.yml")

	if env == "local" {
		(*params)["RANCHER_ACCESS_KEY"] = state.Provider["virtualbox"].AccessKey
		(*params)["RANCHER_SECRET_KEY"] = state.Provider["virtualbox"].SecretKey
	} else {
		(*params)["RANCHER_ACCESS_KEY"] = state.Provider["amazonec2"].AccessKey
		(*params)["RANCHER_SECRET_KEY"] = state.Provider["amazonec2"].SecretKey
	}
}

/***********************************************************************
 * Check for required variables
 */
func checkRequiredUberstackVariables(uberstack model.Uberstack, params utils.Environment) {

	for i := range uberstack.Required {
		required := uberstack.Required[i]
		_, ok := params[required]
		if !ok {
			log.Fatal("Required parameter: ", required)
			os.Exit(1)
		}
	}
}

/***********************************************************************
 * Process any referenced Uberstacks
 */
func (p DefaultProvider) ProcessUberstack(uberHome string, uberstack model.Uberstack, env string, cmd string, exclude_stack string) {

	fmt.Println("process_uberstack", uberstack.Name)
	for i := 0; i < len(uberstack.Uberstacks); i++ {
		name := uberstack.Uberstacks[i]
		inner_uberstack := p.GetUberstack(uberHome, name)
		p.ProcessUberstack(uberHome, inner_uberstack, env, cmd, exclude_stack)
	}

	for i := range uberstack.Stacks {
		name := uberstack.Stacks[i]
		if name == exclude_stack {
			continue
		}
		project := name
		stack := name

		s := strings.SplitN(name, ":", 2)
		if len(s) == 2 {
			project = s[0]
			stack = s[1]
		}
		command := fmt.Sprintf(`rancher-compose --file %s/stacks/%s/docker-compose.yml \
                        --rancher-file %s/stacks/%s/rancher-compose.yml \
                        --project-name %s \
                        %s`,
			uberHome, stack, uberHome, stack, project, cmd)
		env := getParametersFor(uberstack, env, "/state/state.yml")
		utils.Execute(command, env, "")
	}
}
