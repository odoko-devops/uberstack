package apps

import (
	"fmt"
	"utils"
)

func Virtualbox_Configure(ip, netmask, broadcast string) {

	commands := []string{
		"sudo mkdir /mnt/sda1/var/lib/rancher",
		"echo 'sudo mkdir /mnt/sda1/var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile",
		"echo 'sudo mount -r /mnt/sda1/var/lib/rancher /var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile",
		fmt.Sprintf("echo '%s netmask %s broadcast %s' | sudo tee /etc/ip.cfg", ip, netmask, broadcast),
		"echo 'sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill' | sudo tee -a /var/lib/boot2docker/bootsync.sh",
		"echo 'sudo ifconfig eth1 $(cat /etc/ip.cfg) up' | sudo tee -a /var/lib/boot2docker/bootsync.sh",
		"sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill",
		"sudo ifconfig eth1 $(cat /etc/ip.cfg) up",
	}

	for _, command := range commands {
		fmt.Printf("Execute: %s\n", command)
		utils.Execute(command, nil, "")
	}
}



