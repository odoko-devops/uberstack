package amazonec2

var terraformManagementSecurityGroup =`

resource "aws_eip_association" "management" {
  instance_id = "${var.instance_id}"
  allocation_id = "${var.elastic_ip_allocation}"
}

output "management_public_ip" {
  value = "${aws_eip_association.management.public_ip}"
}

resource "aws_security_group" "management" {
  name = "management-tools"
  description = "Security group for our Development Tools"
  vpc_id = "${aws_vpc.default.id}"

  ingress {
    from_port   = "22"
    to_port     = "22"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = "80"
    to_port     = "80"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = "443"
    to_port     = "443"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = "2376"
    to_port     = "2376"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = "1194"
    to_port     = "1194"
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = "22"
    to_port     = "22"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = "443"
    to_port     = "443"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  egress {
    from_port   = "80"
    to_port     = "80"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags { 
    Name = "management-tools" 
  }
}
`

