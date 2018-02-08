package config

import (
	"log"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
)

type HostProviderBase struct {
	Type         string
	Name         string
	Impl         string
	PublicSSHKey string
	SshUser      string `yaml:"ssh-user"`
	Config       map[string]string
	State        *State
}

type HostBase struct {
	Name                 string
	HostProviderFilename string `yaml:"host-provider"`
	Config               map[string]string
	HostProvider         HostProvider
	HostName             string `yaml:"host-name"` // this is the name or IP via which the host is accessible
	Required             []string
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

	CreateHost(Host, ExecutionEnvironment) (map[string]string, map[string]string, error)
	DeleteHost(Host) (error)

	UploadFile(host Host, filename string, destination string, env ExecutionEnvironment) error
	UploadScript(host Host, script string, destination string, env ExecutionEnvironment) error
	Execute(host Host, command string, env ExecutionEnvironment) ([]byte, error)
	ExecuteWithRetrieve(host Host, command string) (string, error)
}

type Host interface {
	GetName() string
	GetHostProvider() HostProvider
	GetHostName() string // this is a name or IP via which the host is accessible
	ConfirmRequired(env ExecutionEnvironment) error
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
func (p *HostProviderBase) Execute(host Host, command string, env ExecutionEnvironment) ([]byte, error) {
	hostName := p.Resolve(host.GetHostName(), env)
	log.Printf("Resolved hostname %s to %s", host.GetHostName(), hostName)
	signer, err := getKeyFile()
	if err != nil {
		return nil, err
	}

	log.Printf("Executing %s", command)
	output, err := executeBySSH(hostName, p.SshUser, signer, command)
	return output, err
}

func (p *HostProviderBase) ExecuteWithRetrieve(host Host, command string) (string, error) {
	return "", nil
}

func (p *HostProviderBase) UploadFile(host Host, filename string, destination string, env ExecutionEnvironment) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Failed to load %s: %s", filename, err)
	}
	return p.UploadScript(host, string(data), destination, env)
}

func (p *HostProviderBase) UploadScript(host Host, script string, destination string, env ExecutionEnvironment) error {
	hostName := p.Resolve(host.GetHostName(), env)
	log.Printf("Resolved hostname %s to %s", host.GetHostName(), hostName)
	signer, err := getKeyFile()
	if err != nil {
		return err
	}
	uploadViaSCP(hostName, p.SshUser, signer, script, destination)
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

func (h *HostBase) ConfirmRequired(env ExecutionEnvironment) error {
	missing := []string{}
	for _, variable := range h.Required {
		ok := false
		if env != nil {
			_, ok = env[variable]
		}
		envValue := os.Getenv(variable)
		if !ok && envValue == "" {
			missing = append(missing, variable)
		}
	}
	if len(missing) == 0 {
		return nil
	} else {
		return fmt.Errorf("Please provide required variables: %s", strings.Join(missing, ", "))
	}

	return nil
}
