package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"
	"utils"
	"installer/model"
)

const path = "/usr/local/bin:/bin:/usr/bin"
const config_file = "/config.yml"
const state_file = "/state/state.yml"

	/*
	  launcher provider up amazonec2
	  launcher provider destroy amazonec2
	  launcher host up management
	  launcher host destroy management
	  launcher host up docker01
	  lancher provider up virtualbox
	  launcher host up local-management
	  lancher host up local-docker
	 */

func processProvider(config model.Config, state model.State, args []string) {
	action := args[0]
	providerName := args[1]

	switch action {
	case "up":
		provider := GetProvider(config, providerName)
		provider.InfrastructureUp()
	case "destroy":
		provider := GetProvider(config, providerName)
		provider.InfrastructureDestroy()
	case "env":
		provider := GetProvider(config, providerName)
		GetProviderEnvironment(config, provider)
	}
}

func processHost(config model.Config, state model.State, args []string) {
	action := args[0]
	hostName := args[1]
	hostConfig := GetHost(config, hostName)
	provider := GetHostProvider(config, hostConfig)

	switch action {
	case "up":
		CreateHost(config, state, provider, hostConfig)
	case "destroy":
		DestroyHost(config, provider, hostConfig)
	case "env":
		GetHostEnvironment(config, hostConfig)
	}
}

func main() {

	config := model.Config{}
	state := model.State{}

        utils.ReadYaml(config_file, &config)
	utils.ReadYaml(state_file, &state)

	uber_home := os.Getenv("UBER_HOME")
	if uber_home == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		utils.Check(err)
		uber_home = dir
	}

	rancher_url := os.Getenv("RANCHER_URL")
	fmt.Printf("Using UBER_HOME=%v\n", uber_home)
	fmt.Printf("Using RANCHER_URL=%v\n", rancher_url)

	flag.Parse()

	group := flag.Arg(1)
	switch group {
	case "provider":
		processProvider(config, state, flag.Args()[1:])
	case "host":
		processHost(config, state, flag.Args()[1:])
	default:
		fmt.Printf("Unknown group: %s\n", group)
	}
}
