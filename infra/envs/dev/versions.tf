terraform {
  backend "s3" {
    bucket         = "gitops-platform-tfstate-290172088615"
    key            = "dev/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "gitops-platform-tf-locks"
    encrypt        = true
  }
}
