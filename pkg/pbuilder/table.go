package pbuilder

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type TableFunctionPermision struct {
	Function   string
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
	data    Data
	table   Table
	created bool
	withNew bool
}

func NewTableBuilder(data Data) *tableBuilder {
	return &tableBuilder{
		data: data,
	}
}

func (tb *tableBuilder) New() *tableBuilder {
	tb.withNew = true
	tb.created = false
	tb.table = Table{
		Permissions: []TableFunctionPermision{},
	}
	return tb
}

func (tb *tableBuilder) WithName(name string) *tableBuilder {
	tb.table.Name = name
	return tb
}

func (tb *tableBuilder) WithPermission(f string, p []string) *tableBuilder {
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
			cname := filepath.Join(sdir, fname, "main.go")
			read, err := os.ReadFile(cname)
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
				tb.WithPermission(fname, permissions)
			}

		}
	}
}
