variable "region" {
  description = "AWS region"

  type = string
}

variable "environment" {
  description = "Name of the environment"

  type = string
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC"

  type = string
}

variable "private_subnets_cidrs" {
  description = "CIDR blocks of private subnets to create inside the VPC"

  type = list(string)
}

variable "public_subnets_cidrs" {
  description = "CIDR blocks of public subnets to create inside the VPC"

  type = list(string)
}
