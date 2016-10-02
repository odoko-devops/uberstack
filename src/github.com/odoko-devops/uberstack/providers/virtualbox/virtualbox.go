package virtualbox

import (
	"github.com/odoko-devops/uberstack/model"
	"fmt"
	"github.com/odoko-devops/uberstack/utils"
	"github.com/odoko-devops/uberstack/providers/defaultProvider"
	"github.com/odoko-devops/uberstack/config"
)

type VirtualBoxHostProvider struct {
	config.HostProviderBase

	Boot2DockerImage string `yaml:"boot2docker_image"`
	Netmask string
	Broadcast string
}

type VirtualBoxHost struct {
	config.Host

	DiskSize string
	ram string
	ip string
}

func LoadHostProvider(filename string) (config.HostProvider, error) {
	provider := new(VirtualBoxHostProvider)
	err := utils.ReadYamlFile(filename, provider)
	if (err != nil) {
		return nil, err
	}
	return provider, nil
}

func (p *VirtualBoxHostProvider) LoadHost(filename string) (config.Host, error) {
	host := new(VirtualBoxHost)
	err := utils.ReadYamlFile(filename, host)
	if (err != nil) {
		return nil, err
	}
	return host, nil
}






/**************************************************************************************************************/
/**************************************************************************************************************/
/**************************************************************************************************************/
/**************************************************************************************************************/
/**************************************************************************************************************/

type VirtualBox struct {
	Boot2DockerImage string
	NetMask string
	Broadcast string
}

func (v *VirtualBox) Configure(config model.Config, state *model.State, provider model.ProviderConfig) error {
	v.Boot2DockerImage = provider.Config["boot2docker-image"]
	v.NetMask = provider.Config["netmask"]
	v.Broadcast = provider.Config["broadcast"]
	return nil
}

func (v *VirtualBox) SampleConfiguration() error {
	return nil
}

func (v *VirtualBox) InfrastructureUp() error {
	println("Nothing to do for Virtualbox Infrastructure")
	return nil
}
func (v *VirtualBox) InfrastructureDestroy() error {
	println("Nothing to do for Virtualbox Infrastructure")
	return nil
}
func (v *VirtualBox) HostUp(host model.HostConfig, state *model.State) error {

	disk := host.Config["disk-size"]
	memory := host.Config["ram"]
	image := v.Boot2DockerImage
	command := fmt.Sprintf(`docker-machine -s %s/machine create %s \
             --driver virtualbox \
             --virtualbox-cpu-count -1 \
             --virtualbox-disk-size %s \
             --virtualbox-memory %s \
             --virtualbox-boot2docker-url=%s
             `, utils.GetUberState(), host.Name, disk, memory, image)
	utils.Execute(command, nil, "")

	defaultProvider := defaultProvider.DefaultProvider{}
	defaultProvider.UploadSelf(host)
	command = fmt.Sprintf("./uberstack-remote-agent -ip=%s -broadcast=%s -netmask=%s virtualbox",
		host.Config["ip"], v.Broadcast, v.NetMask)
	utils.ExecuteRemote(host.Name, command, nil, "")

	command = fmt.Sprintf("docker-machine -s %s/machine regenerate-certs -f %s", utils.GetUberState(), host.Name)
	utils.Execute(command, nil, "")

	hostState := state.HostState[host.Name]
	if hostState == nil {
		hostState = model.HostState{}
	}
	hostState["public-ip"] = host.Config["ip"]
	state.HostState[host.Name] = hostState

	return nil
}

func (v *VirtualBox) HostDestroy(host model.HostConfig, state *model.State) (bool, error) {
	return false, nil
}
