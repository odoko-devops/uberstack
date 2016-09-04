package main

import (
	"os"
	"fmt"
	"strings"
	"utils"
)

func main() {
	installerHome := os.Getenv("INSTALLER_HOME")
	command := fmt.Sprintf(`
	   docker run --rm \
	   -v %s:/state \
	   odoko/docker-stack %s`, installerHome, strings.Join(os.Args[1:], " "))
	utils.Execute(command, nil, "")
}