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

locals {
  cluster_name = "cluster1"
  region       = "ap-northeast-2"
  azs          = ["ap-northeast-2a", "ap-northeast-2b"]
}
