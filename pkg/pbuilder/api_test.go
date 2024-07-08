package pbuilder

import (
	"fmt"
	"testing"
)

func TestApiBuilder(t *testing.T) {
	sname := "testservice"

	dir := getTestDataDirectory()
	data := NewData(dir)

	sb := NewServiceBuilder(data)
	service := sb.New().WithName(sname).Build()

	ab := NewApiBuilder(data).WithService(service)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		ab.Build()
		t.Error("Expected panic")
	}()

	// regular call
	api := ab.New().WithService(service).Build()
	t.Log("api: ", api)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in api", r)
			}
		}()
		ab.Build()
		t.Error("Expected panic")
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		// Build without name
		ab.New().Build()
		t.Error("Expected panic")
	}()
}

func TestApiString(t *testing.T) {
	sname := "testservice"

	dir := getTestDataDirectory()
	data := NewData(dir)

	sb := NewServiceBuilder(data)
	service := sb.New().WithName(sname).Build()

	ab := NewApiBuilder(data)

	api := ab.New().WithService(service).Build()

	expect := "(Service: testservice, methods: [(Function: test_func, method: POST)])"

	if api.String() != expect {
		t.Errorf("Expected %s, got %s", expect, api.String())
	}
}
