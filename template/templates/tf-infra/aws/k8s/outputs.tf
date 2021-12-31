output "cluster_id" {
  description = "EKS cluster ID"

  value = module.env.cluster_id
}

output "worker_iam_role_arn" {
  description = "Worker nodes IAM role ARN"

  value = module.env.worker_iam_role_arn
}

output "external_dns_role_arn" {
  description = "IAM role for ExternalDNS service account ARN"

  value = module.env.external_dns_role_arn
}

output "vpc_cidr_block" {
  description = "EKS VPC CIDR block"

  value = module.env.vpc_cidr_block
}
