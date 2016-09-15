package model

import (
	"utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type State struct {
	Provider      map[string]RancherAccess
	ProviderState map[string]ProviderState
	TerraformState map[string]ProviderState
	HostState     map[string]HostState
}

type RancherAccess struct {
	RancherUrl string `yaml:"rancher-url"`
	AccessKey string `yaml:"access-key"`
	SecretKey string `yaml:"secret-key"`
}

type ProviderState map[string]string
type HostState map[string]string

func LoadState(state_file string) *State {
	state := State{}
	bytes, err := ioutil.ReadFile(state_file)
	utils.Check(err)
	err = yaml.Unmarshal(bytes, &state)
	utils.Check(err)
	return &state
}

func SaveState(state_file string, state *State) {
	bytes, err := yaml.Marshal(state)
	utils.Check(err)
	err = ioutil.WriteFile(state_file, bytes, 0644)
	utils.Check(err)
}