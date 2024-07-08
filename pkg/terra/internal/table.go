package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type table struct {
	t                        pbuilder.Table
	tableName                string
	service                  *service
	tableFunctionPermissions []*tableFunctionPermission
}

func newTable(service *service, t pbuilder.Table) *table {
	table := &table{
		t:                        t,
		tableName:                t.Name,
		service:                  service,
		tableFunctionPermissions: []*tableFunctionPermission{},
	}

	for _, permission := range t.Permissions {
		ok, function := service.getFunctionByName(permission.Function.Name)
		if ok {
			p := newTableFunctionPermission(table, function, permission.Permisions)
			table.tableFunctionPermissions = append(table.tableFunctionPermissions, p)
		}
	}
	return table
}

func (t *table) name() string {
	return t.tableName
}

func (t *table) tableResourceBlockPath() []string {
	return []string{"aws_dynamodb_table", t.name()}
}

func (t *table) resourceIdentifier() string {
	return "[aws_dynamodb_table." + t.name() + ".arn]"
}

func (t *table) build() {
	t.buildHCL()
	for _, permission := range t.tableFunctionPermissions {
		if len(permission.permissions) > 0 {
			permission.buildHCL()
		}
	}
}

func (t *table) buildHCL() {
	table := t.service.project.body.AppendNewBlock("resource", t.tableResourceBlockPath())
	table.Body().SetAttributeValue("name", cty.StringVal(t.name()))
	table.Body().SetAttributeValue("billing_mode", cty.StringVal("PROVISIONED"))
	table.Body().SetAttributeValue("read_capacity", cty.NumberIntVal(5))
	table.Body().SetAttributeValue("write_capacity", cty.NumberIntVal(5))
	table.Body().SetAttributeValue("hash_key", cty.StringVal("id"))
	table.Body().AppendNewline()

	attribute := table.Body().AppendNewBlock("attribute", nil)
	attribute.Body().SetAttributeValue("name", cty.StringVal("id"))
	attribute.Body().SetAttributeValue("type", cty.StringVal("S"))

	t.service.project.body.AppendNewline()
}

type tableFunctionPermission struct {
	table       *table
	function    *function
	permissions []string
}

func newTableFunctionPermission(t *table, f *function, permissions []string) *tableFunctionPermission {
	p := &tableFunctionPermission{
		table:       t,
		function:    f,
		permissions: permissions,
	}
	return p
}

func (t *tableFunctionPermission) buildHCL() {
	policyDocumentId := t.policyDocumentPath()
	blockPermission := t.table.service.project.body.AppendNewBlock("data", policyDocumentId)
	statement := blockPermission.Body().AppendNewBlock("statement", nil)
	statement.Body().SetAttributeValue("effect", cty.StringVal("Allow"))

	l := []cty.Value{}
	for _, g := range t.permissions {
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
	principals.Body().SetAttributeRaw("identifiers", hclwrite.TokensForIdentifier(t.principalIdentifier()))
	statement.Body().AppendNewline()

	statement.Body().SetAttributeRaw("resources", hclwrite.TokensForIdentifier(t.table.resourceIdentifier()))
	t.table.service.project.body.AppendNewline()

	policy := t.table.service.project.body.AppendNewBlock("resource", t.policyResourceBlockPath())
	policy.Body().SetAttributeTraversal("resource_arn", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_dynamodb_table",
		},
		hcl.TraverseAttr{
			Name: t.table.tableName,
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
			Name: t.policyDocumentName(),
		},
		hcl.TraverseAttr{
			Name: "json",
		},
	})
	t.table.service.project.body.AppendNewline()
}

func (t *tableFunctionPermission) policyDocumentName() string {
	return t.table.tableName + "_" + t.function.lambdaName() + "_policy_document"
}

func (t *tableFunctionPermission) policyDocumentPath() []string {
	return []string{"aws_iam_policy_document", t.policyDocumentName()}
}

func (f *tableFunctionPermission) principalIdentifier() string {
	return "[aws_iam_role." + f.function.lambdaName() + "_role.arn]"
}

func (f *tableFunctionPermission) policyName() string {
	return f.table.tableName + "_" + f.function.lambdaName() + "_policy"
}

func (f *tableFunctionPermission) policyResourceBlockPath() []string {
	return []string{"aws_dynamodb_resource_policy", f.policyName()}
}
