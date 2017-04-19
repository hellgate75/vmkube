package action

import (
	"fmt"
	"vmkube/vmio"
	"os"
	"vmkube/utils"
	"vmkube/model"
	"errors"
)

type ProjectActions interface {
	CheckProject() bool
	CreateProject() (Response, error)
	AlterProject() (Response, error)
	InfoProject() (Response, error)
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
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) AlterProject() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) InfoProject() (Response, error) {
	if request.SubType == List {
		//List of elements
		defines := vmio.ListProjectTypeDefines()
		for _,define := range defines {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", utils.StrPad(define.Name,15), define.Description)
		}
	} else {
		//Detail
		TypeVal := ""
		Sample := ""
		for _,option := range request.Arguments.Options {
			if "sample" == CorrectInput(option[0]) {
				Sample = CorrectInput(option[1])
			} else if "elem-type" == CorrectInput(option[0]) {
				TypeVal = option[1]
			}
		}
		if TypeVal == "" {
			fmt.Printf("Element Type not provided ...\n")
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return Response{
				Message: "",
				Status: true,}, nil
		} else {
			if Sample != "" {
				if "json" == Sample || "xml" == Sample {
					defines := vmio.ListProjectTypeDefines()
					for _,define := range defines {
						if define.Name == TypeVal {
							if "json" == Sample {
								bytes, err := utils.GetJSONFromElem(define.Sample, true)
								if err != nil {
									fmt.Printf("Error printing in output format : %s\n", Sample)
									PrintCommandHelper(request.TypeStr, request.SubTypeStr)
									return Response{
										Message: "",
										Status: true,}, nil
								}
								fmt.Fprintf(os.Stdout, "%s\n", bytes)
								return Response{
									Message: "",
									Status: true,}, nil
							} else {
								bytes, err := utils.GetXMLFromElem(define.Sample, true)
								if err != nil {
									fmt.Printf("Error printing in output format : %s\n", Sample)
									PrintCommandHelper(request.TypeStr, request.SubTypeStr)
									return Response{
										Message: "",
										Status: true,}, nil
								}
								fmt.Fprintf(os.Stdout, "%s\n", bytes)
								return Response{
									Message: "",
									Status: true,}, nil
							}
						}
					}
					fmt.Printf("Unable to find Type : %s\n",TypeVal)
					PrintCommandHelper(request.TypeStr, request.SubTypeStr)
					return Response{
						Message: "",
						Status: true,}, nil
				}
				fmt.Printf("Wrong output format : %s\n", Sample)
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return Response{
					Message: "",
					Status: true,}, nil
			} else {
				defines := vmio.ListProjectTypeDefines()
				for _,define := range defines {
					if define.Name == TypeVal {
						fields, err := model.DescribeStruct(define.Sample)
						if err == nil {
							model.PrintFieldsHeader(len(fields)>0)
							model.PrintFieldsRecursively(fields, 0)
							return Response{
								Message: "",
								Status: true,}, nil
						}
						fmt.Printf("Unable to describe Type : %s\n",TypeVal)
						PrintCommandHelper(request.TypeStr, request.SubTypeStr)
						return Response{
							Message: "",
							Status: true,}, nil
					}
				}
				fmt.Printf("Unable to find Type : %s\n",TypeVal)
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return Response{
					Message: "",
					Status: true,}, nil
			}
			return Response{
				Message: "",
				Status: true,}, nil
		}
	}
	response := Response{Status: true}
	return  response, nil
}

func (request *CmdRequest) DeleteProject() (Response, error) {
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) ListProjects() (Response, error) {
	indexes, error := vmio.LoadIndex()
	if error != nil {
		response := Response{
			Status: false,
			Message: error.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	if len(indexes.Projects) > 0 {
		fmt.Printf("%s  %s  %s  %s  %s\n", utils.StrPad("Project Id", 20), utils.StrPad("Project Name", 20), utils.StrPad("Open", 4), utils.StrPad("Infrastructure Name", 20), utils.StrPad("Active", 6))
	} else {
		fmt.Printf("No Projects found\n")
	}
	for _,index := range indexes.Projects {
		open := "no"
		active := "no"
		if index.Open {
			open = "yes"
		}
		if index.Active {
			active = "yes"
		}
		fmt.Printf("%s  %s  %s  %s  %s\n", utils.StrPad(index.Id, 20), utils.StrPad(index.Name, 20), utils.StrPad(" "+open, 4), utils.StrPad(index.InfraName, 20), utils.StrPad("  "+active, 6))
		
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) StatusProject() (Response, error) {
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) ImportProject() (Response, error) {
	response := Response{}
	//Name := ""
	//File := ""
	//Format := ""
	//FullImport := true
	//ElementType := ""
	//for _,option := range request.Arguments.Options {
	//	if "name" == CorrectInput(option[0]) {
	//		Name = CorrectInput(option[1])
	//	} else if "elem-type" == CorrectInput(option[0]) {
	//		ElementType = option[1]
	//	}
	//}
	return  response, nil
}

func (request *CmdRequest) ExportProject() (Response, error) {
	response := Response{}
	return  response, nil
}
