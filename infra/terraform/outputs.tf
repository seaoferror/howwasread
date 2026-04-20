# CloudFront Key Pair ID and PEM File
output "cloudfront_key_pair_id" {
  value       = aws_cloudfront_public_key.signer_public_key.id
  description = "The ID of the CloudFront Public Key (Key Pair ID)"
}

output "cloudfront_private_pem_key" {
  value       = tls_private_key.cloudfront_signer_key.private_key_pem
  description = "The raw private PEM key to be used by your backend to sign URLs"
  sensitive   = true
}

