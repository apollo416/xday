package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type api struct {
	project *project
	methods []*apiMethod
	a       pbuilder.Api
}

func newApi(project *project, a pbuilder.Api) *api {
	ap := &api{
		project: project,
		methods: []*apiMethod{},
		a:       a,
	}

	for _, method := range a.Methods {
		m := &apiMethod{
			api:      ap,
			function: newFunction(project, method.Function),
		}

		ap.methods = append(ap.methods, m)
	}

	return ap
}

func (a *api) build() {
	ap := a.project.body.AppendNewBlock("resource", []string{"aws_api_gateway_rest_api", a.a.Service.Name})
	ap.Body().SetAttributeValue("name", cty.StringVal(a.a.Service.Name))
	ap.Body().AppendNewline()

	endpointConfig := ap.Body().AppendNewBlock("endpoint_configuration", nil)
	endpointConfig.Body().SetAttributeValue("types", cty.ListVal([]cty.Value{cty.StringVal("REGIONAL")}))
	ap.Body().AppendNewline()

	lifecycle := ap.Body().AppendNewBlock("lifecycle", nil)
	lifecycle.Body().SetAttributeValue("create_before_destroy", cty.True)

	a.project.body.AppendNewline()

	// TODO: fazer direito
	resource := a.project.body.AppendNewBlock("resource", []string{"aws_api_gateway_resource", "crops"})
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
	a.project.body.AppendNewline()
}

type apiMethod struct {
	api      *api
	function *function
}

func (m *apiMethod) build() {
	method := m.api.project.body.AppendNewBlock("resource", []string{"aws_api_gateway_method", "crops_post"})
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
	m.api.project.body.AppendNewline()

	integration := m.api.project.body.AppendNewBlock("resource", []string{"aws_api_gateway_integration", "crops_post"})
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
	m.api.project.body.AppendNewline()

	permission := m.api.project.body.AppendNewBlock("resource", []string{"aws_lambda_permission", "crops_post"})
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
	m.api.project.body.AppendNewline()
}
