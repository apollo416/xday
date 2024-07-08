package pbuilder

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ApiMethod struct {
	Function Function
	Method   string
}

type Api struct {
	Service Service
	Methods []ApiMethod
}

func (t Api) isValid() (bool, error) {
	if _, err := os.Stat(t.Service.SourcePath); os.IsNotExist(err) {
		return false, errors.New("Service does not exist")
	}
	return true, nil
}

type apiBuilder struct {
	data            Data
	api             Api
	functionBuilder *functionBuilder
	created         bool
	withNew         bool
}

func NewApiBuilder(data Data) *apiBuilder {
	return &apiBuilder{
		data: data,
	}
}

func (ab *apiBuilder) New() *apiBuilder {
	ab.withNew = true
	ab.created = false
	ab.api = Api{
		Methods: []ApiMethod{},
	}
	ab.functionBuilder = NewFunctionBuilder(ab.data)
	return ab
}

func (ab *apiBuilder) WithService(service Service) *apiBuilder {
	ab.api.Service = service
	return ab
}

func (ab *apiBuilder) Build() Api {
	if !ab.withNew {
		panic("table New() not called before")
	}

	if ab.created {
		panic("Table already created. use New() first to create a new table")
	}

	if ok, err := ab.api.isValid(); !ok {
		panic(err)
	}

	ab.loadMethods()

	ab.created = true
	return ab.api
}

func (ab *apiBuilder) loadMethods() {
	functions := listFunctions(ab.api.Service.SourcePath)
	for _, fname := range functions {
		f := ab.functionBuilder.New().WithService(ab.api.Service.Name).WithName(fname).Build()
		read, err := os.ReadFile(filepath.Join(f.SourcePath, "main.go"))
		if err != nil {
			// TODO: add description to error
			log.Fatal(err)
		}

		if strings.Contains(string(read), "api.POST") {
			ab.api.Methods = append(ab.api.Methods, ApiMethod{Function: f, Method: "POST"})
		}
	}
}

func (a *ApiMethod) String() string {
	return fmt.Sprintf("(Function: %s, method: %s)", a.Function.Name, a.Method)
}

func (a *Api) String() string {
	str := fmt.Sprintf("(Service: %s", a.Service.Name)
	str += ", methods: ["
	methods := []string{}
	for _, m := range a.Methods {
		methods = append(methods, m.String())
	}
	str += strings.Join(methods, ", ")
	str += "])"
	return str
}
