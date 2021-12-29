output "vpc_id" {
  description = "VPC ID"

  value = module.vpc.vpc_id
}

output "vpc_cidr" {
  description = "The CIDR block for the VPC"

  value = var.vpc_cidr
}

output "vpc_private_subnet_ids" {
  description = "IDs of the VPC private subnets"

  value = module.vpc.private_subnets
}

output "vpc_private_route_table_ids" {
  description = "IDs of the route tables associated with private subnets in VPC"

  value = module.vpc.private_route_table_ids
}

output "vpc_default_security_group_id" {
  description = "ID of the VPC default security group"

  value = module.vpc.default_security_group_id
}

output "client_vpn_endpoint_id" {
  description = "ID of the Client VPN endpoint"

  value = aws_ec2_client_vpn_endpoint.this.id
}

output "client_vpn_subnet_id" {
  description = "ID of the subnet attached to the Client VPN"

  value = aws_ec2_client_vpn_network_association.this.subnet_id
}

output "client_vpn_security_group_id" {
  description = "ID of the security group applied to the target network association"

  value = aws_security_group.vpn_access.id
}

output "route53_hosted_zone_id" {
  description = "Route53 Hosted Zone ID"

  value = aws_route53_zone.this.zone_id
}

