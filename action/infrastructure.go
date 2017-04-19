package action

import (
	"fmt"
	"errors"
	"vmkube/utils"
	"vmkube/vmio"
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
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) AlterInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) DeleteInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) StartInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) StopInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) RestartInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) ListInfras() (Response, error) {
	indexes, error := vmio.LoadIndex()
	if error != nil {
		response := Response{
			Status: false,
			Message: error.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	if len(indexes.Projects) > 0 {
		fmt.Printf("%s  %s  %s\n", utils.StrPad("Infrastructure Id", 20), utils.StrPad("Infrastructure Name", 20), utils.StrPad("Active", 6))
	} else {
		fmt.Printf("No Infrastructures found\n")
	}
	for _,index := range indexes.Projects {
		active := "no"
		if index.Active {
			active = "yes"
		}
		fmt.Printf("%s  %s  %s\n", utils.StrPad(index.InfraId, 20), utils.StrPad(index.InfraName, 20), utils.StrPad("  "+active, 6))
		
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) StatusInfra() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}
