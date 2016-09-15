package model

import (
	"utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"text/template"
	"bytes"
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
	Terraform []string `yaml:"terraform-resources"`
	TerraformOutputs []string `yaml:"terraform-outputs"`
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
	RancherInterface string `yaml:"rancher-interface"`
	Config map[string]string
	RancherAgent bool `yaml:"rancher-agent"`
	Apps []AppConfig
	Labels map[string] string
	TerraformBefore []string `yaml:"terraform-resources-before"`
	TerraformAfter []string `yaml:"terraform-resources-after"`
	TerraformOutputsBefore []string `yaml:"terraform-outputs-before"`
	TerraformOutputsAfter []string `yaml:"terraform-outputs-after"`
}

type AppConfig struct {
	Type string
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
	return GetProviderConfig(config, host.Provider)
}

func GetHostConfig(config Config, name string) HostConfig {
	for i := range config.Hosts {
		if config.Hosts[i].Name == name {
			return config.Hosts[i]
		}
	}
	return HostConfig{}
}

func GetAuthRealm(config Config, realmName string) RealmConfig {
	for _, realm := range config.Authentication {
		if realm.Name == realmName {
			return realm
		}
	}
	return RealmConfig{}
}

func GetHostState(state *State, name string) HostState {
	hostState := state.HostState[name]
	if hostState == nil {
		hostState = make(map[string]string)
		state.HostState[name] = hostState
	}
	return hostState
}

func GetDockerEnvironment(state *State, host HostConfig) utils.Environment {
	ip := state.HostState[host.Name]["public-ip"]
	env := utils.Environment{
		"DOCKER_TLS_VERIFY": "1",
		"DOCKER_HOST": fmt.Sprintf("tcp://%s:2376", ip),
		"DOCKER_CERT_PATH": fmt.Sprintf("%s/machine/machines/%s", utils.GetUberState(), host.Name),
		"DOCKER_MACHINE_NAME": host.Name,
	}
	return env
}

func GetRancherEnvironment(state *State, provider ProviderConfig) utils.Environment {
	providerState := state.Provider[provider.Name]
	env := utils.Environment{
		"RANCHER_URL": providerState.RancherUrl,
		"RANCHER_ACCESS_KEY": providerState.AccessKey,
		"RANCHER_SECRET_KEY": providerState.SecretKey,
	}
	return env
}

func GetHostConfigValue(host HostConfig, state *State, name string) string {
	value := host.Config[name]
	configTemplate, err := template.New(name).Parse(value)
	utils.Check(err)
	buf := bytes.Buffer{}
	params := map[string]interface{}{
		"terraform": state.TerraformState[host.Provider],
	}
	configTemplate.Execute(&buf, params)
	return buf.String()
}

func GetProviderConfigValue(provider ProviderConfig, state *State, name string) string {
	value := provider.Config[name]
	configTemplate, err := template.New(name).Parse(value)
	utils.Check(err)
	buf := bytes.Buffer{}
	params := map[string]interface{}{
		"terraform": state.TerraformState[provider.Name],
	}
	configTemplate.Execute(&buf, params)
	return buf.String()
}

func SetTerraformState(state *State, provider, output, value string) {
	terraformState := state.TerraformState[provider]
	if terraformState == nil {
		terraformState = ProviderState{}
	}
	terraformState[output] = value
	state.TerraformState[provider] = terraformState
}

type Uberstack struct {
	Name         string
	Stacks       []string
	Uberstacks   []string
	Required     []string
	Environments map[string]UberEnv
}

type UberEnv struct {
	Provider string
	Environment utils.Environment
}
