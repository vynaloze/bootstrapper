resource "aws_route53_zone" "this" {
  name = local.domain_name

  vpc {
    vpc_id = module.vpc.vpc_id
  }

  tags = local.common_tags

  # Required to manage additional VPC associations with aws_route53_zone_association resource.
  # See: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone_association
  lifecycle {
    ignore_changes = [vpc]
  }
}