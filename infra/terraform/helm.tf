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

resource "helm_release" "strimzi" {
  name             = "my-strimzi-cluster-operator"
  repository       = "oci://quay.io/strimzi-helm"
  chart            = "strimzi-kafka-operator"
  namespace        = "kafka-system"
  create_namespace = true

  # depends_on = [
  #   module.cluster1.eks_managed_node_groups
  # ]
}

resource "helm_release" "k8ssandra" {
  name             = "k8ssandra"
  repository       = "https://helm.k8ssandra.io/stable"
  chart            = "k8ssandra-operator"
  namespace        = "k8ssandra"
  create_namespace = true


  depends_on = [helm_release.cert_manager]

  # set {
  #   name  = "cassandra.version"
  #   value = "4.0.0"
  # }
}

resource "helm_release" "flink_operator" {
  name             = "flink-operator"
  repository       = "https://downloads.apache.org/flink/flink-kubernetes-operator-1.15.0/"
  chart            = "flink-kubernetes-operator"
  namespace        = "flink"
  create_namespace = true

  depends_on = [helm_release.cert_manager]

  set {
    name  = "webhook.create"
    value = "true"
  }

  set_list {
    name  = "watchNamespaces"
    value = ["backend", "kafka-system", "k8ssandra"]
  }
}

resource "helm_release" "trust_manager" {
  name       = "trust-manager"
  repository = "https://charts.jetstack.io"
  chart      = "trust-manager"
  namespace  = "cert-manager"

  set {
    name  = "app.trust.namespace"
    value = "cert-manager"
  }
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

resource "helm_release" "alloy" {
  name             = "alloy"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "alloy"
  namespace        = "monitoring"
  create_namespace = true

  # Passing the complete Alloy River configuration via YAML values
  values = [
    <<-EOF
    alloy:
      envFrom:
        - secretRef:
            name: grafanacloud

      configMap:
        create: true
        content: |
          // 2. Read the environment variable in the Loki config
          loki.write "grafanacloud" {
            endpoint {
              url = "https://logs-prod-030.grafana.net/loki/api/v1/push"

              basic_auth {
                username = "1654261"
                password = sys.env("GRAFANA_CLOUD_TOKEN")
              }
            }
          }

          // ==========================================
          // 2. KUBERNETES DISCOVERY & BASE LABELS
          // ==========================================
          discovery.kubernetes "pods" {
            role = "pod"
          }

          discovery.relabel "base_labels" {
            targets = discovery.kubernetes.pods.targets

            rule {
              source_labels = ["__meta_kubernetes_namespace"]
              target_label  = "namespace"
            }
            rule {
              source_labels = ["__meta_kubernetes_pod_name"]
              target_label  = "pod"
            }
            rule {
              source_labels = ["__meta_kubernetes_pod_container_name"]
              target_label  = "container"
            }
          }

          // ==========================================
          // 3. Default
          // ==========================================
          discovery.relabel "none" {
            targets = discovery.relabel.base_labels.output
            rule {
              action        = "keep"
              // Targets the new label we added to your Helm template
              source_labels = ["__meta_kubernetes_pod_label_framework"]
              regex         = "none"
            }
          }

          loki.source.kubernetes "default_logs" {
            targets    = discovery.relabel.none.output
            forward_to = [loki.write.grafanacloud.receiver]
          }

          // ==========================================
          // 4. SPRING BOOT APPLICATIONS
          // ==========================================
          discovery.relabel "spring_apps" {
            targets = discovery.relabel.base_labels.output
            rule {
              action        = "keep"
              source_labels = ["__meta_kubernetes_pod_label_framework"]
              regex         = "spring-boot"
            }
          }

          loki.source.kubernetes "spring_logs" {
            targets    = discovery.relabel.spring_apps.output
            // Route to the multiline processor first
            forward_to = [loki.process.spring_multiline.receiver]
          }

          loki.process "spring_multiline" {
            // Forward processed logs to Loki
            forward_to = [loki.write.grafanacloud.receiver]

            stage.multiline {
              // Matches standard Spring Boot timestamp format: YYYY-MM-DD HH:MM:SS.mmm
              firstline     = "^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}"
              max_wait_time = "3s"
            }
          }
    EOF
  ]
}

resource "helm_release" "valkey_operator" {
  name             = "valkey-operator"
  repository       = "https://valkey.io/valkey-helm"
  chart            = "valkey-operator"
  namespace        = "valkey-system"
  create_namespace = true

  # values = [
  #   yamlencode({
  #     nodeSelector = {
  #       "node.kubernetes.io/instance-type" = local.idle_node
  #     }
  #   })
  # ]

  # depends_on = [
  #   module.cluster0.eks_managed_node_groups
  # ]
}

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
