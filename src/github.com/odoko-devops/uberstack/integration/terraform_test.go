package integration

import (
	"github.com/DATA-DOG/godog"
	"log"
	"os"
	"errors"
	"strings"
	"bytes"
	u "github.com/odoko-devops/uberstack/uberstack"
	"fmt"
	"runtime/debug"
)

var response string

func uberHomeIsSet() error {
	if os.Getenv("UBER_HOME") == "" {
		return errors.New("UBER_HOME is not set")
	}
	return nil
}

func iExecute(command string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime error: %s", r)
			debug.PrintStack()
		}
	}()
	log.Printf("Executing '%s'\n", command)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	err = u.ProcessHost(strings.Split(command, " "))
	log.SetOutput(os.Stderr)
	response = buf.String()
	log.Println(response)
	if err != nil {
		return err
	}
	return nil
}

func theOutputContains(text string) error {
	if !strings.Contains(response, text) {
		return fmt.Errorf("Response does not contain '%s'", text)
	}
	return nil
}

func isHostCreated() error {
	return nil
}

func aRunningHost(host string) error {
	return nil
}

func aKnownSshKey() error {
	return nil
}

func iExecuteViaSSHOnHost(command, host string) error {
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^UBER_HOME is set$`, uberHomeIsSet)
	s.Step(`^I execute '(.+)'$`, iExecute)
	s.Step(`^(.*) has been created$`, isHostCreated)
	s.Step(`^the output contains '(.+)'$`, theOutputContains)

	s.Step(`a running host '(.+)'`, aRunningHost)
	s.Step(`a known SSH key`, aKnownSshKey)
	s.Step(`I execute '(.+)' via ssh on host '(.+)'$`, iExecuteViaSSHOnHost)
}
