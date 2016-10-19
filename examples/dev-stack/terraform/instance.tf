variable "ami" {}
variable "instance_type" {}
variable "elastic_ip_allocation" {}

resource "aws_instance" "dev_host" {
  ami = "${var.ami}"
  subnet_id = "${aws_subnet.public.id}"
  instance_type = "${var.instance_type}"
  vpc_security_group_ids = ["${aws_security_group.dev_stack.id}"]
  key_name = "${aws_key_pair.default.key_name}"
  tags {
    Name = "dev_host"
  }
}

resource "aws_eip_association" "dev_host" {
  instance_id = "${aws_instance.dev_host.id}"
  allocation_id = "${var.elastic_ip_allocation}"
}

output "dev_host_id" {
  value = "${aws_instance.dev_host.id}"
}

output "dev_host_ip" {
  value = "${aws_eip_association.dev_host.public_ip}"
}
