package config

import (
	"os"
	"fmt"
	"strings"
	"log"
)

type AppProviderBase struct {
	Type                string
	Name                string
	Impl                string
	AppProviderFilename string `yaml:"app-provider"`
	Config              map[string]string
	AppProvider         AppProvider
	State               *State
}

type AppProvider interface {
	LoadApp(filename string) (App, error)
	Resolve(text string, env ExecutionEnvironment) string
	SetState(state *State)

	ConnectHost(host Host) error
	DisconnectHost(host Host) error

	StartApp(app App, envName string, env ExecutionEnvironment) error
	StopApp(app App, envName string) error

	StartDependentApps(app App, envName string, env ExecutionEnvironment) error
}

func (p *AppProviderBase) LoadApp(filename string) (App, error) {
	return nil, nil
}

func (p *AppProviderBase) Resolve(text string, env ExecutionEnvironment) string {
	return p.State.Resolve(text, env)
}

func (p *AppProviderBase) SetState(state *State) {
	p.State = state
}

func (p *AppProviderBase) StartDependentApps(app App, envName string, env ExecutionEnvironment) error {
	env = app.GetEnvironment(envName, env)
	err := app.ConfirmRequired(env)
	if err != nil {
		return err
	}
	log.Printf("Dependent Apps: %s", app.GetApps())
	for _, childApp := range app.GetApps() {
		provider := app.GetAppProvider()
		err := provider.StartApp(childApp, envName, env)
		if err != nil {
			return err
		}
	}
	return nil
}

type ExecutionEnvironment map[string]string
type DeploymentEnvironment struct {
	HostProviderName string `yaml:"host-provider"`
	HostProvider HostProvider
	Environment  ExecutionEnvironment
}

type AppBase struct {
	Name                string
	AppProviderFilename string `yaml:"app-provider"`
	Config              map[string]string
	AppProvider         AppProvider
	HostName            string `yaml:"host-name"`
	Host                Host

	Stacks              []string
	Required            []string
	Environments        map[string]DeploymentEnvironment
	Apps                []App
}

type App interface {
	GetName() string
	GetAppProvider() AppProvider
	GetStacks() []string
	GetEnvironment(envName string, env ExecutionEnvironment) ExecutionEnvironment
	GetEnvironments() map[string]DeploymentEnvironment
	ConfirmRequired(env ExecutionEnvironment) error
	GetHostName() string
	GetHost() Host
	SetHost(host Host)
	AddApp(app App)
	GetApps() []App
}

func (a *AppBase) GetAppProvider() AppProvider {
	return a.AppProvider
}

func (a *AppBase) GetName() string {
	return a.Name
}

func (a *AppBase) GetStacks() []string {
	return a.Stacks
}

func (a *AppBase) GetEnvironment(envName string, env ExecutionEnvironment) ExecutionEnvironment {
	if env == nil {
		return a.Environments[envName].Environment
	} else {
		for k,v := range a.Environments[envName].Environment {
			env[k]=v
		}
		return env
	}
}

func (a *AppBase) GetEnvironments() map[string]DeploymentEnvironment {
	return a.Environments
}

func (a *AppBase) ConfirmRequired(env ExecutionEnvironment) error {
	missing := []string{}
	for _, variable := range a.Required {
		_, ok := env[variable]
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

func (a *AppBase) GetHost() Host {
	return a.Host
}

func (a *AppBase) SetHost(host Host) {
	a.Host = host
}

func (a *AppBase) GetHostName() string {
	return a.HostName
}

func (a *AppBase) AddApp(app App) {
	if a.Apps == nil {
		a.Apps = []App{}
	}
	a.Apps = append(a.Apps, app)
}

func (a *AppBase) GetApps() []App {
	return a.Apps
}