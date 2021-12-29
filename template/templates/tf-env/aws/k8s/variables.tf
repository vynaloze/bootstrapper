# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "environment" {
  description = "Name of the environment"

  type = string
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC"

  type = string
}

variable "base_vpc_id" {
  description = "Base VPC ID (to create peering connection with)"

  type = string
}

variable "base_vpc_cidr" {
  description = "The CIDR block for the base VPC (to create peering connection with)"

  type = string
}

variable "base_vpc_private_route_table_ids" {
  description = "IDs of the route tables associated with private subnets in base VPC (to create peering connection with)"

  type = set(string)
}

variable "client_vpn_endpoint_id" {
  description = "ID of the Client VPN endpoint"

  type = string
}

variable "client_vpn_subnet_id" {
  description = "ID of the subnet attached to the Client VPN"

  type = string
}

variable "route53_hosted_zone_id" {
  description = "Route53 Hosted Zone ID"

  type = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

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

variable "cluster_log_retention_days" {
  description = "Number of days until cluster logs expire"

  default = 7
  type    = number
}

variable "external_dns_service_account_namespace" {
  description = "ExternalDNS service account namespace"

  default = "default"
  type    = string
}

variable "external_dns_service_account_name" {
  description = "ExternalDNS service account name"

  default = "external-dns"
  type    = string
}

variable "worker_general_purpose_ami_type" {
  description = "General purpose worker AMI type"

  default = "AL2_x86_64"
  type    = string
}

variable "worker_general_purpose_node_type" {
  description = "General purpose worker node type"

  default = "t3.medium"
  type    = string
}

variable "worker_general_purpose_disk_size" {
  description = "General purpose  worker disk size in GB"

  default = 20
  type    = number
}

variable "worker_general_purpose_group_size" {
  description = "Initial value of general purpose worker groups"

  default = 2
  type    = number
}

variable "worker_general_purpose_max_group_size" {
  description = "Max size of general purpose worker groups"

  default = 3
  type    = number
}
