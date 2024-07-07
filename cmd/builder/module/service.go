package module

import (
	"log"
	"os"
	"path/filepath"
)

const (
	ServicesDir = "./services"
)

type Service struct {
	Name string
}

func (s Service) SourcePath() string {
	return "." + string(filepath.Separator) + filepath.Join(ServicesDir, s.Name)
}

// This function list all folders inside the services directory
func ListServices() []Service {
	services := []Service{}
	entries, err := os.ReadDir(ServicesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		services = append(services, Service{Name: e.Name()})
	}

	return services
}
