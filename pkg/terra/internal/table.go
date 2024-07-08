package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type functionPermission struct {
	function    function
	permissions []string
}

type table struct {
	table               pbuilder.Table
	functionPermissions []functionPermission
}

func newTable(dataDir string, t pbuilder.Table) table {
	functionpermissions := []functionPermission{}
	for _, permission := range t.Permissions {
		p := functionPermission{
			function:    newFunction(dataDir, permission.Function),
			permissions: permission.Permisions,
		}
		functionpermissions = append(functionpermissions, p)
	}
	return table{
		table:               t,
		functionPermissions: functionpermissions,
	}
}

func (t table) createTable(root *hclwrite.Body) {
	table := root.AppendNewBlock("resource", []string{"aws_dynamodb_table", t.table.Name})
	table.Body().SetAttributeValue("name", cty.StringVal(t.table.Name))
	table.Body().SetAttributeValue("billing_mode", cty.StringVal("PROVISIONED"))
	table.Body().SetAttributeValue("read_capacity", cty.NumberIntVal(5))
	table.Body().SetAttributeValue("write_capacity", cty.NumberIntVal(5))
	table.Body().SetAttributeValue("hash_key", cty.StringVal("id"))
	table.Body().AppendNewline()

	attribute := table.Body().AppendNewBlock("attribute", nil)
	attribute.Body().SetAttributeValue("name", cty.StringVal("id"))
	attribute.Body().SetAttributeValue("type", cty.StringVal("S"))

	root.AppendNewline()

	if len(t.functionPermissions) > 0 {
		for _, permission := range t.functionPermissions {
			policyDocumentName := t.table.Name + "_" + permission.function.LambdaName() + "_policy_document"
			policyDocumentId := []string{"aws_iam_policy_document", policyDocumentName}
			blockPermission := root.AppendNewBlock("data", policyDocumentId)
			statement := blockPermission.Body().AppendNewBlock("statement", nil)
			statement.Body().SetAttributeValue("effect", cty.StringVal("Allow"))

			l := []cty.Value{}
			for _, g := range permission.permissions {
				if g == "get" {
					l = append(l, cty.StringVal("dynamodb:GetItem"))
				}
				if g == "put" {
					l = append(l, cty.StringVal("dynamodb:PutItem"))
				}
			}
			statement.Body().SetAttributeValue("actions", cty.ListVal(l))
			statement.Body().AppendNewline()

			principals := statement.Body().AppendNewBlock("principals", nil)
			principals.Body().SetAttributeValue("type", cty.StringVal("AWS"))
			principalIdentifier := "[aws_iam_role." + permission.function.LambdaName() + "_role.arn]"
			principals.Body().SetAttributeRaw("identifiers", hclwrite.TokensForIdentifier(principalIdentifier))
			statement.Body().AppendNewline()

			resourceIdentifier := "[aws_dynamodb_table." + t.table.Name + ".arn]"
			statement.Body().SetAttributeRaw("resources", hclwrite.TokensForIdentifier(resourceIdentifier))
			root.AppendNewline()

			policyName := t.table.Name + "_" + permission.function.LambdaName() + "_policy"
			policyNameId := []string{"aws_dynamodb_resource_policy", policyName}
			policy := root.AppendNewBlock("resource", policyNameId)
			policy.Body().SetAttributeTraversal("resource_arn", hcl.Traversal{
				hcl.TraverseRoot{
					Name: "aws_dynamodb_table",
				},
				hcl.TraverseAttr{
					Name: t.table.Name,
				},
				hcl.TraverseAttr{
					Name: "arn",
				},
			})

			policy.Body().SetAttributeTraversal("policy", hcl.Traversal{
				hcl.TraverseRoot{
					Name: "data",
				},
				hcl.TraverseAttr{
					Name: "aws_iam_policy_document",
				},
				hcl.TraverseAttr{
					Name: policyDocumentName,
				},
				hcl.TraverseAttr{
					Name: "json",
				},
			})

		}
		root.AppendNewline()
	}
}
