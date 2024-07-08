package pbuilder

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type Function struct {
	Name       string
	Service    string
	SourcePath string
}

func (f Function) isValid() (bool, error) {
	if f.Name == "" {
		return false, errors.New("function name is required")
	}

	if f.Service == "" {
		return false, errors.New("service name is required")
	}

	if f.SourcePath == "" {
		return false, errors.New("source path is required")
	}

	if _, err := os.Stat(f.SourcePath); os.IsNotExist(err) {
		return false, errors.New("sourcePath does not exist")
	}

	return true, nil
}

type functionBuilder struct {
	data     Data
	function Function
	created  bool
}

func NewFunctionBuilder(data Data) *functionBuilder {
	return &functionBuilder{data: data}
}

func (fb *functionBuilder) New() *functionBuilder {
	fb.created = false
	fb.function = Function{}
	return fb
}

func (fb *functionBuilder) WithName(fname string) *functionBuilder {
	fb.function.Name = fname
	return fb
}

func (fb *functionBuilder) WithService(sname string) *functionBuilder {
	fb.function.Service = sname
	return fb
}

func (fb *functionBuilder) Build() Function {
	if fb.created {
		panic("Function already created. use New() first to create a new Function")
	}

	fb.function.SourcePath = filepath.Join(fb.data.ServicesDir, fb.function.Service, fb.function.Name)

	if ok, err := fb.function.isValid(); !ok {
		panic(err)
	}

	fb.created = true
	return fb.function
}

func listFunctions(dir string) []string {
	functions := []string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			functions = append(functions, e.Name())
		}
	}

	return functions
}

// TODO: Verificar a necessidade de especificar um glob para pegar todos os arquivos do diretório
// Atualmente retorna apenas o diretório
// func (f Function) SourcePath() string {
// 	return "." + string(filepath.Separator) + filepath.Join(f.Service.SourcePath(), f.Name)
// }

// func (f Function) BinaryPath() string {
// 	return "." + string(filepath.Separator) + filepath.Join(mdata.FunctionDataDir, f.Service.Name, f.Name, functionHandlerName)
// }

// func (f Function) ZipPath() string {
// 	return f.BinaryPath() + ".zip"
// }
