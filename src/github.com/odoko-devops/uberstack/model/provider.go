package model

type Provider interface {
	Configure(config Config, state *State, provider ProviderConfig) error

	SampleConfiguration() error

	InfrastructureUp() error
	InfrastructureDestroy() error

	HostUp(host HostConfig, state *State) error
	HostDestroy(host HostConfig, state *State) (bool, error)
}

