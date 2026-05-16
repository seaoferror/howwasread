terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.28"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }

    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0" # Keeps you on the stable 2.x releases
    }
  }
}
provider "aws" {
  region = local.region
}

provider "helm" {
  kubernetes {
    host = module.cluster1.cluster_endpoint

    cluster_ca_certificate = base64decode(module.cluster1.cluster_certificate_authority_data)

    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", module.cluster1.cluster_name]
    }
  }
}
locals {
  cluster_name = "cluster1"
  region       = "ap-northeast-2"
  azs          = ["ap-northeast-2a", "ap-northeast-2b"]
}
