package main

import (
	"fmt"
	"github.com/odoko-devops/uberstack/config"
	"github.com/odoko-devops/uberstack/providers/amazonec2"
	"github.com/odoko-devops/uberstack/providers/virtualbox"
	"github.com/odoko-devops/uberstack/utils"
)

func loadHostProvider(filename string) (config.HostProvider, error) {
	baseProvider := config.HostProviderBase{}

	err := utils.ReadYamlFile(filename, &baseProvider)
	if err != nil {
		return nil, err
	}

	if baseProvider.Type != "host-provider" {
		return nil, fmt.Errorf("%s is not a host provider configuration", filename)
	}

	var provider config.HostProvider
	switch baseProvider.Impl {
	case "amazonec2":
		provider, err = amazonec2.LoadHostProvider(filename)
	case "virtualbox":
		provider, err = virtualbox.LoadHostProvider(filename)
	default:
		err = fmt.Errorf("Provider not known: %s", baseProvider.Impl)
	}
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func LoadHost(filename string) (config.Host, error) {
	baseHost := config.HostBase{}

	err := utils.ReadYamlFile(filename, &baseHost)
	if err != nil {
		return nil, err
	}

	provider, err := loadHostProvider(baseHost.ProviderFilename)
	if err != nil {
		return nil, err
	}
	host, err := provider.LoadHost(filename)
	if err != nil {
		return nil, err
	}
	return host, nil
}