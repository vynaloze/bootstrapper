# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

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
  default     = []

  type = list(string)
}

variable "public_subnets_cidrs" {
  description = "CIDR blocks of public subnets to create inside the VPC"
  default     = []

  type = list(string)
}

# VPN

variable "client_cidr_block" {
  description = "The IPv4 address range from which to assign client IP addresses"

  type = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "region_short" {
  description = "Optional short region name override"
  default     = ""

  type = string
}

variable "vpn_mtls_key_algorithm" {
  description = "Algorithm used for the keys for VPN mutual TLS authentication"
  default     = "RSA"

  type = string
}

variable "vpn_mtls_key_size" {
  description = "Size of the keys for VPN mutual TLS authentication"
  default     = 2048

  type = number
}

variable "vpn_mtls_certificate_validity" {
  description = "The number of hours after issuing that the VPN certificates will become invalid"
  default     = 24 * 365 * 10

  type = number
}
