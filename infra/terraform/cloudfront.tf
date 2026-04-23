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

resource "aws_cloudfront_response_headers_policy" "security_headers" {
  name    = "standard-security-headers-policy"
  comment = "Applies standard security headers to all responses"

  security_headers_config {
    frame_options {
      frame_option = "DENY"
      override     = true
    }
    content_type_options {
      override = true
    }
    referrer_policy {
      referrer_policy = "strict-origin-when-cross-origin"
      override        = true
    }
    xss_protection {
      protection = true
      mode_block = true
      override   = true
    }
  }
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

    response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers.id

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
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
