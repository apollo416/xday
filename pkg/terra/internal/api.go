package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func createApi(root *hclwrite.Body, p pbuilder.Project) {
	for _, service := range p.Services {
		api := root.AppendNewBlock("resource", []string{"aws_api_gateway_rest_api", service.Name})
		api.Body().SetAttributeValue("name", cty.StringVal(service.Name))
		api.Body().AppendNewline()

		endpointConfig := api.Body().AppendNewBlock("endpoint_configuration", nil)
		endpointConfig.Body().SetAttributeValue("types", cty.ListVal([]cty.Value{cty.StringVal("REGIONAL")}))
		api.Body().AppendNewline()

		lifecycle := api.Body().AppendNewBlock("lifecycle", nil)
		lifecycle.Body().SetAttributeValue("create_before_destroy", cty.True)

		root.AppendNewline()
	}

	// TODO: fazer direito
	resource := root.AppendNewBlock("resource", []string{"aws_api_gateway_resource", "crops"})
	resource.Body().SetAttributeTraversal("rest_api_id", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_rest_api",
		},
		hcl.TraverseAttr{
			Name: "crops",
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	})

	resource.Body().SetAttributeTraversal("parent_id", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_rest_api",
		},
		hcl.TraverseAttr{
			Name: "crops",
		},
		hcl.TraverseAttr{
			Name: "root_resource_id",
		},
	})
	resource.Body().SetAttributeValue("path_part", cty.StringVal("crops"))
	root.AppendNewline()

	method := root.AppendNewBlock("resource", []string{"aws_api_gateway_method", "crops_post"})
	method.Body().SetAttributeTraversal("rest_api_id", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_rest_api",
		},
		hcl.TraverseAttr{
			Name: "crops",
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	})
	method.Body().SetAttributeTraversal("resource_id", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_resource",
		},
		hcl.TraverseAttr{
			Name: "crops",
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	})
	method.Body().SetAttributeValue("http_method", cty.StringVal("POST"))
	method.Body().SetAttributeValue("authorization", cty.StringVal("NONE"))
	root.AppendNewline()

	integration := root.AppendNewBlock("resource", []string{"aws_api_gateway_integration", "crops_post"})
	integration.Body().SetAttributeTraversal("rest_api_id", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_rest_api",
		},
		hcl.TraverseAttr{
			Name: "crops",
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	})
	integration.Body().SetAttributeTraversal("resource_id", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_resource",
		},
		hcl.TraverseAttr{
			Name: "crops",
		},
		hcl.TraverseAttr{
			Name: "id",
		},
	})
	integration.Body().SetAttributeTraversal("http_method", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_api_gateway_method",
		},
		hcl.TraverseAttr{
			Name: "crops_post",
		},
		hcl.TraverseAttr{
			Name: "http_method",
		},
	})
	integration.Body().SetAttributeValue("integration_http_method", cty.StringVal("POST"))
	integration.Body().SetAttributeValue("type", cty.StringVal("AWS_PROXY"))
	integration.Body().SetAttributeTraversal("uri", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_lambda_function",
		},
		hcl.TraverseAttr{
			Name: "crops_crop_add",
		},
		hcl.TraverseAttr{
			Name: "invoke_arn",
		},
	})
	root.AppendNewline()

	permission := root.AppendNewBlock("resource", []string{"aws_lambda_permission", "crops_post"})
	permission.Body().SetAttributeValue("statement_id", cty.StringVal("AllowExecution"+"_crops_post_"+"FromAPI"))
	permission.Body().SetAttributeValue("action", cty.StringVal("lambda:InvokeFunction"))
	permission.Body().SetAttributeTraversal("function_name", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_lambda_function",
		},
		hcl.TraverseAttr{
			Name: "crops_crop_add",
		},
		hcl.TraverseAttr{
			Name: "function_name",
		},
	})
	permission.Body().SetAttributeValue("principal", cty.StringVal("apigateway.amazonaws.com"))
	lambdaIdentifier := "[aws_lambda_function.crops_crop_add]"
	permission.Body().SetAttributeRaw("depends_on", hclwrite.TokensForIdentifier(lambdaIdentifier))
	root.AppendNewline()
}

// resource "aws_lambda_permission" "allow_api" {
// 	# checkov:skip=CKV_AWS_364:Ensure that AWS Lambda function permissions delegated to AWS services are limited by SourceArn or SourceAccount
// 	action        = "lambda:InvokeFunction"
// 	function_name = aws_lambda_function.this.function_name
// 	principal     = "apigateway.amazonaws.com"
// 	#source_arn    = var.api
// 	depends_on = [aws_lambda_function.this]
//   }
