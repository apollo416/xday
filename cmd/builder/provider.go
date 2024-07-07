package main

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func provider(root *hclwrite.Body) {
	providerAWS := root.AppendNewBlock("provider", []string{"aws"})
	providerAWSBody := providerAWS.Body()
	providerAWSBody.SetAttributeValue("region", cty.StringVal("us-east-1"))
	providerAWSBody.AppendNewline()

	assumeRole := providerAWSBody.AppendNewBlock("assume_role", nil)
	assumeRoleBody := assumeRole.Body()
	assumeRoleBody.SetAttributeTraversal("role_arn", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "var",
		},
		hcl.TraverseAttr{
			Name: "workspace_iam_role",
		},
	})
	providerAWSBody.AppendNewline()

	defaultTags := providerAWSBody.AppendNewBlock("default_tags", nil)
	defaultTagsBody := defaultTags.Body()
	defaultTagsBody.SetAttributeValue("tags", cty.ObjectVal(map[string]cty.Value{
		"Project":            cty.StringVal("xday"),
		"TerraformWorkspace": cty.StringVal("xday"),
	}))
	root.AppendNewline()
}
