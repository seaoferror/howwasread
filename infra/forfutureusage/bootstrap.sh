helm install eg oci://docker.io/envoyproxy/gateway-helm -n envoy-gateway-system --create-namespace

helm install my-strimzi-cluster-operator oci://quay.io/strimzi-helm/strimzi-kafka-operator -n kafka-system --create-namespace

helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true

helm repo add argo https://argoproj.github.io/argo-helm

helm install argocd -n argocd argo/argo-cd \
  --set redis-ha.enabled=false \
  --set controller.replicas=1 \
  --set server.autoscaling.enabled=false \
  --set server.replicas=1 \
  --set repoServer.autoscaling.enabled=false \
  --set repoServer.replicas=1 \
  --set applicationSet.replicaCount=1 \
  --set notifications.enabled=false \
  --set server.service.type=ClusterIP \
  --set server.extraArgs=\{--insecure\} \
  --create-namespace


kubectl port-forward svc/argocd-server -n argocd 8080:80

