package apps

/*
func Vpn_Install(config model.Config, state *model.State, hostConfig model.HostConfig, app model.AppConfig) error {

	log.Println("Deploy VPN Service")

	env := model.GetDockerEnvironment(state, hostConfig)

	authRealm := model.GetAuthRealm(config, app.Config["auth-realm"])
	hostState := state.HostState[hostConfig.Name]
	username := authRealm.Users[0].Username
	command := fmt.Sprintf("./uberstack-remote-agent -cidr=%s -publicIp=%s -username=%s -password=%s vpn-server",
		app.Config["cidr"],
		hostState["public-ip"],
		authRealm.Users[0].Username,
		authRealm.Users[0].Password)

	utils.ExecuteRemote(hostConfig.Name, command, nil, "")

	command = fmt.Sprintf("docker-machine -s %s/machine ssh %s cat %s.ovpn",
		utils.GetUberState(), hostConfig.Name, username)

	ovpn, err := utils.ExecuteAndRetrieve(command, env, "")
	utils.Check(err)
	filename := fmt.Sprintf("%s/%s.ovpn", utils.GetUberState(), username)
	err = ioutil.WriteFile(filename, []byte(ovpn), 0644)
	utils.Check(err)
	fmt.Printf("Changed ovpn\n")

	return nil
}

var vpn_installScript = `#!/bin/sh
sudo mkdir -p /etc/openvpn
sudo iptables -t nat -A POSTROUTING -j MASQUERADE
echo 1 | sudo tee /proc/sys/net/ipv4/conf/all/forwarding > /dev/null
`

var vpn_installCommands = []string{
	"/tmp/uberstack-tmp",
	"docker run --name ovpn-data -v /etc/openvpn:/etc/openvpn busybox",
	"docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_genconfig -p {{.cidr}} -u udp://{{.publicIp}}",
	"{{.password}}\n{{.password}}\n\n{{.password}}%docker run --volumes-from ovpn-data --rm -it gosuri/openvpn ovpn_initpki",
	"docker run --volumes-from ovpn-data -d -p 1194:1194/udp --cap-add=NET_ADMIN gosuri/openvpn",
	"{{.password}}%docker run --volumes-from ovpn-data --rm -it gosuri/openvpn easyrsa build-client-full {{.username}} nopass",
	}

func Vpn_RemoteInstall(cidr, publicIp, username, password string) error {

	params := map[string]string{
		"cidr": cidr,
		"publicIp": publicIp,
		"username": username,
		"password": password,
	}

	ioutil.WriteFile("/tmp/uberstack-tmp", []byte(vpn_installScript), 0755)

	for _, command := range vpn_installCommands {
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
			utils.ExecuteWithInput(cmd, input, nil, "")
		} else {
			fmt.Printf("------- %s\n", command2)
			utils.Execute(command2, nil, "")
		}
	}

	command := "docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_getclient " + params["username"]
	ovpn, err := utils.ExecuteAndRetrieve(command, nil, "")
	utils.Check(err)
	filename := fmt.Sprintf("%s.ovpn", params["username"])
	err = ioutil.WriteFile(filename, []byte(ovpn), 0644)
	utils.Check(err)
	fmt.Println("Completed remote VPN install")
	return nil
}
*/