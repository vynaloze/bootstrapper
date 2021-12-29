locals {
  vpc_name     = "k8s-${var.environment}"
  cluster_name = "eks-${var.environment}"

  common_tags = {
    ManagedBy   = "terraform"
    Environment = lower(var.environment)
  }
}