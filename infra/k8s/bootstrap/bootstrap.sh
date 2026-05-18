kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# for off-charging aws
terraform destroy -target="module.cluster1"

terraform apply -target="module.vpc" # remove nat gateway

# ariadne thread
terraform apply -target="module.vpc"

terraform apply