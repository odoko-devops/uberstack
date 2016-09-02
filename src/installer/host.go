package main

import (
	"installer/model"
	"log"
	"installer/providers/virtualbox"
	"installer/providers/amazonec2"
	"installer/providers/defaultProvider"
)

func GetProvider(config model.Config, state model.State, name string) model.Provider {

	var provider model.Provider

	switch name {
	case "amazonec2":
		provider = amazonec2.Amazonec2{}
	case "virtualbox":
		provider = virtualbox.VirtualBox{}
	default:
		log.Panic("Unknown provider: ", name)
	}
	providerConfig := model.ProviderConfig{}
	for i := range config.Providers {
		if config.Providers[i].Name == name {
			providerConfig = config.Providers[i]
			break
		}
	}
	provider, _ = provider.Configure(config, state, providerConfig)
	return provider
}

func GetHost(config model.Config, name string) model.HostConfig {
	for i := range config.Hosts {
		if config.Hosts[i].Name == name {
			return config.Hosts[i]
		}
	}
	return model.HostConfig{}
}

func GetHostProvider(config model.Config, state model.State, hostConfig model.HostConfig) model.Provider {
	return GetProvider(config, state, hostConfig.Provider)
}

func CreateHost(config model.Config, state model.State, provider model.Provider, hostConfig model.HostConfig) {
	defaultProvider := defaultProvider.DefaultProvider{}
	provider.HostUp(hostConfig, state)

	defaultProvider.AddUbuntuToDockerGroup(hostConfig)
	defaultProvider.RegenerateCerts(hostConfig)
	defaultProvider.StartApps(config, state, hostConfig)
	defaultProvider.StartRancherAgent(config, state, hostConfig)

}

func DestroyHost(config model.Config, state model.State, provider model.Provider, host model.HostConfig) {
	completed, _ := provider.HostDestroy(host, state)
	if !completed {
		defaultProvider := defaultProvider.DefaultProvider{}
		defaultProvider.HostDestroy(host)
	}
}

func GetProviderEnvironment(config model.Config, provider model.Provider) {

}

func GetHostEnvironment(config model.Config, host model.HostConfig) {

}