package pbuilder

import (
	"fmt"
	"testing"
)

func TestTableBuilder(t *testing.T) {
	tname := "products"

	dir := getTestDataDirectory()
	data := NewData(dir)

	tb := NewTableBuilder(data).WithName(tname)

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
	table := tb.New().WithName(tname).Build()
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
