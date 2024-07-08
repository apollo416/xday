package pbuilder

import (
	"path/filepath"
	"runtime"
)

const (
	servicesDir     = "./services"
	functionDataDir = "./data/functions"
)

type Data struct {
	DataDir         string
	ServicesDir     string
	FunctionDataDir string
}

func NewData(basePath string) Data {
	return Data{
		ServicesDir:     filepath.Join(basePath, servicesDir),
		FunctionDataDir: filepath.Join(basePath, functionDataDir),
	}
}

var DefaultData Data

func init() {
	_, filename, _, _ := runtime.Caller(0)
	projectData := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	DefaultData = NewData(projectData)
}
