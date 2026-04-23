# Generate a secure RSA key pair locally via Terraform
resource "tls_private_key" "cloudfront_signer_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

# Upload the Public Key to CloudFront
resource "aws_cloudfront_public_key" "signer_public_key" {
  encoded_key = tls_private_key.cloudfront_signer_key.public_key_pem
  name        = "my-cloudfront-signer-key"
}

# Create a Key Group containing the public key, you can attach more public key id
resource "aws_cloudfront_key_group" "signer_key_group" {
  items = [aws_cloudfront_public_key.signer_public_key.id]
  name  = "my-trusted-signers"
}

# Create the Origin Access Control (OAC)
resource "aws_cloudfront_origin_access_control" "s3_oac" {
  name                              = "s3-oac"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# Create the CloudFront Distribution
resource "aws_cloudfront_distribution" "cloudfront_distribution" {
  enabled = true

  origin {
    domain_name              = aws_s3_bucket.chat_bucket.bucket_regional_domain_name
    origin_id                = "myS3Origin"
    origin_access_control_id = aws_cloudfront_origin_access_control.s3_oac.id
  }

  default_cache_behavior {
    target_origin_id       = "myS3Origin"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]

    # THIS enforces that users MUST use a Signed URL or Cookie
    trusted_key_groups = [aws_cloudfront_key_group.signer_key_group.id]

    response_headers_policy_id = "67f7725c-6f97-4210-82d7-5512b31e9d03" # Managed-SecurityHeadersPolicy

    cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6" # Managed-CachingOptimized
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}
