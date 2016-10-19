package docker_machine

/*
Stuff kept here in case it is needed in producing a docker-machine host-provider...

func (p DefaultProvider) RegenerateCerts(host model.HostConfig) {
	log.Printf("Regenerating certificates for %s\n", host.Name)
	command := fmt.Sprintf("docker-machine -s %s/machine regenerate-certs -f %s", utils.GetUberState(), host.Name)
	utils.Execute(command, nil, "")
}


func (p DefaultProvider) StartRancherAgent(config model.Config, state *model.State, provider model.ProviderConfig, host model.HostConfig) {

	labels := make([]string, len(host.Labels))
	i:=0
	for k, v := range host.Labels {
		labels[i] = k + "=" + v
		i++
	}

	providerState := state.Provider[provider.Name]
	command := fmt.Sprintf(`./uberstack-remote-agent \
	                   -interface=%s \
	                   -rancher=%s \
	                   -access_key=%s \
	                   -secret_key=%s \
	                   -labels=%s \
	                   rancher-agent`,
		host.RancherInterface,
		providerState.RancherUrl,
		providerState.AccessKey,
		providerState.SecretKey,
		strings.Join(labels, ","))
	utils.ExecuteRemote(host.Name, command, nil, "")
}

func (p DefaultProvider) ProcessUberstack(config model.Config, state *model.State, uberHome string,
		uberstack model.Uberstack, env string, cmd string, exclude_stack string, doTerraform bool) {

	for i := 0; i < len(uberstack.Uberstacks); i++ {
		name := uberstack.Uberstacks[i]
		inner_uberstack := p.GetUberstack(uberHome, name)
		p.ProcessUberstack(config, state, uberHome, inner_uberstack, env, cmd, exclude_stack, doTerraform)
	}

	uberEnv := uberstack.Environments[env]
	params := utils.Environment{}
	for k,v := range uberEnv.TerraformConfig {
		params[k] = v
	}

	if doTerraform && len(uberEnv.TerraformBefore)>0 {
		utils.TerraformApply(uberEnv.Provider, uberEnv.TerraformBefore, params)
	}

	for i := range uberstack.Stacks {
		name := uberstack.Stacks[i]
		if name == exclude_stack {
			continue
		}
		project := name
		stack := name

		s := strings.SplitN(name, ":", 2)
		if len(s) == 2 {
			project = s[0]
			stack = s[1]
		}
		command := fmt.Sprintf(`rancher-compose --file %s/stacks/%s/docker-compose.yml \
                        --rancher-file %s/stacks/%s/rancher-compose.yml \
                        --project-name %s \
                        %s`,
			uberHome, stack, uberHome, stack, project, cmd)
		env := getParametersFor(uberstack, env, state)
		utils.Execute(command, env, "")
	}
	if doTerraform && len(uberEnv.TerraformAfter)>0 {
		utils.TerraformApply(uberEnv.Provider, uberEnv.TerraformAfter, params)
	}
}

*/