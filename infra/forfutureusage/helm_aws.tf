resource "helm_release" "external_dns" {
  name             = "external-dns"
  repository       = "https://kubernetes-sigs.github.io/external-dns/"
  chart            = "external-dns"
  namespace        = "external-dns"
  create_namespace = true

  values = [
    yamlencode({
      registry = "txt"

      sources = [
        "gateway-httproute",
      ]

      txtOwnerId = module.cluster1.cluster_name

      serviceAccount = {
        annotations = {
          "eks.amazonaws.com/role-arn" = aws_iam_role.external_dns_role.arn
        }
      }

      provider = {
        name = "aws" //aws is default but in case of migration to other vendor
      }

      managedRecordTypes = [
        "A"
      ]

      policy = "sync"

      # nodeSelector = {
      #   "node.kubernetes.io/instance-type" = local.idle_node
      # }
    })
  ]

  depends_on = [aws_iam_role_policy_attachment.external_dns_attach,  module.cluster1.eks_managed_node_groups]
}


locals {
  cluster1 = {
    oidc_url = replace(module.cluster1.oidc_provider_arn, "/^(.*provider/)/", "")
  }
}

data "aws_iam_policy_document" "lbc_trust_policy" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [module.cluster1.oidc_provider_arn]
    }

    condition {
      test     = "StringEquals"
      variable = "${local.cluster1.oidc_url}:sub"
      values   = ["system:serviceaccount:kube-system:aws-load-balancer-controller"]
    }

    condition {
      test     = "StringEquals"
      variable = "${local.cluster1.oidc_url}:aud"
      values   = ["sts.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lbc_role" {
  name               = "aws-load-balancer-controller-role"
  assume_role_policy = data.aws_iam_policy_document.lbc_trust_policy.json
}

data "http" "lbc_iam_policy_json" {
  url = "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json"
}

resource "aws_iam_policy" "lbc_policy" {
  name        = "AWSLoadBalancerControllerIAMPolicy"
  description = "Permissions for the AWS Load Balancer Controller"
  policy      = data.http.lbc_iam_policy_json.response_body
}

resource "aws_iam_role_policy_attachment" "lbc_attach" {
  role       = aws_iam_role.lbc_role.name
  policy_arn = aws_iam_policy.lbc_policy.arn
}

resource "helm_release" "aws_load_balancer_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"

  set {
    name  = "clusterName"
    value = module.cluster1.cluster_name
  }

  set {
    name  = "serviceAccount.create"
    value = "true"
  }

  set {
    name  = "serviceAccount.name"
    value = "aws-load-balancer-controller"
  }

  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = aws_iam_role.lbc_role.arn
  }

  set {
    name  = "vpcId"
    value = module.vpc.vpc_id
  }

  set {
    name  = "region"
    value = local.region
  }

  # values = [
  #   yamlencode({
  #     nodeSelector = {
  #       "node.kubernetes.io/instance-type" = local.idle_node
  #     }
  #   })
  # ]

  depends_on = [module.cluster1.eks_managed_node_groups, module.cluster1.cluster_addons, aws_iam_role_policy_attachment.lbc_attach]
}
