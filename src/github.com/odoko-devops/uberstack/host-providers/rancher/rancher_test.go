package rancher

import (
	"testing"
)

const (
	expectedRegion = "us-east-1"
	expectedInstanceType = "t2.small"
)
func TestLoadHostProvider(t *testing.T) {

	p, err := LoadHostProvider("rancher")
	if err != nil {
		t.Errorf("Failed with error %v", err)
	}
	p2 := p.(*RancherHostProvider)
	if p2.Name != "rancher-aws" {
		t.Errorf("Expected rancher provider name rancher-aws, got %s", p2.Name)
	}
}


func TestLoadHost(t *testing.T) {
	p, err := LoadHostProvider("rancher")
	if err != nil {
		t.Errorf("Failed with error %v", err)
	}
	h, err := p.LoadHost("rancher-host01")
	if err != nil {
		t.Errorf("Failed with error %v", err)
	}
	tfHost := h.(*RancherHost)
	if tfHost.Name != "rancher01" {
		t.Errorf("Expected host configured with rancher-host01 to have the name rancher01, was %s", tfHost.Name)
	}
}