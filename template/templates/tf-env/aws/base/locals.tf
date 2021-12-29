data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id

  region_short = join("", regex("([a-z]{2})-([a-z]{1})[a-z]+-([0-9]{1})", data.aws_region.current.name))

  domain_name = "${var.environment}.${coalesce(var.region_short, local.region_short)}.aws.${var.base_domain}"

  vpn_name = "vpn-${var.environment}"

  common_tags = {
    ManagedBy   = "terraform"
    Environment = lower(var.environment)
  }
}
