package model

type State struct {
	Provider map[string]RancherAccess
	ProviderState map[string]ProviderState
	HostState map[string]HostState
}

type RancherAccess struct {
	AccessKey string
	SecretKey string
}

type ProviderState map[string]string
type HostState map[string]string