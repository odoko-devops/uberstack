/* Define our vpc */
resource "aws_vpc" "default" {
  cidr_block = "${var.vpc_cidr}"
  enable_dns_hostnames = true
  tags { 
    Name = "odoko-vpc" 
  }
}

output "vpc_id" {
  value = "${aws_vpc.default.id}"
}
