package main

import (
	"flag"
	"utils"
	"strings"
	"remote/apps"
)

func main() {
	rancherHostname := flag.String("rancher", "", "Rancher hostname")
	rancherAccessKey := flag.String("access_key", "", "Rancher Access Key")
	rancherSecretKey := flag.String("secret_key", "", "Rancher Secret Key")
	networkInterface := flag.String("interface", "", "Network Interface used to access Rancher Server")
	labelsPtr := flag.String("labels", "", "Comman separated set of NAME=VALUE labels for this host")
	flag.Parse()

	action := flag.Arg(0)
	switch action {
	case "rancher-agent":
		utils.Required(*rancherHostname, "-rancher required")
		utils.Required(*rancherAccessKey, "-access_key required")
		utils.Required(*rancherSecretKey, "-secret_key required")
		utils.Required(*networkInterface, "-interface required")
		labels := strings.Replace(*labelsPtr, ",", "&", -1)
		apps.RancherInstallAgent(*rancherHostname, *rancherAccessKey, *rancherSecretKey, *networkInterface, labels)

	case "docker-registry":
	case "jenkins-server":
	case "http-proxy":
	case "rancher-server":
	case "vpn-server":
	}
}