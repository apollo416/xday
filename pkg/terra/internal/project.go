package internal

import (
	"fmt"

	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func project(root *hclwrite.Body, dataDir string, p pbuilder.Project) {
	for _, service := range p.Services {
		for _, function := range service.Functions {
			f := newFunction(dataDir, function)
			f.createFunction(root)
		}
		for _, table := range service.Tables {
			fmt.Println(table)
		}
	}
}
