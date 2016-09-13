package main

import (
	"fmt"
	"flag"
	"installer/model"
	"log"
	"installer/providers/defaultProvider"
	"installer/providers/amazonec2"
	"installer/providers/virtualbox"
	"strings"
	"os"
	"utils"
)

/*
  Sample usage (where binary is called 'launcher'):
  	  uberstack init
	  uberstack provider up amazonec2
	  uberstack provider destroy amazonec2

	  uberstack host up management
	  uberstack host destroy management
	  uberstack host up docker01

	  uberstack provider up virtualbox
	  uberstack host up local-management
	  uberstack host up local-docker

	  uberstack app up myapp local
	  uberstack app up myapp dev
 */

func CreateHost(config model.Config, state *model.State, provider model.Provider,
		hostConfig model.HostConfig, skip *model.SkipList) {

	log.Printf("Creating host %s\n", hostConfig.Name)
	defaultProvider := defaultProvider.DefaultProvider{}
	providerConfig := model.GetProviderConfigForHost(config, hostConfig)

	if !skip.Avoid(model.SkipHost) {
		provider.HostUp(hostConfig, state)
		defaultProvider.AddUbuntuToDockerGroup(hostConfig)
		defaultProvider.RegenerateCerts(hostConfig)
	}

	if !skip.Avoid(model.SkipUpload) {
		defaultProvider.UploadSelf(hostConfig)
	}

	if !skip.Avoid(model.SkipApps) {
		defaultProvider.StartApps(config, state, hostConfig, skip)
	}
	if hostConfig.RancherAgent && !skip.Avoid(model.SkipRancherAgent) {
		defaultProvider.StartRancherAgent(config, state, providerConfig, hostConfig)
	}

}

func DestroyHost(config model.Config, state *model.State, provider model.Provider, host model.HostConfig) {

	completed, _ := provider.HostDestroy(host, state)
	if !completed {
		defaultProvider := defaultProvider.DefaultProvider{}
		defaultProvider.HostDestroy(host)
	}
	log.Printf("Destroyed host %s\n", host.Name)
}

func ListHosts() {
	defaultProvider := defaultProvider.DefaultProvider{}
	defaultProvider.ListHosts()
}

func GetProviderEnvironment(state *model.State, provider model.ProviderConfig) {
	env := model.GetRancherEnvironment(state, provider)
	for k,v := range env {
		fmt.Printf("export %s=%s\n", k, v)
	}
}


func GetHostEnvironment(state *model.State, host model.HostConfig) {
	env := model.GetDockerEnvironment(state, host)
	for k,v := range env {
		fmt.Printf("export %s=%s\n", k, v)
	}
}

func GetProvider(config model.Config, state *model.State, name string) model.Provider {

	var provider model.Provider

	switch name {
	case "amazonec2":
		provider = &amazonec2.Amazonec2{}
	case "virtualbox":
		provider = &virtualbox.VirtualBox{}
	default:
		log.Panic("Unknown provider: ", name)
	}
	providerConfig := model.GetProviderConfig(config, name)
	provider.Configure(config, state, providerConfig)
	return provider
}

func GetHostProvider(config model.Config, state *model.State, hostConfig model.HostConfig) model.Provider {
	return GetProvider(config, state, hostConfig.Provider)
}

func processProvider(config model.Config, state *model.State, args []string, skip *model.SkipList) {
	action := args[0]
	providerName := args[1]

	switch action {
	case "up":
		provider := GetProvider(config, state, providerName)
		provider.InfrastructureUp()
	case "destroy":
		provider := GetProvider(config, state, providerName)
		provider.InfrastructureDestroy()
	case "env":
		providerConfig := model.GetProviderConfig(config, providerName)
		GetProviderEnvironment(state, providerConfig)
	default:
		log.Printf("Unknown action: %s\n", action)
	}
}

func processHost(config model.Config, state *model.State, args []string, skip *model.SkipList) {
	action := args[0]

	switch action {
	case "up":
		hostName := args[1]
		hostConfig := model.GetHostConfig(config, hostName)
		provider := GetHostProvider(config, state, hostConfig)
		CreateHost(config, state, provider, hostConfig, skip)
	case "destroy", "rm":
		hostName := args[1]
		hostConfig := model.GetHostConfig(config, hostName)
		provider := GetHostProvider(config, state, hostConfig)
		DestroyHost(config, state, provider, hostConfig)
	case "replace", "re":
		hostName := args[1]
		hostConfig := model.GetHostConfig(config, hostName)
		provider := GetHostProvider(config, state, hostConfig)
		DestroyHost(config, state, provider, hostConfig)
		CreateHost(config, state, provider, hostConfig, skip)
	case "ls", "list":
		ListHosts()
	case "env":
		hostName := args[1]
		hostConfig := model.GetHostConfig(config, hostName)
		GetHostEnvironment(state, hostConfig)
	default:
		log.Printf("Unknown host action: %s\n", action)
	}
}

func processApp(config model.Config, state *model.State, args []string, skip *model.SkipList) {
	uberHome := os.Getenv("UBER_HOME")
	if uberHome == "" {
		println("Please set UBER_HOME.")
		os.Exit(1)
	}

	if len(args) < 3 {
		println("Usage: app <action> <uberstack-name> <environment name>")
		os.Exit(1)
	}

	action := args[0]
	uberstackName := args[1]
	environment := args[2]

	cmd := ""
	desc := ""
	switch action {
	case "up":
		cmd = "up -d"
		desc = "Installing"
	case "upgrade":
		cmd = "up --upgrade --pull -d " + strings.Join(flag.Args()[3:], " ")
		desc = "Upgrading"
	case "confirm-upgrade":
		cmd = "up --upgrade --confirm-upgrade"
		desc = "Confirming"
	case "rollback":
		cmd = "up --upgrade --rollback"
		desc = "Rolling back"
	case "rm":
		if environment != "local" {
			var answer string
			fmt.Print("Retype uberstack name to confirm deletion: ")
			fmt.Scanln(&answer)
			if answer != uberstackName {
				fmt.Println("Confirmation failed, quitting")
				os.Exit(1)
			}
		}
		cmd = "rm --force"
		desc = "Removing"
	default:
		fmt.Printf("Unknown action: %s", action)
		os.Exit(1)
	}
	defaultProvider := defaultProvider.DefaultProvider{}
	uberstack := defaultProvider.GetUberstack(uberHome, uberstackName)
	fmt.Printf("%s uberstack %s\n", desc, uberstack.Name)

	defaultProvider.ProcessUberstack(config, state, uberHome, uberstack, environment, cmd, "")

}

func processInit(config model.Config, state *model.State, args []string, skip *model.SkipList) {
	utils.Download("docker-machine")
	utils.Download("rancher-compose")
	utils.Download("terraform")
}

func main() {

	skipString := flag.String("skip", "", "Process to skip")
	flag.Parse()

	uberState := utils.GetUberState()
	configFile := uberState + "/config.yml"
	stateFile := uberState + "/state.yml"
	config := model.LoadConfig(configFile)
	state := model.LoadState(stateFile)

	skipOptions := new(model.SkipList)
	skipOptions = skipOptions.Configure(*skipString)

	group := flag.Arg(0)
	switch group {
	case "init":
		processInit(config, state, flag.Args()[1:], skipOptions)
	case "provider":
		processProvider(config, state, flag.Args()[1:], skipOptions)
	case "host":
		processHost(config, state, flag.Args()[1:], skipOptions)
	case "app":
		processApp(config, state, flag.Args()[1:], skipOptions)
	default:
		fmt.Printf("Unknown group: %s\n", group)
	}

	model.SaveState(stateFile, state)
}
