package action

import (
	"fmt"
)

type InfrastructureActions interface {
	CheckInfra() bool
	CreateInfra() (Response, error)
	AlterInfra() (Response, error)
	DeleteInfra() (Response, error)
	StartInfra() (Response, error)
	StopInfra() (Response, error)
	RestartInfra() (Response, error)
	ListInfras() (Response, error)
	StatusInfra() (Response, error)
}

func (request *CmdRequest) CheckInfra() bool {
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

func (request *CmdRequest) CreateInfra() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) AlterInfra() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) DeleteInfra() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) StartInfra() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) StopInfra() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) RestartInfra() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) ListInfras() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) StatusInfra() (Response, error) {
	response := Response{}
	return  response, nil
}
