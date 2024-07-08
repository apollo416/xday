package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type project struct {
	p        pbuilder.Project
	dataDir  string
	body     *hclwrite.Body
	apis     []*api
	services []*service
}

func newProject(body *hclwrite.Body, dataDir string, p pbuilder.Project) *project {
	project := &project{
		p:        p,
		dataDir:  dataDir,
		body:     body,
		apis:     []*api{},
		services: []*service{},
	}

	for _, service := range project.p.Services {
		s := newService(project, service)
		project.services = append(project.services, s)
	}

	for _, ap := range project.p.Apis {
		a := newApi(project, ap)
		project.apis = append(project.apis, a)
	}

	return project
}

func (p *project) build() {
	for _, s := range p.services {
		s.build()
	}

	for _, a := range p.apis {
		a.build()
	}
}
