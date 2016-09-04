package model

import (
	"utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Providers []ProviderConfig
	Authentication []RealmConfig
	Hosts []HostConfig
}

type ProviderConfig struct {
	Type string
	Name string
	Config map[string]string
}

type RealmConfig struct {
	Name string
	Users []User
}

type User struct {
	Username string
	Password string
	Email string
}

type HostConfig struct {
	Name string
	Provider string
	Config map[string]string
	RancherAgent bool
	Apps []AppConfig
	Labels map[string] string
}

type AppConfig struct {
	Name string
	Config map[string]string
}

func LoadConfig(config_file string) Config {
	config := Config{}
	bytes, err := ioutil.ReadFile(config_file)
	utils.Check(err)
	err = yaml.Unmarshal(bytes, &config)
	utils.Check(err)
	return config
}

func GetProviderConfig(config Config, providerName string) ProviderConfig {
	for i := range config.Providers {
		if config.Providers[i].Name == providerName {
			return config.Providers[i]
		}
	}
	return ProviderConfig{}
}

func GetProviderConfigForHost(config Config, host HostConfig) ProviderConfig {
	return GetProviderConfig(config, host.Name)
}

func GetHostConfig(config Config, name string) HostConfig {
	for i := range config.Hosts {
		if config.Hosts[i].Name == name {
			return config.Hosts[i]
		}
	}
	return HostConfig{}
}

