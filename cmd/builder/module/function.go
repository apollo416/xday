package module

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const (
	RoleSufix = "_function_role"
)

const (
	functionHandlerName = "bootstrap"
)

type Function struct {
	Name    string
	Service Service
}

func (f Function) LambdaName() string {
	return f.Service.Name + "_" + f.Name
}

func (f Function) LambdaRoleName() string {
	return f.LambdaName() + RoleSufix
}

func (f Function) LambdaRolePolicy() string {
	return f.LambdaRoleName() + "_policy"
}

// TODO: Verificar a necessidade de especificar um glob para pegar todos os arquivos do diretório
// Atualmente retorna apenas o diretório
func (f Function) SourcePath() string {
	return "." + string(filepath.Separator) + filepath.Join(f.Service.SourcePath(), f.Name)
}

func (f Function) BinaryPath() string {
	return "." + string(filepath.Separator) + filepath.Join(FunctionDataDir, f.Service.Name, f.Name, functionHandlerName)
}

func (f Function) ZipPath() string {
	return f.BinaryPath() + ".zip"
}

func (f Function) CreateFunction(root *hclwrite.Body) {
	f.buildFunction()
	f.zip()
	f.createHCL(root)
	f.clean()
}

func (f Function) createHCL(root *hclwrite.Body) {
	function := root.AppendNewBlock("resource", []string{"aws_lambda_function", f.LambdaName()})
	functionBody := function.Body()
	functionBody.SetAttributeValue("filename", cty.StringVal(f.ZipPath()))
	functionBody.SetAttributeValue("function_name", cty.StringVal(f.LambdaName()))
	functionBody.SetAttributeValue("runtime", cty.StringVal("python3.12"))
	functionBody.SetAttributeValue("handler", cty.StringVal(functionHandlerName))
	functionBody.SetAttributeValue("timeout", cty.NumberIntVal(10))
	functionBody.SetAttributeValue("memory_size", cty.NumberIntVal(128))
	functionBody.SetAttributeValue("publish", cty.True)
	functionBody.SetAttributeValue("reserved_concurrent_executions", cty.NumberIntVal(-1))
	functionBody.SetAttributeValue("architectures", cty.ListVal([]cty.Value{cty.StringVal("arm64")}))
	functionBody.SetAttributeValue("source_code_hash", cty.StringVal(f.fileBase64SHA256()))

	functionBody.SetAttributeTraversal("role", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_iam_role",
		},
		hcl.TraverseAttr{
			Name: f.LambdaRoleName(),
		},
		hcl.TraverseAttr{
			Name: "arn",
		},
	})
	root.AppendNewline()

	role := root.AppendNewBlock("resource", []string{"aws_iam_role", f.LambdaRoleName()})
	roleBody := role.Body()
	roleBody.SetAttributeValue("name", cty.StringVal(f.LambdaRoleName()))
	roleBody.SetAttributeTraversal("assume_role_policy", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "data",
		},
		hcl.TraverseAttr{
			Name: "aws_iam_policy_document",
		},
		hcl.TraverseAttr{
			Name: f.LambdaRoleName(),
		},
		hcl.TraverseAttr{
			Name: "json",
		},
	})
	root.AppendNewline()

	data := root.AppendNewBlock("data", []string{"aws_iam_policy_document", f.LambdaRoleName()})
	dataBody := data.Body()

	statement := dataBody.AppendNewBlock("statement", nil)
	statementBody := statement.Body()
	statementBody.SetAttributeValue("actions", cty.ListVal([]cty.Value{cty.StringVal("sts:AssumeRole")}))

	principals := statementBody.AppendNewBlock("principals", nil)
	principalsBody := principals.Body()
	principalsBody.SetAttributeValue("type", cty.StringVal("Service"))
	principalsBody.SetAttributeValue("identifiers", cty.ListVal([]cty.Value{cty.StringVal("lambda.amazonaws.com")}))

	root.AppendNewline()
}

func (f Function) buildFunction() {
	cmd := exec.Command("go", "build", "-o", f.BinaryPath(), f.SourcePath())
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to build function: %v", string(bytes))
	}
}

func (f Function) zip() {
	archive, err := os.Create(f.ZipPath())
	if err != nil {
		log.Fatalf("Failed to create archive: %v", err)
	}

	defer func() {
		if err := archive.Close(); err != nil {
			log.Fatalf("Failed to close archive: %v", err)
		}
	}()

	zipWriter := zip.NewWriter(archive)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			log.Fatalf("Failed to close zipWriter: %v", err)
		}
	}()

	bootstrap, err := os.Open(f.BinaryPath())
	if err != nil {
		log.Fatalf("Failed to open bootstrap: %v", err)
	}

	defer func() {
		if err := bootstrap.Close(); err != nil {
			log.Fatalf("Failed to close bootstrap: %v", err)
		}
	}()

	bootstrapZip, err := zipWriter.Create(functionHandlerName)
	if err != nil {
		log.Fatalf("Failed to create bootstrapZip: %v", err)
	}

	if _, err := io.Copy(bootstrapZip, bootstrap); err != nil {
		log.Fatalf("Failed to copy content: %v", err)
	}
}

func (f Function) clean() {
	if err := os.Remove(f.BinaryPath()); err != nil {
		log.Fatalf("Failed to remove binary: %v", err)
	}
}

func (f Function) fileBase64SHA256() string {
	return filebase64sha256(f)
}

func ListFunctions(services []Service) []Function {
	var functions []Function

	for _, service := range services {
		entries, err := os.ReadDir(service.SourcePath())
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			if e.IsDir() {
				functions = append(functions, Function{Name: e.Name(), Service: service})
			}
		}
	}

	return functions
}

func filebase64sha256(f Function) string {
	arch, err := os.Open(f.BinaryPath())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := arch.Close(); err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, arch); err != nil {
		log.Fatalf("Failed to copy file content to the hash: %v", err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
