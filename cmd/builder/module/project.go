package module

// type Project struct {
// 	Services []Service
// }

// func (p Project) String() string {
// 	str := "Project"
// 	str += "\n  Services"
// 	for _, s := range p.Services {
// 		str += "\n    # " + s.Name
// 		str += "\n      Tables:"
// 		for _, t := range s.Tables {
// 			str += "\n        - " + t.Name
// 		}
// 		str += "\n      Functions:"
// 		for _, f := range s.Functions {
// 			str += "\n        - " + f.Name
// 		}
// 	}
// 	return str
// }

// func LoadProject() Project {
// 	project := Project{}

// 	services := ListServices()
// 	for _, s := range services {
// 		service := LoadService(s)
// 		project.Services = append(project.Services, service)
// 	}

// 	return project
// }
