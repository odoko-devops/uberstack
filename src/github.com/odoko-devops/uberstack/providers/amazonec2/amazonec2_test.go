package amazonec2

import (
	"testing"
)

func TestLoadHostProvider(t *testing.T) {
	p, err := LoadHostProvider("/go/tests/amazonec2")
	if err != nil {
		t.Errorf("Failed with error %v", err)
	}
	p2 := p.(*Amazonec2HostProvider)
	if p2.Region != "us-east-1" {
		t.Errorf("Region should be 'us-east-1'")
	}
	if p2.Zone != "b" {
		t.Errorf("Zone should be b")
	}
/*
	thing: provider
me: amazonec2
impl: amazonec2
region: us-east-1
zone: b
vpc_cidr: 10.128.0.0/16
public_cidr: 10.128.0.0/24
private_cidr: 10.128.1.0/24
ssh_keypath: id_rsa
access_key: AKIAJ2MGU55FFVEILEXA
secret_key: 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
terraform-resources:
  - aws_vpc.default
  - aws_internet_gateway.default
  - aws_subnet.public
  - aws_route_table.public
  - aws_route_table_association.public
terraform-outputs:
  - vpc_id
  - public_subnet_id
*/
}


func TestLoadHost(t *testing.T) {
	/*(p *Amazonec2HostProvider)
} LoadHost(filename string) (config.Host, error) {
	host := new(Amazonec2Host)
	err := utils.ReadYamlFile(filename, host)
	if (err != nil) {
		return nil, err
	}
	return host, nil */
}
