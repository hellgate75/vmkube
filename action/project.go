package action

import (
	"fmt"
	"vmkube/vmio"
	"os"
	"vmkube/utils"
	"vmkube/model"
	"errors"
	"time"
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
	Name := ""
	InputFile := ""
	InputFormat := ""
	Force := false
	DestroyInfra := false
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		}
		if "input-file" == CorrectInput(option[0]) {
			InputFile = option[1]
		}
		if "input-format" == CorrectInput(option[0]) {
			InputFormat = option[1]
		}
		if "force" == CorrectInput(option[0]) {
			Force = CorrectInput(option[0]) == "true"
		}
		if "destroy-infra" == CorrectInput(option[0]) {
			DestroyInfra = CorrectInput(option[0]) == "true"
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Project Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	
	AllowProjectDeletion := Force
	AllowInfraDeletion := DestroyInfra
	AllowProjectBackup := Force
	AllowInfraBackup := DestroyInfra
	
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err == nil {
		if ! AllowProjectDeletion {
			AllowProjectDeletion = utils.RequestConfirmation("Do you want delete Project named '"+descriptor.Name+"'?")
			if ! AllowProjectDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if AllowProjectDeletion && ! AllowProjectBackup {
			AllowProjectBackup = utils.RequestConfirmation("Do you want backup Project named'"+descriptor.Name+"'?")
		}
	}
	
	if err == nil && !AllowProjectDeletion {
		response := Response{
			Status: false,
			Message: "Project '"+Name+"' already exists and no force clause specified ...",
		}
		return  response, errors.New("Unable to execute task")
	}
	
	if descriptor.InfraId != "" {
		if ! AllowInfraDeletion {
			AllowInfraDeletion = utils.RequestConfirmation("Do you want delete Infrastrcuture named '"+descriptor.InfraName+"'?")
			if ! AllowInfraDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if AllowInfraDeletion && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation("Do you want backup Infrastrcuture named'"+descriptor.InfraName+"'?")
		}
	}
	
	if descriptor.InfraId != "" && ! AllowInfraDeletion {
		response := Response{
			Status: false,
			Message: "Project '"+Name+"' already build in Infra '"+descriptor.InfraName+"' and no infrastructure destroy clause specified ...",
		}
		return  response, errors.New("Unable to execute task")
	}
	existsProject := (err == nil)
	existsInfrastructure := (descriptor.InfraId != "")
	existanceClause := "n't"
	existanceClause2 := "proceding with definition of new project"
	if existsProject {
		existanceClause = ""
		existanceClause2 = "proceding with overwrite of existing project"
	}
	fmt.Printf("\nProject: %s does%s exist, now %s...\n", Name, existanceClause, existanceClause2 )
	project := model.Project{}
	if InputFile != "" && InputFormat != "" {
		fmt.Printf("\nLoading project %s from file '%s' using format '%s'...\n", Name, InputFile, InputFormat )
		project, err = vmio.ImportUserProject(InputFile, InputFormat)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		project.LastMessage = "Project imported from file " + InputFile + " in format " + InputFormat
		project.Name = Name
	} else {
		fmt.Printf("\nDefining new empty project %s...\n", Name )
		project.Id = NewUUIDString()
		project.Name = Name
		project.Created = time.Now()
		project.Modified = time.Now()
		project.Open = true
		project.LastMessage = "Empty Project Creation"
		project.Domains = append(project.Domains, model.ProjectDomain{
			Id: NewUUIDString(),
			Name: "Default Domain",
			Options: [][]string{},
			Networks: []model.ProjectNetwork{},
		})
		project.Domains[0].Networks = append(project.Domains[0].Networks, model.ProjectNetwork{
			Id: NewUUIDString(),
			Name: "Defaut Network",
			Options: [][]string{},
			CServers: []model.ProjectCloudServer{},
			Servers: []model.ProjectServer{},
			Installations: []model.InstallationPlan{},
		})
	}
	
	
	if ErrorList := project.Validate(); len(ErrorList) > 0 {
		errorValue := "Imported Project is invalid, clause(s) :"
		for _,err := range ErrorList {
			errorValue += fmt.Sprintf("\n%s", err.Error())
		}
		response := Response{
			Status: false,
			Message: errorValue,
		}
		return  response, errors.New("Unable to execute task")
	}
	InfraBackup := ""
	ProjectBackup := ""

	if existsInfrastructure {
		if AllowInfraBackup {
			InfraBackup = fmt.Sprintf("%s%s.infra-%s-%s.json",model.GetEmergencyFolder(),string(os.PathSeparator),descriptor.InfraId, descriptor.InfraName)
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				vmio.ExportInfrastructure(infra,InfraBackup,"json",true)
				fmt.Printf("Emergency Infrastructure backup at : %s\n", InfraBackup)
			}
		}
		response, err := request.DeleteInfra()
		if err != nil {
			return response, err
		}
	}
	
	if existsProject {
		if AllowProjectBackup {
			ProjectBackup = fmt.Sprintf("%s%s.project-%s-%s.json",model.GetEmergencyFolder(),string(os.PathSeparator),descriptor.Id, descriptor.Name)
			project, err := vmio.LoadProject(descriptor.Id)
			if err == nil {
				vmio.ExportUserProject(project,ProjectBackup,"json",true)
				fmt.Printf("Emergency Project backup at : %s\n", ProjectBackup)
			}
		}
		response, err := request.DeleteProject()
		if err != nil {
			return response, err
		}
	}
	
	err = vmio.SaveProject(project)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		if existsProject {
			if existsInfrastructure {
				return  response, errors.New("Unable to execute task, Project '"+Name+"' and Infrastructure "+descriptor.InfraName+" no longer exist, no rollback available, check emergency backups in logs!!")
			} else {
				return  response, errors.New("Unable to execute task, Project '"+Name+"' no longer exist, no rollback available, check emergency backup in logs!!")
			}
		} else {
			return  response, errors.New("Unable to execute task")
		}
	}
	
	index, err := vmio.LoadIndex()
	
	index.Projects = append(index.Projects, model.ProjectsDescriptor{
		Id: project.Id,
		Name: project.Name,
		Open: project.Open,
		Synced: true,
		Active: false,
		InfraId: "",
		InfraName: "",
	})
	
	
	err = vmio.SaveIndex(index)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		request.DeleteProject()
		return  response, errors.New("Unable to execute task")
	}
	
	if ProjectBackup != "" {
		fmt.Printf("Removing Project backup at : %s\n", ProjectBackup)
		os.Remove(ProjectBackup)
	}
	
	if InfraBackup != "" {
		fmt.Printf("Removing Infrastructure backup at : %s\n", InfraBackup)
		os.Remove(InfraBackup)
	}
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
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
	Name := ""
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
			break
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Project Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	project, err := vmio.LoadProject(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	open := "no"
	if project.Open {
		open = "yes"
	}
	errors := "no"
	if project.Errors {
		errors = "yes"
	}
	fmt.Printf("Id: %s\nProject: %s\nOpen: %s\n", project.Id,project.Name, open)
	fmt.Printf("Created : %d-%02d-%02d %02d:%02d:%02d\n",
		project.Created.Year(), project.Created.Month(), project.Created.Day(),
		project.Created.Hour(), project.Created.Minute(), project.Created.Second())
	fmt.Printf("Modified : %d-%02d-%02d %02d:%02d:%02d\n",
		project.Modified.Year(), project.Modified.Month(), project.Modified.Day(),
		project.Modified.Hour(), project.Modified.Minute(), project.Modified.Second())
	fmt.Printf("Errors: %s\nLast Message: %s\n", errors,project.LastMessage)
	fmt.Printf("Domains: %d\n", len(project.Domains))
	for _,domain := range project.Domains {
		num, options := vmio.StripOptions(domain.Options)
		fmt.Printf("Domain: %s (Id: %s) - Options [%d] :%s\n", domain.Name, domain.Id, num, options)
		fmt.Printf("Networks: %d\n", len(domain.Networks))
		for _,network := range domain.Networks {
			num, options := vmio.StripOptions(network.Options)
			fmt.Printf("   Network: %s (Id: %s) - Options [%d] :%s\n", network.Name, network.Id, num, options)
			fmt.Printf("   Servers: %d\n", len(network.Servers))
			serversMap := make(map[string]string)
			for _,server := range network.Servers {
				serversMap[server.Id] = server.Name
				fmt.Printf("      Server: %s (Id: %s) - Driver: %s - OS : %s:%s\n", server.Name, server.Id, server.Driver, server.OSType, server.OSVersion)
			}
			fmt.Printf("   Cloud Servers: %d\n", len(network.CServers))
			for _,server := range network.CServers {
				serversMap[server.Id] = server.Name
				num, options := vmio.StripOptions(server.Options)
				fmt.Printf("      Server: %s (Id: %s) - Driver: %s - Options [%d] :%s\n", server.Name, server.Id, server.Driver, num, options)
			}
			fmt.Printf("   Installation Plans: %d\n", len(network.Installations))
			for _,installation := range network.Installations {
				serverName,ok := serversMap[installation.ServerId]
				if !ok {
					serverName = "<invalid>"
				}
				cloud := "no"
				if installation.IsCloud {
					cloud = "yes"
				}
				fmt.Printf("      Plan: Id: %s - Server: %s [Id: %s] - Cloud: %s - Envoronment : %s  Role: %s  Type: %s\n", installation.Id, serverName, installation.ServerId, cloud, installation.Environment, installation.Role, installation.Type)
			}
		}
	}
	return Response{
		Message: "Success",
		Status: true,}, nil
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
