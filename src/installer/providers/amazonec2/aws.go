package amazonec2

import (
	"utils"
	"log"
	"installer/model"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"strings"
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
	if aws.sshKeyPath[0:1] != "/" {
		aws.sshKeyPath = fmt.Sprintf("%s/%s", utils.GetUberState(), aws.sshKeyPath)
	}

	aws.hosts = make(map[string]Amazonec2Host, len(config.Hosts))

	for i := range config.Hosts {
		host := config.Hosts[i]
		if host.Provider == provider.Name {
			awsHost := Amazonec2Host{}
			awsHost.host = host
			awsHost.instanceType = host.Config["instance-type"]
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
	aws.makeElasticIPAssociation(state, awsHost)

	hostState := aws.state.HostState[hostConfig.Name]
	if hostState == nil {
		hostState = make(map[string]string)
		aws.state.HostState[hostConfig.Name] = hostState
	}
	hostState["instanceId"] = awsHost.instanceId
	state.HostState[hostConfig.Name] = hostState
	return nil
}

func (aws *Amazonec2) HostDestroy(host model.HostConfig, state *model.State) (bool, error) {
	return false, nil
}


var terraformEipAssociation string =`
provider "aws" {
  access_key = "{{.accessKey}}"
  secret_key = "{{.secretKey}}"
  region     = "{{.region}}"
}

resource "aws_eip_association" "{{.hostName}}" {
  instance_id = "{{.instanceId}}"
  allocation_id = "{{.allocationId}}"
}

output "public_id" {
  value = "${aws_eip_association.{{.hostName}}.public_ip}"
}

`
func (aws *Amazonec2) makeElasticIPAssociation(state *model.State, awsHost Amazonec2Host) error {
	hostState := model.GetHostState(state, awsHost.host.Name)
	if awsHost.eipAlloc != "" {
		log.Println("Associate predefined EIP with Docker Host")

		params := map[string]string{
			"accessKey": aws.accessKey,
			"secretKey": aws.secretKey,
			"region": aws.region,
			"hostName": awsHost.host.Name,
			"instanceId": awsHost.instanceId,
			"allocationId": awsHost.eipAlloc,
		}
		fmt.Printf("PARAMS: %v\n", params)
		createCommandTemplate, err := template.New("terraformEip").Parse(terraformEipAssociation)

		dir, err := ioutil.TempDir("", "terraform")
		fmt.Printf("Created %s\n", dir)
		utils.Check(err)

		f, err := os.Create(dir + "/eip.tf")
		utils.Check(err)

		err = createCommandTemplate.Execute(f, params)
		utils.Check(err)

		f.Close()
		cwd := dir
		utils.Execute("terraform apply", nil, cwd)
		cmd := "terraform output public_id"
		hostState["public-ip"] = utils.ExecuteAndRetrieve(cmd, nil, cwd)
		state.HostState[awsHost.host.Name] = hostState

		os.RemoveAll(dir)
	} else {
		command := fmt.Sprintf("docker-machine -s %s/machine inspect %s -f '{{.Driver.IPAddress}}'",
			utils.GetUberState(), awsHost.host.Name)
		hostState["public-ip"] = utils.ExecuteAndRetrieve(command, nil, "")
	}
	return nil
}

func (aws *Amazonec2) createHost(host Amazonec2Host)  {
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
		aws.subnetId,
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