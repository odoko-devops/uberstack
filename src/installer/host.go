package main

import (
	"installer/model"
	"log"
	"installer/providers/virtualbox"
	"installer/providers/amazonec2"
	"installer/providers/defaultProvider"
)

func GetProvider(config model.Config, name string) model.Provider {
	provider := model.Provider{}
	switch name {
	case "amazonaws":
		provider = amazonec2.Amazonec2{}
	case "virtualbox":
		provider = virtualbox.VirtualBox{}
	default:
		log.Panic("Unknown provider: %s", name)
	}
	providerConfig := model.ProviderConfig{}
	for i := range config.Provider {
		if config.Provider[i].Name == name {
			providerConfig = config.Provider[i]
		}
	}

	provider.Configure(config, providerConfig)
	return provider
}

func GetHost(config model.Config, name string) model.HostConfig {
	for i := range config.Hosts {
		if config.Hosts[i].Name == name {
			return config.Hosts[i]
		}
	}
	return nil
}

func GetHostProvider(config model.Config, hostConfig model.HostConfig) model.Provider {
	return GetProvider(config, hostConfig.Provider)
}

func CreateHost(config model.Config, state model.State, provider model.Provider, hostConfig model.HostConfig) {
	defaultProvider := defaultProvider.DefaultProvider{}
	provider.HostUp(hostConfig)

	defaultProvider.AddUbuntuToDockerGroup(hostConfig)
	defaultProvider.RegenerateCerts(hostConfig)
	defaultProvider.StartApps(config, state, hostConfig)
	defaultProvider.StartRancherAgent(config, state, hostConfig)

}

func DestroyHost(config model.Config, provider model.Provider, host model.HostConfig) {
	completed := provider.HostDestroy(host)
	if !completed {
		defaultProvider := defaultProvider.DefaultProvider{}
		defaultProvider.HostDestroy(host)
	}
}

func GetProviderEnvironment(config model.Config, provider model.Provider) {

}

func GetHostEnvironment(config model.Config, host model.HostConfig) {

}