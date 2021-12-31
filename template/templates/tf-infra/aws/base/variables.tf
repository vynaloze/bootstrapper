variable "region" {
  description = "AWS region"

  type = string
}

variable "environment" {
  description = "Name of the environment"

  type = string
}

variable "base_domain" {
  description = "Base domain, e.g. example.com"

  type = string
}

# VPC

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

# VPN

variable "client_cidr_block" {
  description = "The IPv4 address range from which to assign client IP addresses"

  type = string
}

