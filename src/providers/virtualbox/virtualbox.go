package virtualbox

import (
	"model"
	"fmt"
	"utils"
	"providers/defaultProvider"
)

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
	command = fmt.Sprintf("./remote -ip=%s -broadcast=%s -netmask=%s virtualbox",
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
