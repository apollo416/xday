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
		services: []*service{},
		apis:     []*api{},
	}

	for _, service := range project.p.Services {
		s := newService(project, service)
		project.services = append(project.services, s)

		api := newApi(project)
		project.apis = append(project.apis, api)
	}

	return project
}

func (p *project) build() {
	for _, s := range p.services {
		s.build()
	}
}
