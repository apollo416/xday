package terra

import (
	"errors"
	"os"

	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/apollo416/xday/pkg/terra/internal"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type TerraConfig struct {
	Project pbuilder.Project
	DataDir string
}

func (tc TerraConfig) isValid() (bool, error) {
	if _, err := os.Stat(tc.DataDir); os.IsNotExist(err) {
		return false, errors.New("data dir does not exist")
	}
	return true, nil
}

type Terra struct {
	hcl         *hclwrite.File
	TerraConfig TerraConfig
	builded     bool
}

func New(tc TerraConfig) *Terra {
	if ok, err := tc.isValid(); !ok {
		panic(err)
	}

	return &Terra{
		TerraConfig: tc,
	}
}

func (t *Terra) Build() {
	t.hcl = internal.Main(t.TerraConfig.DataDir, t.TerraConfig.Project)
}

func (t *Terra) String() string {
	return string(t.hcl.Bytes())
}
