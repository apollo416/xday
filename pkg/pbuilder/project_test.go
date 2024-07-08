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
