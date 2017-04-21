package action

import (
	"fmt"
	"vmkube/vmio"
	"os"
	"vmkube/utils"
	"vmkube/model"
	"errors"
	"time"
	"strings"
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
			Force = GetBoolean(option[1])
		}
		if "destroy-infra" == CorrectInput(option[0]) {
			DestroyInfra = GetBoolean(option[1])
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
			AllowProjectDeletion = utils.RequestConfirmation("Do you want proceed with deletion process for Project named '"+descriptor.Name+"'?")
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
			AllowInfraDeletion = utils.RequestConfirmation("Do you want proceed with deletion process for Infrastrcuture named '"+descriptor.InfraName+"'?")
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
		_, errorValue := vmio.StripErrorMessages("Imported Project is invalid, clause(s) :", ErrorList)
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
	
	iFaceProject := vmio.IFaceProject{
		Id: descriptor.Id,
	}
	iFaceProject.WaitForUnlock()
	
	
	vmio.LockProject(project)
	
	err = vmio.SaveProject(project)
	
	vmio.UnlockProject(project)
	
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
	
	indexes, err := vmio.LoadIndex()
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()
	
	indexes, err = vmio.LoadIndex()
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	vmio.LockIndex(indexes)
	
	indexes.Projects = append(indexes.Projects, model.ProjectsDescriptor{
		Id: project.Id,
		Name: project.Name,
		Open: project.Open,
		Synced: true,
		Active: false,
		InfraId: "",
		InfraName: "",
	})
	
	
	err = vmio.SaveIndex(indexes)
	
	vmio.UnlockIndex(indexes)
	
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
	Name := ""
	Force := false
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		}
		if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		}
	}
	AllowProjectDeletion := Force
	
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	if err == nil {
		if ! AllowProjectDeletion {
			if descriptor.InfraId == "" {
				AllowProjectDeletion = utils.RequestConfirmation("Do you want delete Project named '"+descriptor.Name+"'?")
			} else {
				AllowProjectDeletion = utils.RequestConfirmation("Do you want delete Project named '"+descriptor.Name+" and Infrastrcuture named'"+descriptor.InfraName+"'?")
			}
			if ! AllowProjectDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
	}
	existsInfrastructure := (descriptor.InfraId != "")
	existanceClause := ""
	if existsInfrastructure {
		existanceClause = " and proceding with deletion of existing Infrastrcuture named'"+descriptor.InfraName+"'"
	}
	fmt.Printf("\nProceding with deletion of Project named  '%s'%s...\n", descriptor.Name, existanceClause )
	

	indexes,err := vmio.LoadIndex()

	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	if existsInfrastructure {
		resp, err := request.DeleteInfra()
		if err != nil {
			return resp, err
		}
	}
	
	projectMeta := model.Project{
		Id: descriptor.Id,
	}
	
	iFaceProject := vmio.IFaceProject{
		Id: descriptor.Id,
	}
	iFaceProject.WaitForUnlock()
	
	vmio.LockProject(projectMeta)
	
	info := vmio.ProjectInfo{
		Format: "",
		Project: model.Project{
			Id: descriptor.Id,
		},
	}
	
	err = info.Delete()

	vmio.UnlockProject(projectMeta)

	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	
	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()
	
	indexes, err = vmio.LoadIndex()
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	vmio.LockIndex(indexes)
	
	
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _,prj := range indexes.Projects {
		if CorrectInput(prj.Name) != Name {
			NewIndexes = append(NewIndexes, )
		}
	}
	SaveIndex := len(indexes.Projects) > len(NewIndexes)
	if SaveIndex {
		indexes.Projects = NewIndexes
		err = vmio.SaveIndex(indexes)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			vmio.UnlockIndex(indexes)
			return  response, errors.New("Unable to execute task")
		}
	}
	vmio.UnlockIndex(indexes)
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
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
		fmt.Printf("%s  %s  %s  %s  %s\n", utils.StrPad("Project Id", 40), utils.StrPad("Project Name", 40), utils.StrPad("Open", 4), utils.StrPad("Infrastructure Name", 40), utils.StrPad("Active", 6))
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
		fmt.Printf("%s  %s  %s  %s  %s\n", utils.StrPad(index.Id, 40), utils.StrPad(index.Name, 40), utils.StrPad(open, 4), utils.StrPad(index.InfraName, 40), utils.StrPad("  "+active, 6))
		
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
	iFaceProject := vmio.IFaceProject{
		Id: descriptor.Id,
	}
	iFaceProject.WaitForUnlock()
	
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
	Name := ""
	File := ""
	Format := ""
	FullImport := true
	var ElementType CmdElementType = NoElement
	var err error
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = CorrectInput(option[1])
		} else if "file" == CorrectInput(option[0]) {
			File = option[1]
		} else if "format" == CorrectInput(option[0]) {
			Format = option[1]
		} else if "full-export" == CorrectInput(option[0]) {
			FullImport = GetBoolean(option[1])
		} else if "elem-type" == CorrectInput(option[0]) {
			ElementType, err = CmdParseElement(option[0],option[1])
			if err != nil {
				ElementType = NoElement
			}
		}
	}
	if strings.TrimSpace(Name) == "" {
		response := Response{
			Status: false,
			Message: "Project Name Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(File) == "" {
		response := Response{
			Status: false,
			Message: "Import File Path Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(Format) == "" {
		response := Response{
			Status: false,
			Message: "Import Format Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	fmt.Printf("Import File Path : %s, Format: %s\n", File, Format)
	Full := "no"
	if FullImport {
		Full = "yes"
	}
	fmt.Printf("Import Full : %s\n", Full)
	if ! FullImport {
		ImportDomains := "no"
		if ElementType == SDomain || ElementType == SProject {
			ImportDomains = "yes"
		}
		ImportNetworks := "no"
		if ElementType == SNetwork || ElementType == SProject {
			ImportNetworks = "yes"
		}
		ImportServers := "no"
		if ElementType == LServer || ElementType == SProject {
			ImportServers = "yes"
		}
		ImportCServers := "no"
		if ElementType == CLServer || ElementType == SProject {
			ImportCServers = "yes"
		}
		ImportPlans := "no"
		if ElementType == SPlan || ElementType == SProject {
			ImportPlans = "yes"
		}
		fmt.Printf("Import Domains : %s\nImport Networks : %s\nImport Servers : %s\nImport Cloud Servers : %s\nImport Plans : %s\n",
			ImportDomains, ImportNetworks, ImportServers, ImportCServers, ImportPlans)
	}
	if ! FullImport && ElementType == NoElement {
		response := Response{
			Status: false,
			Message: "Nothing to import, application quits ...",
		}
		return  response, errors.New("Unable to execute task")
	}
	model.DeleteIfExists(File)
	if FullImport || ElementType == SProject {
		project, err := vmio.LoadProject(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		project, err = vmio.ImportUserProject(File, Format)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		project.PostImport()
		errorList := project.Validate()
		if len(errorList) > 0 {
			_, errorValue := vmio.StripErrorMessages("Project import is invalid, clause(s) :", errorList)
			response := Response{
				Status: false,
				Message: errorValue,
			}
			return  response, errors.New("Unable to execute task")
		}
	} else {
		project, err := vmio.LoadProject(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		switch ElementType {
		case SDomain:
			domains := ExportImportDomains{
				Domains: []model.ProjectDomain{},
			}
			domains, err = UserImportDomains(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			for _, domain := range domains.Domains {
				domain.PostImport()
				errorList := domain.Validate()
				if len(errorList) > 0 {
					_, errorValue := vmio.StripErrorMessages("Domains import is invalid, clause(s) :", errorList)
					response := Response{
						Status: false,
						Message: errorValue,
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			fmt.Printf("Successfully imported Project '%s' Domains to file '%s' in format '%s' ", Name, File, Format)
			break
		case SNetwork:
			var networks []ExportImportNetwork
			networks, err = UserImportNetworks(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			for _,domain := range project.Domains {
				network := ExportImportNetwork{
					Domain: DomainReference{
						DomainId: domain.Id,
						DomainName: domain.Name,
					},
					Networks: domain.Networks,
				}
				for _, network := range network.Networks {
					network.PostImport()
					errorList := network.Validate()
					if len(errorList) > 0 {
						_, errorValue := vmio.StripErrorMessages("Networks import is invalid, clause(s) :", errorList)
						response := Response{
							Status: false,
							Message: errorValue,
						}
						return  response, errors.New("Unable to execute task")
					}
				}
				networks = append(networks, network)
			}
			fmt.Printf("Successfully imported Project '%s' Networks to file '%s' in format '%s' ", Name, File, Format)
			break
		case  LServer:
			var servers []ExportImportLocalServers
			servers, err = UserImportLocalServers(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			for _,domain := range project.Domains {
				for _,network := range domain.Networks {
					server := ExportImportLocalServers{
						Network: NetworkReference{
							DomainId: domain.Id,
							DomainName: domain.Name,
							NetworkId: network.Id,
							NetworkName: network.Name,
						},
						Servers: network.Servers,
					}
					for _, server := range server.Servers {
						server.PostImport()
						errorList := server.Validate()
						if len(errorList) > 0 {
							_, errorValue := vmio.StripErrorMessages("Local Servers import is invalid, clause(s) :", errorList)
							response := Response{
								Status: false,
								Message: errorValue,
							}
							return  response, errors.New("Unable to execute task")
						}
					}
					servers = append(servers, server)
				}
			}
			fmt.Printf("Successfully imported Project '%s' Local Servers to file '%s' in format '%s' ", Name, File, Format)
			break
		case CLServer:
			var servers []ExportImportCloudServers
			servers, err = UserImportCloudServers(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			for _,domain := range project.Domains {
				for _,network := range domain.Networks {
					server := ExportImportCloudServers{
						Network: NetworkReference{
							DomainId: domain.Id,
							DomainName: domain.Name,
							NetworkId: network.Id,
							NetworkName: network.Name,
						},
						Servers: network.CServers,
					}
					for _, server := range server.Servers {
						server.PostImport()
						errorList := server.Validate()
						if len(errorList) > 0 {
							_, errorValue := vmio.StripErrorMessages("Cloud Servers import is invalid, clause(s) :", errorList)
							response := Response{
								Status: false,
								Message: errorValue,
							}
							return  response, errors.New("Unable to execute task")
						}
					}
					servers = append(servers, server)
				}
			}
			fmt.Printf("Successfully imported Project '%s' Cloud Servers to file '%s' in format '%s' ", Name, File, Format)
			break
		default:
			//Plans
			var plans []ExportImportPlans
			plans, err = UserImportPlans(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			for _,domain := range project.Domains {
				for _,network := range domain.Networks {
					plan := ExportImportPlans{
						Network: NetworkReference{
							DomainId: domain.Id,
							DomainName: domain.Name,
							NetworkId: network.Id,
							NetworkName: network.Name,
						},
						Plans: network.Installations,
					}
					for _,plan := range plan.Plans {
						plan.PostImport()
						errorList := plan.Validate()
						if len(errorList) > 0 {
							_, errorValue := vmio.StripErrorMessages("Plans import is invalid, clause(s) :", errorList)
							response := Response{
								Status: false,
								Message: errorValue,
							}
							return  response, errors.New("Unable to execute task")
						}
					}
					plans = append(plans, plan)
				}
			}
			fmt.Printf("Successfully imported Project '%s' Installation Plans to file '%s' in format '%s' ", Name, File, Format)
		}
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) ExportProject() (Response, error) {
	Name := ""
	File := ""
	Format := ""
	FullExport := true
	var ElementType CmdElementType = NoElement
	var err error
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = CorrectInput(option[1])
		} else if "file" == CorrectInput(option[0]) {
			File = option[1]
		} else if "format" == CorrectInput(option[0]) {
			Format = option[1]
		} else if "full-export" == CorrectInput(option[0]) {
			FullExport = GetBoolean(option[1])
		} else if "elem-type" == CorrectInput(option[0]) {
			ElementType, err = CmdParseElement(option[0],option[1])
			if err != nil {
				ElementType = NoElement
			}
		}
	}
	if strings.TrimSpace(Name) == "" {
		response := Response{
			Status: false,
			Message: "Project Name Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(File) == "" {
		response := Response{
			Status: false,
			Message: "Export File Path Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(Format) == "" {
		response := Response{
			Status: false,
			Message: "Export Format Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	fmt.Printf("Export File Path : %s, Format: %s\n", File, Format)
	Full := "no"
	if FullExport {
		Full = "yes"
	}
	fmt.Printf("Export Full : %s\n", Full)
	if ! FullExport {
		ExportDomains := "no"
		if ElementType == SDomain || ElementType == SProject {
			ExportDomains = "yes"
		}
		ExportNetworks := "no"
		if ElementType == SNetwork || ElementType == SProject {
			ExportNetworks = "yes"
		}
		ExportServers := "no"
		if ElementType == LServer || ElementType == SProject {
			ExportServers = "yes"
		}
		ExportCServers := "no"
		if ElementType == CLServer || ElementType == SProject {
			ExportCServers = "yes"
		}
		ExportPlans := "no"
		if ElementType == SPlan || ElementType == SProject {
			ExportPlans = "yes"
		}
		fmt.Printf("Export Domains : %s\nExport Networks : %s\nExport Servers : %s\nExport Cloud Servers : %s\nExport Plans : %s\n",
			ExportDomains, ExportNetworks, ExportServers, ExportCServers, ExportPlans)
	}
	if ! FullExport && ElementType == NoElement {
		response := Response{
			Status: false,
			Message: "Nothing to export, application quits ...",
		}
		return  response, errors.New("Unable to execute task")
	}
	model.DeleteIfExists(File)
	if FullExport || ElementType == SProject {
		project, err := vmio.LoadProject(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		err = vmio.ExportUserProject(project, File, Format, true)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
	} else {
		project, err := vmio.LoadProject(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		switch ElementType {
		case SDomain:
			domains := ExportImportDomains{
				Domains: []model.ProjectDomain{},
			}
			domains.Domains = append(domains.Domains, project.Domains...)
			err = utils.ExportStructureToFile(File, Format, domains)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			fmt.Printf("Successfully exported Project '%s' Domains to file '%s' in format '%s' ", Name, File, Format)
			break
		case SNetwork:
			var networks []ExportImportNetwork = make([]ExportImportNetwork, 0)
			for _,domain := range project.Domains {
				network := ExportImportNetwork{
					Domain: DomainReference{
						DomainId: domain.Id,
						DomainName: domain.Name,
					},
					Networks: domain.Networks,
				}
				networks = append(networks, network)
			}
			err = utils.ExportStructureToFile(File, Format, networks)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			fmt.Printf("Successfully exported Project '%s' Networks to file '%s' in format '%s' ", Name, File, Format)
			break
		case  LServer:
			var servers []ExportImportLocalServers = make([]ExportImportLocalServers, 0)
				for _,domain := range project.Domains {
					for _,network := range domain.Networks {
						server := ExportImportLocalServers{
							Network: NetworkReference{
								DomainId: domain.Id,
								DomainName: domain.Name,
								NetworkId: network.Id,
								NetworkName: network.Name,
							},
							Servers: network.Servers,
						}
						servers = append(servers, server)
					}
			}
			err = utils.ExportStructureToFile(File, Format, servers)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			fmt.Printf("Successfully exported Project '%s' Local Servers to file '%s' in format '%s' ", Name, File, Format)
			break
		case CLServer:
			var servers []ExportImportCloudServers = make([]ExportImportCloudServers, 0)
			for _,domain := range project.Domains {
				for _,network := range domain.Networks {
					server := ExportImportCloudServers{
						Network: NetworkReference{
							DomainId: domain.Id,
							DomainName: domain.Name,
							NetworkId: network.Id,
							NetworkName: network.Name,
						},
						Servers: network.CServers,
					}
					servers = append(servers, server)
				}
			}
			err = utils.ExportStructureToFile(File, Format, servers)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			fmt.Printf("Successfully exported Project '%s' Cloud Servers to file '%s' in format '%s' ", Name, File, Format)
			break
		default:
			//Plans
			var plans []ExportImportPlans = make([]ExportImportPlans, 0)
			for _,domain := range project.Domains {
				for _,network := range domain.Networks {
					plan := ExportImportPlans{
						Network: NetworkReference{
							DomainId: domain.Id,
							DomainName: domain.Name,
							NetworkId: network.Id,
							NetworkName: network.Name,
						},
						Plans: network.Installations,
					}
					plans = append(plans, plan)
				}
			}
			err = utils.ExportStructureToFile(File, Format, plans)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			fmt.Printf("Successfully exported Project '%s' Installation Plans to file '%s' in format '%s' ", Name, File, Format)
		}
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}
