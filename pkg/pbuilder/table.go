package pbuilder

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	tablesFileName = "tables.json"
)

type TableFunctionPermision struct {
	Function   Function
	Permisions []string
}

type Table struct {
	Name        string
	Permissions []TableFunctionPermision
}

func (t Table) isValid() (bool, error) {
	if t.Name == "" {
		return false, errors.New("table name not defined")
	}
	return true, nil
}

type tableBuilder struct {
	data            Data
	table           Table
	functionBuilder *functionBuilder
	created         bool
	withNew         bool
}

func NewTableBuilder(data Data) *tableBuilder {
	return &tableBuilder{
		data: data,
	}
}

func (tb *tableBuilder) New() *tableBuilder {
	tb.withNew = true
	tb.created = false
	tb.functionBuilder = NewFunctionBuilder(tb.data)
	tb.table = Table{
		Permissions: []TableFunctionPermision{},
	}
	return tb
}

func (tb *tableBuilder) WithName(name string) *tableBuilder {
	tb.table.Name = name
	return tb
}

func (tb *tableBuilder) WithPermission(f Function, p []string) *tableBuilder {
	tb.table.Permissions = append(tb.table.Permissions, TableFunctionPermision{Function: f, Permisions: p})
	return tb
}

func (tb *tableBuilder) Build() Table {
	if !tb.withNew {
		panic("table New() not called before")
	}

	if tb.created {
		panic("Table already created. use New() first to create a new table")
	}

	if ok, err := tb.table.isValid(); !ok {
		panic(err)
	}

	tb.loadTablePermissions()

	tb.created = true
	return tb.table
}

func (tb *tableBuilder) loadTablePermissions() {
	services := listServices(tb.data.ServicesDir)
	for _, sname := range services {
		sdir := filepath.Join(tb.data.ServicesDir, sname)
		functions := listFunctions(sdir)
		for _, fname := range functions {
			f := tb.functionBuilder.New().WithService(sname).WithName(fname).Build()
			read, err := os.ReadFile(filepath.Join(f.SourcePath, "main.go"))
			if err != nil {
				panic(err)
			}

			permissions := []string{}

			if strings.Contains(string(read), "dyna.GetItem") {
				permissions = append(permissions, "get")
			}

			if strings.Contains(string(read), "dyna.PutItem") {
				permissions = append(permissions, "put")
			}

			if len(permissions) > 0 {
				tb.WithPermission(f, permissions)
			}

		}
	}
}

func listTables(dir string) []string {
	tables := []string{}

	filename := filepath.Join(dir, tablesFileName)
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return tables
	}

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var jtables []struct {
		Name string `json:"name"`
	}

	json.Unmarshal(byteValue, &jtables)
	for _, t := range jtables {
		tables = append(tables, t.Name)
	}

	return tables
}
