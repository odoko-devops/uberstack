package main

import (
	"fmt"
	"flag"
	"utils"
	"installer/model"
	"log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
		provider := GetProvider(config, state, providerName)
		provider.InfrastructureUp()
	case "destroy":
		provider := GetProvider(config, state, providerName)
		provider.InfrastructureDestroy()
	case "env":
		provider := GetProvider(config, state, providerName)
		GetProviderEnvironment(config, provider)
	default:
		log.Printf("Unknown action: %s\n", action)
	}
}

func processHost(config model.Config, state model.State, args []string) {
	action := args[0]
	hostName := args[1]
	hostConfig := GetHost(config, hostName)
	provider := GetHostProvider(config, state, hostConfig)

	switch action {
	case "up":
		CreateHost(config, state, provider, hostConfig)
	case "destroy":
		DestroyHost(config, state, provider, hostConfig)
	case "env":
		GetHostEnvironment(config, hostConfig)
	}
}

func main() {

	config := model.Config{}
	state := model.State{}

	bytes, err := ioutil.ReadFile(config_file)
	utils.Check(err)
	err = yaml.Unmarshal(bytes, &config)
	utils.Check(err)
	utils.ReadYaml(state_file, &state)

	fmt.Println("AUTHENTICATION")
	fmt.Println(config.Authentication)
	fmt.Println("HOSTS")
	fmt.Println(config.Hosts)

	flag.Parse()

	group := flag.Arg(0)
	switch group {
	case "provider":
		processProvider(config, state, flag.Args()[1:])
	case "host":
		processHost(config, state, flag.Args()[1:])
	default:
		fmt.Printf("Unknown group: %s\n", group)
	}
}
