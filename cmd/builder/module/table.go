package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const (
	TablesFile = "tables.json"
)

type Table struct {
	Name             string
	AllowedFunctions []Function
}

func ListTables(services []Service) []Table {

	var tables []Table

	for _, service := range services {
		filename := filepath.Join(service.SourcePath(), TablesFile)

		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			continue
		}

		jsonFile, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := jsonFile.Close(); err != nil {
				log.Fatalf("Failed to close file: %v", err)
			}
		}()

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatal(err)
		}
		var inner_tables []struct {
			Name             string   `json:"name"`
			AllowedFunctions []string `json:"allowed_functions"`
		}
		json.Unmarshal(byteValue, &inner_tables)
		fmt.Println(inner_tables)

		for _, it := range inner_tables {
			itt := Table{Name: it.Name, AllowedFunctions: []Function{}}
			for _, f := range it.AllowedFunctions {
				itt.AllowedFunctions = append(itt.AllowedFunctions, Function{Name: f, Service: service})
			}
			tables = append(tables, itt)
		}
	}

	return tables
}

func (t Table) CreateTable(root *hclwrite.Body) {
	block := root.AppendNewBlock("resource", []string{"aws_dynamodb_table", t.Name})
	blockBody := block.Body()
	blockBody.SetAttributeValue("name", cty.StringVal(t.Name))
	blockBody.SetAttributeValue("billing_mode", cty.StringVal("PROVISIONED"))
	blockBody.SetAttributeValue("read_capacity", cty.NumberIntVal(5))
	blockBody.SetAttributeValue("write_capacity", cty.NumberIntVal(5))
	blockBody.SetAttributeValue("hash_key", cty.StringVal("id"))
	blockBody.AppendNewline()

	attribute := blockBody.AppendNewBlock("attribute", nil)
	attributeBody := attribute.Body()
	attributeBody.SetAttributeValue("name", cty.StringVal("id"))
	attributeBody.SetAttributeValue("type", cty.StringVal("S"))
	blockBody.AppendNewline()

	tokens := []*hclwrite.Token{}
	tokens = append(
		tokens,
		&hclwrite.Token{Type: hclsyntax.TokenOBrack, Bytes: []byte("[")},
		&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	)
	for _, f := range t.AllowedFunctions {
		tokens = append(
			tokens,
			&hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte("aws_lambda_function." + f.LambdaName())},
			&hclwrite.Token{Type: hclsyntax.TokenComma, Bytes: []byte(",")},
			&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
		)
	}
	tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCBrack, Bytes: []byte("]")})
	blockBody.SetAttributeRaw("depends_on", tokens)

	root.AppendNewline()
}
