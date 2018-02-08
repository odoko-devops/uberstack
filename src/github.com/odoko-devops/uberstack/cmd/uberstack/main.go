package main

import (
	"fmt"
	"flag"
	"github.com/odoko-devops/uberstack/uberstack"
)

/*
  Sample usage:

	  uberstack host up management
	  uberstack host destroy management
	  uberstack host up docker01

	  uberstack host up local-management
	  uberstack host up local-docker

	  uberstack app up myapp local
	  uberstack app up myapp dev
 */

type EnvFiles []string

func (e *EnvFiles) String() string {
	return "STUFF"
}

func (e *EnvFiles) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func main() {
	var envFiles EnvFiles

	flag.Var(&envFiles, "e", "Environment variable file to use")

	flag.Parse()

	actionType := flag.Arg(0)
	switch actionType {
	case "host":
		err := uberstack.ProcessHost(flag.Args(), envFiles)
		if err != nil {
			fmt.Println(err)
		}
	case "app":
		err := uberstack.ProcessApp(flag.Args(), envFiles)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", actionType)
	}
}