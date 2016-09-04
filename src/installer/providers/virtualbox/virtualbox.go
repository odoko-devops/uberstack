package virtualbox

import "installer/model"

type VirtualBox struct {

}

func (v *VirtualBox) Configure(config model.Config, state *model.State, provider model.ProviderConfig) error {
	return nil
}
func (v *VirtualBox) InfrastructureUp() error {
	return nil
}
func (v *VirtualBox) InfrastructureDestroy() error {
	return nil
}
func (v *VirtualBox) HostUp(host model.HostConfig, state *model.State) error {
	return nil
}

func (v *VirtualBox) HostDestroy(host model.HostConfig, state *model.State) (bool, error) {
	return false, nil
}

/*
def create_local_host(name, disk, memory, image):
  return '''docker-machine create %s \
             --driver virtualbox \
             --virtualbox-cpu-count -1 \
             --virtualbox-disk-size %s \
             --virtualbox-memory %s \
             --virtualbox-boot2docker-url=%s
             ''' % (name, disk, memory, image)

def create_local_rancher_host(config):
  local = config["local"]
  rancher = local["rancher"]
  script = ["#!/bin/sh"]
  script.append(create_local_host("rancher", rancher["disk-size"], rancher["ram"], local["boot2docker-image"]))
  script.extend(set_ip("rancher", local, rancher))
  script.append('docker-machine ssh rancher "docker run -d --restart=always -p 80:8080 rancher/server"')

  write_script("/state/run", script)
  ask("state/run")

def set_ip(name, local, host):
  return ["docker-machine ssh %s \"echo '%s netmask %s broadcast %s' | sudo tee /etc/ip.cfg\"" %
             (name, host["ip"], local["netmask"], local["broadcast"]),
          "docker-machine ssh %s \"echo 'sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill' | sudo tee -a /var/lib/boot2docker/bootsync.sh\"" % name,
          "docker-machine ssh %s \"echo 'sudo ifconfig eth1 \$(cat /etc/ip.cfg) up' | sudo tee -a /var/lib/boot2docker/bootsync.sh\"" % name,
          "docker-machine ssh %s \"sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill\"" % name,
          "docker-machine ssh %s \"sudo ifconfig eth1 \$(cat /etc/ip.cfg) up\"" % name,
          "docker-machine regenerate-certs -f %s" % name
         ]

def make_local_rancher_host_links(host):
  return ["docker-machine ssh %s \"sudo mkdir /mnt/sda1/var/lib/rancher\"" % host,
          "docker-machine ssh %s \"echo 'sudo mkdir /var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile\"" % host,
          "docker-machine ssh %s \"echo 'sudo mount -r /mnt/sda1/var/lib/rancher /var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile\"" % host
         ]

def create_local_docker_host(config):
  local = config["local"]
  docker = local["docker-host"]
  create_local_host("docker", docker["disk-size"], rancher["ram"], local["boot2docker-image"])
  set_ip("docker", local, docker)
  make_local_rancher_host_links("docker")


*/
