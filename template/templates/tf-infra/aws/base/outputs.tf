output "vpc_id" {
  description = "VPC ID"

  value = module.env.vpc_id
}

output "vpc_cidr" {
  description = "The CIDR block for the VPC"

  value = module.env.vpc_cidr
}

output "vpc_private_route_table_ids" {
  description = "IDs of the route tables associated with private subnets in VPC"

  value = module.env.vpc_private_route_table_ids
}

output "client_vpn_endpoint_id" {
  description = "ID of the Client VPN endpoint"

  value = module.env.client_vpn_endpoint_id
}

output "client_vpn_subnet_id" {
  description = "ID of the subnet attached to the Client VPN"

  value = module.env.client_vpn_subnet_id
}

output "client_vpn_security_group_id" {
  description = "ID of the security group applied to the target network association"

  value = module.env.client_vpn_security_group_id
}

output "route53_hosted_zone_id" {
  description = "Route53 Hosted Zone ID"

  value = module.env.route53_hosted_zone_id
}

output "vpc_private_subnet_ids" {
  description = "IDs of the VPC private subnets"

  value = module.env.vpc_private_subnet_ids
}

output "vpc_default_security_group_id" {
  description = "ID of the VPC default security group"

  value = module.env.vpc_default_security_group_id
}
