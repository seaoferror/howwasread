module "cluster1" {
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

  enable_irsa = true

  eks_managed_node_groups = {
    default = {
      instance_types = ["t3.medium"]
      min_size     = 2
      max_size     = 5
      desired_size = 2
    }
  }
}

provider "helm" {
  kubernetes {
    host = module.cluster1.cluster_endpoint

    cluster_ca_certificate = base64decode(module.cluster1.cluster_certificate_authority_data)

    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args = ["eks", "get-token", "--cluster-name", module.cluster1.cluster_name]
    }
  }
}

provider "kubernetes" {
  host                   = module.cluster1.cluster_endpoint
  cluster_ca_certificate = base64decode(module.cluster1.cluster_certificate_authority_data)
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", module.cluster1.cluster_name]
  }
}

resource "helm_release" "external_dns" {
  name             = "external-dns"
  repository       = "https://kubernetes-sigs.github.io/external-dns/"
  chart            = "external-dns"
  namespace        = "external-dns"
  create_namespace = true

  values = [
    yamlencode({
      registry = "txt"

      source = [
        "gateway-httproute"
      ]

      aws = {
        zoneType = "public"
      }

      txtOwnerId = module.cluster1.cluster_name

      serviceAccount = {
        create = true
        annotations = {
          "eks.amazonaws.com/role-arn" = aws_iam_role.external_dns_role.arn
        }
      }

      policy = "sync"
    })
  ]
}

resource "helm_release" "envoy_gateway" {
  name             = "eg"
  repository       = "oci://docker.io/envoyproxy"
  chart            = "gateway-helm"
  namespace        = "envoy-gateway-system"
  create_namespace = true
}

resource "helm_release" "strimzi" {
  name             = "my-strimzi-cluster-operator"
  repository       = "oci://quay.io/strimzi-helm"
  chart            = "strimzi-kafka-operator"
  namespace        = "kafka-system"
  create_namespace = true
}

resource "helm_release" "cert_manager" {
  name             = "cert-manager"
  repository       = "https://charts.jetstack.io"
  chart            = "cert-manager"
  namespace        = "cert-manager"
  create_namespace = true

  set {
    name  = "installCRDs"
    value = "true"
  }
}

resource "helm_release" "argocd" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = "argocd"
  create_namespace = true

  values = [
    yamlencode({
      server = {
        extraArgs = ["--insecure"]
        service = {
          type = "ClusterIP"
        }
        replicas    = 1
        autoscaling = { enabled = false }
      }
      repoServer = {
        replicas    = 1
        autoscaling = { enabled = false }
      }
      controller = {
        replicas = 1
      }
      applicationSet = {
        replicaCount = 1
      }
      redis-ha = {
        enabled = false
      }
      notifications = {
        enabled = false
      }
    })
  ]
}