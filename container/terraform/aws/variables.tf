variable "aws_access_key" { 
  description = "AWS access key"
}

variable "aws_secret_key" { 
  description = "AWS secret access key"
}

variable "aws_region"     { 
  description = "AWS region to host your network"
  default     = "us-east-1" 
}

variable "vpc_cidr" {
  description = "CIDR for VPC"
  default     = "10.128.0.0/16"
}

variable "public_subnet_cidr" {
  description = "CIDR for public subnet"
  default     = "10.128.0.0/24"
}

variable "private_subnet_cidr" {
  description = "CIDR for private subnet"
  default     = "10.128.1.0/24"
}
