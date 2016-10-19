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

func (p *TerraformHostProvider) CreateHost(h config.Host) (map[string]string, map[string]string, error) {
	host := h.(*TerraformHost)

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

	env := utils.Environment{}
	for k, v := range p.Variables {
		env["TF_VAR_" + k]= os.ExpandEnv(v)
	}
	for k, v := range host.Variables {
		env["TF_VAR_" + k]= os.ExpandEnv(v)
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
	err = utils.Execute(command, env, filepath)
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
	}
	resolvedHostOutputs := map[string]string{}
	for _, output := range hostOutputList {
		command := fmt.Sprintf("terraform output %s", output)
		resolvedHostOutputs[output], err = utils.ExecuteAndRetrieve(command, nil, filepath)
		if err != nil {
			return nil, nil, err
		}
	}

	return resolvedOutputs, resolvedHostOutputs, nil
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

	env := utils.Environment{}
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
	err = utils.Execute(command, env, filepath)
	if err != nil {
		return err
	}
	return nil
}