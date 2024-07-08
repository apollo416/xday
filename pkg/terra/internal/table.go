package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
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
	block := root.AppendNewBlock("resource", []string{"aws_dynamodb_table", t.table.Name})
	blockBody := block.Body()
	blockBody.SetAttributeValue("name", cty.StringVal(t.table.Name))
	blockBody.SetAttributeValue("billing_mode", cty.StringVal("PROVISIONED"))
	blockBody.SetAttributeValue("read_capacity", cty.NumberIntVal(5))
	blockBody.SetAttributeValue("write_capacity", cty.NumberIntVal(5))
	blockBody.SetAttributeValue("hash_key", cty.StringVal("id"))
	blockBody.AppendNewline()

	attribute := blockBody.AppendNewBlock("attribute", nil)
	attributeBody := attribute.Body()
	attributeBody.SetAttributeValue("name", cty.StringVal("id"))
	attributeBody.SetAttributeValue("type", cty.StringVal("S"))
	root.AppendNewline()

	// if len(t.functionPermissions) > 0 {
	// 	for _, permission := range t.functionPermissions {
	// 		blockPermission := root.AppendNewBlock("data", []string{"aws_iam_policy_document", t.table.Name + "_" + permission.function.LambdaName() + "_policy_document"})
	// 		blockPermissionBody := blockPermission.Body()
	// 		statement := blockPermissionBody.AppendNewBlock("statement", nil)
	// 		statementBody := statement.Body()
	// 		statementBody.SetAttributeValue("effect", cty.StringVal("Allow"))

	// 		l := []cty.Value{}
	// 		for _, g := range permission.permissions {
	// 			if g == "get" {
	// 				l = append(l, cty.StringVal("dynamodb:GetItem"))
	// 			}
	// 			if g == "put" {
	// 				l = append(l, cty.StringVal("dynamodb:PutItem"))
	// 			}
	// 		}
	// 		statementBody.SetAttributeValue("actions", cty.ListVal(l))
	// 		statementBody.AppendNewline()

	// 		principals := statementBody.AppendNewBlock("principals", nil)
	// 		principalsBody := principals.Body()
	// 		principalsBody.SetAttributeValue("type", cty.StringVal("AWS"))
	// 		principalsBody.SetAttributeRaw("identifiers", hclwrite.TokensForIdentifier("[aws_lambda_function."+permission.function.LambdaName()+".arn]"))
	// 		statementBody.AppendNewline()

	// 		statementBody.SetAttributeRaw("resources", hclwrite.TokensForIdentifier("[aws_dynamodb_table."+t.table.Name+".arn]"))
	// 		root.AppendNewline()

	// 		policy := root.AppendNewBlock("resource", []string{"aws_dynamodb_resource_policy", t.table.Name + "_" + permission.function.LambdaName() + "_policy"})
	// 		policyBody := policy.Body()
	// 		policyBody.SetAttributeTraversal("resource_arn", hcl.Traversal{
	// 			hcl.TraverseRoot{
	// 				Name: "aws_dynamodb_table",
	// 			},
	// 			hcl.TraverseAttr{
	// 				Name: t.table.Name,
	// 			},
	// 			hcl.TraverseAttr{
	// 				Name: "arn",
	// 			},
	// 		})

	// 		policyBody.SetAttributeTraversal("policy", hcl.Traversal{
	// 			hcl.TraverseRoot{
	// 				Name: "data",
	// 			},
	// 			hcl.TraverseAttr{
	// 				Name: "aws_iam_policy_document",
	// 			},
	// 			hcl.TraverseAttr{
	// 				Name: "json",
	// 			},
	// 		})
	// 	}
	// }
}

// func loadServiceTables(s Service) []Table {
// 	tables := []Table{}

// 	filename := filepath.Join(s.SourcePath(), TablesFile)
// 	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
// 		return tables
// 	}

// 	jsonFile, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer func() {
// 		if err := jsonFile.Close(); err != nil {
// 			log.Fatalf("Failed to close file: %v", err)
// 		}
// 	}()

// 	byteValue, err := io.ReadAll(jsonFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var jtables []struct {
// 		Name             string   `json:"name"`
// 		AllowedFunctions []string `json:"allowed_functions"`
// 	}
// 	json.Unmarshal(byteValue, &jtables)

// 	for _, it := range jtables {
// 		itt := Table{Name: it.Name, AllowedFunctions: []Function{}}
// 		for _, f := range it.AllowedFunctions {
// 			itt.AllowedFunctions = append(itt.AllowedFunctions, Function{Name: f, Service: s})
// 		}
// 		tables = append(tables, itt)
// 	}

// 	return tables
// }
