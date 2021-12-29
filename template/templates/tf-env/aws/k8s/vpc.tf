data "aws_availability_zones" "available" {}

locals {
  azs               = slice(data.aws_availability_zones.available.names, 0, 2)
  private_subnets_a = [for idx, subnet_id in module.eks_vpc.private_subnets : subnet_id if idx % 2 == 0]
  private_subnets_b = [for idx, subnet_id in module.eks_vpc.private_subnets : subnet_id if idx % 2 == 1]
}

module "eks_vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 2.48"

  name                 = local.vpc_name
  cidr                 = var.vpc_cidr
  azs                  = local.azs
  public_subnets       = var.public_subnets_cidrs
  private_subnets      = var.private_subnets_cidrs
  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
  enable_dns_support   = true

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                      = "1"
  }
  private_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"             = "1"
  }
  tags = local.common_tags
  vpc_tags = {
    Name = local.vpc_name
  }
}
