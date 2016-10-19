package config

import (
	"testing"
)

type XAppProvider struct {
	AppProviderBase
}

func (p *XAppProvider) GetConfig(name string) string {
	return p.Config[name] + "--"
}

func (p *XAppProvider) LoadApp(filename string) (App, error) {return nil, nil}

func TestSetAppConfig(t *testing.T) {
}