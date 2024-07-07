package main

import (
	"fmt"
	"log"
	"os"

	"github.com/apollo416/xday/cmd/builder/module"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func main() {
	tfFile, err := os.Create("terraform.tf")
	defer func() {
		if err := tfFile.Close(); err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("Error creating the file: %v", err)
		return
	}

	hclFile := hclwrite.NewEmptyFile()
	rootBody := hclFile.Body()

	rootVariables(rootBody)
	provider(rootBody)
	terraform(rootBody)

	services := module.ListServices()
	for _, service := range services {
		fmt.Println("Service: ", service)
	}

	tables := module.ListTables(services)
	fmt.Println("Tables: ", tables)
	for _, t := range tables {
		t.CreateTable(rootBody)
	}

	functions := module.ListFunctions(services)
	fmt.Println("Functions: ", functions)
	for _, f := range functions {
		f.CreateFunction(rootBody)
	}

	tfFile.Write(hclFile.Bytes())
}
