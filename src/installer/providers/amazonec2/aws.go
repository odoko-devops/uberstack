package amazonec2

import (
	"utils"
	"log"
	"installer/model"
	"fmt"
)

type Amazonec2 struct {
	config     model.Config
	state      *model.State
	provider   model.ProviderConfig

	name       string
	cidr       string
	accessKey  string
	secretKey  string
	region     string
	zone       string
	sshKeyPath string

	hosts      map[string]Amazonec2Host

	vpcId      string
	subnetId   string
}

type Amazonec2Host struct {
	host          model.HostConfig
	instanceType  string
	elasticIp     string
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
	aws.cidr = provider.Config["public-cidr"]
	aws.accessKey = provider.Config["access-key"]
	aws.secretKey = provider.Config["secret-key"]
	aws.region = provider.Config["region"]
	aws.zone = provider.Config["zone"]
	aws.sshKeyPath = provider.Config["ssh-keypath"]

	aws.hosts = make(map[string]Amazonec2Host, len(config.Hosts))

	for i := range config.Hosts {
		host := config.Hosts[i]
		if host.Provider == provider.Name {
			awsHost := Amazonec2Host{}
			awsHost.host = host
			awsHost.instanceType = host.Config["instance-type"]
			awsHost.elasticIp = host.Config["elastic-ip"]
			awsHost.eipAlloc = host.Config["elastic-ip-allocation"]
			awsHost.securityGroup = host.Config["security-group"]

			hostState := aws.state.HostState[host.Name]
			if hostState != nil {
				awsHost.instanceId = hostState["instanceId"]
			}
			aws.hosts[host.Name] = awsHost
		}
	}

	return nil
}

func (aws *Amazonec2) InfrastructureUp() error {
	log.Println("Create AWS VPC Environment")
	cwd := "terraform/aws"
	env := utils.Environment{
		"TF_VAR_aws_access_key": aws.accessKey,
		"TF_VAR_aws_secret_key": aws.secretKey,
	}

	utils.Execute("terraform apply -state=/state/terraform.tfstate", env, cwd)

	aws.vpcId = utils.ExecuteAndRetrieve("terraform output -state=/state/terraform.tfstate vpc_id", env, cwd)
	aws.subnetId = utils.ExecuteAndRetrieve("terraform output -state=/state/terraform.tfstate subnet_id", env, cwd)

	providerState := aws.state.ProviderState[aws.provider.Name]
	providerState["vpcId"] = aws.vpcId
	providerState["subnetId"] = aws.subnetId

	return nil
}

func (aws *Amazonec2) InfrastructureDestroy() error {
	log.Print("Destroy AWS VPC Environment")

	cwd := "terraform/aws"
	env := utils.Environment{
		"TF_VAR_aws_access_key": aws.accessKey,
		"TF_VAR_aws_secret_key": aws.secretKey,
	}

	utils.Execute("terraform destroy -state=/state/terraform.tfstate -force", env, cwd)
	return nil

}

func (aws *Amazonec2) HostUp(hostConfig model.HostConfig, state *model.State) error {
	awsHost := aws.hosts[hostConfig.Name]
	aws.createHost(awsHost)
	awsHost.instanceId = aws.getInstanceId(awsHost)
	aws.makeElasticIPAssociation(awsHost)

	hostState := aws.state.HostState[hostConfig.Name]
	hostState["instanceId"] = awsHost.instanceId
	return nil
}

func (aws *Amazonec2) HostDestroy(host model.HostConfig, state *model.State) (bool, error) {
	return false, nil
}

func (aws *Amazonec2) makeElasticIPAssociation(awsHost Amazonec2Host) error {
	if awsHost.eipAlloc != "" {
		log.Println("Associate predefined EIP with Docker Host")
		env := utils.Environment{
			"TF_VAR_aws_access_key": aws.accessKey,
			"TF_VAR_aws_secret_key": aws.secretKey,
			"TF_VAR_instance_id": awsHost.instanceId,
			"TF_VAR_allocation_id": awsHost.eipAlloc,
		}
		cwd := "terraform/aws-eip"
		utils.Execute("terraform apply", env, cwd)
	}
	return nil
}

func (aws *Amazonec2) createHost(host Amazonec2Host)  {
	log.Printf("Create host %s\n", host.host.Name)
	command := fmt.Sprintf(`docker-machine -s /state/machine create --driver amazonec2 \
           --amazonec2-access-key=%s \
           --amazonec2-secret-key=%s \
               --amazonec2-vpc-id=%s \
               --amazonec2-instance-type=%s \
               --amazonec2-security-group=%s \
               --amazonec2-region=%s \
               --amazonec2-zone=%s \
               --amazonec2-subnet-id=%s \
               --amazonec2-tags name=%s \
               --amazonec2-ssh-keypath=%s \
               %s`, aws.accessKey,
		aws.secretKey,
		aws.vpcId,
		host.instanceType,
		host.securityGroup,
		aws.region,
		aws.zone,
		aws.subnetId,
		host.host.Name,
		aws.sshKeyPath,
		host.host.Name)
	fmt.Println(command)
	utils.Execute(command, nil, "")
}

func (aws *Amazonec2) getInstanceId(host Amazonec2Host) string {
	command := fmt.Sprintf("docker-machine -s /state/machine inspect %s -f '{{.Driver.InstanceId}}'", host.host.Name)
	instanceId := utils.ExecuteAndRetrieve(command, nil, "")
	return instanceId
}

/*

def create_docker_host_with_rancher_cli(config):
step("Create Management Docker Host")
aws = config["aws"]
host = aws["docker-host"]
rancher = config["rancher"]
execute('''rancher --url http://%s/v1 \
--access-key %s \
--secret-key %s \
host create \
--driver amazonec2 \
--amazonec2-access-key %s \
--amazonec2-secret-key %s \
--amazonec2-vpc-id %s \
--amazonec2-instance-type %s \
--amazonec2-security-group management-tools \
--amazonec2-region %s \
--amazonec2-zone %s \
--amazonec2-subnet-id %s \
--amazonec2-tags name= management-tools \
--amazonec2-ssh-keypath %s \
''' % (config["apps"]["rancher"]["name"],
rancher["api-access-key"],
rancher["api-secret-key"],
aws["access-key"],
aws["secret-key"],
vpc_id,
host["instance-type"],
aws["region"],
aws["zone"],
subnet_id,
"/id_rsa"))
*/