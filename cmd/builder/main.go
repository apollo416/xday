package main

import (
	"fmt"

	"github.com/apollo416/xday/pkg/pbuilder"
	"github.com/apollo416/xday/pkg/terra"
)

func main() {

	data := pbuilder.DefaultData
	pb := pbuilder.NewProjectBuilder(data)
	p := pb.New().WithName("xday").Build()

	config := terra.TerraConfig{
		Project:      p,
		DataDir:      "./data",
		S3Name:       "tf-state-hday",
		DynamoDBName: "terraform_state",
	}

	t := terra.New(config)
	t.Build()
	fmt.Println(t)
}
