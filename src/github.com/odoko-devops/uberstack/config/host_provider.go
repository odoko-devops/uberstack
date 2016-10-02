package config

type HostProviderBase struct {
	Type string
	Name string
	Impl string
	Config map[string]string
}

type HostBase struct {
	Name string
	ProviderFilename string `yaml:"provider"`
	Config map[string]string
	HostProvider HostProvider
}

type HostProvider interface {
	LoadHost(filename string) (Host, error)

	GetConfig(name string) string
	SetConfig(name, value string)
}

type Host interface {
}

func (p *HostProviderBase) GetConfig(name string) string {
	if (p.Config == nil) {
		return ""
	} else {
		return p.Config[name]
	}
}

func (p *HostProviderBase) SetConfig(name, value string) {
	if (p.Config == nil) {
		p.Config = make(map[string]string)
	}
	p.Config[name] = value
}
