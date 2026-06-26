module "cluster1" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  name               = local.cluster_name
  region             = local.region
  kubernetes_version = "1.36"

  create_cloudwatch_log_group = false
  enabled_log_types           = []

  compute_config = {
    enabled    = false
    node_pools = []
  }

  addons = {
    coredns = {}
    eks-pod-identity-agent = {
      before_compute = true
    }
    kube-proxy = {}
    vpc-cni = {
      before_compute = true
      configuration_values = jsonencode({
        env = {
          ENABLE_PREFIX_DELEGATION = "true"
          WARM_PREFIX_TARGET       = "1"
        }
      })
    }
    # aws-ebs-csi-driver = {
    #   most_recent              = true
    #   service_account_role_arn = aws_iam_role.ebs_csi_role.arn
    # }
  }

  vpc_id                   = module.vpc.vpc_id
  subnet_ids               = module.vpc.private_subnets
  control_plane_subnet_ids = module.vpc.private_subnets

  endpoint_public_access                   = true
  endpoint_private_access                  = true
  enable_cluster_creator_admin_permissions = true

  enable_irsa = true

  eks_managed_node_groups = {
    medium = {
      instance_types = ["t3.medium"]
      min_size       = 1
      max_size       = 1
      desired_size   = 1
    }
    not_medium = {
      instance_types = [local.idle_node]
      min_size       = 1
      max_size       = 1
      desired_size   = 1
    }

  }

  node_security_group_additional_rules = {
    ingress_cluster_8080 = {
      description                   = "Allow Control Plane to talk to 8080 port where Sealed Secrets Controller locate in order to use kubeseal"
      protocol                      = "tcp"
      from_port                     = 8080
      to_port                       = 8080
      type                          = "ingress" //ingress  egress 's ingress
      source_cluster_security_group = true
    }
  }
}
#
# data "aws_iam_policy_document" "ebs_csi_trust_policy" {
#   statement {
#     actions = ["sts:AssumeRoleWithWebIdentity"]
#     effect  = "Allow"
#
#     principals {
#       type        = "Federated"
#       identifiers = [module.cluster1.oidc_provider_arn]
#     }
#
#     # Restrict the role to ONLY the specific namespace and service account used by the driver
#     condition {
#       test = "StringEquals"
#       # Note: We use module.eks.oidc_provider here (the URL), NOT the ARN
#       variable = "${module.cluster1.oidc_provider}:sub"
#       values   = ["system:serviceaccount:kube-system:ebs-csi-controller-sa"]
#     }
#
#     condition {
#       test     = "StringEquals"
#       variable = "${module.cluster1.oidc_provider}:aud"
#       values   = ["sts.amazonaws.com"]
#     }
#   }
# }
#
# resource "aws_iam_role" "ebs_csi_role" {
#   name               = "ebs-csi-driver-role"
#   assume_role_policy = data.aws_iam_policy_document.ebs_csi_trust_policy.json
# }
#
# resource "aws_iam_role_policy_attachment" "ebs_csi_policy_attachment" {
#   role       = aws_iam_role.ebs_csi_role.name
#   policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
# }
