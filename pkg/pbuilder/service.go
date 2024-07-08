package pbuilder

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type Service struct {
	Name       string
	SourcePath string
	Functions  []Function
	Tables     []Table
}

func (s Service) isValid() (bool, error) {
	if s.Name == "" {
		return false, errors.New("service name not defined")
	}

	if _, err := os.Stat(s.SourcePath); os.IsNotExist(err) {
		return false, errors.New("sourcePath does not exist")
	}

	return true, nil
}

type serviceBuilder struct {
	data            Data
	service         Service
	functionBuilder *functionBuilder
	created         bool
	withNew         bool
}

func NewServiceBuilder(data Data) *serviceBuilder {
	return &serviceBuilder{
		data: data,
	}
}

func (sb *serviceBuilder) New() *serviceBuilder {
	sb.withNew = true
	sb.created = false
	sb.functionBuilder = NewFunctionBuilder(sb.data)
	sb.service = Service{
		Functions: []Function{},
		Tables:    []Table{},
	}
	return sb
}

func (sb *serviceBuilder) WithName(name string) *serviceBuilder {
	sb.service.Name = name
	return sb
}

func (sb *serviceBuilder) Build() Service {
	if !sb.withNew {
		panic("New() not called before")
	}
	if sb.created {
		panic("Service already created. use New() first to create a new service")
	}

	sourcePath := filepath.Join(sb.data.ServicesDir, sb.service.Name)
	sb.service.SourcePath = sourcePath

	if ok, err := sb.service.isValid(); !ok {
		panic(err)
	}

	functions := listFunctions(sb.service.SourcePath)
	for _, fname := range functions {
		f := sb.functionBuilder.New().WithService(sb.service.Name).WithName(fname).Build()
		sb.service.Functions = append(sb.service.Functions, f)
	}

	tables := listTables(sb.service.SourcePath)
	for _, tname := range tables {
		t := NewTableBuilder(sb.data).New().WithName(tname).Build()
		sb.service.Tables = append(sb.service.Tables, t)
	}

	sb.created = true
	return sb.service
}

func listServices(dir string) []string {
	services := []string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			services = append(services, e.Name())
		}
	}

	return services
}
