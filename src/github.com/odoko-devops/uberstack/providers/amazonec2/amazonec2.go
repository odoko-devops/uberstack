package amazonec2

import (
	"github.com/odoko-devops/uberstack/config"
	"github.com/odoko-devops/uberstack/utils"
)

type Amazonec2HostProvider struct {
	*config.HostProviderBase

	Region             string
	Zone               string
	VpcCidr            string `yaml:"vpc_cidr"`
	SshKeypath         string `yaml:"ssh_keypath"`
	AccessKey          string `yaml:"access_key"`
	SecretKey          string `yaml:"secret_key"`
	TerraformResources []string `yaml:"terraform-resources"`
	TerraformOutputs   []string `yaml:"terraform-outputs"`
}

type Amazonec2Host struct {
	config.Host

	InstanceType string
	SecurityGroup string
	Subnet string
	IamRole string `yaml:"iam"`
}

func LoadHostProvider(filename string) (config.HostProvider, error) {
	provider := new(Amazonec2HostProvider)
	err := utils.ReadYamlFile(filename, provider)
	if (err != nil) {
		return nil, err
	}
	return provider, nil
}

func (p *Amazonec2HostProvider) LoadHost(filename string) (config.Host, error) {
	host := new(Amazonec2Host)
	err := utils.ReadYamlFile(filename, host)
	if (err != nil) {
		return nil, err
	}
	return host, nil
}