package main

import "utils"

/*
    This is a dummy package that creates a simple command line app which has access
    to all of the classes behind Uberstack. It is intended for experimentation
    when exploring new functionality within Uberstack.
 */
func main() {
	command := "./foo.sh"
	input := "hello"
	utils.ExecuteWithInput(command, input, nil, "")
}