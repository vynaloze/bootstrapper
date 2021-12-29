output "cluster_id" {
  description = "EKS cluster ID"

  value = module.eks.cluster_id
}

output "vpc_cidr_block" {
  description = "EKS VPC CIDR block"

  value = module.eks_vpc.vpc_cidr_block
}

output "worker_iam_role_arn" {
  description = "Worker nodes IAM role ARN"

  value = module.eks.worker_iam_role_arn
}

output "external_dns_role_arn" {
  description = "IAM role for ExternalDNS service account ARN"

  value = aws_iam_role.external_dns_role.arn
}
