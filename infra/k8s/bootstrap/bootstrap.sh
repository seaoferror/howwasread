# for off-charging aws
kubectl get secret -n kube-system \
  -l sealedsecrets.bitnami.com/sealed-secrets-key \
  -o yaml > sealed-secrets-key.yaml

kubectl delete gateway cluster-gateway -n envoy-gateway-system

terraform destroy -target="module.cluster1"

# remove nat gateway
terraform apply -target="module.vpc"

# ariadne thread
terraform apply -target="module.vpc"

terraform apply

kubectl rollout restart deployment sealed-secrets-controller -n kube-system

kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# k8ssandra

# do this before generate cert
kubectl create secret generic k8ssandra-jks-password \
-n k8ssandra \
--from-literal=password="$(openssl rand -base64 24)"

# do this before provisioning cluster
kubectl get secret k8ssandra-jks -n k8ssandra -o json | \
jq '.data.keystore = .data["keystore.jks"] | .data.truststore = .data["truststore.jks"]' | \
kubectl apply -f -

PW=$(kubectl get secret k8ssandra-jks-password -n k8ssandra -o jsonpath='{.data.password}' | base64 -d)

kubectl get secret k8ssandra-jks -n k8ssandra -o json | \
  jq --arg pw "$PW" '
    .data.keystore = .data["keystore.jks"]
    | .data.truststore = .data["truststore.jks"]
    | .data["keystore-password"] = ($pw | @base64)
    | .data["truststore-password"] = ($pw | @base64)
  ' | kubectl apply -f -

# do this after provisioning cluster

kubectl annotate secret k8ssandra-cluster-superuser -n k8ssandra \
reflector.v1.k8s.emberstack.com/reflection-allowed="true" \
reflector.v1.k8s.emberstack.com/reflection-auto-enabled="true" \
reflector.v1.k8s.emberstack.com/reflection-allowed-namespaces="backend" \
--overwrite

CASS_USERNAME=$(kubectl get secret k8ssandra-cluster-superuser -n k8ssandra -o=jsonpath='{.data.username}' | base64 --decode)
echo $CASS_USERNAME
CASS_PASSWORD=$(kubectl get secret k8ssandra-cluster-superuser -n k8ssandra -o=jsonpath='{.data.password}' | base64 --decode)
echo $CASS_PASSWORD
kubectl exec -it k8ssandra-cluster-dc1-default-sts-0 -n k8ssandra -c cassandra -- env SSL_VALIDATE=false cqlsh --ssl -u $CASS_USERNAME -p $CASS_PASSWORD -e "CREATE KEYSPACE production WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': 1};"
