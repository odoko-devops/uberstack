package apps

import (
	"model"
	"log"
	"utils"
	"fmt"
	"strings"
	"text/template"
	"bytes"
	"io/ioutil"
)

var vpn_install_commands = []string{
	"docker run --name ovpn-data -v /etc/openvpn:/etc/openvpn busybox",
	"docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_genconfig -p {{.cidr}} -u udp://{{.publicIp}}",
	"{{.password}}\n{{.password}}\n\n{{.password}}%docker run --volumes-from ovpn-data --rm -it gosuri/openvpn ovpn_initpki",
	"docker run --volumes-from ovpn-data -d -p 1194:1194/udp --cap-add=NET_ADMIN gosuri/openvpn",
	"{{.password}}%docker run --volumes-from ovpn-data --rm -it gosuri/openvpn easyrsa build-client-full {{.username}} nopass",
	}


func Vpn_Install(config model.Config, state *model.State, hostConfig model.HostConfig, app model.AppConfig) error {

	log.Println("Deploy VPN Service")

	env := model.GetDockerEnvironment(state, hostConfig)

	authRealm := model.GetAuthRealm(config, app.Config["auth-realm"])
	hostState := state.HostState[hostConfig.Name]
	params := map[string]string{
		"cidr": app.Config["cidr"],
		"publicIp": hostState["public-ip"],
		"username": authRealm.Users[0].Username,
		"password": authRealm.Users[0].Password,
		"uberstate": utils.GetUberState(),
	}

	utils.ExecuteRemote(hostConfig.Name, "sudo mkdir -p /etc/openvpn", nil, "")

	for _, command := range vpn_install_commands {
		println("-------")
		commandTemplate, err := template.New("vpnCommand").Parse(command)
		utils.Check(err)

		var buf bytes.Buffer
		err = commandTemplate.Execute(&buf, params)
		utils.Check(err)

		command2 := buf.String()
		if strings.Contains(command2, "%") {
			parts := strings.Split(command2, "%")
			input := parts[0]
			cmd := parts[1]
			fmt.Printf("%s ---> %s\n", input, cmd)
			utils.ExecuteWithInput(cmd, input, env, "")
		} else {
			fmt.Printf("------- %s\n", command2)
			utils.Execute(command2, env, "")
		}
	}

	command := "docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_getclient " + params["username"]
	ovpn := utils.ExecuteAndRetrieve(command, env, "")
	filename := fmt.Sprintf("%s/%s", utils.GetUberState(), params["username"])
	err := ioutil.WriteFile(filename, []byte(ovpn), 0644)
	utils.Check(err)
	return nil
}
