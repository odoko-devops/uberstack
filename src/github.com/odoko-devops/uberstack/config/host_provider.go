package config

import (
	"log"
	"io/ioutil"
	"fmt"
)

type HostProviderBase struct {
	Type string
	Name string
	Impl string
	PublicSSHKey string
	Config map[string]string
	State *State
}

type HostBase struct {
	Name string
	HostProviderFilename string `yaml:"host-provider"`
	Config map[string]string
	HostProvider HostProvider
	HostName string `yaml:"host-name"`// this is the name or IP via which the host is accessible
}

type HostProvider interface {
	LoadHost(filename string) (Host, error)
	Resolve(text string, env ExecutionEnvironment) string
	SetState(state *State)
	GetState() *State

	GetConfig(name string) string
	SetConfig(name, value string)
	GetName() string
	GetType() string
	GetImpl() string

	CreateHost(Host) (map[string]string, map[string]string, error)
	DeleteHost(Host) (error)

	UploadFile(host Host, filename string, destination string) error
	UploadScript(host Host, script string, destination string) error
	Execute(host Host, command string, env ExecutionEnvironment) error
	ExecuteWithRetrieve(host Host, command string) (string, error)
	InstallDockerOnUbuntu(host Host) error
}

type Host interface {
	GetName() string
	GetHostProvider() HostProvider
	GetHostName() string // this is a name or IP via which the host is accessible
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

func (p *HostProviderBase) GetName() string {
	return p.Name
}

func (p *HostProviderBase) Resolve(text string, env ExecutionEnvironment) string {
	return p.State.Resolve(text, env)
}

func (p *HostProviderBase) SetState(state *State) {
	p.State = state
}

func (p *HostProviderBase) GetState() *State {
	return p.State
}
func (p *HostProviderBase) Execute(host Host, command string, env ExecutionEnvironment) error {
	hostName := p.Resolve(host.GetHostName(), env)
	log.Printf("Resolved hostname %s to %s", host.GetHostName(), hostName)
	signer, err := getKeyFile()
	if err != nil {
		return err
	}

	log.Printf("Executing %s", command)
	return executeBySSH(hostName, "ubuntu", signer, command)
}

func (p *HostProviderBase) ExecuteWithRetrieve(host Host, command string) (string, error) {
	return "", nil
}

func (p *HostProviderBase) UploadFile(host Host, filename string, destination string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Failed to load %s: %s", filename, err)
	}
	return p.UploadScript(host, string(data), destination)
}

func (p *HostProviderBase) UploadScript(host Host, script string, destination string) error {
	hostName := p.Resolve(host.GetHostName(), nil)
	log.Printf("Resolved hostname %s to %s", host.GetHostName(), hostName)
	signer, err := getKeyFile()
	if err != nil {
		return err
	}
	uploadViaSCP(hostName, "ubuntu", signer, script, destination)
	return nil
}

func (p *HostProviderBase) GetType() string {
	return p.Type
}

func (p *HostProviderBase) GetImpl() string {
	return p.Impl
}

func (p *HostBase) GetHostProvider() HostProvider {
	return p.HostProvider
}

func (h *HostBase) GetName() string {
	return h.Name
}

func (h *HostBase) GetHostName() string {
	return h.HostName
}

func (p *HostProviderBase) InstallDockerOnUbuntu(host Host) error {
	log.Println("Installing docker on", host.GetName())
	var ubuntuDockerInstallScript=`
		sudo apt-get update &&
		sudo apt-get install -y apt-transport-https ca-certificates &&
		sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D &&
		echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list &&
		sudo apt-get update &&
		sudo apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual &&
		sudo apt-get install -y docker-engine &&
		sudo service docker start &&
		sudo gpasswd -a ubuntu docker
		`
	return p.Execute(host, ubuntuDockerInstallScript, nil)
}