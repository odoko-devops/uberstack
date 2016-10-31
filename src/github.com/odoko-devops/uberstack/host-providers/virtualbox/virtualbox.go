package virtualbox

import (
	"github.com/odoko-devops/uberstack/utils"
	"github.com/odoko-devops/uberstack/config"
	"fmt"
	"io/ioutil"
)

type VirtualBoxHostProvider struct {
	config.HostProviderBase `yaml:",inline"`

	Boot2DockerImage string `yaml:"boot2docker_image"`
	Netmask string
	Broadcast string
}

type VirtualBoxHost struct {
	config.HostBase `yaml:",inline"`

	DiskSize string `yaml:"disk-size"`
	Ram string
	Ip string
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
	host.HostProvider = p
	return host, nil
}


func (p *VirtualBoxHostProvider) DeleteHost(host config.Host) (error) {
	return nil
}

var virtualBoxInstallScript=`#!/bin/sh
sudo mkdir /mnt/sda1/var/lib/rancher
echo 'sudo mkdir /mnt/sda1/var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile
echo 'sudo mount -r /mnt/sda1/var/lib/rancher /var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile
echo '${IP} netmask ${NETMASK} broadcast ${BROADCAST}' | sudo tee /etc/ip.cfg
echo 'sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill' | sudo tee -a /var/lib/boot2docker/bootsync.sh
echo 'sudo ifconfig eth1 $$(cat /etc/ip.cfg) up' | sudo tee -a /var/lib/boot2docker/bootsync.sh
sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill
sudo ifconfig eth1 $$(cat /etc/ip.cfg) up
`

func (p *VirtualBoxHostProvider) CreateHost(h config.Host) (map[string]string, map[string]string, error) {
	host := h.(*VirtualBoxHost)

	image := p.Boot2DockerImage
	if image != "" {
		image = "--virtualbox-boot2docker-url=" + image
	}
	command := fmt.Sprintf(`docker-machine create %s
             --driver virtualbox
             --virtualbox-cpu-count -1
             --virtualbox-disk-size %s
             --virtualbox-memory %s
             %s
             `, host.HostName, host.DiskSize, host.Ram, image)
	_, err := utils.Execute(command, nil, "")
	if err != nil {
		return nil, nil, err
	}

	env := config.ExecutionEnvironment{}
	env["IP"] = host.Ip
	env["NETMASK"] = p.Netmask
	env["BROADCAST"] = p.Broadcast

	script := p.Resolve(virtualBoxInstallScript, env)
	err = ioutil.WriteFile("/tmp/uberstack-vb-install", []byte(script), 0755)
	if err != nil {
		return nil, nil, err
	}
	command = fmt.Sprintf(`docker-machine scp /tmp/uberstack-vb-install %s:/tmp/`, host.HostName)
	_, err = utils.Execute(command, nil, "")
	if err != nil {
		return nil, nil, err
	}
	command = fmt.Sprintf(`docker-machine ssh %s /tmp/uberstack-vb-install`, host.HostName)
	_, err = utils.Execute(command, nil, "")
	if err != nil {
		return nil, nil, err
	}
	command = fmt.Sprintf("docker-machine regenerate-certs -f %s", host.HostName)
	_, err = utils.Execute(command, nil, "")

	return nil, nil, err
}
