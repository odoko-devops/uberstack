package uberstack

import (
	"testing"
	"fmt"
)

/*
 * Can I successfully load a provider, and get one of a particular type?
 */
func TestLoadProvider(t *testing.T) {
	provider, err := LoadHostProvider("terraform")
	if err != nil {
		t.Error(err)
	}

	if (provider == nil) {
		t.Fatalf("Expected a non-nil provider")
	}

	providerInterface := fmt.Sprintf("%T", provider)
	if providerInterface != "*terraform.TerraformHostProvider" {
		t.Fatalf("Expected to load terraform.yml into an *terraform.TerraformHostProvider object, got a %s",
			providerInterface)
	}
	if (provider.GetName() != "terraform-aws") {
		t.Errorf("Expected name=terraform when loading terraform.yml, got %s", provider.GetName())
	}
	if (provider.GetImpl() != "terraform") {
		t.Errorf("Expected impl=terraform when loading terraform.yml, got %s", provider.GetImpl())
	}
}

func TestLoadHost(t *testing.T) {

}