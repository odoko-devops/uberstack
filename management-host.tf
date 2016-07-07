/* App servers */
resource "aws_instance" "management" {
  ami = "${lookup(var.amis, var.aws_region)}"
  instance_type = "t2.small"
  subnet_id = "${aws_subnet.public.id}"
  security_groups = ["${aws_security_group.management.id}"]
  key_name = "${aws_key_pair.sshkey.key_name}"
  source_dest_check = false
  associate_public_ip_address = true
  tags = { 
    Name = "management"
  }
  connection {
    user = "ubuntu"
    key_file = "~/.ssh/id_rsa"
  }
  provisioner "remote-exec" {
     script = "install-docker.sh"
  }

  provisioner "file" {
    source = "docker-compose.yml"
    destination = "/tmp/docker-compose.yml"
  }

  provisioner "file" {
    source = "install-management-tools.sh"
    destination = "/tmp/install-management-tools.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "bash /tmp/install-management-tools.sh ${var.docker_domain} ${var.email} ${var.username} ${var.password} ${var.rancher_domain} ${var.jenkins_domain}"
    ]
  }

}

resource "aws_eip" "management" {
  vpc = true
}

resource "aws_eip_association" "management" {
  instance_id = "${aws_instance.management.id}"
  allocation_id = "${aws_eip.management.id}"
}

output "management.ip" {
  value = "${aws_eip.management.public_ip}"
}

resource "aws_security_group" "management" {
  name = "docker-registry"
  description = "Security group for our Private Docker Registry"
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
    Name = "docker-registry" 
  }
}

