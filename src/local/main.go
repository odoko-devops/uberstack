package main

import (
	"os"
	"fmt"
	"strings"
	"utils"
)

func main() {
	uberHome := os.Getenv("UBER_HOME")
	uberState := os.Getenv("UBER_STATE")

	uberHomeEnv := ""
	if uberHome != "" {
		uberHomeEnv = fmt.Sprintf("-e UBER_HOME=/uberhome -v %s:/uberhome", uberHome)
	}

	if uberState == "" {
		uberState = uberHome + "/state"
	}
	uberStateEnv := fmt.Sprintf("-e UBER_STATE=/state -v %s:/state", uberState)

	args := strings.Join(os.Args[1:], " ")
	command := fmt.Sprintf("docker run --rm %s %s odoko/docker-stack %s", uberHomeEnv, uberStateEnv, args)
	utils.Execute(command, nil, "")

	if utils.DoesUberScriptExist() {
		utils.ExecuteUberScript()
		utils.RemoveUberScript()
	}
}