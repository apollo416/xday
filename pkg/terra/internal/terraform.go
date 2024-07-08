package internal

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func terraform(root *hclwrite.Body) {
	terraform := root.AppendNewBlock("terraform", nil)
	terraformBody := terraform.Body()
	terraformBody.SetAttributeValue("required_version", cty.StringVal("1.9.1"))
	terraformBody.AppendNewline()

	requiredProviders := terraformBody.AppendNewBlock("required_providers", nil)
	requiredProvidersBody := requiredProviders.Body()
	requiredProvidersBody.SetAttributeValue("aws", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal("hashicorp/aws"),
		"version": cty.StringVal("~> 5.0"),
	}))

	root.AppendNewline()
}
