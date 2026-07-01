resource "aws_route53_zone" "main" {
  name = "mikekim1032.shop"
}

data "aws_iam_policy_document" "external_dns_policy_doc" {
  statement {
    effect = "Allow"
    actions = [
      "route53:ChangeResourceRecordSets"
    ]
    resources = [
      "arn:aws:route53:::hostedzone/${aws_route53_zone.main.zone_id}"
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "route53:ListHostedZones",
      "route53:ListResourceRecordSets"
    ]
    resources = [
      "*"
    ]
  }
}

resource "aws_iam_policy" "external_dns" {
  name        = "ExternalDNSRoute53Policy"
  description = "Allows ExternalDNS to manage Route53"
  policy      = data.aws_iam_policy_document.external_dns_policy_doc.json
}

data "aws_iam_policy_document" "external_dns_trust" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [module.cluster0.oidc_provider_arn]
    }

    condition {
      test     = "StringEquals"
      # Strips https:// from the OIDC URL
      variable = "${replace(module.cluster0.cluster_oidc_issuer_url, "https://", "")}:sub"
      # Restricts assumption to ONLY the external-dns service account
      values   = ["system:serviceaccount:external-dns:external-dns"]
    }
  }
}

resource "aws_iam_role" "external_dns_role" {
  name               = "ExternalDNS_IRSA_Role"
  assume_role_policy = data.aws_iam_policy_document.external_dns_trust.json
}

resource "aws_iam_role_policy_attachment" "external_dns_attach" {
  role       = aws_iam_role.external_dns_role.name
  policy_arn = aws_iam_policy.external_dns.arn
}