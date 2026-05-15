terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.28"
    }
  }
}
provider "aws" {
  region = local.region
}

locals {
  cluster_name = "cluster0"
  region       = "ap-northeast-2"
  azs          = ["ap-northeast-2a", "ap-northeast-2b"]
}
