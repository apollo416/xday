variable "workspace_iam_role" {
  type        = string
  description = "the iam role to assume for the workspace"
  nullable    = false
  sensitive   = true
}

provider "aws" {
  region = "us-east-1"

  assume_role {
    role_arn = var.workspace_iam_role
  }

  default_tags {
    tags = {
      Project            = "xday"
      TerraformWorkspace = "xday"
    }
  }
}

terraform {
  required_version = "1.9.1"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

