data "aws_iam_policy_document" "external_dns_role_assume_role_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [module.eks.oidc_provider_arn]
    }

    condition {
      test     = "StringEquals"
      variable = "${trimprefix(module.eks.cluster_oidc_issuer_url, "https://")}:sub"
      values   = ["system:serviceaccount:${var.external_dns_service_account_namespace}:${var.external_dns_service_account_name}"]
    }
  }
}

resource "aws_iam_role" "external_dns_role" {
  name               = "${var.environment}-external-dns"
  assume_role_policy = data.aws_iam_policy_document.external_dns_role_assume_role_policy.json

  tags = local.common_tags
}

data "aws_iam_policy_document" "external_dns_role_policy" {
  statement {
    effect = "Allow"
    actions = [
      "route53:ChangeResourceRecordSets",
    ]
    resources = ["arn:aws:route53:::hostedzone/*"]
  }

  statement {
    effect = "Allow"
    actions = [
      "route53:ListHostedZones",
      "route53:ListResourceRecordSets",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "external_dns_account_role" {
  name   = "${var.environment}-external-dns"
  role   = aws_iam_role.external_dns_role.id
  policy = data.aws_iam_policy_document.external_dns_role_policy.json
}
