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
	"runtime"
	"vmkube/operations"
	"vmkube/tasks"
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
		for _, option := range request.Arguments.Helper.Options {
			if option.Mandatory {
				//Mandatory Option
				found := false
				for _, argument := range request.Arguments.Options {
					if CorrectInput(argument[0]) == option.Option {
						found = true
						break
					}
				}
				if !found {
					correctness = false
					utils.PrintlnError(fmt.Sprintf("Option '--%s' is mandatory!!", option.Option))
				}
			}
		}
		if !correctness {
			return false
		}
	}
	return true
}

func (request *CmdRequest) CreateProject() (Response, error) {
	utils.NO_COLORS = false
	Name := ""
	InputFile := ""
	Backup := false
	BackupDir := ""
	InputFormat := "json"
	Force := false
	OverrideInfra := false
	for _, option := range request.Arguments.Options {
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
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "backup-dir" == CorrectInput(option[0]) {
			BackupDir = option[1]
		} else if "backup" == CorrectInput(option[0]) {
			Backup = GetBoolean(option[1])
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Project Name not provided",
			Status: false, }, errors.New("Unable to execute task")
	}
	
	AllowProjectDeletion := Force
	AllowInfraDeletion := OverrideInfra && Force
	AllowProjectBackup := Backup && Force
	AllowInfraBackup := Backup && OverrideInfra && Force
	
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
	} else {
		descriptor.Open = true
	}
	
	if err == nil {
		if ! AllowProjectDeletion {
			
			AllowProjectDeletion = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with deletion process for Project named '%s'?", descriptor.Name))
			if ! AllowProjectDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if Backup && AllowProjectDeletion && ! AllowProjectBackup {
			
			AllowProjectBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Project named '%s'?", descriptor.Name))
		}
	}
	
	if err == nil && !AllowProjectDeletion {
		response := Response{
			Status: false,
			Message: "Project '" + Name + "' already exists and no force clause specified ...",
		}
		return response, errors.New("Unable to execute task")
	}
	
	if descriptor.InfraId != "" && existsInfrastructure && OverrideInfra {
		if ! AllowInfraDeletion {
			AllowInfraDeletion = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with deletion process for Infrastructure named '%s'?", descriptor.InfraName))
			if ! AllowInfraDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if Backup && AllowInfraDeletion && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?",descriptor.InfraName))
		}
	}
	
	if descriptor.InfraId != "" && ! AllowInfraDeletion &&  existsInfrastructure && OverrideInfra {
		response := Response{
			Status: false,
			Message: "Project '" + Name + "' already build in Infra '" + descriptor.InfraName + "' and no infrastructure destroy clause specified ...",
		}
		return response, errors.New("Unable to execute task")
	}
	
	existanceClause := "n't"
	existanceClause2 := "proceding with definition of new project"
	if existsProject {
		existanceClause = ""
		existanceClause2 = "proceding with overwrite of existing project"
	}
	utils.PrintlnWarning(fmt.Sprintf("\nProject: %s does%s exist, now %s...", Name, existanceClause, existanceClause2))
	project := model.Project{}
	if InputFile != "" && InputFormat != "" {
		utils.PrintlnWarning(fmt.Sprintf("\nLoading project %s from file '%s' using format '%s'...", Name, InputFile, InputFormat))
		project, err = vmio.ImportUserProject(InputFile, InputFormat)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		project.LastMessage = "Project imported from file " + InputFile + ", format " + InputFormat
		project.Name = Name
	} else {
		utils.PrintlnWarning(fmt.Sprintf("\nDefining new empty project %s...", Name))
		project.Id = NewUUIDString()
		project.Name = Name
		project.Created = time.Now()
		project.Modified = time.Now()
		project.Open = true
		project.LastMessage = "Empty Project Creation"
		project.Domains = append(project.Domains, model.MachineDomain{
			Id: NewUUIDString(),
			Name: "Default Domain",
			Options: [][]string{},
			Networks: []model.MachineNetwork{},
		})
		project.Domains[0].Networks = append(project.Domains[0].Networks, model.MachineNetwork{
			Id: NewUUIDString(),
			Name: "Defaut Network",
			Options: [][]string{},
			CloudMachines: []model.CloudMachine{},
			LocalMachines: []model.LocalMachine{},
			Installations: []model.InstallationPlan{},
		})
	}
	
	if ErrorList := project.Validate(); len(ErrorList) > 0 {
		_, errorValue := vmio.StripErrorMessages("Imported Project is invalid, clause(s) :", ErrorList)
		response := Response{
			Status: false,
			Message: errorValue,
		}
		return response, errors.New("Unable to execute task")
	}
	InfraBackup := ""
	ProjectBackup := ""
	
	if existsInfrastructure {
		if AllowInfraBackup {
			folder := model.GetEmergencyFolder()
			if BackupDir != "" {
				folder = BackupDir
			}
			if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
				folder += string(os.PathSeparator)
			}
			InfraBackup = fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name), utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure backup at : %s", InfraBackup))
			}
		}
		request.Arguments.Options = append(request.Arguments.Options, []string{"infra-name", descriptor.InfraName})
		request.Type = DestroyInfrastructure
		request.TypeStr = "delete-infra"
		request.SubType = NoSubCommand
		request.SubTypeStr = ""
		response, err := request.DeleteInfra()
		if err != nil {
			return response, err
		}
	}
	if existsProject {
		if AllowProjectBackup {
			folder := model.GetEmergencyFolder()
			if BackupDir != "" {
				folder = BackupDir
			}
			if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
				folder += string(os.PathSeparator)
			}
			ProjectBackup = fmt.Sprintf("%s.project-%s-%s.json", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name))
			project, err := vmio.LoadProject(descriptor.Id)
			if err == nil {
				vmio.ExportUserProject(project, ProjectBackup, "json", true)
				utils.PrintlnImportant(fmt.Sprintf("Emergency Project backup at : %s", ProjectBackup))
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
			folder := model.GetEmergencyFolder()
			if BackupDir != "" {
				folder = BackupDir
			}
			if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
				folder += string(os.PathSeparator)
			}
			
			InfraBackup = fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name), utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure backup at : %s", InfraBackup))
			}
		}
		if (AllowInfraDeletion) {
			request.Arguments.Options = append(request.Arguments.Options, []string{"infra-name", descriptor.InfraName})
			request.Type = DestroyInfrastructure
			request.TypeStr = "delete-infra"
			request.SubType = NoSubCommand
			request.SubTypeStr = ""
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
				return response, errors.New("Unable to execute task, Project '" + Name + "' and Infrastructure " + descriptor.InfraName + " no longer exist, no rollback available, check emergency backups in logs!!")
			} else {
				return response, errors.New("Unable to execute task, Project '" + Name + "' no longer exist, no rollback available, check emergency backup in logs!!")
			}
		} else {
			return response, errors.New("Unable to execute task")
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
		return response, errors.New("Unable to execute task")
	}
	
	if ProjectBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project backup at : %s", ProjectBackup))
		os.Remove(ProjectBackup)
	}
	
	if InfraBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Infrastructure backup at : %s", InfraBackup))
		os.Remove(InfraBackup)
	}
	
	if existsProject && existsInfrastructure && OverrideInfra {
		AddProjectChangeActions(descriptor.Id, ActionDescriptor{
			Id: NewUUIDString(),
			Date: time.Now(),
			DropAction: true,
			ElementType: SProject,
			ElementId: descriptor.Id,
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
	return response, nil
}

func (request *CmdRequest) AlterProject() (Response, error) {
	utils.NO_COLORS = false
	Name := ""
	File := ""
	Format := "json"
	Force := true
	OverrideInfra := false
	BackupDir := ""
	Backup := false
	var ElementType CmdElementType = NoElement
	ElementName := ""
	ElementId := ""
	var AnchorElementType CmdElementType = NoElement
	AnchorElementName := ""
	AnchorElementId := ""
	Sample := false
	SampleFormat := "json"
	var err error
	for _, option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
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
		} else if "sample" == CorrectInput(option[0]) {
			Sample = GetBoolean(option[1])
		} else if "sample-format" == CorrectInput(option[0]) {
			SampleFormat = option[1]
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "backup-dir" == CorrectInput(option[0]) {
			BackupDir = option[1]
		} else if "backup" == CorrectInput(option[0]) {
			Backup = GetBoolean(option[1])
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
		return response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(File) == "" && ! Sample && request.SubType != Open && request.SubType != Close && request.SubType != Remove {
		response := Response{
			Status: false,
			Message: "Input File Path Field is mandatory in Alter Project  Create or Alter operations",
		}
		return response, errors.New("Unable to execute task")
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
		return response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(ElementName) == "" && strings.TrimSpace(ElementId) == "" && ! Sample && request.SubType != Open && request.SubType != Close {
		response := Response{
			Status: false,
			Message: "Element Name Field or Id is mandatory, use project-import for massive changes",
		}
		return response, errors.New("Unable to execute task")
	}
	
	if Sample && request.SubType != Open && request.SubType != Close {
		if CorrectInput(SampleFormat) != "json" && CorrectInput(SampleFormat) != "xml" {
			response := Response{
				Status: false,
				Message: "Sample Format '" + SampleFormat + "' not supported!!",
			}
			return response, errors.New("Unable to execute task")
		}
		if ElementType == SProject {
			if CorrectInput(SampleFormat) == "json" {
				fmt.Printf("%s\n", utils.GetJSONFromObj(vmio.ProjectSample, true))
			} else {
				fmt.Printf("%s\n", utils.GetXMLFromObj(vmio.ProjectSample, true))
			}
		} else if ElementType != NoElement {
			switch ElementType {
			case SDomain:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(vmio.DomainSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(vmio.DomainSample, true))
				}
				break
			case SNetwork:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(vmio.NetworkSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(vmio.NetworkSample, true))
				}
				break
			case LMachine:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(vmio.MachineSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(vmio.MachineSample, true))
				}
				break
			case CLMachine:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(vmio.CloudMachineSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(vmio.CloudMachineSample, true))
				}
				break
			default:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(vmio.InstallationPlanSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(vmio.InstallationPlanSample, true))
				}
			}
		} else {
			response := Response{
				Status: false,
				Message: "Infrastructure Element not supported!!",
			}
			return response, errors.New("Unable to execute task")
		}
		response := Response{
			Status: true,
			Message: "Success",
		}
		return response, nil
	}
	
	if AnchorElementType == NoElement && request.SubType == Create {
		response := Response{
			Status: false,
			Message: "Anchor Element Type Field is mandatory in Alter Project Create command",
		}
		return response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(AnchorElementName) == "" && strings.TrimSpace(AnchorElementId) == ""  && request.SubType == Create {
		response := Response{
			Status: false,
			Message: "Anchor Element Name or Id Fields is mandatory in Alter Project Create command",
		}
		return response, errors.New("Unable to execute task")
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
			return response, errors.New("Unable to execute task")
		default:
			// Any but Project Type
			Id, err := operations.AddElementToProject(project, int(ElementType), ElementName, int(AnchorElementType), AnchorElementName, AnchorElementId, File, Format)
			ElemRealName := ""
			Names := strings.Split(Id, ":")
			Id = Names[0]
			ElemRealName = Names[1]
			SubId := ""
			Ids := strings.Split(Id, ":")
			if len(Ids) > 1 {
				Id = Ids[0]
				SubId = Ids[1]
			}
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}
			
			
			if existsInfrastructure {
				AddProjectChangeActions(descriptor.Id, ActionDescriptor{
					Id: NewUUIDString(),
					Date: time.Now(),
					DropAction: false,
					ElementType: ElementType,
					ElementId: Id,
					ElementName: ElemRealName,
					RelatedId: SubId,
					FullProject: false,
					JSONImage: ProjectJSON,
					Request: request.Type,
					SubRequest: request.SubType,
				})
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
			return response, errors.New("Unable to execute task")
		default:
			// Any but Project Type
			Id, err := operations.AlterElementInProject(project, int(ElementType), ElementName, ElementId, File, Format)
			ElemRealName := ""
			Names := strings.Split(Id, ":")
			Id = Names[0]
			ElemRealName = Names[1]
			SubId := ""
			Ids := strings.Split(Id, ":")
			if len(Ids) > 1 {
				Id = Ids[0]
				SubId = Ids[1]
			}
			
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}
			
			if existsInfrastructure {
				AddProjectChangeActions(descriptor.Id, ActionDescriptor{
					Id: NewUUIDString(),
					Date: time.Now(),
					DropAction: false,
					ElementType: ElementType,
					ElementId: Id,
					ElementName: ElemRealName,
					RelatedId: SubId,
					FullProject: false,
					JSONImage: ProjectJSON,
					Request: request.Type,
					SubRequest: request.SubType,
				})
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
			return response, errors.New("Unable to execute task")
		default:
			// Any but Project Type
			Id, err := operations.DeleteElementInProject(project, int(ElementType), ElementName, ElementId)
			ElemRealName := ""
			Names := strings.Split(Id, ":")
			Id = Names[0]
			ElemRealName = Names[1]
			SubId := ""
			Ids := strings.Split(Id, ":")
			if len(Ids) > 1 {
				Id = Ids[0]
				SubId = Ids[1]
			}

			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}

			if existsInfrastructure {
				AddProjectChangeActions(descriptor.Id, ActionDescriptor{
					Id: NewUUIDString(),
					Date: time.Now(),
					DropAction: true,
					ElementType: ElementType,
					ElementId: Id,
					ElementName: ElemRealName,
					RelatedId: SubId,
					FullProject: false,
					JSONImage: ProjectJSON,
					Request: request.Type,
					SubRequest: request.SubType,
				})
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
			return response, errors.New("Unable to execute task")
		}
		descriptor.Active = !project.Open
		break
	case Close:
		project, err = operations.CloseProject(project)
		
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		descriptor.Active = !project.Open
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
		AllowProjectOverwrite = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with alter process for Project named '%s'?", descriptor.Name))
		if ! AllowProjectOverwrite {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
		
	}
	
	AllowInfraDeletion := OverrideInfra && Force
	AllowInfraBackup := Backup && OverrideInfra && Force
	InfraBackup := ""
	
	if descriptor.InfraId != "" && OverrideInfra {
		if ! AllowInfraDeletion {
			AllowInfraDeletion = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with deletion process for Infrastructure named '%s'?", descriptor.InfraName))
			if ! AllowInfraDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if Backup && AllowInfraDeletion && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?",descriptor.InfraName))
		}
	}
	
	if existsInfrastructure && OverrideInfra {
		if AllowInfraBackup {
			folder := model.GetEmergencyFolder()
			if BackupDir != "" {
				folder = BackupDir
			}
			if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
				folder += string(os.PathSeparator)
			}
			InfraBackup = fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name), utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure backup at : %s", InfraBackup))
			}
		}
		if (AllowInfraDeletion) {
			request.Arguments.Options = append(request.Arguments.Options, []string{"infra-name", descriptor.InfraName})
			request.Type = DestroyInfrastructure
			request.TypeStr = "delete-infra"
			request.SubType = NoSubCommand
			request.SubTypeStr = ""
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
			return response, errors.New("Unable to execute task, Infrastructure " + descriptor.InfraName + " no longer exist, no rollback available, check emergency backups in logs!!")
		} else {
			return response, errors.New("Unable to execute task")
		}
	}
	
	if InfraBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Infrastructure backup at : %s", InfraBackup))
		os.Remove(InfraBackup)
	}
	
	if existsInfrastructure {
		request.Arguments.Options = append(request.Arguments.Options, []string{"rebuild", "true"})
		request.BuildProject()
	}
	
	if existsProject {
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
			return response, errors.New("Unable to execute task")
		}
	}
	
	if existsInfrastructure && OverrideInfra {
		request.Arguments.Options = append(request.Arguments.Options, []string{"rebuild", "true"})
		request.BuildProject()
	}
	
	utils.PrintlnSuccess("Alter Project : Command Executed correctly!!")
	response := Response{
		Status: true,
		Message: "Success",
	}
	return response, nil
	
}

func (request *CmdRequest) InfoProject() (Response, error) {
	utils.NO_COLORS = false
	if request.SubType == List {
		//List of elements
		defines := vmio.ListProjectTypeDefines()
		for _, define := range defines {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", utils.StrPad(define.Name, 15), define.Description)
		}
	} else {
		//Detail
		TypeVal := ""
		Sample := ""
		for _, option := range request.Arguments.Options {
			if "sample" == CorrectInput(option[0]) {
				Sample = CorrectInput(option[1])
			} else if "elem-type" == CorrectInput(option[0]) {
				TypeVal = option[1]
			} else if "no-colors" == CorrectInput(option[0]) {
				utils.NO_COLORS = GetBoolean(option[1])
			}
		}
		if TypeVal == "" {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr)
			return Response{
				Message: "Element Type not provided",
				Status: false, }, nil
		} else {
			if Sample != "" {
				if "json" == Sample || "xml" == Sample {
					defines := vmio.ListProjectTypeDefines()
					for _, define := range defines {
						if define.Name == TypeVal {
							if "json" == Sample {
								bytes, err := utils.GetJSONFromElem(define.Sample, true)
								if err != nil {
									utils.PrintlnError(fmt.Sprintf("Error: Output format '%s' not provided", Sample))
									PrintCommandHelper(request.TypeStr, request.SubTypeStr)
									return Response{
										Message: "",
										Status: true, }, nil
								}
								fmt.Fprintf(os.Stdout, "%s\n", bytes)
								return Response{
									Message: "",
									Status: true, }, nil
							} else {
								bytes, err := utils.GetXMLFromElem(define.Sample, true)
								if err != nil {
									utils.PrintlnError(fmt.Sprintf("Error: Output format '%s' not provided", Sample))
									PrintCommandHelper(request.TypeStr, request.SubTypeStr)
									return Response{
										Message: "",
										Status: true, }, nil
								}
								fmt.Fprintf(os.Stdout, "%s\n", bytes)
								return Response{
									Message: "",
									Status: true, }, nil
							}
						}
					}
					utils.PrintlnError(fmt.Sprintf("Project Type '%s' not provided", TypeVal))
					PrintCommandHelper(request.TypeStr, request.SubTypeStr)
					return Response{
						Message: "",
						Status: true, }, nil
				}
				utils.PrintlnError(fmt.Sprintf("Error: Output format '%s' not provided", Sample))
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return Response{
					Message: "",
					Status: true, }, nil
			} else {
				defines := vmio.ListProjectTypeDefines()
				for _, define := range defines {
					if define.Name == TypeVal {
						fields, err := model.DescribeStruct(define.Sample)
						if err == nil {
							model.PrintFieldsHeader(len(fields) > 0)
							model.PrintFieldsRecursively(fields, 0)
							return Response{
								Message: "",
								Status: true, }, nil
						}
						utils.PrintlnError(fmt.Sprintf("Project Type '%s' not provided", TypeVal))
						PrintCommandHelper(request.TypeStr, request.SubTypeStr)
						return Response{
							Message: "",
							Status: true, }, nil
					}
				}
				utils.PrintlnError(fmt.Sprintf("Project Type '%s' not provided", TypeVal))
				PrintCommandHelper(request.TypeStr, request.SubTypeStr)
				return Response{
					Message: "",
					Status: true, }, nil
			}
			return Response{
				Message: "",
				Status: true, }, nil
		}
	}
	response := Response{Status: true}
	return response, nil
}

func (request *CmdRequest) DeleteProject() (Response, error) {
	utils.NO_COLORS = false
	Name := ""
	Force := false
	SkipIndexes := false
	for _, option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "skip-indexes" == CorrectInput(option[0]) {
			SkipIndexes = GetBoolean(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
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
				AllowProjectDeletion = utils.RequestConfirmation(fmt.Sprintf("Do you want delete Project named '%s'?", descriptor.Name))
			} else {
				AllowProjectDeletion = utils.RequestConfirmation(fmt.Sprintf("Do you want delete Project named '%s' and Infrastructure named '%s'?", descriptor.Name, descriptor.InfraName))
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
		existanceClause = " and proceding with deletion of existing Infrastructure named '" + descriptor.InfraName + "'"
	}
	utils.PrintlnWarning(fmt.Sprintf("\nProceding with deletion of Project named  '%s'%s...", descriptor.Name, existanceClause))
	
	indexes, err := vmio.LoadIndex()
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	if existsInfrastructure {
		request.Arguments.Options = append(request.Arguments.Options, []string{"infra-name", descriptor.InfraName})
		request.Type = DestroyInfrastructure
		request.TypeStr = "delete-infra"
		request.SubType = NoSubCommand
		request.SubTypeStr = ""
		resp, err := request.DeleteInfra()
		if err != nil {
			return resp, err
		}
	}
	
	actionIndex, err := LoadProjectActionIndex(descriptor.Id)

	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	actionIndexMeta := IFaceProjectActionIndex{
		ProjectId: actionIndex.ProjectId,
		Id: actionIndex.Id,
	}
	
	actionIndexMeta.WaitForUnlock()
	
	actionInfo := ProjectActionIndexInfo{
		Format: "",
		Index: actionIndex,
	}
	
	err = actionInfo.Delete()
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
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
		return response, errors.New("Unable to execute task")
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
		return response, errors.New("Unable to execute task")
	}
	vmio.UnlockIndex(indexes)
	if ! SkipIndexes {
		utils.PrintlnWarning(fmt.Sprintf("\nProceding with deletion of indexes for Project named '%s'...", descriptor.Name))
		err := DeleteProjectActionChanges(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		err = DeleteProjectRollBackIndex(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
	}
	utils.PrintlnSuccess(fmt.Sprintf("Project %s deleted successfully!!", Name))
	response := Response{
		Status: true,
		Message: "Success",
	}
	
	return response, nil
}

func (request *CmdRequest) ListProjects() (Response, error) {
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
		if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		}
	}
	indexes, error := vmio.LoadIndex()
	if error != nil {
		response := Response{
			Status: false,
			Message: error.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if len(indexes.Projects) > 0 {
		utils.PrintlnImportant(fmt.Sprintf("%s  %s  %s  %s  %s  %s", utils.StrPad("Project Id", 40), utils.StrPad("Project Name", 40), utils.StrPad("Open", 4), utils.StrPad("Infrastructure Name", 40), utils.StrPad("Active", 6), utils.StrPad("Synced", 6)))
	} else {
		utils.PrintlnImportant("No Projects found")
	}
	for _, index := range indexes.Projects {
		open := "no"
		active := "no"
		synced := "no"
		if index.Open {
			open = "yes"
		}
		if index.Active {
			active = "yes"
		}
		if index.Synced {
			synced = "yes"
		}
		fmt.Printf("%s  %s  %s  %s  %s  %s\n", utils.StrPad(index.Id, 40), utils.StrPad(index.Name, 40), utils.StrPad(open, 4), utils.StrPad(index.InfraName, 40), utils.StrPad("  " + active, 6), utils.StrPad("  " + synced, 6) )
		
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) StatusProject() (Response, error) {
	utils.NO_COLORS = false
	Name := ""
	Details := false
	Format := "json"
	for _, option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "show-full" == CorrectInput(option[0]) {
			Details = GetBoolean(option[1])
		} else if "format" == CorrectInput(option[0]) {
			Format = CorrectInput(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Project Name not provided",
			Status: false, }, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
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
		return response, errors.New("Unable to execute task")
	}
	if Details {
		var bytesArray  []byte = make([]byte, 0)
		var err error
		if "json" == Format {
			bytesArray, err = utils.GetJSONFromElem(project, true)
		} else if "xml" == Format {
			bytesArray, err = utils.GetXMLFromElem(project, true)
		} else {
			response := Response{
				Status: false,
				Message: "Sample Format '" + Format + "' not supported!!",
			}
			return response, errors.New("Unable to execute task")
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
		fmt.Printf("Id: %s\nProject: %s\nOpen: %s\n", project.Id, project.Name, open)
		fmt.Printf("Created : %d-%02d-%02d %02d:%02d:%02d\n",
			project.Created.Year(), project.Created.Month(), project.Created.Day(),
			project.Created.Hour(), project.Created.Minute(), project.Created.Second())
		fmt.Printf("Modified : %d-%02d-%02d %02d:%02d:%02d\n",
			project.Modified.Year(), project.Modified.Month(), project.Modified.Day(),
			project.Modified.Hour(), project.Modified.Minute(), project.Modified.Second())
		if descriptor.InfraId != "" {
			fmt.Printf("Infrastructure: %s (Id: %s)\n", descriptor.InfraName, descriptor.InfraId)
		} else {
			fmt.Println("Infrastructure: No Infrastructure\n")
		}
		fmt.Printf("Errors: %s\nLast Message: %s\n", errors, project.LastMessage)
		fmt.Printf("Domains: %d\n", len(project.Domains))
		for _, domain := range project.Domains {
			num, options := vmio.StripOptions(domain.Options)
			fmt.Printf("Domain: %s (Id: %s) - Options [%d] :%s\n", domain.Name, domain.Id, num, options)
			fmt.Printf("Networks: %d\n", len(domain.Networks))
			for _, network := range domain.Networks {
				num, options := vmio.StripOptions(network.Options)
				fmt.Printf("   Network: %s (Id: %s) - Options [%d] :%s\n", network.Name, network.Id, num, options)
				fmt.Printf("   Local Machines: %d\n", len(network.LocalMachines))
				machinesMap := make(map[string]string)
				for _, machine := range network.LocalMachines {
					machinesMap[machine.Id] = machine.Name
					fmt.Printf("      Local Machine: %s (Id: %s) - Driver: %s - OS : %s:%s\n", machine.Name, machine.Id, machine.Driver, machine.OSType, machine.OSVersion)
				}
				fmt.Printf("   Cloud Machines: %d\n", len(network.CloudMachines))
				for _, machine := range network.CloudMachines {
					machinesMap[machine.Id] = machine.Name
					num, options := vmio.StripOptions(machine.Options)
					fmt.Printf("      Cloud Machine: %s (Id: %s) - Driver: %s - Options [%d] :%s\n", machine.Name, machine.Id, machine.Driver, num, options)
				}
				fmt.Printf("   Installation Plans: %d\n", len(network.Installations))
				for _, installation := range network.Installations {
					machineName, ok := machinesMap[installation.MachineId]
					if !ok {
						machineName = "<invalid>"
					}
					cloud := "no"
					if installation.IsCloud {
						cloud = "yes"
					}
					fmt.Printf("      Plan: Id: %s - Machine: %s [Id: %s] - Cloud: %s - Envoronment : %s  Role: %s  Type: %s\n", installation.Id, machineName, installation.MachineId, cloud, installation.Environment, installation.Role, installation.Type)
				}
			}
		}
		fmt.Println("")
		fmt.Println("Changes :")
		changes, err := LoadProjectActionChanges(descriptor.Id)
		if err != nil {
			fmt.Println("Unable to load changes")
		} else {
			if len(changes.Actions) > 0 {
				for _, change := range changes.Actions {
					isFull := "no"
					if change.FullProject {
						isFull = "yes"
					}
					fmt.Printf("Change: Id: %s - Request: %s - Sub-Request: %s - Element Type: %s - Element Name : %s (Id: %s) - Full-Impact : %s\n", change.Id, CmdRequestDescriptors[change.Request], CmdSubRequestDescriptors[change.SubRequest], CmdElementTypeDescriptors[change.ElementType], change.ElementName, change.ElementId, isFull)
					fmt.Printf("JSON : %s\n", change.JSONImage)
				}
			} else {
				utils.PrintlnImportant("No changes available")
			}
		}
	}
	return Response{
		Message: "Success",
		Status: true,
	}, nil
}

func (request *CmdRequest) BuildProject() (Response, error) {
	Name := ""
	InfraName := ""
	Backup := false
	BackupDir := ""
	Force := false
	Rebuild := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "infra-name" == CorrectInput(option[0]) {
			InfraName = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "rebuild" == CorrectInput(option[0]) {
			Rebuild = GetBoolean(option[1])
		} else if "overclock" == CorrectInput(option[0]) {
			Overclock = GetBoolean(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "threads" == CorrectInput(option[0]) {
			Threads = GetInteger(option[1], Threads)
		} else if "backup-dir" == CorrectInput(option[0]) {
			BackupDir = option[1]
		} else if "backup" == CorrectInput(option[0]) {
			Backup = GetBoolean(option[1])
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Project Name not provided",
			Status: false, }, errors.New("Unable to execute task")
	}
	if InfraName == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastructure Name not provided",
			Status: false,
		}, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	infraDescriptor, err := vmio.GetInfrastructureProjectDescriptor(InfraName)
	if err == nil {
		if infraDescriptor.Id != descriptor.Id {
			return Response{
				Message: fmt.Sprintf("Infrastructure Name already used in Project named %s", infraDescriptor.Name),
				Status: false,
			}, errors.New("Unable to execute task")
		}
	}
	existsInfrastructure := descriptor.InfraId != ""
	ForceRebuild := Force && existsInfrastructure
	AllowInfraBackup := Backup && Force && existsInfrastructure
	InfraBackup := ""
	
	project, err := vmio.LoadProject(descriptor.Id)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	ValidProject := false
	for _, domain := range project.Domains {
		if ValidProject {
			continue
		}
		for _, network := range domain.Networks {
			if len(network.LocalMachines) > 0 {
				ValidProject = true
			} else if len(network.CloudMachines) > 0 {
				ValidProject = true
			}
		}
	}
	
	if !ValidProject {
		response := Response{
			Status: false,
			Message: "Project not valid, please define some machines and eventually plans before to build it!!",
		}
		return response, errors.New("Unable to execute task")
	}
	
	Infrastructure, err := ProjectToInfrastructure(project)
	Infrastructure.Name = InfraName
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnImportant("Infrastructure JSON:")
	fmt.Printf("%s\n", string(utils.GetJSONFromObj(Infrastructure, true)))
	
	if existsInfrastructure {
		if ! Rebuild {
			response := Response{
				Status: false,
				Message: fmt.Sprintf("Infrastructure Named : %s already exists and no rebuild clause provided!!", descriptor.InfraName),
			}
			return response, errors.New("Unable to execute task")
		} else if ! ForceRebuild {
			ForceRebuild = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with override process for existing Infrastructure named '%s'?",descriptor.InfraName))
			if ! ForceRebuild {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		
		if Backup && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?",descriptor.InfraName))
		}
		
	}
	
	if Rebuild {
		if AllowInfraBackup {
			folder := model.GetEmergencyFolder()
			if BackupDir != "" {
				folder = BackupDir
			}
			if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
				folder += string(os.PathSeparator)
			}
			InfraBackup = fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name), utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure backup at : %s", InfraBackup))
			}
		}
		request.Arguments.Options = append(request.Arguments.Options, []string{"infra-name", descriptor.InfraName})
		request.Type = DestroyInfrastructure
		request.TypeStr = "delete-infra"
		request.SubType = NoSubCommand
		request.SubTypeStr = ""
		response, err := request.DeleteInfra()
		if err != nil {
			return response, err
		}
	}
	
	var actionIndex ProjectActionIndex
	actionCouples, err := make([]tasks.ActivityCouple, 0), errors.New("Unknown Error")
	if ! existsInfrastructure {
		creationCouples, err := tasks.GetTaskActivities(project, Infrastructure, tasks.CreateMachine)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		
		inspectCouples, err := tasks.GetTaskActivities(project, Infrastructure, tasks.MachineInspect)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		ipAddressCouples, err := tasks.GetTaskActivities(project, Infrastructure, tasks.MachineIPAddress)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		stopCouples, err := tasks.GetTaskActivities(project, Infrastructure, tasks.StopMachine)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		extendsDiskCouples, err := tasks.GetTaskActivities(project, Infrastructure, tasks.MachineExtendsDisk)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		actionCouples = append(actionCouples, creationCouples...)
		actionCouples = append(actionCouples, inspectCouples...)
		actionCouples = append(actionCouples, ipAddressCouples...)
		actionCouples = append(actionCouples, stopCouples...)
		actionCouples = append(actionCouples, extendsDiskCouples...)
	} else {
		var exclusionIdList []string = make([]string, 0)
		creationCouples, err := tasks.GetPostBuildTaskActivities(Infrastructure, tasks.CreateMachine, exclusionIdList)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		actionIndex, err = LoadProjectActionIndex(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		// If no changes in the action logs nothing to do
		if len(actionIndex.Actions) == 0 {
			response := Response{
				Status: false,
				Message: fmt.Sprintf("No changes for Project '%s' Infrastructure '%s'!!",descriptor.Name,descriptor.InfraName),
			}
			return response, errors.New("Unable to execute task")
		}
		// If no effective DDL into action logs nothing to do
		if len(creationCouples) == 0 {
			response := Response{
				Status: false,
				Message: fmt.Sprintf("No effective changes for Project '%s' Infrastructure '%s' to justify a build!!",descriptor.Name,descriptor.InfraName),
			}
			return response, errors.New("Unable to execute task")
		}
		exclusionIdList, creationCouples = FilterCreationBasedOnProjectActions(actionIndex, creationCouples)
		inspectCouples, err := tasks.GetPostBuildTaskActivities(Infrastructure, tasks.MachineInspect, exclusionIdList)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		ipAddressCouples, err := tasks.GetPostBuildTaskActivities(Infrastructure, tasks.MachineIPAddress, exclusionIdList)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		stopCouples, err := tasks.GetPostBuildTaskActivities(Infrastructure, tasks.StopMachine, exclusionIdList)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		extendsDiskCouples, err := tasks.GetPostBuildTaskActivities(Infrastructure, tasks.MachineExtendsDisk, exclusionIdList)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		actionCouples = append(actionCouples, creationCouples...)
		actionCouples = append(actionCouples, inspectCouples...)
		actionCouples = append(actionCouples, ipAddressCouples...)
		actionCouples = append(actionCouples, stopCouples...)
		actionCouples = append(actionCouples, extendsDiskCouples...)
		request.StopInfra()
	}
	
	utils.PrintlnImportant("Now Proceding with machine creation ...!!")
	NumThreads := Threads
	if runtime.NumCPU() - 1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	
	var errorsList []error = make([]error, 0)
	var fixInfraValue int = len(actionCouples)
	utils.PrintlnImportant(fmt.Sprintf("Number of scheduled processes : %d", fixInfraValue))

	errorsList = ExecuteInfrastructureActions(Infrastructure, actionCouples, NumThreads,func(task tasks.ScheduleTask){
		go func(task tasks.ScheduleTask) {
			for i := 0; i < len(task.Jobs); i++ {
				response := strings.Split(fmt.Sprintf("%s",task.Jobs[i].GetRunnable().Response()),"|")
				if len(response) > 3 {
					if len(response) > 4 {
						if response[0] == "ip" {
							instanceId := response[1]
							ipAddress := response[2]
							json := ""
							log := response[3] + "\n" + response[4]
							FixInfrastructureElementValue(&Infrastructure, instanceId, ipAddress, json, log)
						} else if response[0] == "json" {
							instanceId := response[1]
							ipAddress := ""
							json := response[2]
							log := response[3] + "\n" + response[4]
							FixInfrastructureElementValue(&Infrastructure, instanceId, ipAddress, json, log)
						}
					} else {
						instanceId := response[1]
						ipAddress := ""
						json := ""
						log := response[2] + "\n" + response[3]
						FixInfrastructureElementValue(&Infrastructure, instanceId, ipAddress, json, log)
					}
				}
				fixInfraValue--
			}
		}(task)
	})
	
	reOpt := ""
	if existsInfrastructure {
		reOpt = "Re"
	}

	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete %sBuild of project '%s' : Errors building Infrastructure : '%s'!!", reOpt, Name, InfraName))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error building infrastructure -> '%s' : ", InfraName), errorsList)
		response := Response{
			Status: false,
			Message: message,
		}
		exclusionList, _ := FilterForExistState(Infrastructure)
		rollbackActions, err := tasks.GetPostBuildTaskActivities(Infrastructure, tasks.DestroyMachine, exclusionList)

		if err != nil {
			return response, errors.New("Unable to execute task")
		}
		utils.PrintlnWarning(fmt.Sprintf("Executing rollback for Project '%s' Infrastrcucture '%s'...", Name, InfraName))
		if existsInfrastructure && AllowInfraBackup {
			utils.PrintlnImportant("Check logs for backup activities, you can use for recovery...")
		} else {
			utils.PrintlnImportant("No backup activities, you can not use recovery utils...")
		}
		if ! utils.NO_COLORS {
			time.Sleep(4*time.Second)
		}
		ExecuteInfrastructureActions(Infrastructure, rollbackActions, NumThreads,func(task tasks.ScheduleTask){})
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Waiting for Instance recovery information in Project '%s' Infrastrcucture '%s'", Name, InfraName))
	
	for fixInfraValue > 0 {
		time.Sleep(1*time.Second)
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Update Project '%s' indexes linking New Infrastrcucture '%s'...", Name, InfraName))
	

	err = UpdateIndexWithInfrastructure(Infrastructure)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Saving new Project '%s' Infrastrcucture '%s' descriptors...", Name, InfraName))

	vmio.LockInfrastructureById(descriptor.Id, descriptor.InfraId)
	
	err = vmio.SaveInfrastructure(Infrastructure)
	
	vmio.UnlockInfrastructureById(descriptor.Id, descriptor.InfraId)

	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if existsInfrastructure {
		utils.PrintlnWarning(fmt.Sprintf("Migrating Project '%s' Infrastrcucture '%s' action logs to Rollback Segments...", Name, InfraName))
		err = MigrateProjectActionsToRollbackSegments(actionIndex)

		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
	}
	
	utils.PrintlnSuccess(fmt.Sprintf("Project '%s' Infrastrcucture '%s' %sBuild successful!!", Name, InfraName, reOpt))

	if InfraBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Infrastructure backup file : %s", InfraBackup))
		err = model.DeleteIfExists(InfraBackup)
	}


	response := Response{
		Status: true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) ImportProject() (Response, error) {
	utils.NO_COLORS = false
	Name := ""
	File := ""
	Backup := false
	BackupDir := ""
	Format := "json"
	FullImport := true
	Force := false
	Sample := false
	SampleFormat := "json"
	OverrideInfra := false
	var ElementType CmdElementType = NoElement
	var err error
	for _, option := range request.Arguments.Options {
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
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "backup-dir" == CorrectInput(option[0]) {
			BackupDir = option[1]
		} else if "backup" == CorrectInput(option[0]) {
			Backup = GetBoolean(option[1])
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
		return response, errors.New("Unable to execute task")
	}
	if ! Sample {
		if strings.TrimSpace(File) == "" {
			response := Response{
				Status: false,
				Message: "Import File Path Field is mandatory",
			}
			return response, errors.New("Unable to execute task")
		}
		if strings.TrimSpace(Format) == "" {
			response := Response{
				Status: false,
				Message: "Import Format Field is mandatory",
			}
			return response, errors.New("Unable to execute task")
		}
	} else {
		if CorrectInput(SampleFormat) != "json" && CorrectInput(SampleFormat) != "xml" {
			response := Response{
				Status: false,
				Message: "Sample Format '" + SampleFormat + "' not supported!!",
			}
			return response, errors.New("Unable to execute task")
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
			case LMachine:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(ImportLocalMachineSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(ImportLocalMachineSample, true))
				}
				break
			case CLMachine:
				if CorrectInput(SampleFormat) == "json" {
					fmt.Printf("%s\n", utils.GetJSONFromObj(ImportCloudMachineSample, true))
				} else {
					fmt.Printf("%s\n", utils.GetXMLFromObj(ImportCloudMachineSample, true))
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
			return response, errors.New("Unable to execute task")
		}
		response := Response{
			Status: true,
			Message: "Success",
		}
		return response, nil
	}
	descriptor, err := vmio.GetProjectDescriptor(Name)
	if err != nil && ! FullImport && ElementType != SProject {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if descriptor.Id !="" && !descriptor.Open {
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
		ImportMachines := "no"
		if ElementType == LMachine {
			ImportMachines = "yes"
		}
		ImportCMachines := "no"
		if ElementType == CLMachine {
			ImportCMachines = "yes"
		}
		ImportPlans := "no"
		if ElementType == SPlan {
			ImportPlans = "yes"
		}
		fmt.Printf("Import Domains : %s\nImport Networks : %s\nImport Machines : %s\nImport Cloud Machines : %s\nImport Plans : %s\n",
			ImportDomains, ImportNetworks, ImportMachines, ImportCMachines, ImportPlans)
	}
	if ! FullImport && ElementType == NoElement {
		response := Response{
			Status: false,
			Message: "Nothing to import, application quits ...",
		}
		return response, errors.New("Unable to execute task")
	}
	var project model.Project
	Found := false
	ProjectBackup := ""
	
	if ! FullImport && ElementType == NoElement {
		response := Response{
			Status: false,
			Message: "No Full Import and Element Type not supported, nothing to import, VMKube quits!!",
		}
		return response, errors.New("Unable to execute task")
	}
	
	project, err = vmio.LoadProject(descriptor.Id)
	
	existsProject := descriptor.Id != ""
	existsInfrastructure := existsProject && descriptor.InfraId != ""
	
	var ProjectJSON string = ""
	
	if err == nil && project.Id != "" {
		existsProject = true
		Found = true
		ProjectJSON = string(utils.GetJSONFromObj(project, true))
		if FullImport || ElementType == SProject {
			AllowProjectDeletion := Force
			AllowProjectBackup := Backup && Force
			if err == nil {
				if ! AllowProjectDeletion {
					AllowProjectDeletion = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with deletion process for Project named '%s'?", descriptor.Name))
					if ! AllowProjectDeletion {
						response := Response{
							Status: false,
							Message: "User task interruption",
						}
						return response, errors.New("Unable to execute task")
					}
				}
				if Backup && AllowProjectDeletion && ! AllowProjectBackup {
					AllowProjectBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Project named '%s'?", descriptor.Name))
				}
			}
			
			if !AllowProjectDeletion {
				response := Response{
					Status: false,
					Message: "Project '" + Name + "' already exists and no force clause specified ...",
				}
				return response, errors.New("Unable to execute task")
			}
			
			if AllowProjectBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				ProjectBackup = fmt.Sprintf("%s.project-%s-%s.json", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name))
				project, err := vmio.LoadProject(descriptor.Id)
				if err == nil {
					vmio.ExportUserProject(project, ProjectBackup, "json", true)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Project backup at : %s", ProjectBackup))
				}
			}
			
			request.Arguments.Options = append(request.Arguments.Options, []string{"skip-indexes", "true"})
			response, err := request.DeleteProject()
			if err != nil {
				return response, err
			}
		}
	} else {
		descriptor.Open = true
	}
	
	AllowInfraDeletion := OverrideInfra && Force
	AllowInfraBackup := Backup && OverrideInfra && Force
	InfraBackup := ""
	
	if descriptor.InfraId != "" && OverrideInfra && (FullImport || ElementType == SProject) {
		if ! AllowInfraDeletion {
			AllowInfraDeletion = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with override process for Infrastructure named '%s'?", descriptor.InfraName))
			if ! AllowInfraDeletion {
				response := Response{
					Status: false,
					Message: "User task interruption",
				}
				return response, errors.New("Unable to execute task")
			}
		}
		if Backup && AllowInfraDeletion && ! AllowInfraBackup {
			AllowInfraBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?",descriptor.InfraName))
		}
	}
	
	if existsInfrastructure && OverrideInfra && (FullImport || ElementType == SProject) {
		if AllowInfraBackup {
			folder := model.GetEmergencyFolder()
			if BackupDir != "" {
				folder = BackupDir
			}
			if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
				folder += string(os.PathSeparator)
			}
			InfraBackup = fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name), utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
			infra, err := vmio.LoadInfrastructure(descriptor.Id)
			if err == nil {
				infra.Save(InfraBackup)
				utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure backup at : %s", InfraBackup))
			}
		}
		if (AllowInfraDeletion) {
			request.Arguments.Options = append(request.Arguments.Options, []string{"infra-name", descriptor.InfraName})
			request.Type = DestroyInfrastructure
			request.TypeStr = "delete-infra"
			request.SubType = NoSubCommand
			request.SubTypeStr = ""
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
			return response, errors.New("Unable to execute task")
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
			return response, errors.New("Unable to execute task")
		}
		err = UpdateIndexWithProjectStates(project, false, false)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		err = vmio.SaveProject(project)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			request.DeleteProject()

			return response, errors.New("Unable to execute task")
		}
		if existsInfrastructure {
			AddProjectChangeActions(descriptor.Id, ActionDescriptor{
				Id: NewUUIDString(),
				Date: time.Now(),
				DropAction: false,
				ElementType: SProject,
				ElementId: descriptor.Id,
				ElementName: descriptor.Name,
				FullProject: true,
				JSONImage: ProjectJSON,
				Request: request.Type,
				SubRequest: request.SubType,
			})
		}
		utils.PrintlnSuccess(fmt.Sprintf("Successfully imported Project '%s' from file '%s' in format '%s'", Name, File, Format))
	} else {
		if ! Found {
			response := Response{
				Status: false,
				Message: "Project '" + Name + "' doesn't exist, no suitable project for new elements",
			}
			return response, errors.New("Unable to execute task")
		}
		switch ElementType {
		case SDomain:
			domains := ExportImportDomains{
				Domains: []model.MachineDomain{},
			}
			domains, err = UserImportDomains(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
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
					return response, errors.New("Unable to execute task")
				}
				project.Domains = append(project.Domains, domain)
				if existsInfrastructure {
					AddProjectChangeActions(descriptor.Id, ActionDescriptor{
						Id: NewUUIDString(),
						Date: time.Now(),
						DropAction: false,
						ElementType: SDomain,
						ElementId: domain.Id,
						ElementName: domain.Name,
						FullProject: false,
						JSONImage: ProjectJSON,
						Request: request.Type,
						SubRequest: request.SubType,
					})
				}
			}
			if len(domains.Domains) > 0 {
				project.LastMessage = fmt.Sprintf("Domains (no. %d) imported from file %s, format %s", len(domains.Domains), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return response, errors.New("Unable to execute task")
				}
				utils.PrintlnSuccess(fmt.Sprintf("Successfully imported Project '%s' Domains from file '%s' in format '%s'", Name, File, Format))
			} else {
				utils.PrintlnWarning(fmt.Sprintf("Warning: No Domain to import from file '%s' in format '%s'", File, Format))
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
				return response, errors.New("Unable to execute task")
			}
			for _, networkImport := range networks {
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
								return response, errors.New("Unable to execute task")
							}
							project.Domains[i].Networks = append(project.Domains[i].Networks, network)
							if existsInfrastructure {
								AddProjectChangeActions(descriptor.Id, ActionDescriptor{
									Id: NewUUIDString(),
									Date: time.Now(),
									DropAction: false,
									ElementType: SNetwork,
									ElementId: network.Id,
									ElementName: network.Name,
									FullProject: false,
									JSONImage: ProjectJSON,
									Request: request.Type,
									SubRequest: request.SubType,
								})
							}
						}
					}
				}
				if ! Found {
					response := Response{
						Status: false,
						Message: fmt.Sprintf("Errors during Networks Import : Import target Domain not found by name : '%s' or by id : '%s'", DomainName, DomainId),
					}
					return response, errors.New("Unable to execute task")
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
					return response, errors.New("Unable to execute task")
				}
				utils.PrintlnSuccess(fmt.Sprintf("Successfully imported Project '%s' Networks from file '%s' in format '%s'", Name, File, Format))
			} else {
				utils.PrintlnWarning(fmt.Sprintf("Warning: No Network to import from file '%s' in format '%s'", File, Format))
			}
			break
		case LMachine:
			var machines []ExportImportLocalMachines
			machines, err = UserImportLocalMachines(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}
			for _, machineImport := range machines {
				DomainId := machineImport.Network.DomainId
				DomainName := machineImport.Network.DomainName
				NetworkId := machineImport.Network.NetworkId
				NetworkName := machineImport.Network.NetworkName
				Found := false
				for i := 0; i < len(project.Domains); i++ {
					if project.Domains[i].Id == DomainId || project.Domains[i].Name == DomainName {
						for j := 0; j < len(project.Domains[i].Networks); j++ {
							if project.Domains[i].Networks[j].Id == NetworkId || project.Domains[i].Networks[j].Name == NetworkName {
								Found = true
								for _, machine := range machineImport.Machines {
									machine.PostImport()
									errorList := machine.Validate()
									if len(errorList) > 0 {
										_, errorValue := vmio.StripErrorMessages("Machines import is invalid, clause(s) :", errorList)
										response := Response{
											Status: false,
											Message: errorValue,
										}
										return response, errors.New("Unable to execute task")
									}
									project.Domains[i].Networks[j].LocalMachines = append(project.Domains[i].Networks[j].LocalMachines, machine)
									if existsInfrastructure {
										AddProjectChangeActions(descriptor.Id, ActionDescriptor{
											Id: NewUUIDString(),
											Date: time.Now(),
											DropAction: false,
											ElementType: LMachine,
											ElementId: machine.Id,
											ElementName: machine.Name,
											FullProject: false,
											JSONImage: ProjectJSON,
											Request: request.Type,
											SubRequest: request.SubType,
										})
									}
								}
							}
						}
					}
				}
				if ! Found {
					response := Response{
						Status: false,
						Message: fmt.Sprintf("Errors during Machines Import : Import target Domain not found by name : '%s' or by id : '%s and target Network not found by name : '%s' or by id : '%s '", DomainName, DomainId, NetworkName, NetworkId),
					}
					return response, errors.New("Unable to execute task")
				}
			}
			if len(machines) > 0 {
				project.LastMessage = fmt.Sprintf("Machines (no. %d) imported from file %s, format %s", len(machines), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return response, errors.New("Unable to execute task")
				}
				utils.PrintlnSuccess(fmt.Sprintf("Successfully imported Project '%s' Local Machines from file '%s' in format '%s'", Name, File, Format))
			} else {
				utils.PrintlnWarning(fmt.Sprintf("Warning: No Machine to import from file '%s' in format '%s'", File, Format))
			}
			break
		case CLMachine:
			var machines []ExportImportCloudMachines
			machines, err = UserImportCloudMachines(File, Format)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}
			for _, machineImport := range machines {
				DomainId := machineImport.Network.DomainId
				DomainName := machineImport.Network.DomainName
				NetworkId := machineImport.Network.NetworkId
				NetworkName := machineImport.Network.NetworkName
				Found := false
				for i := 0; i < len(project.Domains); i++ {
					if project.Domains[i].Id == DomainId || project.Domains[i].Name == DomainName {
						for j := 0; j < len(project.Domains[i].Networks); j++ {
							if project.Domains[i].Networks[j].Id == NetworkId || project.Domains[i].Networks[j].Name == NetworkName {
								Found = true
								for _, machine := range machineImport.Machines {
									machine.PostImport()
									errorList := machine.Validate()
									if len(errorList) > 0 {
										_, errorValue := vmio.StripErrorMessages("Machines import is invalid, clause(s) :", errorList)
										response := Response{
											Status: false,
											Message: errorValue,
										}
										return response, errors.New("Unable to execute task")
									}
									project.Domains[i].Networks[j].CloudMachines = append(project.Domains[i].Networks[j].CloudMachines, machine)
									if existsInfrastructure {
										AddProjectChangeActions(descriptor.Id, ActionDescriptor{
											Id: NewUUIDString(),
											Date: time.Now(),
											DropAction: false,
											ElementType: CLMachine,
											ElementId: machine.Id,
											ElementName: machine.Name,
											FullProject: false,
											JSONImage: ProjectJSON,
											Request: request.Type,
											SubRequest: request.SubType,
										})
									}
								}
							}
						}
					}
				}
				if ! Found {
					response := Response{
						Status: false,
						Message: fmt.Sprintf("Errors during Cloud Machines Import : Import target Domain not found by name : '%s' or by id : '%s and target Network not found by name : '%s' or by id : '%s '", DomainName, DomainId, NetworkName, NetworkId),
					}
					return response, errors.New("Unable to execute task")
				}
			}
			if len(machines) > 0 {
				project.LastMessage = fmt.Sprintf("Cloud Machines (no. %d) imported from file %s, format %s", len(machines), File, Format)
				err = vmio.SaveProject(project)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return response, errors.New("Unable to execute task")
				}
				utils.PrintlnSuccess(fmt.Sprintf("Successfully imported Project '%s' Cloud Machines from file '%s' in format '%s'", Name, File, Format))
			} else {
				utils.PrintlnWarning(fmt.Sprintf("Warning: No Machine to import from file '%s' in format '%s'", File, Format))
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
				return response, errors.New("Unable to execute task")
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
										return response, errors.New("Unable to execute task")
									}
									FoundMachine := false
									if plan.IsCloud {
										for _, machine := range project.Domains[i].Networks[j].CloudMachines {
											if plan.MachineId == machine.Id || CorrectInput(plan.MachineId) == CorrectInput(machine.Name) {
												FoundMachine = true
												plan.MachineId = machine.Id
												break
											}
										}
									} else {
										for _, machine := range project.Domains[i].Networks[j].LocalMachines {
											if plan.MachineId == machine.Id || CorrectInput(plan.MachineId) == CorrectInput(machine.Name) {
												FoundMachine = true
												plan.MachineId = machine.Id
												break
											}
										}
									}
									if ! FoundMachine {
										IsCloud := "no"
										if plan.IsCloud {
											IsCloud = "yes"
										}
										response := Response{
											Status: false,
											Message: fmt.Sprintf("Errors during Plans Import : Plan Machine id/name: %s from Cloud: %s not found in domain id: %s, name: %s and network id: %s, name: %s", plan.MachineId, IsCloud, DomainName, DomainId, NetworkName, NetworkId),
										}
										return response, errors.New("Unable to execute task")
									}
									project.Domains[i].Networks[j].Installations = append(project.Domains[i].Networks[j].Installations, plan)
									if existsInfrastructure {
										AddProjectChangeActions(descriptor.Id, ActionDescriptor{
											Id: NewUUIDString(),
											Date: time.Now(),
											DropAction: false,
											ElementType: SPlan,
											ElementId: plan.Id,
											ElementName: fmt.Sprintf("Plan for machine id : %s", plan.MachineId),
											RelatedId: plan.MachineId,
											FullProject: false,
											JSONImage: ProjectJSON,
											Request: request.Type,
											SubRequest: request.SubType,
										})
									}
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
					return response, errors.New("Unable to execute task")
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
					return response, errors.New("Unable to execute task")
				}
				utils.PrintlnSuccess(fmt.Sprintf("Successfully imported Project '%s' Installation Plans from file '%s' in format '%s'", Name, File, Format))
			} else {
				utils.PrintlnWarning(fmt.Sprintf("Warning: No Plan to import from file '%s' in format '%s'", File, Format))
			}
		}
		err = UpdateIndexWithProjectStates(project, false, false)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
	}
	if ProjectBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project backup file : %s", ProjectBackup))
		err = model.DeleteIfExists(ProjectBackup)
	}

	if InfraBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Infrastructure backup file : %s", InfraBackup))
		err = model.DeleteIfExists(InfraBackup)
	}

	if existsInfrastructure && OverrideInfra {
		request.Arguments.Options = append(request.Arguments.Options, []string{"rebuild", "true"})
		request.BuildProject()
	}
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) ExportProject() (Response, error) {
	utils.NO_COLORS = false
	Name := ""
	File := ""
	Format := "json"
	FullExport := true
	var ElementType CmdElementType = NoElement
	var err error
	for _, option := range request.Arguments.Options {
		if "name" == CorrectInput(option[0]) {
			Name = option[1]
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
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		}
	}
	if strings.TrimSpace(Name) == "" {
		response := Response{
			Status: false,
			Message: "Project Name Field is mandatory",
		}
		return response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(File) == "" {
		response := Response{
			Status: false,
			Message: "Export File Path Field is mandatory",
		}
		return response, errors.New("Unable to execute task")
	}
	if strings.TrimSpace(Format) == "" {
		response := Response{
			Status: false,
			Message: "Export Format Field is mandatory",
		}
		return response, errors.New("Unable to execute task")
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
		ExportMachines := "no"
		if ElementType == LMachine {
			ExportMachines = "yes"
		}
		ExportCMachines := "no"
		if ElementType == CLMachine {
			ExportCMachines = "yes"
		}
		ExportPlans := "no"
		if ElementType == SPlan {
			ExportPlans = "yes"
		}
		fmt.Printf("Export Domains : %s\nExport Networks : %s\nExport Machines : %s\nExport Cloud Machines : %s\nExport Plans : %s\n",
			ExportDomains, ExportNetworks, ExportMachines, ExportCMachines, ExportPlans)
	}
	if ! FullExport && ElementType == NoElement {
		response := Response{
			Status: false,
			Message: "No Full Import and Element Type not supported, nothing to import, VMKube quits!!",
		}
		return response, errors.New("Unable to execute task")
	}
	model.DeleteIfExists(File)
	if FullExport || ElementType == SProject {
		project, err := vmio.LoadProject(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		err = vmio.ExportUserProject(project, File, Format, true)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		utils.PrintlnSuccess(fmt.Sprintf("Successfully exported Project '%s' to file '%s' in format '%s' ", Name, File, Format))
	} else {
		project, err := vmio.LoadProject(descriptor.Id)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		switch ElementType {
		case SDomain:
			domains := ExportImportDomains{
				Domains: []model.MachineDomain{},
			}
			domains.Domains = append(domains.Domains, project.Domains...)
			err = utils.ExportStructureToFile(File, Format, domains)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}
			utils.PrintlnSuccess(fmt.Sprintf("Successfully exported Project '%s' Domains to file '%s' in format '%s' ", Name, File, Format))
			break
		case SNetwork:
			var networks []ExportImportNetwork = make([]ExportImportNetwork, 0)
			for _, domain := range project.Domains {
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
				return response, errors.New("Unable to execute task")
			}
			utils.PrintlnSuccess(fmt.Sprintf("Successfully exported Project '%s' Networks to file '%s' in format '%s'", Name, File, Format))
			break
		case LMachine:
			var machines []ExportImportLocalMachines = make([]ExportImportLocalMachines, 0)
			for _, domain := range project.Domains {
				for _, network := range domain.Networks {
					machine := ExportImportLocalMachines{
						Network: NetworkReference{
							DomainId: domain.Id,
							DomainName: domain.Name,
							NetworkId: network.Id,
							NetworkName: network.Name,
						},
						Machines: network.LocalMachines,
					}
					machines = append(machines, machine)
				}
			}
			err = utils.ExportStructureToFile(File, Format, machines)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}
			utils.PrintlnSuccess(fmt.Sprintf("Successfully exported Project '%s' Local Machines to file '%s' in format '%s'", Name, File, Format))
			break
		case CLMachine:
			var machines []ExportImportCloudMachines = make([]ExportImportCloudMachines, 0)
			for _, domain := range project.Domains {
				for _, network := range domain.Networks {
					machine := ExportImportCloudMachines{
						Network: NetworkReference{
							DomainId: domain.Id,
							DomainName: domain.Name,
							NetworkId: network.Id,
							NetworkName: network.Name,
						},
						Machines: network.CloudMachines,
					}
					machines = append(machines, machine)
				}
			}
			err = utils.ExportStructureToFile(File, Format, machines)
			if err != nil {
				response := Response{
					Status: false,
					Message: err.Error(),
				}
				return response, errors.New("Unable to execute task")
			}
			utils.PrintlnSuccess(fmt.Sprintf("Successfully exported Project '%s' Cloud Machines to file '%s' in format '%s'", Name, File, Format))
			break
		default:
			//Plans
			var plans []ExportImportPlans = make([]ExportImportPlans, 0)
			for _, domain := range project.Domains {
				for _, network := range domain.Networks {
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
				return response, errors.New("Unable to execute task")
			}
			utils.PrintlnSuccess(fmt.Sprintf("Successfully exported Project '%s' Installation Plans to file '%s' in format '%s'", Name, File, Format))
		}
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return response, nil
}
