package internal

import (
	"fmt"

	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func Main(dataDir string, projectName, s3Name, dynamodbName string, p pbuilder.Project) *hclwrite.File {
	hclFile := hclwrite.NewEmptyFile()
	variables(hclFile.Body())
	provider(hclFile.Body(), projectName)
	terraform(hclFile.Body(), s3Name, dynamodbName)

	fmt.Println(p)
	panic("hu")
	project := newProject(hclFile.Body(), dataDir, p)

	project.build()

	return hclFile
}
