package model

type Provider interface {
	Configure(config Config, provider ProviderConfig) bool
	WriteState(config Config, stateFile string)

	InfrastructureUp() bool
	InfrastructureDestroy() bool

	HostUp(host HostConfig) bool
	HostDestroy(host HostConfig) bool
}

