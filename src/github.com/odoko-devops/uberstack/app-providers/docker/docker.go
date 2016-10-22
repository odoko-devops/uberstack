package docker

import (
	"github.com/odoko-devops/uberstack/config"
	"github.com/odoko-devops/uberstack/utils"
	"fmt"
	"log"
	"io/ioutil"
)

type DockerAppProvider struct {
	config.AppProviderBase `yaml:",inline"`
}

type DockerApp struct {
	config.AppBase `yaml:",inline"`

	DockerCompose string `yaml:"docker-compose"`
	Hostname string
}

func LoadAppProvider(filename string) (config.AppProvider, error) {
	provider := new(DockerAppProvider)
	err := utils.ReadYamlFile(filename, provider)
	return provider, err
}

func (p *DockerAppProvider) LoadApp(filename string) (config.App, error) {
	host := new(DockerApp)
	err := utils.ReadYamlFile(filename, host)
	if (err != nil) {
		return nil, err
	}
	host.AppProvider = p
	return host, nil
}

func (p *DockerAppProvider) ConnectHost(host config.Host) error {
	if p.State.GetHostValue(host, "connected") == "" {
		// Here we need to install docker-compose
		provider := host.GetHostProvider()
		provider.Execute(host, "sudo apt-get install -y python-pip", nil)
		provider.Execute(host, "sudo pip install docker-compose", nil)
		p.State.SetHostValue(host, "connected", "connected")
	}
	return nil
}

func (p *DockerAppProvider) DisconnectHost(host config.Host) error {
	// nothing to do.
	return nil
}

func (p *DockerAppProvider) StartApp(a config.App, envName string, env config.ExecutionEnvironment) error {
	app := a.(*DockerApp)

	log.Printf("Starting %s", app.GetName())
	if app.Host == nil {
		return fmt.Errorf("App %s requires a hostname to start", app.GetName())
	}

	if env == nil {
		env = config.ExecutionEnvironment{}
	}
	innerEnv := app.Environments[envName].Environment
	for k, v := range innerEnv {
		env[k] = v
	}
	provider := app.Host.GetHostProvider()
	if app.DockerCompose != "" {
		dockerComposePath, err := utils.Resolve(fmt.Sprintf("%s/docker-compose.yml", app.DockerCompose), false)
		if err != nil {
			return err
		}
		composeBytes, err := ioutil.ReadFile(dockerComposePath)
		if err != nil {
			return err
		}
		compose := provider.Resolve(string(composeBytes), env)
		err = provider.UploadScript(app.Host, compose, "/tmp/docker-compose.yml")
		if err != nil {
			return err
		}
		command := fmt.Sprintf("docker-compose -f /tmp/docker-compose.yml -p %s up -d", app.GetName())
		err = provider.Execute(app.Host, command, env)
		if err != nil {
			return err
		}
	}
	log.Printf("Started %s", a.GetName())

	err := p.StartDependentApps(app, envName, env)
	return err
}

func (p *DockerAppProvider) StopApp(app config.App, envName string) error {
//	err := p.StopDependentApps(app, envName, env)
//	return err
	return nil
}
