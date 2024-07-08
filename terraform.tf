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
      TerraformWorkspace = terraform.workspace
    }
  }
}

terraform {
  required_version = "1.9.1"

  backend "s3" {
    bucket         = "tf-state-hday"
    key            = "terraform.state"
    region         = "us-east-1"
    dynamodb_table = "terraform_state"

  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_lambda_function" "crops_crop_add" {
  filename                       = "data/functions/crops/crop_add/bootstrap.zip"
  function_name                  = "crops_crop_add"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "p85k2hpl51PgJdCW9ZE/r1JBjizTLIL5MlTkWk8udlE="
  role                           = aws_iam_role.crops_crop_add_role.arn
}

resource "aws_iam_role" "crops_crop_add_role" {
  name               = "crops_crop_add_role"
  assume_role_policy = data.aws_iam_policy_document.crops_crop_add_role.json
}

data "aws_iam_policy_document" "crops_crop_add_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_lambda_function" "crops_crop_get" {
  filename                       = "data/functions/crops/crop_get/bootstrap.zip"
  function_name                  = "crops_crop_get"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "7TAak3Vv/s92tlG9h9p1Xc5nA8JMSIonCaVA14rv1To="
  role                           = aws_iam_role.crops_crop_get_role.arn
}

resource "aws_iam_role" "crops_crop_get_role" {
  name               = "crops_crop_get_role"
  assume_role_policy = data.aws_iam_policy_document.crops_crop_get_role.json
}

data "aws_iam_policy_document" "crops_crop_get_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_lambda_function" "products_product_get" {
  filename                       = "data/functions/products/product_get/bootstrap.zip"
  function_name                  = "products_product_get"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "wDoibJ1wkU5X1+OL6A7DUyEFq/FJBBsZpIRNyhEQIWY="
  role                           = aws_iam_role.products_product_get_role.arn
}

resource "aws_iam_role" "products_product_get_role" {
  name               = "products_product_get_role"
  assume_role_policy = data.aws_iam_policy_document.products_product_get_role.json
}

data "aws_iam_policy_document" "products_product_get_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_lambda_function" "products_product_init" {
  filename                       = "data/functions/products/product_init/bootstrap.zip"
  function_name                  = "products_product_init"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "UAGNSEE2JN5ceW8OcX+iSHrnmFtJDAbT5xdip+5Qmq4="
  role                           = aws_iam_role.products_product_init_role.arn
}

resource "aws_iam_role" "products_product_init_role" {
  name               = "products_product_init_role"
  assume_role_policy = data.aws_iam_policy_document.products_product_init_role.json
}

data "aws_iam_policy_document" "products_product_init_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}


