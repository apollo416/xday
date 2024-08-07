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
  runtime                        = "provided.al2023"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "B9lXIj6/AzLRDZExoG01fipM/oQVKNQy1o8gmKFkfV0="
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
  runtime                        = "provided.al2023"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "zGnaCyXOV0cbREm9SXJ8ph4fx5v6Uis/MSaAT3VnNFA="
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
}

data "aws_iam_policy_document" "crops_crops_crop_add_policy_document" {
  statement {
    effect  = "Allow"
    actions = ["dynamodb:PutItem"]

    principals {
      type        = "AWS"
      identifiers = [aws_iam_role.crops_crop_add_role.arn]
    }

    resources = [aws_dynamodb_table.crops.arn]
  }
}

resource "aws_dynamodb_resource_policy" "crops_crops_crop_add_policy" {
  resource_arn = aws_dynamodb_table.crops.arn
  policy       = data.aws_iam_policy_document.crops_crops_crop_add_policy_document.json
}

resource "aws_lambda_function" "products_product_get" {
  filename                       = "data/functions/products/product_get/bootstrap.zip"
  function_name                  = "products_product_get"
  runtime                        = "provided.al2023"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "SZPnDrxKTUnavLEmR6GYnVPx/cZyKZCxN1vZMtTapHw="
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
  runtime                        = "provided.al2023"
  handler                        = "bootstrap"
  timeout                        = 10
  memory_size                    = 128
  publish                        = true
  reserved_concurrent_executions = -1
  architectures                  = ["arm64"]
  source_code_hash               = "m1qPfmGAiVib2s2f3Kf0ptdgL7L2tXnHDDo8AwnnKaI="
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
}

resource "aws_api_gateway_rest_api" "crops" {
  name = "crops"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_resource" "crops" {
  rest_api_id = aws_api_gateway_rest_api.crops.id
  parent_id   = aws_api_gateway_rest_api.crops.root_resource_id
  path_part   = "crops"
}

resource "aws_api_gateway_rest_api" "products" {
  name = "products"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_resource" "crops" {
  rest_api_id = aws_api_gateway_rest_api.crops.id
  parent_id   = aws_api_gateway_rest_api.crops.root_resource_id
  path_part   = "crops"
}


