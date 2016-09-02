package model

type Provider interface {
	Configure(config Config, state State, provider ProviderConfig) (Provider, error)
	WriteState(config Config, stateFile string) error

	InfrastructureUp() error
	InfrastructureDestroy() error

	HostUp(host HostConfig, state State) error
	HostDestroy(host HostConfig, state State) (bool, error)
}

