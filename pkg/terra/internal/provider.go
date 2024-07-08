package internal

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func provider(root *hclwrite.Body, project string) {
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

	tokens := []*hclwrite.Token{}
	tokens = append(
		tokens,
		&hclwrite.Token{Type: hclsyntax.TokenOBrace, Bytes: []byte("{")},
		&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	)
	tokens = append(
		tokens,
		&hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte("Project")},
		&hclwrite.Token{Type: hclsyntax.TokenEqual, Bytes: []byte("=")},
		&hclwrite.Token{Type: hclsyntax.TokenOQuote, Bytes: []byte("\"")},
		&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(project)},
		&hclwrite.Token{Type: hclsyntax.TokenCQuote, Bytes: []byte("\"")},
		&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},

		&hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte("TerraformWorkspace")},
		&hclwrite.Token{Type: hclsyntax.TokenEqual, Bytes: []byte("=")},
		&hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte("terraform.workspace")},
		&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	)

	tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCBrace, Bytes: []byte("}")})

	defaultTagsBody.SetAttributeRaw("tags", tokens)
	root.AppendNewline()
}

// backend "s3" {
//     bucket = "mybucket"
//     key    = "path/to/my/key"
//     region = "us-east-1"
//   }
