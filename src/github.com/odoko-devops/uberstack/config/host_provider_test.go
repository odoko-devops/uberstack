package config

import (
	"testing"
)

type XHostProvider struct {
	HostProviderBase
}

func (p *XHostProvider) GetConfig(name string) string {
	return p.Config[name] + "--"
}

func (p *XHostProvider) LoadHost(filename string) (Host, error) {return nil, nil}


const (
	PARAM_NAME="a"
	PARAM_VALUE="b"
)
func TestSetConfig(t *testing.T) {
	p := HostProviderBase{}
	p.SetConfig(PARAM_NAME, PARAM_VALUE)
	value := p.GetConfig(PARAM_NAME)
	if (value != PARAM_VALUE) {
		t.Errorf("GetConfig('%s') should return %s, but was %s", PARAM_NAME, PARAM_VALUE, value)
	}
}