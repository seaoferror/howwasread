module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  name    = local.cluster_name
  region = local.region
  kubernetes_version= "1.35"

  compute_config = {
    enabled = false
    node_pools = []
  }

  addons = {
    coredns                = {}
    eks-pod-identity-agent = {
      before_compute = true
    }
    kube-proxy             = {}
    vpc-cni                = {
      before_compute = true
    }
  }

  vpc_id                   = aws_vpc.vpc.id
  subnet_ids               = aws_subnet.private[*].id
  control_plane_subnet_ids = aws_subnet.private[*].id

  endpoint_public_access = true
  endpoint_private_access = true
  enable_cluster_creator_admin_permissions = true

  eks_managed_node_groups = {
    default = {
      instance_types = ["t3.medium"]
      min_size     = 2
      max_size     = 5
      desired_size = 2
    }
  }
}