package pbuilder

import (
	"errors"
	"strings"
)

type Project struct {
	Name     string
	Services []Service
	Apis     []Api
}

func (p Project) isValid() (bool, error) {
	if p.Name == "" {
		return false, errors.New("project name not defined")
	}
	return true, nil
}

type projectBuilder struct {
	data    Data
	project Project
	created bool
}

func NewProjectBuilder(data Data) *projectBuilder {
	p := &projectBuilder{data: data}
	return p
}

func (pb *projectBuilder) New() *projectBuilder {
	pb.created = false
	pb.project = Project{}
	return pb
}

func (pb *projectBuilder) WithName(name string) *projectBuilder {
	pb.project.Name = name
	return pb
}

func (pb *projectBuilder) Build() Project {
	if pb.created {
		panic("Project already created. use New() first to create a new Project")
	}

	if ok, err := pb.project.isValid(); !ok {
		panic(err)
	}

	services := listServices(pb.data.ServicesDir)
	sb := NewServiceBuilder(pb.data)
	ab := NewApiBuilder(pb.data)
	for _, sname := range services {
		s := sb.New().WithName(sname).Build()
		pb.project.Services = append(pb.project.Services, s)

		a := ab.New().WithService(s).Build()
		pb.project.Apis = append(pb.project.Apis, a)
	}

	pb.created = true
	return pb.project
}

func (p Project) String() string {
	str := "(Project: " + p.Name
	str += ", services: ["
	services := []string{}
	for _, s := range p.Services {
		services = append(services, s.String())
	}
	str += strings.Join(services, ", ") + "], "
	str += "apis: ["
	apis := []string{}
	for _, a := range p.Apis {
		apis = append(apis, a.String())
	}
	str += strings.Join(apis, ", ") + "])"
	return str
}
