# for off-charging aws
kubectl get secret -n kube-system \
  -l sealedsecrets.bitnami.com/sealed-secrets-key \
  -o yaml > sealed-secrets-key.yaml

kubectl delete gateway cluster-gateway -n envoy-gateway-system

terraform destroy -target="module.cluster1"

terraform apply -target="module.vpc" # remove nat gateway

# ariadne thread
terraform apply -target="module.vpc"

terraform apply

kubectl create namespace backend

kubectl apply -f . #in bootstrap

kubectl apply -f . #in sealedsecret

kubectl rollout restart deployment sealed-secrets-controller -n kube-system

helm upgrade --install holiday . -n backend

kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d