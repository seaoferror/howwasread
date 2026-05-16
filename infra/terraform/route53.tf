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
      identifiers = [module.cluster1.oidc_provider_arn]
    }

    condition {
      test     = "StringEquals"
      # Strips https:// from the OIDC URL
      variable = "${replace(module.cluster1.cluster_oidc_issuer_url, "https://", "")}:sub"
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

data "aws_lb" "envoy_gateway_nlb" {
  name = "k8s-envoygat-envoyenv-b904f39fc7"
}

resource "aws_route53_record" "backend_alias" {
  zone_id = aws_route53_zone.main.id
  name    = "backend.mikekim1032.shop"
  type    = "A"
  alias {
    name                   = data.aws_lb.envoy_gateway_nlb.dns_name
    zone_id                = data.aws_lb.envoy_gateway_nlb.zone_id
    evaluate_target_health = true
  }
}

# 4. Create the ALIAS record for argocd.mikekim1032.shop
resource "aws_route53_record" "argocd_alias" {
  zone_id = aws_route53_zone.main.id
  name    = "argocd.mikekim1032.shop"
  type    = "A"

  alias {
    name                   = data.aws_lb.envoy_gateway_nlb.dns_name
    zone_id                = data.aws_lb.envoy_gateway_nlb.zone_id
    evaluate_target_health = true
  }
}