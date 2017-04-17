package action

import "fmt"

type ProjectActions interface {
	CheckProject() bool
	CreateProject() (Response, error)
	AlterProject() (Response, error)
	DeleteProject() (Response, error)
	ListProjects() (Response, error)
	StatusProject() (Response, error)
	ImportProject() (Response, error)
	ExportProject() (Response, error)
}

func (request *CmdRequest) CheckProject() bool {
	if len(request.Arguments.Helper.Options) > 0 {
		correctness := true
		for _,option := range request.Arguments.Helper.Options {
			if "true" == option[3] {
				//Mandatory Option
				found := false
				for _,argument := range request.Arguments.Options {
					if CorrectInput(argument[0]) == option[0] {
						found = true
						break
					}
				}
				if !found {
					correctness = false
					fmt.Printf("Option '--%s' is mandatory!!\n", option[0])
				}
			}
		}
		if !correctness {
			return  false
		}
	}
	return  true
}

func (request *CmdRequest) CreateProject() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) AlterProject() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) DeleteProject() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) ListProjects() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) StatusProject() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) ImportProject() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) ExportProject() (Response, error) {
	response := Response{}
	return  response, nil
}
