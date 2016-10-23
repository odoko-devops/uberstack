package rancher

import (
	"github.com/odoko-devops/uberstack/config"
	"github.com/odoko-devops/uberstack/utils"
	"fmt"
	"log"
)

type RancherAppProvider struct {
	config.AppProviderBase `yaml:",inline"`
	RancherHost string `yaml:"rancher-host"`
	AccessKey   string `yaml:"access-key"`
	SecretKey   string `yaml:"secret-key"`
}

type RancherApp struct {
	config.AppBase `yaml:",inline"`

	RancherCompose string `yaml:"rancher-compose"`
}

func LoadAppProvider(filename string) (config.AppProvider, error) {
	provider := new(RancherAppProvider)
	err := utils.ReadYamlFile(filename, provider)
	return provider, err
}

func (p *RancherAppProvider) LoadApp(filename string) (config.App, error) {
	host := new(RancherApp)
	err := utils.ReadYamlFile(filename, host)
	if (err != nil) {
		return nil, err
	}
	host.AppProvider = p
	return host, nil
}

func (p *RancherAppProvider) ConnectHost(host config.Host) error {
	return nil ; // this is now handled by the rancher-host-provider
	if p.State.GetHostValue(host, "connected") == "" {
		// Here we need to install Rancher-agent
		rancherHostname := p.Resolve(p.RancherHost, nil)
		accessKey := p.Resolve(p.AccessKey, nil)
		secretKey := p.Resolve(p.SecretKey, nil)
		networkInterface := ""
		labels := ""
		ipAddress, err := identifyIpAddress(networkInterface)
		if err != nil {
			return err
		}
		log.Printf("IP address: %s\n", ipAddress)
		rancherEnvironment := identifyRancherEnvironment(rancherHostname, accessKey, secretKey)
		log.Printf("Environment: %s\n", rancherEnvironment)
		registrationUrl, err := identifyRegistrationUrl(rancherHostname, accessKey, secretKey, rancherEnvironment)
		if err != nil {
			return err
		}
		log.Printf("Registration url: %s\n", registrationUrl)
		installRancherAgent(ipAddress, labels, registrationUrl)
		p.State.SetHostValue(host, "connected", "connected")
	}
	return nil
}

func (p *RancherAppProvider) DisconnectHost(host config.Host) error {
	provider := host.GetHostProvider()
	provider.Execute(host, "docker rm -f rancher-agent", nil)
	return nil
}

func (p *RancherAppProvider) StartApp(a config.App, envName string, env config.ExecutionEnvironment) error {
	app := a.(*RancherApp)


	if env == nil {
		env = config.ExecutionEnvironment{}
	}
	env = app.GetEnvironment(envName, env)
	err := app.ConfirmRequired(env)
	if err != nil {
		return err
	}
	env = p.ResolveEnvironment(env)
	log.Printf("Rancher compose environment: %s", env)

	err = p.StartDependentApps(app, envName, env)
	if err != nil {
		return err
	}

	if app.RancherCompose != "" {
		composePath, err := utils.Resolve(app.RancherCompose, false)

		command := fmt.Sprintf(
			`rancher-compose --file %s/docker-compose.yml
                        --rancher-file %s/rancher-compose.yml
                        --project-name %s
                        --url http://%s/
                        --access-key %s
                        --secret-key %s
                        up -d
                        `, composePath,
			   composePath,
			app.GetName(),
			p.Resolve(p.RancherHost, nil),
			p.Resolve(p.AccessKey, nil),
			p.Resolve(p.SecretKey, nil),
		        )

		_, err = utils.Execute(command, env, "")
		if err != nil {
			return err
		}
	}
	log.Printf("Started %s", a.GetName())

	return nil
}

func (p *RancherAppProvider) StopApp(a config.App, envName string) error {
	app := a.(*RancherApp)

	log.Printf("Stoping %s", app.GetName())

	if app.RancherCompose != "" {
		composePath, err := utils.Resolve(app.RancherCompose, false)

		command := fmt.Sprintf(
			`rancher-compose --file %s/docker-compose.yml
                        --rancher-file %s/rancher-compose.yml
                        --project-name %s down
                        `, composePath, composePath, app.GetName())

		_, err = utils.Execute(command, nil, "")
		if err != nil {
			return err
		}
	}
	err := p.StopDependentApps(app, envName)
	return err
}
