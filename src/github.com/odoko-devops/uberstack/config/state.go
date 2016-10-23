package config

import (
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type HostState map[string]*string

type State struct {
	Values map[string]*string
	Hosts map[string]HostState
	VariableMap map[string]*string
}

func (s *State) SetValue(name, value string) {
	if s.Values == nil {
		s.Values = map[string]*string{}
	}
	if s.VariableMap == nil {
		s.VariableMap = map[string]*string{}
	}
	s.Values[name] = &value
	s.VariableMap[name] = &value
}

func (s *State) SetHostValue(host Host, name, value string) {
	if s.Hosts == nil {
		s.Hosts = map[string]HostState{}
	}
	if s.Hosts[host.GetName()] == nil {
		s.Hosts[host.GetName()] = HostState{}
	}
	if s.VariableMap == nil {
		s.VariableMap = map[string]*string{}
	}
	s.Hosts[host.GetName()][name] = &value
	variableName := fmt.Sprintf("host.%s.%s", host.GetName(), name)
	log.Printf("Setting %s to %s at %s", name, value, variableName)
	s.VariableMap[ variableName ] = &value
}

func (s *State) GetHostValue(host Host, name string) string {
	val, ok := s.Hosts[host.GetName()][name]
	if ok {
		return *val
	} else {
		return ""
	}
}

func (s *State) Get(name string) string {
	val, ok := s.VariableMap[name]
	if ok {
		return *val
	} else {
		return ""
	}
}

func (s *State) Resolve(text string, env ExecutionEnvironment) string {
	for ; strings.Contains(text, "${"); {
		text = os.Expand(text, func(name string) string {
			if envVal, ok := env[name]; ok {
				return envVal
			}
			if varVal, ok := s.VariableMap[name]; ok {
				return *varVal
			}
			return os.Getenv(name)
		})
	}
	return text
}

func (s *State) GetStateFilename() string {
	uberState := os.Getenv("UBER_STATE")
	if uberState == "" {
		uberState = os.Getenv("UBER_HOME")
	}
	return fmt.Sprintf("%s/state.yml", uberState)
}

func (s *State) Load() error {
	stateFile := s.GetStateFilename()
	if _, err := os.Stat(stateFile) ;os.IsNotExist(err) {
		log.Println("Not loading state, can't find file at ", stateFile)
		return nil // it isn't an error if a state file does not exist.
	}
	log.Println("Loading state from", stateFile)
	bytes, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, s)
	return err
}

func (s *State) Save() error {
	stateFile := s.GetStateFilename()
	log.Printf("Saving state as %s from %s", stateFile, *s)
	bytes, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(stateFile, bytes, 0644)
	return err
}