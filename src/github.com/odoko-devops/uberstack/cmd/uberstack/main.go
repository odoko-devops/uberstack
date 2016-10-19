package main

import (
	"fmt"
	"flag"
	u "github.com/odoko-devops/uberstack/uberstack"
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

func main() {
	flag.Parse()
	actionType := flag.Arg(0)
	switch actionType {
	case "host":
		err := u.ProcessHost(flag.Args())
		if err != nil {
			fmt.Println(err)
		}
	case "app":
		err := u.ProcessApp(flag.Args())
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", actionType)
	}
}