package defaultProvider

var sampleVariables = []string{
	"AWS_ACCESS_KEY",
	"AWS_SECRET_KEY",
	"USERNAME",
	"PASSWORD",
	"EMAIL",
	"DOMAIN_NAME",
	"MANAGEMENT_HOST_ELASTIC_IP_ALLOCATION",
}
var sampleConfiguration = `
providers:
  - type: docker-machine
    name: amazonec2
    config:
      region: us-east-1
      zone: b
      vpc_cidr: 10.128.0.0/16
      public_cidr: 10.128.0.0/24
      private_cidr: 10.128.1.0/24
      ssh_keypath: id_rsa
      access_key: %AWS_ACCESS_KEY%
      secret_key: %AWS_SECRET_KEY%
    terraform-resources:
      - aws_vpc.default
      - aws_internet_gateway.default
      - aws_subnet.public
      - aws_route_table.public
      - aws_route_table_association.public
    terraform-outputs:
      - vpc_id
      - public_subnet_id
  - type: docker-machine
    name: virtualbox
    config:
      boot2docker-image: http://www.odoko.co.uk/boot2docker.iso
      netmask: 255.255.255.0
      broadcast: 192.168.99.255
  - type: rancher
    name: amazonec2-rancher

authentication:
  - name: default
    users:
    - username: %USERNAME%
      password: %PASSWORD%
      email: %EMAIL%

hosts:
  - name: management
    provider: amazonec2
    rancher-interface: eth0
    config:
      instance_type: t2.small
      elastic_ip_allocation: %MANAGEMENT_HOST_ELASTIC_IP_ALLOCATION%
      security_group: management-tools
      subnet: "{{.terraform.public_subnet_id}}"
    rancher-agent: false
    apps:
      - type: vpn
        config:
          host: vpn.%DOMAIN_NAME%
          auth-realm: default
          cidr: 10.128.0.0/16
      - type: http-proxy
        config:
          port: 80
          jenkins-host: jenkins.%DOMAIN_NAME%
          rancher-host: rancher.%DOMAIN_NAME%
      - type: registry
        config:
          host: docker.%DOMAIN_NAME%
          ssl: true
          auth: basic
          auth-realm: default
      - type: rancher-server
        config:
          host: rancher.%DOMAIN_NAME%
          docker-host: docker.%DOMAIN_NAME%
          ssl: false
          auth-type: local
          auth-realm: default 
      - type: jenkins
        config:
          host: jenkins.%DOMAIN_NAME%
          docker-host: docker.%DOMAIN_NAME%
          ssl: false
          auth: basic
          auth-realm: default
    terraform-resources-before:
      - aws_security_group.management
    terraform-resources-after:
      - aws_subnet.private
      - aws_route_table.private
      - aws_route_table_association.private
      - aws_eip_association.management
    terraform-outputs-after:
      - private_subnet_id
  - name: docker01
    provider: amazonec2
    rancher-interface: eth0
    config:
      instance_type: t2.small
      security_group: docker-hosts
      subnet: "{{.terraform.public_subnet_id}}"
    rancher-agent: true
    labels:
      service.zookeeper: "true"
      service.solrcloud: "true"
    terraform-resources-before:
      - aws_security_group.docker
  - name: local-management
    provider: virtualbox
    rancher-interface: eth1
    config:
      disk-size: 8000
      ram: 512
      ip: 192.168.99.99
    rancher-agent: false
    apps:
      - type: http-proxy
        config:
          port: 80
          jenkins-host: jenkins
          rancher-host: rancher
      - type: rancher-server
        config:
          host: rancher
          docker-host: docker.odoko.org
          ssl: false
          auth-type: local
          auth-realm: default 

  - name: local-docker
    provider: virtualbox
    rancher-interface: eth1
    config:
      disk-size: 40000
      ram: 2048
      ip: 192.168.99.101
    rancher-agent: true
    labels:
      service.zookeeper: "true"
      service.solrcloud: "true"
`