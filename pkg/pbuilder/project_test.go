package pbuilder

import (
	"fmt"
	"testing"
)

func TestProjectBuilder(t *testing.T) {

	pjname := "test_project"
	dir := getTestDataDirectory()
	data := NewData(dir)

	pb := NewProjectBuilder(data)
	pb.WithName(pjname)
	project := pb.Build()
	if project.Name != pjname {
		t.Error("project name diff")
	}

	t.Log("project: ", project)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		pb.Build()
		t.Error("Expected panic")
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	pb.New().Build()
	t.Error("Expected panic")
}

func TestProjectString(t *testing.T) {
	pjname := "test_project"
	dir := getTestDataDirectory()
	data := NewData(dir)

	pb := NewProjectBuilder(data)
	project := pb.New().WithName(pjname).Build()

	expected := "(Project: test_project, services: [(Service: testservice, functions: "
	expected += "[(Function: test_func, service: testservice, sourcePath: "
	expected += "/home/tarcisio/xday/pkg/pbuilder/testdata/services/testservice/test_func)], "
	expected += "tables: [(Name: test_table, Service: testservice, Permissions: "
	expected += "[(Function: test_func, permissions: [get, put])])])], apis: "
	expected += "[(Service: testservice, methods: [(Function: test_func, method: POST)])])"

	if project.String() != expected {
		t.Errorf("Expected %s, got %s", expected, project.String())
	}
}
