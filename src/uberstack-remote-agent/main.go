package main

import (
	"flag"
	"utils"
	"strings"
	"apps"
	"fmt"
)

func main() {
	rancherHostname := flag.String("rancher", "", "Rancher hostname")
	rancherAccessKey := flag.String("access_key", "", "Rancher Access Key")
	rancherSecretKey := flag.String("secret_key", "", "Rancher Secret Key")
	networkInterface := flag.String("interface", "", "Network Interface used to access Rancher Server")
	labelsPtr := flag.String("labels", "", "Comman separated set of NAME=VALUE labels for this host")

        dockerHostname := flag.String("docker", "", "Docker hostname")
	username := flag.String("username", "", "Login username")
	password := flag.String("password", "", "Login password")
	//email    := flag.String("email", "", "Login user's email address")

	cidr := flag.String("cidr", "", "CIDR for VPN")
	publicIp := flag.String("publicIp", "", "Public IP for VPN server")

	jenkinsHostname := flag.String("jenkins", "", "Jenkins hostname")

	ip        := flag.String("ip", "", "IP for current host")
	netmask   := flag.String("netmask", "", "Netmask for local network")
	broadcast := flag.String("broadcast", "", "Broadcast address for local network")

	flag.Parse()

	action := flag.Arg(0)
	switch action {
	case "rancher-agent":
		utils.Required(*rancherHostname,  "-rancher required")
		utils.Required(*rancherAccessKey, "-access_key required")
		utils.Required(*rancherSecretKey, "-secret_key required")
		utils.Required(*networkInterface, "-interface required")
		labels := strings.Replace(*labelsPtr, ",", "&", -1)
		apps.RancherInstallAgent(*rancherHostname, *rancherAccessKey, *rancherSecretKey, *networkInterface, labels)
	case "docker-registry":

	case "jenkins-server":
		utils.Required(*dockerHostname, "-docker required")
		utils.Required(*username, "-username required")
		utils.Required(*password, "-password required")
		//apps.Jenkins_RemoteInstall(*dockerHostname, *username, *password)

	case "http-proxy":
		utils.Required(*rancherHostname, "-rancher required")
		utils.Required(*jenkinsHostname, "-jenkins required")
		//apps.Proxy_RemoteInstall(*jenkinsHostname, *rancherHostname)

	case "vpn-server":
		utils.Required(*cidr, "-cidr required")
		utils.Required(*publicIp, "-publicIp required")
		utils.Required(*username, "-username required")
		utils.Required(*password, "-password required")
		apps.Vpn_RemoteInstall(*cidr, *publicIp, *username, *password)
fmt.Println("Remote vpn-server complete")
	case "virtualbox":
		utils.Required(*ip, "-ip required")
		utils.Required(*netmask, "-netmask required")
		utils.Required(*broadcast, "-broadcast required")
		apps.Virtualbox_Configure(*ip, *netmask, *broadcast)
	default:
		println("Unknown actiom: " + action)
	}
}