package main

import (
	"fmt"
	"flag"
	"installer/model"
	"log"
	"installer/providers/defaultProvider"
	"installer/providers/amazonec2"
	"installer/providers/virtualbox"
)

const path = "/usr/local/bin:/bin:/usr/bin"
const config_file = "/state/config.yml"
const state_file = "/state/state.yml"

/*
  Sample usage (where binary is called 'launcher'):
	  launcher provider up amazonec2
	  launcher provider destroy amazonec2

	  launcher host up management
	  launcher host destroy management
	  launcher host up docker01

	  launcher provider up virtualbox
	  launcher host up local-management
	  launcher host up local-docker
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

func List() {
	defaultProvider := defaultProvider.DefaultProvider{}
	defaultProvider.List()
}

func GetProviderEnvironment(state *model.State, provider model.ProviderConfig) {
	defaultProvider := defaultProvider.DefaultProvider{}
	defaultProvider.GetRancherEnvironment(state, provider)
}


func GetHostEnvironment(config model.Config, host model.HostConfig) {

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
	case "destroy":
		hostName := args[1]
		hostConfig := model.GetHostConfig(config, hostName)
		provider := GetHostProvider(config, state, hostConfig)
		DestroyHost(config, state, provider, hostConfig)
	case "ls":
		List()
	case "list":
		List()
	case "env":
		hostName := args[1]
		hostConfig := model.GetHostConfig(config, hostName)
		GetHostEnvironment(config, hostConfig)
	}
}

func main() {

	skipString := flag.String("skip", "", "Process to skip")
	config := model.LoadConfig(config_file)
	state := model.LoadState(state_file)

	flag.Parse()

	skipOptions := new(model.SkipList)
	skipOptions = skipOptions.Configure(skipString)

	group := flag.Arg(0)
	switch group {
	case "provider":
		processProvider(config, state, flag.Args()[1:], skipOptions)
	case "host":
		processHost(config, state, flag.Args()[1:], skipOptions)
	default:
		fmt.Printf("Unknown group: %s\n", group)
	}

	model.SaveState(state_file, state)
}
