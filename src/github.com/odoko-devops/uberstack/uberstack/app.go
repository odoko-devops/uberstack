package uberstack

import (
	"fmt"
	"log"
	"strings"
	"github.com/odoko-devops/uberstack/config"
)

func ProcessHost(args []string) error {
	action := args[1]
	hostName := args[2]

	state := new(config.State)
	state.Load()

	host, err := LoadHost(hostName, state)
	if err != nil {
		return err
	}

	provider := host.GetHostProvider()
	if provider == nil {
		return fmt.Errorf("Host %s does not have a configured host provider", hostName)
	}
	switch action {
	case "up":
		outputs, hostOutputs, err := provider.CreateHost(host)
		if err != nil {
			return err
		}
		for k, v := range outputs {
			state.SetValue(k, v)
		}
		for k,v := range hostOutputs {
			state.SetHostValue(host, k, v)
			log.Printf("Set %s=%s for %s", k, v, host.GetName())
		}
		log.Printf("Host %s created.", host.GetName())
	case "ssh":
		_, err = provider.Execute(host, strings.Join(args[4:], " "), nil)
	case "rm", "destroy":
		err = provider.DeleteHost(host)
		log.Printf("Host %s deleted.", host.GetName())
	default:
		log.Printf("Unknown action:", action)
	}
	if err != nil {
		return err
	}
	err = state.Save()
	return err
}

func ProcessApp(args []string) error {
	action := args[1]
	log.Printf("Action: %s", action)
	appName := args[2]
	log.Printf("App: %s", appName)
	envName := args[3]
	log.Printf("Env: %s", envName)

	state := new(config.State)
	state.Load()

	app, err := LoadApp(appName, state)
	if err != nil {
		return err
	}
	provider := app.GetAppProvider()
	if provider == nil {
		return fmt.Errorf("App %s does not have a configured app provider", appName)
	}
	if len(args)>5 {
		hostName := args[5]
		host, err := LoadHost(hostName, state)
		if err != nil {
			return err
		}
		app.SetHost(host)
	}
	switch action {
	case "up":
		err := provider.ConnectHost(app.GetHost())
		if err != nil {
			return err
		}
		err = provider.StartApp(app, envName, nil)
		if err != nil {
			return err
		}
		log.Printf("App %s started.", app.GetName())
	case "rm", "destroy":
		err := provider.StopApp(app, envName)
		if err != nil {
			return err
		}
		log.Printf("App %s stoped.", app.GetName())
	default:
		log.Printf("Unknown action:", action)
	}
	err = state.Save()
	return err
}