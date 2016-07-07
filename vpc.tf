/* Define our vpc */
resource "aws_vpc" "default" {
  cidr_block = "${var.vpc_cidr}"
  enable_dns_hostnames = true
  tags { 
    Name = "odoko-vpc" 
  }
}

resource "aws_key_pair" "sshkey" {
  key_name = "sshkey"
  public_key = "${var.public_ssh_key}"
}

