package model

import (
	"strings"
	"fmt"
)

type SkipList int

const (
	SkipHost SkipList = 1 << iota
	SkipUpload
	SkipApps
	SkipRancherAgent
	SkipJenkins
	SkipRancherServer
	SkipDockerRegistry
	SkipVpn
	SkipProxy
)

var skipMapping = map[string]SkipList{
	"host": SkipHost,
	"upload": SkipUpload,
	"apps": SkipApps,
	"rancher-agent": SkipRancherAgent,
	"jenkins": SkipJenkins,
	"rancher-server": SkipRancherServer,
	"docker-registry": SkipDockerRegistry,
	"vpn": SkipVpn,
	"proxy": SkipProxy,
}

type SkipFunctions interface {
	Configure(string)
	Avoid(int) bool
}

func (s *SkipList) Configure(skipOptions string) *SkipList {
	if skipOptions == "" {
		return s
	}
	parts := strings.Split(skipOptions, ",")
	for _, part := range parts {
		fmt.Println(part)
		id, ok := skipMapping[part]
		if ok {
			var newList SkipList = *s | id
			s = &newList
		} else {
			fmt.Printf("Unknown skip option: %s\n", part)
		}
	}
	return s
}
func (s *SkipList) Avoid(option SkipList) bool {
	return (*s & option) > 0
}