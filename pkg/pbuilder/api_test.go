package pbuilder

import (
	"fmt"
	"testing"
)

func TestApiBuilder(t *testing.T) {
	sname := "testservice"
	fname := "test_func"

	dir := getTestDataDirectory()
	data := NewData(dir)

	sb := NewServiceBuilder(data)
	service := sb.New().WithName(sname).Build()

	ab := NewApiBuilder(data).WithService(service).WithName(fname)

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
	api := ab.New().WithService(service).WithName(fname).Build()
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
