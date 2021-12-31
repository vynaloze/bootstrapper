provider "aws" {
  region = var.region
}

module "env" {
  source = "git@[[ .GitProvider ]]:[[ .GitProject ]]/[[ .GitRepo ]].git//base?ref=[[ .Ref ]]"

  environment           = var.environment
  base_domain           = var.base_domain
  vpc_cidr              = var.vpc_cidr
  client_cidr_block     = var.client_cidr_block
  private_subnets_cidrs = var.private_subnets_cidrs
  public_subnets_cidrs  = var.public_subnets_cidrs
}
