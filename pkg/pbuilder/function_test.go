package pbuilder

import (
	"fmt"
	"testing"
)

func TestFunctionValidation(t *testing.T) {
	fname := "test_func"
	sname := "testservice"

	dir := getTestDataDirectory()
	data := NewData(dir)

	fref := NewFunctionBuilder(data).WithName(fname).WithService(sname).Build()

	f := Function{}
	if ok, _ := f.isValid(); ok {
		t.Errorf("Expected invalid function")
	}

	f.Name = fref.Name
	if ok, _ := f.isValid(); ok {
		t.Errorf("Expected invalid function")
	}

	f.Service = fref.Service
	if ok, _ := f.isValid(); ok {
		t.Errorf("Expected invalid service")
	}

	f.SourcePath = fref.SourcePath
	if ok, err := f.isValid(); !ok {
		t.Errorf("Expected no errors, get: %v", err)
	}
}

func TestFunctionBuilder(t *testing.T) {

	fname := "test_func"
	sname := "testservice"

	dir := getTestDataDirectory()
	data := NewData(dir)

	fb := NewFunctionBuilder(data)
	fb.WithName(fname)
	fb.WithService(sname)

	f := fb.Build()
	if f.Name != fname {
		t.Errorf("Expected %s, got %s", fname, f.Name)
	}

	if f.Service != sname {
		t.Errorf("Expected %s, got %s", sname, f.Service)
	}

	if ok, _ := f.isValid(); !ok {
		t.Errorf("Expected valid function")
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()

		fb.Build()
		t.Error("Expected panic")
	}()

	func() {
		f2 := fb.New().WithName(fname).WithService("unknow")

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()

		f2.Build()
		t.Errorf("Expected panic")
	}()
}
