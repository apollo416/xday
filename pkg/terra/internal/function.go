package internal

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

const (
	functionHandlerName = "bootstrap"
	dataFunctionsDir    = "functions"
)

type function struct {
	project *project
	f       pbuilder.Function
}

func newFunction(project *project, f pbuilder.Function) *function {
	return &function{
		project: project,
		f:       f,
	}
}

func (f *function) lambdaName() string {
	return f.f.Service + "_" + f.f.Name
}

func (f *function) lambdaRoleName() string {
	return f.lambdaName() + "_role"
}

func (f *function) binaryPath() string {
	return filepath.Join(
		f.project.dataDir,
		dataFunctionsDir,
		f.f.Service,
		f.f.Name,
		functionHandlerName)
}

func (f *function) zipPath() string {
	return filepath.Join(
		f.project.dataDir,
		dataFunctionsDir,
		f.f.Service,
		f.f.Name,
		functionHandlerName+".zip")
}

func (f *function) SourcePath() string {
	return f.f.SourcePath
}

func (f *function) fileBase64SHA256() string {
	return filebase64sha256(f.binaryPath())
}

func (f *function) build() {
	f.buildFunction()
	f.zip()
	f.createHCL()
	f.clean()
}

func (f *function) buildFunction() {
	cmd := exec.Command(
		"go",
		"build",
		"-o", f.binaryPath(),
		"-trimpath",
		"-buildvcs=false",
		"-ldflags=-s -w -buildid=",
		"-tags",
		"lambda.norpc",
		f.SourcePath())
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=arm64", "CGO_ENABLED=0")
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to build function: %v", string(bytes))
	}
}

func (f *function) zip() {
	archive, err := os.Create(f.zipPath())
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

	bootstrap, err := os.Open(f.binaryPath())
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

func (f *function) createHCL() {
	function := f.project.body.AppendNewBlock("resource", []string{"aws_lambda_function", f.lambdaName()})
	function.Body().SetAttributeValue("filename", cty.StringVal(f.zipPath()))
	function.Body().SetAttributeValue("function_name", cty.StringVal(f.lambdaName()))
	function.Body().SetAttributeValue("runtime", cty.StringVal("provided.al2023"))
	function.Body().SetAttributeValue("handler", cty.StringVal(functionHandlerName))
	function.Body().SetAttributeValue("timeout", cty.NumberIntVal(10))
	function.Body().SetAttributeValue("memory_size", cty.NumberIntVal(128))
	function.Body().SetAttributeValue("publish", cty.True)
	function.Body().SetAttributeValue("reserved_concurrent_executions", cty.NumberIntVal(-1))
	function.Body().SetAttributeValue("architectures", cty.ListVal([]cty.Value{cty.StringVal("arm64")}))
	function.Body().SetAttributeValue("source_code_hash", cty.StringVal(f.fileBase64SHA256()))

	function.Body().SetAttributeTraversal("role", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "aws_iam_role",
		},
		hcl.TraverseAttr{
			Name: f.lambdaRoleName(),
		},
		hcl.TraverseAttr{
			Name: "arn",
		},
	})
	f.project.body.AppendNewline()

	role := f.project.body.AppendNewBlock("resource", []string{"aws_iam_role", f.lambdaRoleName()})
	role.Body().SetAttributeValue("name", cty.StringVal(f.lambdaRoleName()))
	role.Body().SetAttributeTraversal("assume_role_policy", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "data",
		},
		hcl.TraverseAttr{
			Name: "aws_iam_policy_document",
		},
		hcl.TraverseAttr{
			Name: f.lambdaRoleName(),
		},
		hcl.TraverseAttr{
			Name: "json",
		},
	})
	f.project.body.AppendNewline()

	data := f.project.body.AppendNewBlock("data", []string{"aws_iam_policy_document", f.lambdaRoleName()})

	statement := data.Body().AppendNewBlock("statement", nil)
	statement.Body().SetAttributeValue("actions", cty.ListVal([]cty.Value{cty.StringVal("sts:AssumeRole")}))

	principals := statement.Body().AppendNewBlock("principals", nil)
	principals.Body().SetAttributeValue("type", cty.StringVal("Service"))
	principals.Body().SetAttributeValue("identifiers", cty.ListVal([]cty.Value{cty.StringVal("lambda.amazonaws.com")}))

	f.project.body.AppendNewline()
}

func (f *function) clean() {
	if err := os.Remove(f.binaryPath()); err != nil {
		log.Fatalf("Failed to remove binary: %v", err)
	}
}
