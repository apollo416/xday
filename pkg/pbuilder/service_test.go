package pbuilder

import (
	"fmt"
	"testing"
)

func TestServiceValidation(t *testing.T) {
	sname := "testservice"

	dir := getTestDataDirectory()
	data := NewData(dir)

	sref := NewServiceBuilder(data).New().WithName(sname).Build()

	s := Service{}
	if ok, _ := s.isValid(); ok {
		t.Error("expected invalid name")
	}

	s.Name = sref.Name
	if ok, _ := s.isValid(); ok {
		t.Error("expected invalid sourcePath")
	}

	s.SourcePath = sref.SourcePath
	if ok, err := s.isValid(); !ok {
		t.Errorf("Expected no errors, get: %v", err)
	}
}

func TestServiceBuilder(t *testing.T) {
	sname := "testservice"
	dir := getTestDataDirectory()
	data := NewData(dir)

	sb := NewServiceBuilder(data).WithName(sname)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		sb.Build()
		t.Error("Expected panic")
	}()

	// regular call
	s := sb.New().WithName(sname).Build()
	t.Log("service: ", s)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		sb.Build()
		t.Error("Expected panic")
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		// Build without name
		sb.New().Build()
		t.Error("Expected panic")
	}()
}

func TestServiceString(t *testing.T) {
	sname := "testservice"
	dir := getTestDataDirectory()
	data := NewData(dir)

	s := NewServiceBuilder(data).New().WithName(sname).Build()

	expected := "(Service: testservice, functions: [(Function: test_func, service: "
	expected += "testservice, sourcePath: /home/tarcisio/xday/pkg/pbuilder/testdata/services/testservice/test_func)], "
	expected += "tables: [(Name: test_table, Service: testservice, Permissions: [(Function: test_func, permissions: [get, put])])])"

	if s.String() != expected {
		t.Errorf("Expected: %s, got: %s", expected, s.String())
	}
}
