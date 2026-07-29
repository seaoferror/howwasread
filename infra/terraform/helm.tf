resource "helm_release" "envoy_gateway" {
  name             = "eg"
  repository       = "oci://docker.io/envoyproxy"
  chart            = "gateway-helm"
  namespace        = "envoy-gateway-system"
  create_namespace = true

  # depends_on = [
  #   module.cluster1.eks_managed_node_groups
  # ]
}

# resource "helm_release" "strimzi" {
#   name             = "my-strimzi-cluster-operator"
#   repository       = "oci://quay.io/strimzi-helm"
#   chart            = "strimzi-kafka-operator"
#   namespace        = "kafka-system"
#   create_namespace = true
#
#   depends_on = [
#     module.cluster1.eks_managed_node_groups
#   ]
# }

# resource "helm_release" "external_dns" {
#   name             = "external-dns"
#   repository       = "https://kubernetes-sigs.github.io/external-dns/"
#   chart            = "external-dns"
#   namespace        = "external-dns"
#   create_namespace = true
#
#   values = [
#     yamlencode({
#       registry = "txt"
#
#       sources = [
#         "gateway-httproute",
#       ]
#
#       txtOwnerId = module.cluster1.cluster_name
#
#       serviceAccount = {
#         annotations = {
#           "eks.amazonaws.com/role-arn" = aws_iam_role.external_dns_role.arn
#         }
#       }
#
#       provider = {
#         name = "aws" //aws is default but in case of migration to other vendor
#       }
#
#       managedRecordTypes = [
#         "A"
#       ]
#
#       policy = "sync"
#
#       # nodeSelector = {
#       #   "node.kubernetes.io/instance-type" = local.idle_node
#       # }
#     })
#   ]
#
#   depends_on = [aws_iam_role_policy_attachment.external_dns_attach,  module.cluster1.eks_managed_node_groups]
# }

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

  set {
    name  = "config.apiVersion"
    value = "controller.config.cert-manager.io/v1alpha1"
  }

  set {
    name  = "config.kind"
    value = "ControllerConfiguration"
  }

  set {
    name  = "config.enableGatewayAPI"
    value = "true"
  }

  # values = [
  #   yamlencode({
  #     nodeSelector = {
  #       "node.kubernetes.io/instance-type" = local.idle_node
  #     }
  #     webhook = {
  #       nodeSelector = {
  #         "node.kubernetes.io/instance-type" = local.idle_node
  #       }
  #     }
  #     cainjector = {
  #       nodeSelector = {
  #         "node.kubernetes.io/instance-type" = local.idle_node
  #       }
  #     }
  #   })
  # ]

  # depends_on = [
  #   module.cluster1.eks_managed_node_groups
  # ]
}

# resource "helm_release" "argocd" {
#   name             = "argocd"
#   repository       = "https://argoproj.github.io/argo-helm"
#   chart            = "argo-cd"
#   namespace        = "argocd"
#   create_namespace = true
#
#   values = [
#     yamlencode({
#       server = {
#         extraArgs = ["--insecure"]
#         service = {
#           type = "ClusterIP"
#         }
#         replicas    = 1
#         autoscaling = { enabled = false }
#       }
#       repoServer = {
#         replicas    = 1
#         autoscaling = { enabled = false }
#       }
#       controller = {
#         replicas = 1
#       }
#       applicationSet = {
#         replicaCount = 1
#       }
#       redis-ha = {
#         enabled = false
#       }
#       notifications = {
#         enabled = false
#       }
#     })
#   ]
#
#   depends_on = [
#     module.cluster1.eks_managed_node_groups
#   ]
# }

# locals {
#   cluster1 = {
#     oidc_url = replace(module.cluster1.oidc_provider_arn, "/^(.*provider/)/", "")
#   }
# }
#
# data "aws_iam_policy_document" "lbc_trust_policy" {
#   statement {
#     effect  = "Allow"
#     actions = ["sts:AssumeRoleWithWebIdentity"]
#
#     principals {
#       type        = "Federated"
#       identifiers = [module.cluster1.oidc_provider_arn]
#     }
#
#     condition {
#       test     = "StringEquals"
#       variable = "${local.cluster1.oidc_url}:sub"
#       values   = ["system:serviceaccount:kube-system:aws-load-balancer-controller"]
#     }
#
#     condition {
#       test     = "StringEquals"
#       variable = "${local.cluster1.oidc_url}:aud"
#       values   = ["sts.amazonaws.com"]
#     }
#   }
# }
#
# resource "aws_iam_role" "lbc_role" {
#   name               = "aws-load-balancer-controller-role"
#   assume_role_policy = data.aws_iam_policy_document.lbc_trust_policy.json
# }
#
# data "http" "lbc_iam_policy_json" {
#   url = "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json"
# }
#
# resource "aws_iam_policy" "lbc_policy" {
#   name        = "AWSLoadBalancerControllerIAMPolicy"
#   description = "Permissions for the AWS Load Balancer Controller"
#   policy      = data.http.lbc_iam_policy_json.response_body
# }
#
# resource "aws_iam_role_policy_attachment" "lbc_attach" {
#   role       = aws_iam_role.lbc_role.name
#   policy_arn = aws_iam_policy.lbc_policy.arn
# }

# resource "helm_release" "aws_load_balancer_controller" {
#   name       = "aws-load-balancer-controller"
#   repository = "https://aws.github.io/eks-charts"
#   chart      = "aws-load-balancer-controller"
#   namespace  = "kube-system"
#
#   set {
#     name  = "clusterName"
#     value = module.cluster1.cluster_name
#   }
#
#   set {
#     name  = "serviceAccount.create"
#     value = "true"
#   }
#
#   set {
#     name  = "serviceAccount.name"
#     value = "aws-load-balancer-controller"
#   }
#
#   set {
#     name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
#     value = aws_iam_role.lbc_role.arn
#   }
#
#   set {
#     name  = "vpcId"
#     value = module.vpc.vpc_id
#   }
#
#   set {
#     name  = "region"
#     value = local.region
#   }
#
#   # values = [
#   #   yamlencode({
#   #     nodeSelector = {
#   #       "node.kubernetes.io/instance-type" = local.idle_node
#   #     }
#   #   })
#   # ]
#
#   depends_on = [module.cluster1.eks_managed_node_groups, module.cluster1.cluster_addons, aws_iam_role_policy_attachment.lbc_attach]
# }

resource "helm_release" "sealed_secrets" {
  name       = "sealed-secrets-controller"
  repository = "oci://registry-1.docker.io/bitnamicharts"
  chart      = "sealed-secrets"
  namespace  = "kube-system"

  # values = [
  #   yamlencode({
  #     nodeSelector = {
  #       "node.kubernetes.io/instance-type" = local.idle_node
  #     }
  #   })
  # ]

  # depends_on = [
  #   module.cluster1.eks_managed_node_groups
  # ]
}

resource "helm_release" "reflector" {
  name       = "reflector-controller"
  repository = "oci://ghcr.io/emberstack/helm-charts"
  chart      = "reflector"
  namespace  = "kube-system"

  # values = [
  #   yamlencode({
  #     nodeSelector = {
  #       "node.kubernetes.io/instance-type" = local.idle_node
  #     }
  #   })
  # ]

  # depends_on = [
  #   module.cluster1.eks_managed_node_groups
  # ]
}

# resource "helm_release" "alloy" {
#   name             = "alloy"
#   repository       = "https://grafana.github.io/helm-charts"
#   chart            = "alloy"
#   namespace        = "monitoring"
#   create_namespace = true
#
#   # Passing the complete Alloy River configuration via YAML values
#   values = [
#     <<-EOF
#     alloy:
#       envFrom:
#         - secretRef:
#             name: grafanacloud
#
#       configMap:
#         create: true
#         content: |
#           // 2. Read the environment variable in the Loki config
#           loki.write "grafanacloud" {
#             endpoint {
#               url = "https://logs-prod-030.grafana.net/loki/api/v1/push"
#
#               basic_auth {
#                 username = "1654261"
#                 password = sys.env("GRAFANA_CLOUD_TOKEN")
#               }
#             }
#           }
#
#           // ==========================================
#           // 2. KUBERNETES DISCOVERY & BASE LABELS
#           // ==========================================
#           discovery.kubernetes "pods" {
#             role = "pod"
#           }
#
#           discovery.relabel "base_labels" {
#             targets = discovery.kubernetes.pods.targets
#
#             rule {
#               source_labels = ["__meta_kubernetes_namespace"]
#               target_label  = "namespace"
#             }
#             rule {
#               source_labels = ["__meta_kubernetes_pod_name"]
#               target_label  = "pod"
#             }
#             rule {
#               source_labels = ["__meta_kubernetes_pod_container_name"]
#               target_label  = "container"
#             }
#           }
#
#           // ==========================================
#           // 3. Default
#           // ==========================================
#           discovery.relabel "none" {
#             targets = discovery.relabel.base_labels.output
#             rule {
#               action        = "keep"
#               // Targets the new label we added to your Helm template
#               source_labels = ["__meta_kubernetes_pod_label_framework"]
#               regex         = "none"
#             }
#           }
#
#           loki.source.kubernetes "default_logs" {
#             targets    = discovery.relabel.none.output
#             forward_to = [loki.write.grafanacloud.receiver]
#           }
#
#           // ==========================================
#           // 4. SPRING BOOT APPLICATIONS
#           // ==========================================
#           discovery.relabel "spring_apps" {
#             targets = discovery.relabel.base_labels.output
#             rule {
#               action        = "keep"
#               source_labels = ["__meta_kubernetes_pod_label_framework"]
#               regex         = "spring-boot"
#             }
#           }
#
#           loki.source.kubernetes "spring_logs" {
#             targets    = discovery.relabel.spring_apps.output
#             // Route to the multiline processor first
#             forward_to = [loki.process.spring_multiline.receiver]
#           }
#
#           loki.process "spring_multiline" {
#             // Forward processed logs to Loki
#             forward_to = [loki.write.grafanacloud.receiver]
#
#             stage.multiline {
#               // Matches standard Spring Boot timestamp format: YYYY-MM-DD HH:MM:SS.mmm
#               firstline     = "^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}"
#               max_wait_time = "3s"
#             }
#           }
#     EOF
#   ]
# }

# resource "helm_release" "valkey_operator" {
#   name             = "valkey-operator
#   repository       = "https://valkey.io/valkey-helm"
#   chart            = "valkey-operator"
#   namespace        = "valkey-system"
#   create_namespace = true
#
#   values = [
#     yamlencode({
#       nodeSelector = {
#         "node.kubernetes.io/instance-type" = local.idle_node
#       }
#     })
#   ]
#
#   depends_on = [
#     module.cluster0.eks_managed_node_groups
#   ]
# }

# resource "helm_release" "metrics_server" {
#   name       = "metrics-server"
#   repository = "https://kubernetes-sigs.github.io/metrics-server/"
#   chart      = "metrics-server"
#   namespace  = "kube-system"
#
#   depends_on = [
#     module.cluster1.eks_managed_node_groups
#   ]
# }
