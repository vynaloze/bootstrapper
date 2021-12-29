module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 17.1"

  cluster_name    = local.cluster_name
  cluster_version = "1.21"
  subnets         = module.eks_vpc.private_subnets
  vpc_id          = module.eks_vpc.vpc_id

  cluster_endpoint_private_access                = true
  cluster_endpoint_private_access_cidrs          = ["0.0.0.0/0"]
  cluster_create_endpoint_private_access_sg_rule = true
  cluster_endpoint_public_access                 = false
  enable_irsa                                    = true
  manage_aws_auth                                = false
  write_kubeconfig                               = false

  node_groups_defaults = {
    ami_type            = var.worker_general_purpose_ami_type
    ami_release_version = "1.21.2-20210826"

    capacity_type    = "ON_DEMAND"
    instance_types   = [var.worker_general_purpose_node_type]
    desired_capacity = var.worker_general_purpose_group_size
    min_capacity     = 1
    max_capacity     = var.worker_general_purpose_max_group_size

    disk_type = "gp2"
    disk_size = var.worker_general_purpose_disk_size

    create_launch_template = true
    enable_monitoring      = true
  }

  node_groups = {
    "gen-purpose-a" = {
      subnets = local.private_subnets_a
    },

    "gen-purpose-b" = {
      subnets = local.private_subnets_b
    },
  }

  tags = local.common_tags

  cluster_enabled_log_types     = ["api", "audit", "authenticator", "controllerManager", "scheduler"]
  cluster_log_retention_in_days = var.cluster_log_retention_days
}
