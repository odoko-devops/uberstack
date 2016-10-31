package rancher

import (
	"github.com/odoko-devops/uberstack/config"
	"github.com/odoko-devops/uberstack/utils"
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"io"
	"log"
)

type RancherHostProvider struct {
	config.HostProviderBase `yaml:",inline"`

	RancherHost string `yaml:"rancher-host"`
	AccessKey   string `yaml:"access-key"`
	SecretKey   string `yaml:"secret-key"`

	AmazonEC2 *AmazonEc2Config `yaml:"amazonec2"`
}

type AmazonEc2Config struct {

	AccessKey     string `yaml:"access-key"`
	SecretKey     string `yaml:"secret-key"`

	Region        string
	Zone          string
	VpcId         string `yaml:"vpc-id"`
	SubnetId      string `yaml:"subnet-id"`

	Ami           string
	RootSize      int
	SecurityGroup string `yaml:"security-group"`
	SshKeyPath    string `yaml:"ssh-keypath"`
	SshUser       string `yaml:"ssh-user"`
	IamProfile    string `yaml:"iam-profile"`
	SpotPrice     string `yaml:"spot-price"`
	InstanceType  string `yaml:"instance-type"`
	DeviceName    string `yaml:"device-name"`
	VolumeType    string `yaml:"volume-type"`
	Type          string

}

type RancherHost struct {
	config.HostBase `yaml:",inline"`

	AmazonEc2 AmazonEc2Config

	Labels           map[string]string
	Interface        string
	Environment      string // default: Default
}

func LoadHostProvider(filename string) (config.HostProvider, error) {
	provider := new(RancherHostProvider)
	err := utils.ReadYamlFile(filename, provider)
	if (err != nil) {
		return nil, err
	}
	return provider, nil
}

func (p *RancherHostProvider) LoadHost(filename string) (config.Host, error) {
	host := new(RancherHost)
	host.AmazonEc2 = AmazonEc2Config{
		DeviceName: "/dev/sda1",
		VolumeType: "gb2",
		SshUser: "rancher",
		SecurityGroup: "rancher-machine",
		RootSize: 16,
		SpotPrice: "",
	}

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

func (p *RancherHostProvider) resolve(providerValue, hostValue string) string {
	if hostValue != "" {
		return p.Resolve(hostValue, nil)
	} else {
		return p.Resolve(providerValue, nil)
	}
}
func (p *RancherHostProvider) CreateHost(h config.Host) (map[string]string, map[string]string, error) {

	host := h.(*RancherHost)

	log.Printf("HOST: %s", host)
	log.Printf("AWS: %s", host.AmazonEc2)
	rancherHost := p.Resolve(p.RancherHost, nil)
	accessKey := p.Resolve(p.AccessKey, nil)
	secretKey := p.Resolve(p.SecretKey, nil)

	err := host.ConfirmRequired(nil)
	if err != nil {
		return nil, nil, err
	}

	environment := host.Environment
	if environment == "" {
		environment = "Default"
	}

	hostData := map[string]interface{}{
		"type": "machine",
		"name": host.GetName(),
		"labels": host.Labels,
	}
	if p.AmazonEC2 != nil {
		pAws := p.AmazonEC2
		hAws := host.AmazonEc2

		rootSize := pAws.RootSize
		if rootSize == 0 {
			rootSize = hAws.RootSize
		}
		hostData["amazonec2Config"] = map[string]interface{} {
			"accessKey": p.resolve(pAws.AccessKey, hAws.AccessKey),
			"secretKey": p.resolve(pAws.SecretKey, hAws.SecretKey),
			"ami": p.resolve(pAws.Ami, hAws.Ami),
			"deviceName": "/dev/sda1",
			"iamInstanceProfile": p.resolve(pAws.IamProfile, hAws.IamProfile),
			"instanceType": p.resolve(pAws.InstanceType, hAws.InstanceType),
			"region": p.resolve(pAws.Region, hAws.Region),
			"rootSize": rootSize,
			"securityGroup": p.resolve(pAws.SecurityGroup, hAws.SecurityGroup),
			"sessionToken": "",
			"spotPrice": "",
			"sshKeypath": p.resolve(pAws.SshKeyPath, hAws.SshKeyPath),
			"sshUser": p.resolve(pAws.SshUser, hAws.SshUser),
			"subnetId": p.resolve(pAws.SubnetId, hAws.SubnetId),
			"tags": "",
			"volumeType": "gp2",
			"vpcId": p.resolve(pAws.VpcId, hAws.VpcId),
			"zone": p.resolve(pAws.Zone, hAws.Zone),
			"type": "amazonec2Config",
		}
	}
	hostBytes, err := json.Marshal(&hostData)
	if err != nil {
		return nil, nil, err
	}
log.Println(string(hostBytes))
	log.Printf("GETTING ENVIRONMENT FROM %s", p.RancherHost)
	rancherEnv := getEnvironmentId(rancherHost, accessKey, secretKey, environment)
	url := fmt.Sprintf("http://%s/v1/projects/%s/machine", rancherHost, rancherEnv)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(hostBytes))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(accessKey, secretKey)
	log.Printf("Creating host %s...", host.GetName())
	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode > 299 {
		return nil, nil, fmt.Errorf("Create node failed: %s", resp.Status)
	}

	type Result struct {
		Id string
	}
	result := Result{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if (err != nil && err != io.EOF) {
		return nil, nil, err
	}
	log.Println("Created host id:%s", result.Id)
	return map[string]string{}, map[string]string{}, nil
}

func (p *RancherHostProvider) DeleteHost(h config.Host) (error) {
	return nil
}