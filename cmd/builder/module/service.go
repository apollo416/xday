package module

// const (
// 	ServicesDir = "./services"
// )

// type Service struct {
// 	Name      string
// 	Tables    []Table
// 	Functions []Function
// }

// func (s Service) SourcePath() string {
// 	return "." + string(filepath.Separator) + filepath.Join(ServicesDir, s.Name)
// }

// func LoadService(s string) Service {
// 	service := Service{Name: s}
// 	service.Tables = loadServiceTables(service)
// 	service.Functions = loadServiceFunctions(service)
// 	return service
// }

// func ListServices() []string {
// 	services := []string{}
// 	entries, err := os.ReadDir(ServicesDir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, e := range entries {
// 		services = append(services, e.Name())
// 	}

// 	return services
// }
