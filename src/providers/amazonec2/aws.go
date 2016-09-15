package amazonec2

import (
	"utils"
	"log"
	"model"
	"fmt"
	"strings"
)

type Amazonec2 struct {
	config      model.Config
	state       *model.State
	provider    model.ProviderConfig

	name        string
	accessKey   string
	secretKey   string
	region      string
	zone        string
	sshKeyPath  string

	hosts       map[string]Amazonec2Host

	vpcId       string
	subnetId    string
}

type Amazonec2Host struct {
	host          model.HostConfig
	instanceType  string
	eipAlloc      string
	securityGroup string

	instanceId    string
}

func (aws *Amazonec2) Configure(config model.Config, state *model.State, provider model.ProviderConfig) error {
	aws.config = config
	aws.state = state
	aws.provider = provider

	providerMap := aws.state.ProviderState[provider.Name]
	if providerMap == nil {
		aws.state.ProviderState[provider.Name] = map[string]string{}
	} else {
		aws.vpcId = providerMap["vpcId"]
		aws.subnetId = providerMap["subnetId"]
	}

	aws.name = provider.Name
	aws.accessKey = provider.Config["access_key"]
	aws.secretKey = provider.Config["secret_key"]
	aws.region = provider.Config["region"]
	aws.zone = provider.Config["zone"]
	aws.sshKeyPath = provider.Config["ssh_keypath"]
	if aws.sshKeyPath[0:1] != "/" {
		aws.sshKeyPath = fmt.Sprintf("%s/%s", utils.GetUberState(), aws.sshKeyPath)
	}

	aws.hosts = make(map[string]Amazonec2Host, len(config.Hosts))

	for i := range config.Hosts {
		host := config.Hosts[i]
		if host.Provider == provider.Name {
			awsHost := Amazonec2Host{}
			awsHost.host = host
			awsHost.instanceType = host.Config["instance_type"]
			awsHost.eipAlloc = host.Config["elastic_ip_allocation"]
			awsHost.securityGroup = host.Config["security_group"]

			hostState := aws.state.HostState[host.Name]
			if hostState != nil {
				awsHost.instanceId = hostState["instanceId"]
			}
			aws.hosts[host.Name] = awsHost
		}
	}

	return nil
}

func (aws *Amazonec2) SampleConfiguration() error {

	for k, v := range terraformConfig {
		fmt.Printf("Exporting configuration %s\n", k)
		utils.TerraformExport(v, aws.provider.Name, k, aws.provider.Config)
	}
	return nil
}

func (aws *Amazonec2) terraformConfig() utils.Environment {
	return utils.Environment{
		"TF_VAR_aws_access_key": aws.accessKey,
		"TF_VAR_aws_secret_key": aws.secretKey,
	}
}

func (aws *Amazonec2) InfrastructureUp() error {
	log.Println("Create AWS VPC Environment")
	utils.TerraformApply(aws.provider.Name, aws.provider.Terraform, aws.terraformConfig())

	aws.vpcId = utils.TerraformOutput(aws.provider.Name, "vpc_id")

	providerState := aws.state.ProviderState[aws.provider.Name]
	providerState["vpcId"] = aws.vpcId
	providerState["subnetId"] = aws.subnetId

	for _, output := range aws.provider.TerraformOutputs {
		model.SetTerraformState(aws.state, aws.provider.Name, output, utils.TerraformOutput(aws.provider.Name, output))
	}
	return nil
}

func (aws *Amazonec2) InfrastructureDestroy() error {
	log.Print("Destroy AWS VPC Environment")

	utils.TerraformDestroy(aws.provider.Name, aws.terraformConfig())
	utils.TerraformRemoveState(aws.provider.Name)
	return nil

}

func (aws *Amazonec2) HostUp(hostConfig model.HostConfig, state *model.State) error {

	env := aws.terraformConfig()
	utils.TerraformApply(aws.provider.Name, hostConfig.TerraformBefore, env)

	for _, output := range hostConfig.TerraformOutputsBefore {
		model.SetTerraformState(state, hostConfig.Provider, output, utils.TerraformOutput(aws.provider.Name, output))
	}

	awsHost := aws.hosts[hostConfig.Name]
	aws.createHost(awsHost)
	awsHost.instanceId = aws.getInstanceId(awsHost)

	hostState := model.GetHostState(state, hostConfig.Name)
	hostState["instanceId"] = awsHost.instanceId
	env["TF_VAR_instance_id"] = awsHost.instanceId
	for k, v := range hostConfig.Config {
		env["TF_VAR_" + k] = v
	}

	utils.TerraformApply(aws.provider.Name, hostConfig.TerraformAfter, env)
	for _, output := range hostConfig.TerraformOutputsAfter {
		model.SetTerraformState(state, hostConfig.Provider, output, utils.TerraformOutput(aws.provider.Name, output))
	}

	if _, ok := hostConfig.Config["elastic_ip_allocation"]; ok {
		fmt.Println("Getting EIP address")
		outputName := fmt.Sprintf("%s_public_ip", hostConfig.Name)
		hostState["public-ip"] = utils.TerraformOutput(hostConfig.Provider, outputName)
		fmt.Printf("Public IP for %s = %s", hostConfig.Name, hostState["public-ip"])
	} else {
		fmt.Println("Retrieving IP from docker-machine")
		command := fmt.Sprintf("docker-machine -s %s/machine inspect %s -f '{{.Driver.IPAddress}}'",
			utils.GetUberState(), awsHost.host.Name)
		hostState["public-ip"] = strings.Replace(utils.ExecuteAndRetrieve(command, nil, ""), "'", "", -1)
	}
	fmt.Printf("Public IP for %s = %s", hostConfig.Name, hostState["public-ip"])
	state.HostState[hostConfig.Name] = hostState

	return nil
}

func (aws *Amazonec2) HostDestroy(host model.HostConfig, state *model.State) (bool, error) {
	return false, nil
}

func (aws *Amazonec2) createHost(host Amazonec2Host) {
	log.Printf("Create host %s\n", host.host.Name)
	command := fmt.Sprintf(`docker-machine -s %s/machine create --driver amazonec2 \
           --amazonec2-access-key=%s \
           --amazonec2-secret-key=%s \
               --amazonec2-vpc-id=%s \
               --amazonec2-instance-type=%s \
               --amazonec2-security-group=%s \
               --amazonec2-region=%s \
               --amazonec2-zone=%s \
               --amazonec2-subnet-id=%s \
               --amazonec2-tags name,%s \
               --amazonec2-ssh-keypath=%s \
               %s`,
		utils.GetUberState(),
		aws.accessKey,
		aws.secretKey,
		aws.vpcId,
		host.instanceType,
		host.securityGroup,
		aws.region,
		aws.zone,
		model.GetHostConfigValue(host.host, aws.state, "subnet"),
		host.host.Name,
		aws.sshKeyPath,
		host.host.Name)
	utils.Execute(command, nil, "")
}

func (aws *Amazonec2) getInstanceId(host Amazonec2Host) string {
	command := fmt.Sprintf("docker-machine -s %s/machine inspect %s -f '{{.Driver.InstanceId}}'",
		utils.GetUberState(), host.host.Name)
	instanceId := strings.Replace(utils.ExecuteAndRetrieve(command, nil, ""), "'", "", -1)
	fmt.Printf("INSTANCE ID = %s\n", instanceId)
	return instanceId
}
