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
	"vmkube/operations"
)

type ProjectActions interface {
	CheckProject() bool
	CreateProject() (Response, error)
	AlterProject() (Response, error)
	InfoProject() (Response, error)
	DeleteProject() (Response, error)
	ListProjects() (Response, error)
	StatusProject() (Response, error)
	BuildProject() (Response, error)
	ImportProject() (Response, error)
	ExportProject() (Response, error)
}

func (request *CmdRequest) CheckProject() bool {
	if len(request.Arguments.Helper.Options) > 0 {
		correctness := true
		for _,option := range request.Arguments.Helper.Options {
			if option.Mandatory {
				//Mandatory Option
				found := false
				for _,argument := range request.Arguments.Options {
					if CorrectInput(argument[0]) == option.Option {
						found = true
						break
					}
				}
				if !found {
					correctness = false
					fmt.Printf("Option '--%s' is mandatory!!\n", option.Option)
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
	InputFormat := "json"
	Force := false
	OverrideInfra := false
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "file" == CorrectInput(option[0]) {
			InputFile = option[1]
		} else if "format" == CorrectInput(option[0]) {
			InputFormat = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
			option[1] = "true"
		} else if "override-infra" == CorrectInput(option[0]) {
			OverrideInfra = GetBoolean(option[1])
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Project Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	
	AllowProjectDeletion := Force
	AllowInfraDeletion := OverrideInfra && Force
	AllowProjectBackup := Force
	AllowInfraBackup := OverrideInfra && Force
	
	descriptor, err := vmio.GetProjectDescriptor(Name)

	var ProjectJSON string = ""
	
	existsProject := (err == nil)
	existsInfrastructure := existsProject && (descriptor.InfraId != "")
	
	
	if err == nil {
		if ! descriptor.Open {
			response := Response{
				Status: false,
				Message: "Project closed!!",
			}
			return response, errors.New("Unable to execute task")
		}
		oldProject, err2 := vmio.LoadProject(descriptor.Id)
		
		if err2 != nil && oldProject.Id != "" {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		ProjectJSON = string(utils.GetJSONFromObj(oldProject, true))
	}
	
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
	
	if descriptor.InfraId != "" && existsInfrastructure && OverrideInfra {
		if ! AllowInfraDeletion {
			AllowInfraDeletion = utils.RequestConfirmation("Do you want proceed with deletion process for Infrastructure named '"+descriptor.InfraName+"'?")
			if ! AllowInfraDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if AllowInfraDeletion && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation("Do you want backup Infrastructure named'"+descriptor.InfraName+"'?")
		}
	}
	
	if descriptor.InfraId != "" && ! AllowInfraDeletion &&  existsInfrastructure && OverrideInfra {
		response := Response{
			Status: false,
			Message: "Project '"+Name+"' already build in Infra '"+descriptor.InfraName+"' and no infrastructure destroy clause specified ...",
		}
		return  response, errors.New("Unable to execute task")
	}
	
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
		project.LastMessage = "Project imported from file " + InputFile + ", format " + InputFormat
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
			InfraBackup = fmt.Sprintf("%s%s.prj-%s-%s-infra-export-%s-%s.vmkube",model.GetEmergencyFolder(),string(os.PathSeparator), utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name),utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
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
			ProjectBackup = fmt.Sprintf("%s%s.project-%s-%s.json",model.GetEmergencyFolder(),string(os.PathSeparator),utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name))
			project, err := vmio.LoadProject(descriptor.Id)
			if err == nil {
				vmio.ExportUserProject(project,ProjectBackup,"json",true)
				fmt.Printf("Emergency Project backup at : %s\n", ProjectBackup)
			}
		}
		request.Arguments.Options = append(request.Arguments.Options, []string{"skip-indexes", "true"})
		response, err := request.DeleteProject()
		if err != nil {
			return response, err
		}
	}
	
	if existsInfrastructure {
		if AllowInfraBackup {
			InfraBackup = fmt.Sprintf("%s%s.prj-%s-%s-infra-export-%s-%s.vmkube",model.GetEmergencyFolder(),string(os.PathSeparator), utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name),utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				fmt.Printf("Emergency Infrastructure backup at : %s\n", InfraBackup)
			}
		}
		if (AllowInfraDeletion) {
			response, err := request.DeleteInfra()
			if err != nil {
				return response, err
			}
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
	
	err = UpdateIndexWithProjectsDescriptor(model.ProjectsDescriptor{
		Id: project.Id,
		Name: project.Name,
		Open: project.Open,
		Synced: true,
		Active: false,
		InfraId: "",
		InfraName: "",
	}, true)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		request.Arguments.Options = append(request.Arguments.Options, []string{"skip-indexes", "true"})
		response2, err := request.DeleteProject()
		if err != nil {
			return response2, err
		}
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
	
	if existsProject && existsInfrastructure && OverrideInfra {
		AddProjectChangeActions(descriptor.Id, ActionDescriptor{
			Id: NewUUIDString(),
			Date: time.Now(),
			DropAction: true,
			ElementType: SProject,
			ElementName: descriptor.Name,
			FullProject: true,
			JSONImage: ProjectJSON,
			Request: request.Type,
			SubRequest: request.SubType,
		})
	}
	
	if existsInfrastructure {
		request.Arguments.Options = append(request.Arguments.Options, []string{"rebuild", "true"})
		request.BuildProject()
	}
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) AlterProject() (Response, error) {
	Name := ""
	File := ""
	Format := "json"
	Force := true
	OverrideInfra := false
	var ElementType CmdElementType = NoElement
	ElementName := ""
	ElementId := ""
	var AnchorElementType CmdElementType = NoElement
	AnchorElementName := ""
	AnchorElementId := ""
	var err error
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = CorrectInput(option[1])
		} else if "file" == CorrectInput(option[0]) {
			File = option[1]
		} else if "format" == CorrectInput(option[0]) {
			Format = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "override-infra" == CorrectInput(option[0]) {
			OverrideInfra = GetBoolean(option[1])
		} else if "elem-type" == CorrectInput(option[0]) {
			ElementType, err = CmdParseElement(option[1])
			if err != nil {
				ElementType = NoElement
			}
		} else if "elem-id" == CorrectInput(option[0]) {
			ElementId = option[1]
		} else if "elem-name" == CorrectInput(option[0]) {
			ElementName = option[1]
		} else if "anchor-elem-type" == CorrectInput(option[0]) {
			AnchorElementType, err = CmdParseElement(option[1])
			if err != nil {
				ElementType = NoElement
			}
		} else if "anchor-elem-id" == CorrectInput(option[0]) {
			AnchorElementId = option[1]
		} else if "anchor-elem-name" == CorrectInput(option[0]) {
			AnchorElementName = option[1]
		}
		//else {
		//	PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		//	return Response{
		//		Message: fmt.Sprintf("Argument %s not provided, please review the help", option[0]),
		//		Status: false,},errors.New("Unable to execute task")
		//}
	}
	if strings.TrimSpace(Name) == "" {
		response := Response{
			Status: false,
			Message: "Project Name Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(File) == "" && request.SubType != Open && request.SubType != Close && request.SubType != Remove  {
		response := Response{
			Status: false,
			Message: "Input File Path Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	//if strings.TrimSpace(Format) == "" && request.SubType != Open && request.SubType != Close && request.SubType != Remove {
	//	response := Response{
	//		Status: false,
	//		Message: "Input File Format Field is mandatory",
	//	}
	//	return  response, errors.New("Unable to execute task")
	//}
	if ElementType == NoElement && request.SubType != Open && request.SubType != Close {
		response := Response{
			Status: false,
			Message: "Element Type Field is mandatory, use project-import for massive changes",
		}
		return  response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(ElementName) == "" && strings.TrimSpace(ElementId) == "" && request.SubType == Open && request.SubType != Close {
		response := Response{
			Status: false,
			Message: "Element Name Field or Id is mandatory, use project-import for massive changes",
		}
		return  response, errors.New("Unable to execute task")
	}
	if AnchorElementType == NoElement && request.SubType == Create {
		response := Response{
			Status: false,
			Message: "Anchor Element Type Field is mandatory in Alter Project Create command",
		}
		return  response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(AnchorElementName) == "" && strings.TrimSpace(AnchorElementId) == ""  && request.SubType == Create {
		response := Response{
			Status: false,
			Message: "Anchor Element Name or Id Fields is mandatory in Alter Project Create command",
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
	if ! descriptor.Open {
		response := Response{
			Status: false,
			Message: "Project closed!!",
		}
		return response, errors.New("Unable to execute task")
	}
	
	var ProjectJSON string = ""
	project, err2 := vmio.LoadProject(descriptor.Id)
	
	if err2 != nil && project.Id != "" {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	ProjectJSON = string(utils.GetJSONFromObj(project, true))
	existsProject := true
	existsInfrastructure := (descriptor.InfraId != "")
	
	switch request.SubType {
		case Create:
			switch ElementType {
				case SProject:
					response := Response{
						Status: false,
						Message: "Entire project define not allowed by alter-project, use import-project or define-project for project import",
					}
					return  response, errors.New("Unable to execute task")
				default:
					// Any but Project Type
					err = operations.AddElementToProject(project, int(ElementType), ElementName, int(AnchorElementType), AnchorElementName, AnchorElementId, File, Format)
					
					if err != nil {
						response := Response{
							Status: false,
							Message: err.Error(),
						}
						return  response, errors.New("Unable to execute task")
					}
			}
			break
		case Alter:
			switch ElementType {
			case SProject:
				response := Response{
					Status: false,
					Message: "Entire project chenge not allowed by alter-project, use import-project or define-project for project replacement",
				}
				return  response, errors.New("Unable to execute task")
			default:
				// Any but Project Type
				err = operations.AlterElementInProject(project, int(ElementType), ElementName, ElementId, File, Format)
				
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			break
		case Remove:
			switch ElementType {
			case SProject:
				response := Response{
					Status: false,
					Message: "Entire project chenge not allowed by alter-project, use delete-project for project removal",
				}
				return  response, errors.New("Unable to execute task")
			default:
				// Any but Project Type
				err = operations.DeleteElementInProject(project, int(ElementType), ElementName, ElementId)
				
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			break
		case Open:
			project, err = operations.OpenProject(project)
			
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			break
		case Close:
			project, err = operations.CloseProject(project)
			
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return  response, errors.New("Unable to execute task")
			}
			break
		default:
			response := Response{
				Status: false,
				Message: fmt.Sprintf("Sub-Command %s not provided!!", request.SubTypeStr),
			}
			return response, errors.New("Unable to execute task")
	}

	AllowProjectOverwrite := Force
	if existsProject && ! AllowProjectOverwrite {
		AllowProjectOverwrite = utils.RequestConfirmation("Do you want proceed with deletion process for Infrastructure named '"+descriptor.InfraName+"'?")
		if ! AllowProjectOverwrite {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
		
	}
	
	AllowInfraDeletion := OverrideInfra && Force
	AllowInfraBackup := OverrideInfra && Force
	InfraBackup := ""
	
	if descriptor.InfraId != "" && OverrideInfra {
		if ! AllowInfraDeletion {
			AllowInfraDeletion = utils.RequestConfirmation("Do you want proceed with deletion process for Infrastructure named '"+descriptor.InfraName+"'?")
			if ! AllowInfraDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if AllowInfraDeletion && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation("Do you want backup Infrastructure named'"+descriptor.InfraName+"'?")
		}
	}
	
	if existsInfrastructure && OverrideInfra {
		if AllowInfraBackup {
			InfraBackup = fmt.Sprintf("%s%s.prj-%s-%s-infra-export-%s-%s.vmkube",model.GetEmergencyFolder(),string(os.PathSeparator), utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name),utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				fmt.Printf("Emergency Infrastructure backup at : %s\n", InfraBackup)
			}
		}
		if (AllowInfraDeletion) {
			response, err := request.DeleteInfra()
			if err != nil {
				return response, err
			}
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
		if existsInfrastructure {
			return  response, errors.New("Unable to execute task, Infrastructure "+descriptor.InfraName+" no longer exist, no rollback available, check emergency backups in logs!!")
		} else {
			return  response, errors.New("Unable to execute task")
		}
	}
	
	if InfraBackup != "" {
		fmt.Printf("Removing Infrastructure backup at : %s\n", InfraBackup)
		os.Remove(InfraBackup)
	}
	
	
	
	if existsInfrastructure {
		request.Arguments.Options = append(request.Arguments.Options, []string{"rebuild", "true"})
		request.BuildProject()
	}
	
	if existsProject && existsInfrastructure {
		AddProjectChangeActions(descriptor.Id, ActionDescriptor{
			Id: NewUUIDString(),
			Date: time.Now(),
			DropAction: true,
			ElementType: SProject,
			ElementName: descriptor.Name,
			FullProject: true,
			JSONImage: ProjectJSON,
			Request: request.Type,
			SubRequest: request.SubType,
		})
	}
	
	if existsProject {
		indexes,err := vmio.LoadIndex()
		
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
		
		err = UpdateIndexWithProjectsDescriptor(model.ProjectsDescriptor{
			Id: descriptor.Id,
			Name: descriptor.Name,
			Open: project.Open,
			Synced: false,
			Active: descriptor.Active,
			InfraId: descriptor.InfraId,
			InfraName: descriptor.InfraName,
		}, true)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			request.Arguments.Options = append(request.Arguments.Options, []string{"skip-indexes", "true"})
			response2, err := request.DeleteProject()
			if err != nil {
				return response2, err
			}
			return  response, errors.New("Unable to execute task")
		}
	}
	
	
	if existsInfrastructure && OverrideInfra {
		request.Arguments.Options = append(request.Arguments.Options, []string{"rebuild", "true"})
		request.BuildProject()
	}
	
	response := Response{
		Status: true,
		Message: "Success",
	}
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
			fmt.Println("Element Type not provided ...\n")
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
	SkipIndexes := false
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		}
		if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		}
		if "skip-indexes" == CorrectInput(option[0]) {
			SkipIndexes = GetBoolean(option[1])
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
				AllowProjectDeletion = utils.RequestConfirmation("Do you want delete Project named '"+descriptor.Name+" and Infrastructure named'"+descriptor.InfraName+"'?")
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
		existanceClause = " and proceding with deletion of existing Infrastructure named'"+descriptor.InfraName+"'"
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
	
	err = UpdateIndexWithProjectsDescriptor(descriptor, false)

	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		vmio.UnlockIndex(indexes)
		return  response, errors.New("Unable to execute task")
	}
	vmio.UnlockIndex(indexes)
	if ! SkipIndexes {
		fmt.Printf("\nProceding with deletion of indexes for Project named  '%s'...\n", descriptor.Name )
		err := DeleteProjectActionChanges(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		err = DeleteProjectRollBackIndex(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
	}
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
		fmt.Println("No Projects found\n")
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
	Details := false
	Format := "json"
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "show-full" == CorrectInput(option[0]) {
			Details = GetBoolean(option[1])
		} else if "format" == CorrectInput(option[0]) {
			Format = CorrectInput(option[1])
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
	if Details {
		var bytesArray  []byte = make([]byte, 0)
		var err         error
		if "json" == Format {
			bytesArray, err = utils.GetJSONFromElem(project, true)
		} else if "xml" == Format {
			bytesArray, err = utils.GetXMLFromElem(project, true)
		} else {
			response := Response{
				Status: false,
				Message: "Sample Format '"+Format+"' not supported!!",
			}
			return  response, errors.New("Unable to execute task")
		}
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		fmt.Printf("%s\n", bytesArray)
	} else {
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
		fmt.Println("")
		fmt.Println("Changes :")
		changes, err := LoadProjectActionChanges(descriptor.Id)
		if err !=nil {
			fmt.Println("Unable to load changes")
		} else {
			if len(changes.Actions) > 0 {
				for _,change := range changes.Actions {
					isFull := "no"
					if change.FullProject {
						isFull = "yes"
					}
					fmt.Printf("Change: Id: %s - Request: %s - Sub-Request: %s - Element Type: %s - Element Name : %s - Full-Impact : %s\n", change.Id, CmdRequestDescriptors[change.Request], CmdSubRequestDescriptors[change.SubRequest], CmdElementTypeDescriptors[change.ElementType], change.ElementName, isFull )
					fmt.Printf("JSON : %s\n", change.JSONImage )
				}
			} else {
				fmt.Println("No changes available")
			}
		}
	}
	return Response{
		Message: "Success",
		Status: true,
	}, nil
}

func (request *CmdRequest) BuildProject() (Response, error) {
	//response := Response{
	//	Status: false,
	//	Message: "Not Implemented",
	//}
	//return  response, errors.New("Unable to execute task")
	Name := ""
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		}
	}
	fmt.Printf("Unable to complete Rebuild of project %s : Build Project not implemented!!\n", Name)
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, nil
}

func (request *CmdRequest) ImportProject() (Response, error) {
	Name := ""
	File := ""
	Format := "json"
	FullImport := true
	Force := false
	Sample := false
	SampleFormat := "json"
	OverrideInfra := false
	var ElementType CmdElementType = NoElement
	var err error
	for _,option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "file" == CorrectInput(option[0]) {
			File = option[1]
		} else if "format" == CorrectInput(option[0]) {
			Format = option[1]
		} else if "full-import" == CorrectInput(option[0]) {
			FullImport = GetBoolean(option[1])
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "override-infra" == CorrectInput(option[0]) {
			OverrideInfra = GetBoolean(option[1])
		} else if "sample" == CorrectInput(option[0]) {
			Sample = GetBoolean(option[1])
		} else if "sample-format" == CorrectInput(option[0]) {
			SampleFormat = option[1]
		} else if "elem-type" == CorrectInput(option[0]) {
			ElementType, err = CmdParseElement(option[1])
			if err != nil {
				ElementType = NoElement
			}
		}
	}
	if CorrectInput(SampleFormat) == "" {
		SampleFormat = "json"
	}
	if strings.TrimSpace(Name) == "" {
		response := Response{
			Status: false,
			Message: "Project Name Field is mandatory",
		}
		return  response, errors.New("Unable to execute task")
	}
	if ! Sample {
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
	} else {
		if CorrectInput(SampleFormat) != "json" && CorrectInput(SampleFormat) != "xml" {
			response := Response{
				Status: false,
				Message: "Sample Format '"+SampleFormat+"' not supported!!",
			}
			return  response, errors.New("Unable to execute task")
		}
		if FullImport || ElementType == SProject {
			if CorrectInput(SampleFormat) == "json" {
				fmt.Printf("%s\n", utils.GetJSONFromObj(vmio.ProjectSample, true))
			} else {
				fmt.Printf("%s\n", utils.GetXMLFromObj(vmio.ProjectSample, true))
			}
		} else if ElementType != NoElement {
			switch ElementType {
			case SDomain:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(ImportDomainSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(ImportDomainSample, true))
				}
				break
			case SNetwork:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(ImportNetworkSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(ImportNetworkSample, true))
				}
				break
			case LServer:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(ImportLocalServerSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(ImportLocalServerSample, true))
				}
				break
			case CLServer:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(ImportCloudServerSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(ImportCloudServerSample, true))
				}
				break
			default:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(ImportPlansSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(ImportPlansSample, true))
				}
			}
		} else {
			response := Response{
				Status: false,
				Message: "Infrastructure Element not supported!!",
			}
			return  response, errors.New("Unable to execute task")
		}
		response := Response{
			Status: true,
			Message: "Success",
		}
		return  response, nil
	}
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil && ! FullImport && ElementType != SProject {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if ! descriptor.Open {
		response := Response{
			Status: false,
			Message: "Project closed!!",
		}
		return response, errors.New("Unable to execute task")
	}
	
	fmt.Printf("Import File Path : %s, Format: %s\n", File, Format)
	Full := "no"
	if FullImport || ElementType == SProject {
		Full = "yes"
	}
	fmt.Printf("Import Full : %s\n", Full)
	if ! FullImport && ElementType != SProject {
		ImportDomains := "no"
		if ElementType == SDomain {
			ImportDomains = "yes"
		}
		ImportNetworks := "no"
		if ElementType == SNetwork {
			ImportNetworks = "yes"
		}
		ImportServers := "no"
		if ElementType == LServer {
			ImportServers = "yes"
		}
		ImportCServers := "no"
		if ElementType == CLServer {
			ImportCServers = "yes"
		}
		ImportPlans := "no"
		if ElementType == SPlan {
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
	var project model.Project
	Found := false
	ProjectBackup := ""
	
	if ! FullImport && ElementType == NoElement {
		response := Response{
			Status: false,
			Message: "No Full Import and Element Type not supported, nothing to import, VMKube quits!!",
		}
		return  response, errors.New("Unable to execute task")
	}
	
	project, err = vmio.LoadProject(descriptor.Id)
	
	existsProject := false
	existsInfrastructure := descriptor.InfraId != ""
	
	var ProjectJSON string = ""

	if err == nil && project.Id != "" {
		existsProject = true
		Found = true
		ProjectJSON = string(utils.GetJSONFromObj(project, true))
		if FullImport || ElementType == SProject {
			AllowProjectDeletion := Force
			AllowProjectBackup := Force
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
			
			if !AllowProjectDeletion {
				response := Response{
					Status: false,
					Message: "Project '"+Name+"' already exists and no force clause specified ...",
				}
				return  response, errors.New("Unable to execute task")
			}
			
			
			if AllowProjectBackup {
				ProjectBackup = fmt.Sprintf("%s%s.project-%s-%s.json",model.GetEmergencyFolder(),string(os.PathSeparator),utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name))
				project, err := vmio.LoadProject(descriptor.Id)
				if err == nil {
					vmio.ExportUserProject(project,ProjectBackup,"json",true)
					fmt.Printf("Emergency Project backup at : %s\n", ProjectBackup)
				}
			}
			
			request.Arguments.Options = append(request.Arguments.Options, []string{"skip-indexes", "true"})
			response, err := request.DeleteProject()
			if err != nil {
				return response, err
			}
		}
	}
	
	AllowProjectOverwrite := Force
	if existsProject && ! AllowProjectOverwrite {
		AllowProjectOverwrite = utils.RequestConfirmation("Do you want proceed with deletion process for Infrastructure named '"+descriptor.InfraName+"'?")
		if ! AllowProjectOverwrite {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
		
	}
	
	AllowInfraDeletion := OverrideInfra && Force
	AllowInfraBackup := OverrideInfra && Force
	InfraBackup := ""
	
	if descriptor.InfraId != "" && OverrideInfra {
		if ! AllowInfraDeletion {
			AllowInfraDeletion = utils.RequestConfirmation("Do you want proceed with deletion process for Infrastructure named '"+descriptor.InfraName+"'?")
			if ! AllowInfraDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if AllowInfraDeletion && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation("Do you want backup Infrastructure named'"+descriptor.InfraName+"'?")
		}
	}
	
	if existsInfrastructure && OverrideInfra {
		if AllowInfraBackup {
			InfraBackup = fmt.Sprintf("%s%s.prj-%s-%s-infra-export-%s-%s.vmkube",model.GetEmergencyFolder(),string(os.PathSeparator), utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name),utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				fmt.Printf("Emergency Infrastructure backup at : %s\n", InfraBackup)
			}
		}
		if (AllowInfraDeletion) {
			response, err := request.DeleteInfra()
			if err != nil {
				return response, err
			}
		}
	}
	
	
	if FullImport || ElementType == SProject {
		
		project, err = vmio.ImportUserProject(File, Format)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		project.Name = Name
		project.LastMessage = "Project imported from file " + File + ", format " + Format
		errorList := project.Validate()
		if len(errorList) > 0 {
			_, errorValue := vmio.StripErrorMessages("Project import is invalid, clause(s) :", errorList)
			response := Response{
				Status: false,
				Message: errorValue,
			}
			return  response, errors.New("Unable to execute task")
		}
		err = UpdateIndexWithProject(project)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return  response, errors.New("Unable to execute task")
		}
		err = vmio.SaveProject(project)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			UpdateIndexWithProjectsDescriptor(descriptor, false)
			return  response, errors.New("Unable to execute task")
		}
		fmt.Printf("Successfully imported Project '%s' from file '%s' in format '%s'\n", Name, File, Format)
	} else {
		if ! Found {
			response := Response{
				Status: false,
				Message: "Project '"+Name+"' doesn't exist, no suitable project for new elements",
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
					_, errorValue := vmio.StripErrorMessages("Domain import is invalid, clause(s) :", errorList)
					response := Response{
						Status: false,
						Message: errorValue,
					}
					return  response, errors.New("Unable to execute task")
				}
				project.Domains = append(project.Domains, domain)
			}
			if len(domains.Domains) > 0 {
				project.LastMessage = fmt.Sprintf("Domains (no. %d) imported from file %s, format %s", len(domains.Domains), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
				fmt.Printf("Successfully imported Project '%s' Domains from file '%s' in format '%s'\n", Name, File, Format)
			} else {
				fmt.Printf("Warning: No Domain to import from file '%s' in format '%s'\n", File, Format)
			}
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
			for _,networkImport := range networks {
				DomainId := networkImport.Domain.DomainId
				DomainName := networkImport.Domain.DomainName
				Found := false
				for i := 0; i < len(project.Domains); i++ {
					if project.Domains[i].Id == DomainId || project.Domains[i].Name == DomainName {
						Found = true
						for _, network := range networkImport.Networks {
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
							project.Domains[i].Networks = append(project.Domains[i].Networks, network)
						}
					}
				}
				if ! Found {
					response := Response{
						Status: false,
						Message: fmt.Sprintf("Errors during Networks Import : Import target Domain not found by name : '%s' or by id : '%s'", DomainName, DomainId),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			if len(networks) > 0 {
				project.LastMessage = fmt.Sprintf("Networks (no. %d) imported from file %s, format %s", len(networks), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
				fmt.Printf("Successfully imported Project '%s' Networks from file '%s' in format '%s'\n", Name, File, Format)
			} else {
				fmt.Printf("Warning: No Network to import from file '%s' in format '%s'\n", File, Format)
			}
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
			for _,serverImport := range servers {
				DomainId := serverImport.Network.DomainId
				DomainName := serverImport.Network.DomainName
				NetworkId := serverImport.Network.NetworkId
				NetworkName := serverImport.Network.NetworkName
				Found := false
				for i := 0; i < len(project.Domains); i++ {
					if project.Domains[i].Id == DomainId || project.Domains[i].Name == DomainName {
						for j := 0; j < len(project.Domains[i].Networks); j++ {
							if project.Domains[i].Networks[j].Id == NetworkId || project.Domains[i].Networks[j].Name == NetworkName {
								Found = true
								for _, server := range serverImport.Servers {
									server.PostImport()
									errorList := server.Validate()
									if len(errorList) > 0 {
										_, errorValue := vmio.StripErrorMessages("Servers import is invalid, clause(s) :", errorList)
										response := Response{
											Status: false,
											Message: errorValue,
										}
										return  response, errors.New("Unable to execute task")
									}
									project.Domains[i].Networks[j].Servers = append(project.Domains[i].Networks[j].Servers, server)
								}
							}
						}
					}
				}
				if ! Found {
					response := Response{
						Status: false,
						Message: fmt.Sprintf("Errors during Servers Import : Import target Domain not found by name : '%s' or by id : '%s and target Network not found by name : '%s' or by id : '%s '", DomainName, DomainId, NetworkName, NetworkId),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			if len(servers) > 0 {
				project.LastMessage = fmt.Sprintf("Servers (no. %d) imported from file %s, format %s", len(servers), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
				fmt.Printf("Successfully imported Project '%s' Local Servers from file '%s' in format '%s'\n", Name, File, Format)
			} else {
				fmt.Printf("Warning: No Server to import from file '%s' in format '%s'\n", File, Format)
			}
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
			for _,serverImport := range servers {
				DomainId := serverImport.Network.DomainId
				DomainName := serverImport.Network.DomainName
				NetworkId := serverImport.Network.NetworkId
				NetworkName := serverImport.Network.NetworkName
				Found := false
				for i := 0; i < len(project.Domains); i++ {
					if project.Domains[i].Id == DomainId || project.Domains[i].Name == DomainName {
						for j := 0; j < len(project.Domains[i].Networks); j++ {
							if project.Domains[i].Networks[j].Id == NetworkId || project.Domains[i].Networks[j].Name == NetworkName {
								Found = true
								for _, server := range serverImport.Servers {
									server.PostImport()
									errorList := server.Validate()
									if len(errorList) > 0 {
										_, errorValue := vmio.StripErrorMessages("Servers import is invalid, clause(s) :", errorList)
										response := Response{
											Status: false,
											Message: errorValue,
										}
										return  response, errors.New("Unable to execute task")
									}
									project.Domains[i].Networks[j].CServers = append(project.Domains[i].Networks[j].CServers, server)
								}
							}
						}
					}
				}
				if ! Found {
					response := Response{
						Status: false,
						Message: fmt.Sprintf("Errors during Cloud Servers Import : Import target Domain not found by name : '%s' or by id : '%s and target Network not found by name : '%s' or by id : '%s '", DomainName, DomainId, NetworkName, NetworkId),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			if len(servers) > 0 {
				project.LastMessage = fmt.Sprintf("Cloud Servers (no. %d) imported from file %s, format %s", len(servers), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
				fmt.Printf("Successfully imported Project '%s' Cloud Servers from file '%s' in format '%s'\n", Name, File, Format)
			} else {
				fmt.Printf("Warning: No Server to import from file '%s' in format '%s'\n", File, Format)
			}
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
			for _, planImport := range plans {
				DomainId := planImport.Network.DomainId
				DomainName := planImport.Network.DomainName
				NetworkId := planImport.Network.NetworkId
				NetworkName := planImport.Network.NetworkName
				Found := false
				for i := 0; i < len(project.Domains); i++ {
					if project.Domains[i].Id == DomainId || project.Domains[i].Name == DomainName {
						for j := 0; j < len(project.Domains[i].Networks); j++ {
							if project.Domains[i].Networks[j].Id == NetworkId || project.Domains[i].Networks[j].Name == NetworkName {
								Found = true
								for _, plan := range planImport.Plans {
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
									FoundServer := false
									if plan.IsCloud {
										for _,server := range project.Domains[i].Networks[j].CServers {
											if plan.ServerId == server.Id || CorrectInput(plan.ServerId) == CorrectInput(server.Name) {
												FoundServer = true
												plan.ServerId = server.Id
												break
											}
										}
									} else {
										for _,server := range project.Domains[i].Networks[j].Servers {
											if plan.ServerId == server.Id || CorrectInput(plan.ServerId) == CorrectInput(server.Name) {
												FoundServer = true
												plan.ServerId = server.Id
												break
											}
										}
									}
									if ! FoundServer {
										IsCloud := "no"
										if plan.IsCloud {
											IsCloud = "yes"
										}
										response := Response{
											Status: false,
											Message: fmt.Sprintf("Errors during Plans Import : Plan Server id/name: %s from Cloud: %s not found in domain id: %s, name: %s and network id: %s, name: %s", plan.ServerId, IsCloud, DomainName, DomainId, NetworkName, NetworkId),
										}
										return  response, errors.New("Unable to execute task")
									}
									project.Domains[i].Networks[j].Installations = append(project.Domains[i].Networks[j].Installations, plan)
								}
							}
						}
					}
				}
				if ! Found {
					response := Response{
						Status: false,
						Message: fmt.Sprintf("Errors during Plans Import : Import target Domain not found by name : '%s' or by id : '%s and target Network not found by name : '%s' or by id : '%s '", DomainName, DomainId, NetworkName, NetworkId),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			if len(plans) > 0 {
				project.LastMessage = fmt.Sprintf("Installation Plans (no. %d) imported from file %s, format %s", len(plans), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
				fmt.Printf("Successfully imported Project '%s' Installation Plans from file '%s' in format '%s'\n", Name, File, Format)
			} else {
				fmt.Printf("Warning: No Plan to import from file '%s' in format '%s'\n", File, Format)
			}
		}
	}
	if ProjectBackup != "" {
		fmt.Printf("Removing Project backup file : %s\n", ProjectBackup)
		err = model.DeleteIfExists(ProjectBackup)
	}
	if existsProject && existsInfrastructure {
		AddProjectChangeActions(descriptor.Id, ActionDescriptor{
			Id: NewUUIDString(),
			Date: time.Now(),
			DropAction: true,
			ElementType: SProject,
			ElementName: descriptor.Name,
			FullProject: true,
			JSONImage: ProjectJSON,
			Request: request.Type,
			SubRequest: request.SubType,
		})
	}
	
	if existsInfrastructure && OverrideInfra {
		request.Arguments.Options = append(request.Arguments.Options, []string{"rebuild", "true"})
		request.BuildProject()
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
	Format := "json"
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
			ElementType, err = CmdParseElement(option[1])
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
	if FullExport || ElementType == SProject {
		Full = "yes"
	}
	fmt.Printf("Export Full : %s\n", Full)
	if ! FullExport {
		ExportDomains := "no"
		if ElementType == SDomain {
			ExportDomains = "yes"
		}
		ExportNetworks := "no"
		if ElementType == SNetwork {
			ExportNetworks = "yes"
		}
		ExportServers := "no"
		if ElementType == LServer {
			ExportServers = "yes"
		}
		ExportCServers := "no"
		if ElementType == CLServer {
			ExportCServers = "yes"
		}
		ExportPlans := "no"
		if ElementType == SPlan {
			ExportPlans = "yes"
		}
		fmt.Printf("Export Domains : %s\nExport Networks : %s\nExport Servers : %s\nExport Cloud Servers : %s\nExport Plans : %s\n",
			ExportDomains, ExportNetworks, ExportServers, ExportCServers, ExportPlans)
	}
	if ! FullExport && ElementType == NoElement {
		response := Response{
			Status: false,
			Message: "No Full Import and Element Type not supported, nothing to import, VMKube quits!!",
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
			fmt.Printf("Successfully exported Project '%s' Domains to file '%s' in format '%s'\n ", Name, File, Format)
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
			fmt.Printf("Successfully exported Project '%s' Networks to file '%s' in format '%s'\n", Name, File, Format)
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
			fmt.Printf("Successfully exported Project '%s' Local Servers to file '%s' in format '%s'\n", Name, File, Format)
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
			fmt.Printf("Successfully exported Project '%s' Cloud Servers to file '%s' in format '%s'\n", Name, File, Format)
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
			fmt.Printf("Successfully exported Project '%s' Installation Plans to file '%s' in format '%s'\n", Name, File, Format)
		}
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}
