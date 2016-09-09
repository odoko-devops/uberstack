package virtualbox

import (
	"installer/model"
	"fmt"
	"utils"
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
	command := fmt.Sprintf(`docker-machine create %s \
             --driver virtualbox \
             --virtualbox-cpu-count -1 \
             --virtualbox-disk-size %s \
             --virtualbox-memory %s \
             --virtualbox-boot2docker-url=%s
             `, host.Name, disk, memory, image)
	utils.ExtendScript(command)

	v.setIp(host.Name)
	v.makeLocalRancherHostLinks(host.Name)
	return nil
}

func (v *VirtualBox) HostDestroy(host model.HostConfig, state *model.State) (bool, error) {
	command := fmt.Sprintf("docker-machine rm -f %s", host.Name)
	utils.ExtendScript(command)
	return true, nil
}

func (v *VirtualBox) setIp(name string) {

	ssh := "docker-machine ssh " + name
	command := fmt.Sprintf("%s \"echo '%s netmask %s broadcast %s' | sudo tee /etc/ip.cfg\"",
		ssh, v.NetMask, v.Broadcast)
	utils.ExtendScript(command)

	command = ssh + "\"echo 'sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill' | sudo tee -a /var/lib/boot2docker/bootsync.sh\""
	utils.ExtendScript(command)

	command = ssh + "\"echo 'sudo ifconfig eth1 \\$(cat /etc/ip.cfg) up' | sudo tee -a /var/lib/boot2docker/bootsync.sh\""
	utils.ExtendScript(command)

	command = ssh + "\"sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill\""
	utils.ExtendScript(command)

	command = ssh + "\"sudo ifconfig eth1 \\$(cat /etc/ip.cfg) up\""
	utils.ExtendScript(command)

	command = "docker-machine regenerate-certs -f " + name
	utils.ExtendScript(command)
}

func (v *VirtualBox) makeLocalRancherHostLinks(name string) {
	ssh := "docker-machine ssh " + name

	command := ssh + "\"sudo mkdir /mnt/sda1/var/lib/rancher\""
	utils.ExtendScript(command)

	command = ssh + "\"echo 'sudo mkdir /var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile\""
	utils.ExtendScript(command)

	command = ssh + "\"echo 'sudo mount -r /mnt/sda1/var/lib/rancher /var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile\""
	utils.ExtendScript(command)
}