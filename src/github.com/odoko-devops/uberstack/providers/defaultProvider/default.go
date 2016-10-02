package defaultProvider

import (
	"github.com/odoko-devops/uberstack/model"
	"log"
	"github.com/odoko-devops/uberstack/utils"
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/odoko-devops/uberstack/apps"
	"os/exec"
)

type DefaultProvider struct {
	Config model.Config
}

func (p DefaultProvider) Configure(config model.Config) {
	p.Config = config
}

func (p DefaultProvider) HostDestroy(host model.HostConfig) (bool, error) {
	log.Printf("Destroy host: %s\n", host.Name)

	command := fmt.Sprintf("docker-machine -s %s/machine rm -f %s", utils.GetUberState(), host.Name)
	utils.Execute(command, nil, "")
	return true, nil
}

func (p DefaultProvider) AddUbuntuToDockerGroup(host model.HostConfig) {
	log.Printf("Add ubuntu user to docker unix group on host %s\n", host.Name)
	command := "sudo gpasswd -a ubuntu docker"
	utils.ExecuteRemote(host.Name, command, nil, "")
}

func (p DefaultProvider) RegenerateCerts(host model.HostConfig) {
	log.Printf("Regenerating certificates for %s\n", host.Name)
	command := fmt.Sprintf("docker-machine -s %s/machine regenerate-certs -f %s", utils.GetUberState(), host.Name)
	utils.Execute(command, nil, "")
}

func (p DefaultProvider) UploadSelf(host model.HostConfig) {
	log.Printf("Upload configuration utility to %s\n", host.Name)
	path, err := exec.LookPath("uberstack-remote-agent")
	if err != nil {
		fmt.Println("Cannot locate uberstack-remote-agent in your PATH. Please correct and retry.")
		os.Exit(1)
	}
	command := fmt.Sprintf("docker-machine -s %s/machine scp %s %s:",
		utils.GetUberState(), path, host.Name)
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
				apps.Vpn_Install(config, state, host, app)
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
	command := fmt.Sprintf(`./uberstack-remote-agent \
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

func (p DefaultProvider) ListHosts() {
	command := fmt.Sprintf("docker-machine -s %s/machine ls", utils.GetUberState())
	utils.Execute(command, nil, "")
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
func getParametersFor(uberstack model.Uberstack, env string, state *model.State) utils.Environment {
	params := getParametersFromEnvironmentAndUberstack(uberstack, env)

	provider := uberstack.Environments[env].Provider
	providerState := state.Provider[provider]
	params["RANCHER_URL"] = fmt.Sprintf("http://%s/", providerState.RancherUrl)
	params["RANCHER_ACCESS_KEY"] = providerState.AccessKey
	params["RANCHER_SECRET_KEY"] = providerState.SecretKey

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

	for k, v := range uberstack.Environments[env].Environment {
		params[k] = v
	}

	return params
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
func (p DefaultProvider) ProcessUberstack(config model.Config, state *model.State, uberHome string,
		uberstack model.Uberstack, env string, cmd string, exclude_stack string, doTerraform bool) {

	for i := 0; i < len(uberstack.Uberstacks); i++ {
		name := uberstack.Uberstacks[i]
		inner_uberstack := p.GetUberstack(uberHome, name)
		p.ProcessUberstack(config, state, uberHome, inner_uberstack, env, cmd, exclude_stack, doTerraform)
	}

	uberEnv := uberstack.Environments[env]
	params := utils.Environment{}
	for k,v := range uberEnv.TerraformConfig {
		params[k] = v
	}

	if doTerraform && len(uberEnv.TerraformBefore)>0 {
		utils.TerraformApply(uberEnv.Provider, uberEnv.TerraformBefore, params)
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
		env := getParametersFor(uberstack, env, state)
		utils.Execute(command, env, "")
	}
	if doTerraform && len(uberEnv.TerraformAfter)>0 {
		utils.TerraformApply(uberEnv.Provider, uberEnv.TerraformAfter, params)
	}
}

/***********************************************************************
 * Generate a sample config to get started
 */
func (p DefaultProvider) GenerateSampleConfiguration() {
	configPath := fmt.Sprintf("%s/config.yml", utils.GetUberState())
	if _, err := os.Stat(configPath); err == nil {
		message := fmt.Sprintf("Configuration exists at %s. Continue? (yes/no)", configPath)
		if !utils.Confirm(message, "yes") {
			fmt.Println("Quitting")
			os.Exit(1)
		}
	}

	configuration := strings.TrimSpace(sampleConfiguration)
	for _, variable := range sampleVariables {
		message := fmt.Sprintf("Value for %s? ", variable)
		value := utils.Ask(message)
		variableName := fmt.Sprintf("%%%s%%", variable)
		configuration = strings.Replace(configuration, variableName, value, -1)
	}
	ioutil.WriteFile(configPath, []byte(configuration), 0644)
}
