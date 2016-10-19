package terraform

import (
	"testing"
)

const (
	expectedRegion = "us-east-1"
	expectedInstanceType = "t2.small"
)
func TestLoadHostProvider(t *testing.T) {

	p, err := LoadHostProvider("terraform")
	if err != nil {
		t.Errorf("Failed with error %v", err)
	}
	p2 := p.(*TerraformHostProvider)
	if p2.Name != "terraform-aws" {
		t.Errorf("Expected terraform provider name terraform-aws, got %s", p2.Name)
	}
	if p2.Variables["region"] != expectedRegion {
		t.Errorf("Region should be '%s', got %s", expectedRegion, p2.Variables["region"])
	}
	if len(p2.Outputs)!=2 {
		t.Errorf("Expected 2 outputs, got %s", p2.Outputs)
	}
}


func TestLoadHost(t *testing.T) {
	p, err := LoadHostProvider("terraform")
	if err != nil {
		t.Errorf("Failed with error %v", err)
	}
	h, err := p.LoadHost("terraform-host01")
	if err != nil {
		t.Errorf("Failed with error %v", err)
	}
	tfHost := h.(*TerraformHost)
	if tfHost.Name != "terraform01" {
		t.Errorf("Expected host configured with aws-host01 to have the name terraform01, was %s", tfHost.Name)
	}
	if tfHost.Variables["instance_type"] != expectedInstanceType {
		t.Errorf("Expected host configured with aws-host01 to have an IAM role of %s, was %s",
			expectedInstanceType,
			tfHost.Variables["iam_role"])
	}
}