package pbuilder

import (
	"fmt"
	"testing"
)

func TestTableBuilder(t *testing.T) {
	sname := "testservice"
	tname := "products"

	dir := getTestDataDirectory()
	data := NewData(dir)

	sb := NewServiceBuilder(data)
	service := sb.New().WithName(sname).Build()

	tb := NewTableBuilder(data).WithService(service).WithName(tname)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		tb.Build()
		t.Error("Expected panic")
	}()

	// regular call
	table := tb.New().WithService(service).WithName(tname).Build()
	t.Log("table: ", table)

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		tb.Build()
		t.Error("Expected panic")
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		// Build without name
		tb.New().Build()
		t.Error("Expected panic")
	}()
}

func TestTableFunctionPermisionString(t *testing.T) {
	sname := "testservice"
	fname := "test_func"

	dir := getTestDataDirectory()
	data := NewData(dir)

	sb := NewServiceBuilder(data)
	service := sb.New().WithName(sname).Build()

	f := NewFunctionBuilder(data).New().WithService(service.Name).WithName(fname).Build()

	tfp := TableFunctionPermision{Function: f, Permisions: []string{"get", "put"}}

	Expected := "(Function: test_func, permissions: [get, put])"
	if tfp.String() != Expected {
		t.Errorf("Expected: %s, got: %s", Expected, tfp.String())
	}
}

func TestTableString(t *testing.T) {
	sname := "testservice"
	tname := "products"

	dir := getTestDataDirectory()
	data := NewData(dir)

	sb := NewServiceBuilder(data)
	service := sb.New().WithName(sname).Build()

	table := NewTableBuilder(data).New().WithService(service).WithName(tname).Build()

	Expected := "(Name: products, Service: testservice, Permissions: [(Function: test_func, permissions: [get, put])])"
	if table.String() != Expected {
		t.Errorf("Expected: %s, got: %s", Expected, table.String())
	}
}
