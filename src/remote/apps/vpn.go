package apps

import (
	"installer/model"
	"log"
	"utils"
	"fmt"
	"strings"
)

var vpn_install_commands = []string{
	"sudo mkdir -p /etc/openvpn",
	"docker run --name ovpn-data -v /etc/openvpn:/etc/openvpn busybox",
	"docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_genconfig -p ${CIDR} -u udp://${PUBLIC_IP}",
	"docker run --volumes-from ovpn-data --rm -it gosuri/openvpn ovpn_initpki",
	"\n%docker run --volumes-from ovpn-data -d -p 1194:1194/udp --cap-add=NET_ADMIN gosuri/openvpn",
	"docker run --volumes-from ovpn-data --rm -it gosuri/openvpn easyrsa build-client-full ${USERNAME} nopass",
	"docker run --volumes-from ovpn-data --rm     gosuri/openvpn ovpn_getclient ${USERNAME} > ${USERNAME}.ovpn"}


func Vpn_Install(config model.Config, state *model.State, hostConfig model.HostConfig, app model.AppConfig) error {
	authRealm := model.GetAuthRealm(config, app.Config["auth-realm"])
	username := authRealm.Users[0].Username

	hostState := state.HostState[hostConfig.Name]
	publicIp := hostState["public-ip"]
	cidr := app.Config["cidr"]

	command := fmt.Sprintf("./remote -cidr=%s -publicip=%s -username=%s -host=%s vpn-server",
		cidr,
		publicIp,
		username,
		hostConfig.Name)
	utils.ExecuteRemote(hostConfig.Name, command, nil, "")
	return nil
}

func Vpn_RemoteInstall(cidr, publicIp, username, hostName string) error {
	log.Println("Deploy VPN Service")
	env := utils.Environment{
		"CIDR": cidr,
		"PUBLIC_IP": publicIp,
		"USERNAME": username,
	}
	for _, command := range vpn_install_commands {
		if strings.Contains(command, "%") {
			parts := strings.Split(command, "%")
			input := parts[0]
			cmd := parts[1]
			utils.ExecuteWithInput(cmd, input, env, "")
		} else {
			utils.Execute(command, env, "")
		}
	}
	return nil
}
