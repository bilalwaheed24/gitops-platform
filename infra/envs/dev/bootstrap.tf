terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 1.7"
}

provider "aws" {
  region = "us-east-1"
}

module "state_backend" {
  source       = "../../modules/s3-state"
  project_name = "gitops-platform"
  account_id   = data.aws_caller_identity.current.account_id
}

data "aws_caller_identity" "current" {}

output "bucket_name"    { value = module.state_backend.bucket_name }
output "dynamodb_table" { value = module.state_backend.dynamodb_table }
