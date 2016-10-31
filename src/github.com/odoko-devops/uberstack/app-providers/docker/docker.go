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
	DockerType string `yaml:"docker-type"`
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
	return nil
}

func (p *DockerAppProvider) DisconnectHost(host config.Host) error {
	// nothing to do.
	return nil
}

func (p *DockerAppProvider) StartApp(a config.App, envName string, env config.ExecutionEnvironment) error {
	app := a.(*DockerApp)

	log.Printf("Starting %s", app.GetName())
	if app.Host == nil && ! (p.DockerType == "local") {
		return fmt.Errorf("App %s requires a hostname to start", app.GetName())
	}
	env = app.GetEnvironment(envName, env)
	err := app.ConfirmRequired(env)
	if err != nil {
		return err
	}
	err = p.StartDependentApps(app, envName, env)
	if err != nil {
		return err
	}
	if app.DockerCompose != "" {
		dockerComposePath, err := utils.Resolve(fmt.Sprintf("%s/docker-compose.yml", app.DockerCompose), false)
		if err != nil {
			return err
		}
		if p.DockerType == "local" {
			command := fmt.Sprintf("docker-compose -f %s -p %s up -d", dockerComposePath, app.GetName())
			log.Printf("LOCAL COMMAND: %s", command)
			output, err := utils.Execute(command, env, "")
			if err != nil {
				return err
			}
			p.ResolveOutputs(app, output)
		} else if p.DockerType == "docker-machine" {
			command := fmt.Sprintf("docker-machine env %s", app.Host.GetHostName())
			script, err := utils.ExecuteAndRetrieve(command, nil, "")
			if err != nil {
				return err
			}
			composeBytes, err := ioutil.ReadFile(dockerComposePath)
			if err != nil {
				return err
			}
			provider := app.Host.GetHostProvider()
			compose := provider.Resolve(string(composeBytes), env)
			ioutil.WriteFile("/tmp/uberstack-docker-compose.yml", []byte(compose), 0644)
			command = fmt.Sprintf("docker-machine scp /tmp/uberstack-docker-compose.yml %s:/tmp/", app.Host.GetHostName())
			_, err = utils.Execute(command, env, "")
			if err != nil {
				return err
			}
			script = fmt.Sprintf(`#!/bin/sh
						%s
						docker-compose -f /tmp/uberstack-docker-compose.yml -p %s up -d`,
				script, app.GetName())
			ioutil.WriteFile("/tmp/uberstack-docker.sh", []byte(script), 0755)

			command = fmt.Sprintf("docker-machine scp /tmp/uberstack-docker.sh %s:/tmp/uberstack-docker.sh", app.Host.GetHostName())
			_, err = utils.Execute(command, env, "")
			if err != nil {
				return err
			}
			command = fmt.Sprintf("docker-machine ssh %s /tmp/uberstack-docker.sh", app.Host.GetHostName())
			log.Printf("MACHINE COMMAND: %s", command)
			output, err := utils.Execute(command, env, "")
			if err != nil {
				return err
			}
			p.ResolveOutputs(app, output)
		} else {
			log.Printf("Running remote app on %s", app.GetHost().GetHostName())
			composeBytes, err := ioutil.ReadFile(dockerComposePath)
			if err != nil {
				return err
			}
			provider := app.Host.GetHostProvider()
			compose := provider.Resolve(string(composeBytes), env)
			err = provider.UploadScript(app.Host, compose, "/tmp/docker-compose.yml")
			if err != nil {
				return err
			}
			command := fmt.Sprintf("docker-compose -f /tmp/docker-compose.yml -p %s up -d", app.GetName())
			output, err := provider.Execute(app.Host, command, env)
			if err != nil {
				return err
			}
			p.ResolveOutputs(app, output)
			return err
		}
	}
	log.Printf("Started %s", a.GetName())
	return nil
}

func (p *DockerAppProvider) StopApp(app config.App, envName string) error {
//	err := p.StopDependentApps(app, envName, env)
//	return err
	return nil
}
