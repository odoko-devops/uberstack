package model

type Config struct {
	Provider []ProviderConfig
	Authentication map[string]RealmConfig
	Hosts []HostConfig
}

type ProviderConfig struct {
	Type string
	Name string
	Config map[string]string
}

type RealmConfig struct {
	Name string
	User []User
}

type User struct {
	Username string
	Password string
	Email string
}

type HostConfig struct {
	Name string
	Provider string
	ProviderConfig map[string]string
	RancherAgent bool
	Apps []map[string]string
	Labels map[string] string
}