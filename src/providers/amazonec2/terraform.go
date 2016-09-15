package amazonec2

var terraformVPCConfig =`
/* Amazon configuration */
provider "aws" {
  access_key = "{{.access_key}}"
  secret_key = "{{.secret_key}}"
  region     = "{{.region}}"
}

/* VPC configuration */
resource "aws_vpc" "default" {
  cidr_block = "{{.vpc_cidr}}"
  enable_dns_hostnames = true
  tags {
    Name = "Uberstack VPC"
  }
}

output "vpc_id" {
  value = "${aws_vpc.default.id}"
}
`

var terraformInfraConfig = `
/* Internet gateway for the public subnet */
resource "aws_internet_gateway" "default" {
  vpc_id = "${aws_vpc.default.id}"
}

/* Public subnet */
resource "aws_subnet" "public" {
  vpc_id            = "${aws_vpc.default.id}"
  cidr_block        = "{{.public_cidr}}"
  availability_zone = "{{.region}}{{.zone}}"
  map_public_ip_on_launch = true
  depends_on = ["aws_internet_gateway.default"]
  tags {
    Name = "public"
  }
}

output "subnet_id" {
  value = "${aws_subnet.public.id}"
}

/* Routing table for public subnet */
resource "aws_route_table" "public" {
  vpc_id = "${aws_vpc.default.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.default.id}"
  }
}

/* Associate the routing table to public subnet */
resource "aws_route_table_association" "public" {
  subnet_id = "${aws_subnet.public.id}"
  route_table_id = "${aws_route_table.public.id}"
}

/* Private subnet */
resource "aws_subnet" "private" {
  vpc_id            = "${aws_vpc.default.id}"
  cidr_block        = "{{.private_cidr}}"
  availability_zone = "{{.region}}{{.zone}}"
  map_public_ip_on_launch = false
  tags {
    Name = "private"
  }
}
`
var terraformManagementHost =`
variable "instance_id" {
  default = "NONE"
}
variable "elastic_ip_allocation" {
  default = "NONE"
}

/* Routing table for private subnet */
resource "aws_route_table" "private" {
  vpc_id = "${aws_vpc.default.id}"
  route {
    cidr_block = "0.0.0.0/0"
    instance_id = "${var.instance_id}"
  }
}

/* Associate the routing table to private subnet */
resource "aws_route_table_association" "private" {
  subnet_id = "${aws_subnet.private.id}"
  route_table_id = "${aws_route_table.private.id}"
}
`

var terraformConfig = map[string]string{
	"vpc.tf": terraformVPCConfig,
	"infrastructure.tf": terraformInfraConfig,
	"management-host.tf": terraformManagementHost,
        "management-sg.tf": terraformManagementSecurityGroup,
	"docker-sg.tf": terraformDockerSecurityGroup,
}
