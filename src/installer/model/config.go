package model

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
	ProviderConfig map[string]string `yaml: "provider-config"`
	RancherAgent bool
	Apps []AppConfig
	Labels map[string] string
}

type AppConfig struct {
	Name string
	Config map[string]string
}
