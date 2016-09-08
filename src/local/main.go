package main

import (
	"os"
	"fmt"
	"strings"
	"utils"
)

func main() {
	uberState := os.Getenv("UBER_STATE")

	if uberState == "" {

		uberHome := os.Getenv("UBER_HOME")
		if uberHome == "" {
			println("Please set either UBER_HOME or UBER_STATE")
			os.Exit(1)
		}
		uberState = fmt.Sprintf("%s/state", uberHome)
	}
	command := fmt.Sprintf(`
	   docker run --rm \
	   -v %s:/state \
	   odoko/docker-stack %s`, uberState, strings.Join(os.Args[1:], " "))
	utils.Execute(command, nil, "")
}