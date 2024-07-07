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

resource "aws_dynamodb_table" "crops" {
  name           = "crops"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  depends_on = [
  ]
}

resource "aws_dynamodb_table" "products" {
  name           = "products"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  depends_on = [
    aws_lambda_function.products_product_get,
    aws_lambda_function.products_product_init,
  ]
}

resource "aws_lambda_function" "crops_crop_add" {
  filename                       = "./data/functions/crops/crop_add/bootstrap.zip"
  function_name                  = "crops_crop_add"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "wKCn/01cCRspCjQtLlDcstFN2NFXf6Y3jH5T1JFOelo="
  role                           = aws_iam_role.crops_crop_add_function_role.arn
}

resource "aws_iam_role" "crops_crop_add_function_role" {
  name               = "crops_crop_add_function_role"
  assume_role_policy = data.aws_iam_policy_document.crops_crop_add_function_role.json
}

data "aws_iam_policy_document" "crops_crop_add_function_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_lambda_function" "crops_crop_get" {
  filename                       = "./data/functions/crops/crop_get/bootstrap.zip"
  function_name                  = "crops_crop_get"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "nIVq/2t6fubILKIkF4WRZSgEsZcFXQNhkVHIQjYNUq0="
  role                           = aws_iam_role.crops_crop_get_function_role.arn
}

resource "aws_iam_role" "crops_crop_get_function_role" {
  name               = "crops_crop_get_function_role"
  assume_role_policy = data.aws_iam_policy_document.crops_crop_get_function_role.json
}

data "aws_iam_policy_document" "crops_crop_get_function_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_lambda_function" "products_product_get" {
  filename                       = "./data/functions/products/product_get/bootstrap.zip"
  function_name                  = "products_product_get"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "OlL5qeVJW48KzzxITY4icDPB1jydTK9nWQQRAkb344s="
  role                           = aws_iam_role.products_product_get_function_role.arn
}

resource "aws_iam_role" "products_product_get_function_role" {
  name               = "products_product_get_function_role"
  assume_role_policy = data.aws_iam_policy_document.products_product_get_function_role.json
}

data "aws_iam_policy_document" "products_product_get_function_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_lambda_function" "products_product_init" {
  filename                       = "./data/functions/products/product_init/bootstrap.zip"
  function_name                  = "products_product_init"
  runtime                        = "python3.12"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "h5nUD+5bX+VTrfw5FRxQqQzuB4FwjEFqL7QpoIu3HeQ="
  role                           = aws_iam_role.products_product_init_function_role.arn
}

resource "aws_iam_role" "products_product_init_function_role" {
  name               = "products_product_init_function_role"
  assume_role_policy = data.aws_iam_policy_document.products_product_init_function_role.json
}

data "aws_iam_policy_document" "products_product_init_function_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

