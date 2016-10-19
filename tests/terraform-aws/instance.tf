variable "ami" {}
variable "instance_type" {}

resource "aws_instance" "terraform_host01" {
  ami = "${var.ami}"
  subnet_id = "${aws_subnet.public.id}"
  instance_type = "${var.instance_type}"
  vpc_security_group_ids = ["${aws_security_group.simple.id}"]
  key_name = "${aws_key_pair.default.key_name}"
}

output "terraform_host01_id" {
  value = "${aws_instance.terraform_host01.id}"
}

output "terraform_host01_ip" {
  value = "${aws_instance.terraform_host01.public_ip}"
}
