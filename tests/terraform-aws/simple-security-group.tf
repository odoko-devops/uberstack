resource "aws_security_group" "simple" {
  name = "simple"
  description = "Simple SG for Uberstack tests"
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
}

output "simple_sg_id" {
  value = "${aws_security_group.simple.id}"
}
