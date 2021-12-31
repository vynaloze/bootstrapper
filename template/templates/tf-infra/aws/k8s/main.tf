provider "aws" {
  region = var.region
}

module "env" {
  source = "git@[[ .GitProvider ]]:[[ .GitProject ]]/[[ .GitRepo ]].git//k8s?ref=[[ .Ref ]]"

  environment = var.environment

  vpc_cidr              = var.vpc_cidr
  private_subnets_cidrs = var.private_subnets_cidrs
  public_subnets_cidrs  = var.public_subnets_cidrs

  base_vpc_cidr                    = data.terraform_remote_state.base.outputs.vpc_cidr
  base_vpc_id                      = data.terraform_remote_state.base.outputs.vpc_id
  base_vpc_private_route_table_ids = data.terraform_remote_state.base.outputs.vpc_private_route_table_ids
  client_vpn_endpoint_id           = data.terraform_remote_state.base.outputs.client_vpn_endpoint_id
  client_vpn_subnet_id             = data.terraform_remote_state.base.outputs.client_vpn_subnet_id
  route53_hosted_zone_id           = data.terraform_remote_state.base.outputs.route53_hosted_zone_id
}
