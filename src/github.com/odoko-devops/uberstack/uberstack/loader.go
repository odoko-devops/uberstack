package uberstack

import (
	"fmt"
	"log"
	"github.com/odoko-devops/uberstack/config"
	"github.com/odoko-devops/uberstack/utils"
	"github.com/odoko-devops/uberstack/host-providers/virtualbox"
	"github.com/odoko-devops/uberstack/host-providers/terraform"
	"github.com/odoko-devops/uberstack/app-providers/docker"
)

func LoadHostProvider(filename string, state *config.State) (config.HostProvider, error) {
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
	case "terraform":
		provider, err = terraform.LoadHostProvider(filename)
	case "virtualbox":
		provider, err = virtualbox.LoadHostProvider(filename)
	default:
		err = fmt.Errorf("Provider not known: %s", baseProvider.Impl)
	}
	if err != nil {
		return nil, err
	}
	provider.SetState(state)
	log.Printf("Provider %s has state %s", provider.GetName(), provider.GetState())
	return provider, nil
}

func LoadHost(filename string, state *config.State) (config.Host, error) {
	log.Printf("Loading host from %s\n", filename)
	baseHost := config.HostBase{}

	err := utils.ReadYamlFile(filename, &baseHost)
	if err != nil {
		return nil, err
	}

	provider, err := LoadHostProvider(baseHost.HostProviderFilename, state)
	if err != nil {
		return nil, err
	}
	host, err := provider.LoadHost(filename)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func LoadAppProvider(filename string, state *config.State) (config.AppProvider, error) {
	baseProvider := config.AppProviderBase{}

	err := utils.ReadYamlFile(filename, &baseProvider)
	if err != nil {
		return nil, err
	}

	if baseProvider.Type != "app-provider" {
		return nil, fmt.Errorf("%s is not an app provider configuration", filename)
	}

	var provider config.AppProvider
	switch baseProvider.Impl {
	case "docker":
		provider, err = docker.LoadAppProvider(filename)
	case "rancher":
	default:
		err = fmt.Errorf("Provider not known: %s", baseProvider.Impl)
	}
	provider.SetState(state)
	return provider, err
}
func LoadApp(filename string, state *config.State) (config.App, error) {
	log.Printf("Loading app from %s\n", filename)
	baseApp := config.AppBase{}

	err := utils.ReadYamlFile(filename, &baseApp)
	if err != nil {
		return nil, err
	}

	provider, err := LoadAppProvider(baseApp.AppProviderFilename, state)
	if err != nil {
		return nil, err
	}
	app, err := provider.LoadApp(filename)
	if err != nil {
		return nil, err
	}

	if app.GetHostName() != "" {
		log.Printf("Loading host %s from app %s", app.GetHostName(), app.GetName())
		host, err := LoadHost(app.GetHostName(), state)
		log.Printf("LOADED HOST")
		if err != nil {
			return nil, err
		}
		p := host.GetHostProvider()
		log.Printf("Host %s has provider %s with state %s", host.GetName(), p.GetName(), p.GetState())
		app.SetHost(host)
	}

	for _, stack := range app.GetStacks() {
		log.Printf("Loading child app %s", stack)
		childApp, err := LoadApp(stack, state)
		if err != nil {
			return nil, err
		}
		if app.GetHost() != nil && childApp.GetHost() == nil {
			log.Printf("Setting host for %s to %s", childApp.GetName(), app.GetHost().GetName())
			childApp.SetHost(app.GetHost())
		}
		app.AddApp(childApp)
	}
	for _, env := range app.GetEnvironments() {
		if env.HostProviderName != "" {
			hostProvider, err := LoadHostProvider(env.HostProviderName, state)
			if err != nil {
				return nil, err
			}
			env.HostProvider = hostProvider
		}
	}
	return app, nil
}