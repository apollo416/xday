package internal

import (
	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func Main(dataDir string, p pbuilder.Project) *hclwrite.File {
	hclFile := hclwrite.NewEmptyFile()
	variables(hclFile.Body())
	provider(hclFile.Body())
	terraform(hclFile.Body())

	project(hclFile.Body(), dataDir, p)

	return hclFile
}
