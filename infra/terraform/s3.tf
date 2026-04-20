resource "aws_s3_bucket" "chat_bucket" {
  bucket = "chat-bucket-14129" # Must be globally unique
}


data "aws_iam_policy_document" "allow_cloudfront" {
  statement {
    sid       = "AllowCloudFrontServicePrincipal"
    effect    = "Allow"
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.chat_bucket.arn}/*"]

    principals {
      type        = "Service"
      identifiers = ["cloudfront.amazonaws.com"]
    }

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceArn"
      values   = [aws_cloudfront_distribution.cloudfront_distribution.arn]
    }
  }
}


resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.chat_bucket.id
  policy = data.aws_iam_policy_document.allow_cloudfront.json
}
