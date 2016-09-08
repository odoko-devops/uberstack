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

        dockerHostname := flag.String("docker", "", "Docker hostname")
	username := flag.String("username", "", "Login username")
	password := flag.String("password", "", "Login password")
	//email    := flag.String("email", "", "Login user's email address")

	cidr := flag.String("cidr", "", "Network CIDR")
	publicIp := flag.String("publicip", "", "Host's Public IP")
	hostName := flag.String("host", "", "Hostname for current host")

	jenkinsHostname := flag.String("jenkins", "", "Jenkins hostname")

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

	case "rancher-server":

	case "vpn-server":
		utils.Required(*cidr, "-cidr required")
		utils.Required(*publicIp, "-publicip required")
		utils.Required(*username, "-username required")
		utils.Required(*hostName, "-host required")
		apps.Vpn_RemoteInstall(*cidr, *publicIp, *username, *hostName)

	default:
		print("Unknown actiom: " + action)
	}
}