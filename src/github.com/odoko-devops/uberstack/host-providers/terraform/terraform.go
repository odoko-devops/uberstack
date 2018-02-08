package terraform

import (
	"github.com/odoko-devops/uberstack/config"
	"github.com/odoko-devops/uberstack/utils"
	"strings"
	"fmt"
	"log"
	"os"
)

type TerraformHostProvider struct {
	config.HostProviderBase `yaml:",inline"`

	TerraformDirectory string `yaml:"terraform-directory"`
	Variables map[string]string
	Resources []string
	Outputs []string
	InstallDocker bool `yaml:"install-docker"`
}

type TerraformHost struct {
	config.HostBase `yaml:",inline"`

	Variables map[string]string
	Resources []string
	Outputs []string
}

func LoadHostProvider(filename string) (config.HostProvider, error) {
	provider := new(TerraformHostProvider)
	err := utils.ReadYamlFile(filename, provider)
	if (err != nil) {
		return nil, err
	}
	return provider, nil
}

func (p *TerraformHostProvider) LoadHost(filename string) (config.Host, error) {
	host := new(TerraformHost)
	err := utils.ReadYamlFile(filename, host)
	if (err != nil) {
		return nil, err
	}
	host.HostProvider = p
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (p *TerraformHostProvider) CreateHost(h config.Host, env config.ExecutionEnvironment) (map[string]string, map[string]string, error) {
	host := h.(*TerraformHost)

	err := host.ConfirmRequired(env)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range env {
		log.Printf("ENV %s=%s", k, v)
	}
	filepath, err := utils.Resolve(p.TerraformDirectory, false)
	if err != nil {
		return nil, nil, err
	}

	resourceTargets := make([]string, len(p.Resources)+len(host.Resources))
	for i, resource := range p.Resources {
		resourceTargets[i] = "-target=" + resource
	}
	for i, resource := range host.Resources {
		resourceTargets[i+len(p.Resources)] = "-target=" + resource
	}
	targets := strings.Join(resourceTargets, " ")

	for k, v := range p.Variables {
		env["TF_VAR_" + k]= p.Resolve(v, env)
	}
	for k, v := range host.Variables {
		env["TF_VAR_" + k]= p.Resolve(v, env)
	}

	outputList := make([]string, len(p.Outputs))
	for i, output := range p.Outputs {
		outputList[i] = output
	}
	hostOutputList := make([]string, len(host.Outputs))
	for i, output := range host.Outputs {
		hostOutputList[i] = output
	}

	command := fmt.Sprintf("terraform apply -refresh=true %s", targets)
	_, err = utils.Execute(command, env, filepath)
	if err != nil {
		return nil, nil, err
	}

	resolvedOutputs := map[string]string{}
	for _, output := range outputList {
		command := fmt.Sprintf("terraform output %s", output)
		resolvedOutputs[output], err = utils.ExecuteAndRetrieve(command, nil, filepath)
		if err != nil {
			return nil, nil, err
		}
		key := fmt.Sprintf("host.%s.%s", host.GetHostName(), output)
		log.Printf("Resolved %s to %s", key, resolvedOutputs[output])
		env[output] = resolvedOutputs[key]
	}
	resolvedHostOutputs := map[string]string{}
	for _, output := range hostOutputList {
		command := fmt.Sprintf("terraform output %s", output)
		resolvedHostOutputs[output], err = utils.ExecuteAndRetrieve(command, nil, filepath)
		if err != nil {
			return nil, nil, err
		}
		key := fmt.Sprintf("host.%s.%s", host.GetHostName(), output)
		log.Printf("Resolved %s to %s", key, resolvedHostOutputs[output])
		env[output] = resolvedHostOutputs[key]
	}

	if p.InstallDocker {
		err = p.InstallDockerOnHost(*host, p.SshUser, env)
		if err != nil {
			return nil, nil, err
		}
	}

	return resolvedOutputs, resolvedHostOutputs, nil
}

func (p *TerraformHostProvider) InstallDockerOnHost(host TerraformHost, sshUser string, env config.ExecutionEnvironment) error {
	var command string

	switch sshUser {
	case "ubuntu":
		command = `
		sudo apt-get update &&
		sudo apt-get install -y apt-transport-https ca-certificates curl &&
		sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D &&
		echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list &&
		sudo apt-get update &&
		sudo apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual &&
		sudo apt-get install -y docker-engine &&
		sudo service docker start &&
		sudo gpasswd -a ubuntu docker
		sudo curl -L "https://github.com/docker/compose/releases/download/1.8.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
		sudo chmod 755 /usr/local/bin/docker-compose`
	default:
		return fmt.Errorf("Unknown Docker installation user: %s", sshUser)
	}
	log.Println("Installing docker on", host.GetName())
	log.Println(env)
	_, err := p.Execute(&host, command, env)
	return err
}

func (p *TerraformHostProvider) DeleteHost(h config.Host) (error) {
	log.Println("HERE")
	host := h.(*TerraformHost)
	log.Println("DELETING HOST")

	filepath, err := utils.Resolve(p.TerraformDirectory, false)
	if err != nil {
		return err
	}

	resourceTargets := make([]string, len(host.Resources))
	for i, resource := range host.Resources {
		resourceTargets[i] = "-target=" + resource
	}
	targets := strings.Join(resourceTargets, " ")
	log.Printf("Targets for %v and %v are %v\n", p.Resources, host.Resources, targets)

	env := config.ExecutionEnvironment{}
	for k, v := range p.Variables {
		env["TF_VAR_" + k]= os.ExpandEnv(v)
	}
	for k, v := range host.Variables {
		env["TF_VAR_" + k]= os.ExpandEnv(v)
	}
	log.Printf("Variables for %v and %v are %v\n", p.Variables, host.Variables, env)
	for k,v := range env {
		log.Printf("export %s=%s\n", k, v)
	}

	command := fmt.Sprintf("terraform destroy -force -refresh=true %s", targets)
	log.Printf("COMMAND: %s\n", command)
	_, err = utils.Execute(command, env, filepath)
	if err != nil {
		return err
	}
	return nil
}