package apps

import (
	"github.com/odoko-devops/uberstack/utils"
	"text/template"
	"io/ioutil"
	"bytes"
)

var virtualBoxInstallScript=`#!/bin/sh
sudo mkdir /mnt/sda1/var/lib/rancher
echo 'sudo mkdir /mnt/sda1/var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile
echo 'sudo mount -r /mnt/sda1/var/lib/rancher /var/lib/rancher' | sudo tee -a /var/lib/boot2docker/profile
echo '{{.ip}} netmask {{.netmask}} broadcast {{.broadcast}}' | sudo tee /etc/ip.cfg
echo 'sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill' | sudo tee -a /var/lib/boot2docker/bootsync.sh
echo 'sudo ifconfig eth1 $(cat /etc/ip.cfg) up' | sudo tee -a /var/lib/boot2docker/bootsync.sh
sudo cat /var/run/udhcpc.eth1.pid | xargs sudo kill
sudo ifconfig eth1 $(cat /etc/ip.cfg) up
`
func Virtualbox_Configure(ip, netmask, broadcast string) {

	virtualboxTemplate, err := template.New("terraformEip").Parse(virtualBoxInstallScript)
	utils.Check(err)

	var buf bytes.Buffer


	params := map[string]string{
		"ip": ip,
		"netmask": netmask,
		"broadcast": broadcast,
	}
	err = virtualboxTemplate.Execute(&buf, params)
	utils.Check(err)
	err = ioutil.WriteFile("/tmp/virtualbox.sh", buf.Bytes(), 0755)
	utils.Check(err)
	utils.Execute("/tmp/virtualbox.sh", nil, "")
}



