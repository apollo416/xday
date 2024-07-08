package module

// type Function struct {
// }

// func (f Function) CreateFunction(root *hclwrite.Body) {
// 	f.buildFunction()
// 	f.zip()
// 	f.createHCL(root)
// 	f.clean()
// }

// func (f Function) createHCL(root *hclwrite.Body) {
// 	function := root.AppendNewBlock("resource", []string{"aws_lambda_function", f.LambdaName()})
// 	functionBody := function.Body()
// 	functionBody.SetAttributeValue("filename", cty.StringVal(f.ZipPath()))
// 	functionBody.SetAttributeValue("function_name", cty.StringVal(f.LambdaName()))
// 	functionBody.SetAttributeValue("runtime", cty.StringVal("python3.12"))
// 	functionBody.SetAttributeValue("handler", cty.StringVal(functionHandlerName))
// 	functionBody.SetAttributeValue("timeout", cty.NumberIntVal(10))
// 	functionBody.SetAttributeValue("memory_size", cty.NumberIntVal(128))
// 	functionBody.SetAttributeValue("publish", cty.True)
// 	functionBody.SetAttributeValue("reserved_concurrent_executions", cty.NumberIntVal(-1))
// 	functionBody.SetAttributeValue("architectures", cty.ListVal([]cty.Value{cty.StringVal("arm64")}))
// 	functionBody.SetAttributeValue("source_code_hash", cty.StringVal(f.fileBase64SHA256()))

// 	functionBody.SetAttributeTraversal("role", hcl.Traversal{
// 		hcl.TraverseRoot{
// 			Name: "aws_iam_role",
// 		},
// 		hcl.TraverseAttr{
// 			Name: f.LambdaRoleName(),
// 		},
// 		hcl.TraverseAttr{
// 			Name: "arn",
// 		},
// 	})
// 	root.AppendNewline()

// 	role := root.AppendNewBlock("resource", []string{"aws_iam_role", f.LambdaRoleName()})
// 	roleBody := role.Body()
// 	roleBody.SetAttributeValue("name", cty.StringVal(f.LambdaRoleName()))
// 	roleBody.SetAttributeTraversal("assume_role_policy", hcl.Traversal{
// 		hcl.TraverseRoot{
// 			Name: "data",
// 		},
// 		hcl.TraverseAttr{
// 			Name: "aws_iam_policy_document",
// 		},
// 		hcl.TraverseAttr{
// 			Name: f.LambdaRoleName(),
// 		},
// 		hcl.TraverseAttr{
// 			Name: "json",
// 		},
// 	})
// 	root.AppendNewline()

// 	data := root.AppendNewBlock("data", []string{"aws_iam_policy_document", f.LambdaRoleName()})
// 	dataBody := data.Body()

// 	statement := dataBody.AppendNewBlock("statement", nil)
// 	statementBody := statement.Body()
// 	statementBody.SetAttributeValue("actions", cty.ListVal([]cty.Value{cty.StringVal("sts:AssumeRole")}))

// 	principals := statementBody.AppendNewBlock("principals", nil)
// 	principalsBody := principals.Body()
// 	principalsBody.SetAttributeValue("type", cty.StringVal("Service"))
// 	principalsBody.SetAttributeValue("identifiers", cty.ListVal([]cty.Value{cty.StringVal("lambda.amazonaws.com")}))

// 	root.AppendNewline()
// }

// func (f Function) buildFunction() {
// 	cmd := exec.Command("go", "build", "-o", f.BinaryPath(), f.SourcePath())
// 	bytes, err := cmd.CombinedOutput()
// 	if err != nil {
// 		log.Fatalf("Failed to build function: %v", string(bytes))
// 	}
// }

// func (f Function) zip() {
// 	archive, err := os.Create(f.ZipPath())
// 	if err != nil {
// 		log.Fatalf("Failed to create archive: %v", err)
// 	}

// 	defer func() {
// 		if err := archive.Close(); err != nil {
// 			log.Fatalf("Failed to close archive: %v", err)
// 		}
// 	}()

// 	zipWriter := zip.NewWriter(archive)
// 	defer func() {
// 		if err := zipWriter.Close(); err != nil {
// 			log.Fatalf("Failed to close zipWriter: %v", err)
// 		}
// 	}()

// 	bootstrap, err := os.Open(f.BinaryPath())
// 	if err != nil {
// 		log.Fatalf("Failed to open bootstrap: %v", err)
// 	}

// 	defer func() {
// 		if err := bootstrap.Close(); err != nil {
// 			log.Fatalf("Failed to close bootstrap: %v", err)
// 		}
// 	}()

// 	bootstrapZip, err := zipWriter.Create(functionHandlerName)
// 	if err != nil {
// 		log.Fatalf("Failed to create bootstrapZip: %v", err)
// 	}

// 	if _, err := io.Copy(bootstrapZip, bootstrap); err != nil {
// 		log.Fatalf("Failed to copy content: %v", err)
// 	}
// }

// func (f Function) clean() {
// 	if err := os.Remove(f.BinaryPath()); err != nil {
// 		log.Fatalf("Failed to remove binary: %v", err)
// 	}
// }

// func (f Function) fileBase64SHA256() string {
// 	return filebase64sha256(f)
// }

// func loadServiceFunctions(s Service) []Function {
// 	var functions []Function

// 	entries, err := os.ReadDir(s.SourcePath())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, e := range entries {
// 		if e.IsDir() {
// 			functions = append(functions, Function{Name: e.Name(), Service: s})
// 		}
// 	}

// 	return functions
// }

// func filebase64sha256(f Function) string {
// 	arch, err := os.Open(f.BinaryPath())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		if err := arch.Close(); err != nil {
// 			log.Fatalf("Failed to close file: %v", err)
// 		}
// 	}()

// 	h := sha256.New()
// 	if _, err := io.Copy(h, arch); err != nil {
// 		log.Fatalf("Failed to copy file content to the hash: %v", err)
// 	}

// 	return base64.StdEncoding.EncodeToString(h.Sum(nil))
// }
