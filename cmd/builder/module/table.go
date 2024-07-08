package module

// const (
// 	TablesFile = "tables.json"
// )

// type Table struct {
// 	Name             string
// 	AllowedFunctions []Function
// }

// func (t Table) CreateTable(root *hclwrite.Body) {
// 	block := root.AppendNewBlock("resource", []string{"aws_dynamodb_table", t.Name})
// 	blockBody := block.Body()
// 	blockBody.SetAttributeValue("name", cty.StringVal(t.Name))
// 	blockBody.SetAttributeValue("billing_mode", cty.StringVal("PROVISIONED"))
// 	blockBody.SetAttributeValue("read_capacity", cty.NumberIntVal(5))
// 	blockBody.SetAttributeValue("write_capacity", cty.NumberIntVal(5))
// 	blockBody.SetAttributeValue("hash_key", cty.StringVal("id"))
// 	blockBody.AppendNewline()

// 	attribute := blockBody.AppendNewBlock("attribute", nil)
// 	attributeBody := attribute.Body()
// 	attributeBody.SetAttributeValue("name", cty.StringVal("id"))
// 	attributeBody.SetAttributeValue("type", cty.StringVal("S"))
// 	blockBody.AppendNewline()

// 	tokens := []*hclwrite.Token{}
// 	tokens = append(
// 		tokens,
// 		&hclwrite.Token{Type: hclsyntax.TokenOBrack, Bytes: []byte("[")},
// 		&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
// 	)
// 	for _, f := range t.AllowedFunctions {
// 		tokens = append(
// 			tokens,
// 			&hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte("aws_lambda_function." + f.LambdaName())},
// 			&hclwrite.Token{Type: hclsyntax.TokenComma, Bytes: []byte(",")},
// 			&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
// 		)
// 	}
// 	tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenCBrack, Bytes: []byte("]")})
// 	blockBody.SetAttributeRaw("depends_on", tokens)

// 	root.AppendNewline()
// }

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
