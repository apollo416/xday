package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
)

type service struct {
	s         pbuilder.Service
	project   *project
	functions []*function
	tables    []*table
	api       *api
	funcRefs  map[string]*function
}

func newService(project *project, s pbuilder.Service) *service {
	service := &service{
		s:         s,
		project:   project,
		functions: []*function{},
		tables:    []*table{},
		funcRefs:  map[string]*function{},
	}

	for _, function := range s.Functions {
		f := newFunction(project, function)
		service.functions = append(service.functions, f)
		service.funcRefs[f.f.Name] = f
	}

	for _, table := range s.Tables {
		t := newTable(service, table)
		service.tables = append(service.tables, t)
	}

	return service
}

func (s *service) build() {
	s.buildFunctions()
	s.buildTables()
}

func (s *service) buildFunctions() {
	for _, function := range s.functions {
		function.build()
	}
}

func (s *service) buildTables() {
	for _, table := range s.tables {
		table.build()
	}
}

func (s *service) buildApi() {
	s.api.build()
}

func (s *service) getFunctionByName(name string) (bool, *function) {
	if function, ok := s.funcRefs[name]; ok {
		return true, function
	}

	return false, nil
}
