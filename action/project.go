package action

import (
	"fmt"
	"vmkube/vmio"
	"os"
	"vmkube/utils"
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
	response := Response{}
	return  response, nil
}

func (request *CmdRequest) AlterProject() (Response, error) {
	response := Response{}
	return  response, nil
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
			fmt.Println(os.Stderr, "Element Type not provided ...")
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
									fmt.Println(os.Stderr, "Error printing in output format : "+Sample)
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
									fmt.Println(os.Stderr, "Error printing in output format : "+Sample)
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
					fmt.Println(os.Stderr, "Unable to find Type : "+TypeVal)
					PrintCommandHelper(request.TypeStr, request.SubTypeStr)
					return Response{
						Message: "",
						Status: true,}, nil
				}
				fmt.Println(os.Stderr, "Wrong output format : "+Sample)
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
