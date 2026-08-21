resource "cloudflare_zone" "main" {
  account_id = "974f94f87ea3d25fca82e9fb1f408b83"
  zone       = "mikekim1032.shop"
  plan       = "free"
}

resource "cloudflare_record" "apex" {
  zone_id = cloudflare_zone.main.id
  name    = "@"
  content = cloudflare_zero_trust_tunnel_cloudflared.kind_cluster.cname
  type    = "CNAME"
  proxied = true
}

# 2. CNAME record for all subdomains (*.mikekim1032.shop)
resource "cloudflare_record" "wildcard" {
  zone_id = cloudflare_zone.main.id
  name    = "*"
  content = cloudflare_zero_trust_tunnel_cloudflared.kind_cluster.cname
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "turn" {
  zone_id = cloudflare_zone.main.id
  name    = "turn"
  content = "127.0.0.1"
  type    = "A"
  proxied = false

  lifecycle {
    ignore_changes = [content]
  }
}

resource "cloudflare_record" "vercel_app" {
  zone_id = cloudflare_zone.main.id
  name    = "app"
  content   = "1c285eb961bd2418.vercel-dns-017.com"
  type    = "CNAME"
  proxied = false
  ttl     = 1
}

resource "random_id" "tunnel_secret" {
  byte_length = 35
}

resource "cloudflare_zero_trust_tunnel_cloudflared" "kind_cluster" {
  account_id = "974f94f87ea3d25fca82e9fb1f408b83"
  name       = "local-kind-cluster"
  secret     = random_id.tunnel_secret.b64_std
}

resource "cloudflare_zero_trust_tunnel_cloudflared_config" "kind_cluster_config" {
  account_id = "974f94f87ea3d25fca82e9fb1f408b83"
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.kind_cluster.id

  config {

    ingress_rule {
      hostname = "mikekim1032.shop"
      service  = "http://static-envoy-ingress.envoy-gateway-system.svc.cluster.local:80"
    }

    ingress_rule {
      hostname = "*.mikekim1032.shop"
      service  = "http://static-envoy-ingress.envoy-gateway-system.svc.cluster.local:80"
    }

    ingress_rule {
      service = "http_status:404"
    }
  }
}

output "cloudflare_nameservers" {
  value       = cloudflare_zone.main.name_servers
  description = "Copy these nameservers and paste them into your Gabia domain management dashboard."
}

output "cloudflare_zone_id" {
  value       = cloudflare_zone.main.id
  description = "The Cloudflare Zone ID for mikekim1032.shop"
}

output "cloudflare_turn_record_id" {
  value       = cloudflare_record.turn.id
  description = "The DNS Record ID for turn.mikekim1032.shop"
}

output "tunnel_token" {
  value       = cloudflare_zero_trust_tunnel_cloudflared.kind_cluster.tunnel_token
  sensitive   = true
  description = "The token required for the cloudflared pod."
}
