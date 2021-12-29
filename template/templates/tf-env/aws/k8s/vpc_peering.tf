resource "aws_vpc_peering_connection" "this" {
  peer_vpc_id = var.base_vpc_id
  vpc_id      = module.eks_vpc.vpc_id
  auto_accept = true

  accepter {
    allow_remote_vpc_dns_resolution = true
  }

  requester {
    allow_remote_vpc_dns_resolution = true
  }

  tags = merge(
    {
      Name = "${local.vpc_name}-to-${var.environment}"
    },
    local.common_tags
  )
}

resource "aws_ec2_client_vpn_route" "this" {
  client_vpn_endpoint_id = var.client_vpn_endpoint_id
  target_vpc_subnet_id   = var.client_vpn_subnet_id
  destination_cidr_block = module.eks_vpc.vpc_cidr_block
}

# count instead of for_each is a workaround for https://github.com/hashicorp/terraform/issues/4149
resource "aws_route" "eks_to_base_vpc" {
  count = length(module.eks_vpc.private_route_table_ids)

  route_table_id = element(module.eks_vpc.private_route_table_ids, count.index)

  destination_cidr_block    = var.base_vpc_cidr
  vpc_peering_connection_id = aws_vpc_peering_connection.this.id
}

resource "aws_route" "base_vpc_to_eks" {
  for_each = var.base_vpc_private_route_table_ids

  route_table_id = each.value

  destination_cidr_block    = module.eks_vpc.vpc_cidr_block
  vpc_peering_connection_id = aws_vpc_peering_connection.this.id
}

resource "aws_route53_zone_association" "this" {
  zone_id = var.route53_hosted_zone_id
  vpc_id  = module.eks_vpc.vpc_id
}
