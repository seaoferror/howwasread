output "cloudfront_key_pair_id" {
  value       = aws_cloudfront_public_key.signer_public_key.id
  description = "The ID of the CloudFront Public Key (Key Pair ID)"
}

output "cloudfront_private_pem_key" {
  value       = tls_private_key.cloudfront_signer_key.private_key_pem
  description = "The raw private PEM key to be used by your backend to sign URLs"
  sensitive   = true
}

output "route53_name_servers" {
  description = "The Name Servers for the Route53 Zone"
  value       = aws_route53_zone.main.name_servers
}

# output "ecr_repository_urls" {
#   description = "A map of repository names to their URLs"
#   value       = { for key, repo in aws_ecr_repository.repos : key => repo.repository_url }
# }

output "configure_kubectl" {
  description = "Run this command to configure kubectl"
  value       = "aws eks update-kubeconfig --region ${local.region} --name ${local.cluster_name}"
}

data "kubernetes_secret" "argocd_admin_pwd" {
  metadata {
    name      = "argocd-initial-admin-secret"
    namespace = "argocd"
  }

  # This ensures Terraform doesn't try to read the secret before the helm chart creates it
  depends_on = [helm_release.argocd]
}

output "argocd_password" {
  value     = data.kubernetes_secret.argocd_admin_pwd.data["password"]
  sensitive = true
}