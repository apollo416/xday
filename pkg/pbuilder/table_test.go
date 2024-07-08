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
