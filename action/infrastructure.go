package action

import (
	"fmt"
	"errors"
	"vmkube/utils"
	"vmkube/vmio"
	"vmkube/model"
	"vmkube/operations"
	"vmkube/scheduler"
	"runtime"
	"strings"
	"time"
	"os"
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
					utils.PrintlnError(fmt.Sprintf("Option '--%s' is mandatory!!\n", option.Option))
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
	//TODO: Alter Infra to : Start/Stop/Recreate/Remove Single Instance or AutoFix instances errors.
	response := Response{
		Status: false,
		Message: "Not Implemented",
	}
	return  response, errors.New("Unable to execute task")
}

func (request *CmdRequest) DeleteInfra() (Response, error) {
	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "overclock" == CorrectInput(option[0]) {
			Overclock = GetBoolean(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "threads" == CorrectInput(option[0]) {
			Threads = GetInteger(option[1], Threads)
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	if ! Force {
		AllowInfraDeletion := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with deletion process for Infrastructure named '%s'?", Name))
		if ! AllowInfraDeletion {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()
	
	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
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
	
	utils.PrintlnImportant("Now Proceding with machines destroy ...!!")
	NumThreads := Threads
	if runtime.NumCPU() - 1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	
	deleteActivities, err := operations.GetPostBuildTaskActivities(infrastructure, operations.DestroyMachine, []string{})
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, deleteActivities, NumThreads,func(task scheduler.ScheduleTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete Destroy machines process : Errors deleting Infrastructure : %s!!",  Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error deleting infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status: false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning("Removing Infrastructure logs ...")
	
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,instance := range network.LocalInstances {
				err = DeleteInfrastructureLogs(instance.Logs)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			for _,instance := range network.CloudInstances {
				err = DeleteInfrastructureLogs(instance.Logs)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
			for _,installation := range network.Installations {
				err = DeleteInfrastructureLogs(installation.Logs)
				if err != nil {
					response := Response{
						Status: false,
						Message: err.Error(),
					}
					return  response, errors.New("Unable to execute task")
				}
			}
		}
	}
	
	utils.PrintlnWarning("Removing Rollback Index and Segments ...")
	
	
	rollbackIndexInfo := ProjectRollbackIndexInfo{
		Format: "",
		Index: RollBackIndex{
			Id: "",
			ProjectId: descriptor.Id,
			IndexList: []RollBackSegmentIndex{},
		},
	}
	
	err = rollbackIndexInfo.Read()
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	iFaceRollbackIndex := IFaceRollBackIndex{
		Id: rollbackIndexInfo.Index.Id,
		ProjectId: descriptor.Id,
	}
	
	iFaceRollbackIndex.WaitForUnlock()
	
	LockRollBackIndex(rollbackIndexInfo.Index)

	err = rollbackIndexInfo.Delete()

	UnlockRollBackIndex(rollbackIndexInfo.Index)

	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning("Removing Infrastructure ...")
	
	iFaceInfra.WaitForUnlock()
	
	vmio.LockInfrastructureById(descriptor.Id, descriptor.InfraId)
	var info vmio.InfrastructureInfo =  vmio.InfrastructureInfo{
		Format: "",
		Infra: model.Infrastructure{
			ProjectId: descriptor.Id,
			Id: descriptor.InfraId,
		},
	}
	err = info.Delete()
	vmio.UnlockInfrastructureById(descriptor.Id, descriptor.InfraId)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	descriptor.InfraId = ""
	descriptor.InfraName = ""

	utils.PrintlnWarning(fmt.Sprintf("Reopening Project '%s' ...", descriptor.Name))
	
	project.Open = true
	
	err = vmio.SaveProject(project)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning("Updating global indexes ...")
	
	descriptor.Open = true

	err = UpdateIndexWithProjectsDescriptor(descriptor, true)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnSuccess(fmt.Sprintf("Delete of Infrastructure named : %s completed successfully!!", Name))
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) BackupInfra() (Response, error) {
	Name := ""
	File := ""
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "file" == CorrectInput(option[0]) {
			File = option[1]
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name field not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	if File == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "File Path field not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()
	
	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	if ! strings.HasSuffix(strings.ToLower(File), ".vmkube") {
		utils.PrintlnWarning("File extension changed, cause it was not standard...")
		File += ".vmkube"
	}
	
	err = infrastructure.Save(File)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnSuccess(fmt.Sprintf("Backup Infrastructure named : %s completed successfully!!", Name))
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) RecoverInfra() (Response, error) {
	//TODO: Test Beckup / Recovery commands
	Name := ""
	File := ""
	Override := false
	Force := false
	ProjectName := ""
	Threads := 1
	Overclock := false
	BackupDir := ""
	Backup := false
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "file" == CorrectInput(option[0]) {
			File = option[1]
		} else if "override" == CorrectInput(option[0]) {
			Override = GetBoolean(option[1])
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "project-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "overclock" == CorrectInput(option[0]) {
			Overclock = GetBoolean(option[1])
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
			Message: "Infrastrcuture Name field not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	if File == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "File Path field not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	
	if ! strings.HasSuffix(strings.ToLower(File), ".vmkube") {
		utils.PrintlnError("File extension wrong, it can cause panic in the system")
		response := Response{
			Status: false,
			Message: fmt.Sprintf("Wrong file type for '%s', expected vmkube extension!!", File),
		}
		return  response, errors.New("Unable to execute task")
	}

	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)

	existsInfrastructure := err == nil && descriptor.InfraId == ""

	if existsInfrastructure {
		utils.PrintlnWarning(fmt.Sprintf("Infrastructure '%s' already exists ...", Name))
	}
	
	DeleteFromDescriptor := false
	
	AllowOverride := false
	
	if existsInfrastructure && CorrectInput(ProjectName) != CorrectInput(descriptor.Name) {
		if ! Override {
			return Response{
				Message: fmt.Sprintf("Project named %s already associated with infrastructure %s, and no override clause specified", descriptor.Name, descriptor.InfraName),
				Status: false,
			},errors.New("Unable to execute task")
		} else {
			if ! Force {
				DeleteFromDescriptor = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with override of existing Project named '%s' and infrastructure '%s'?", descriptor.Name, descriptor.InfraName))
				if ! DeleteFromDescriptor {
					return Response{
						Message: "User task interruption",
						Status: false,
					},errors.New("Unable to execute task")
				}
				AllowOverride = true
			} else {
				utils.PrintlnWarning(fmt.Sprintf("Project named '%s' and related existing '%s' infratructure will be removed ...", descriptor.Name, descriptor.InfraName))
				DeleteFromDescriptor = true
				AllowOverride = true
			}
		}
	}
	
	if existsInfrastructure && !AllowOverride && CorrectInput(ProjectName) == CorrectInput(descriptor.Name) {
		if !Override {
			return Response{
				Message: fmt.Sprintf("Infrastructure named %s already associated with project %s, and no override clause specified", descriptor.InfraName, descriptor.Name),
				Status: false,
			},errors.New("Unable to execute task")
		} else {
			if ! Force {
				AllowOverride = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with override of existing Infrastructure named '%s' associated with project '%s'?", descriptor.InfraName, descriptor.Name))
				if ! AllowOverride {
					return Response{
						Message: "User task interruption",
						Status: false,
					},errors.New("Unable to execute task")
				}
				DeleteFromDescriptor = true
			} else {
				utils.PrintlnWarning(fmt.Sprintf("Infrastructure named %s and Project %s will be replaced ...", descriptor.InfraName, descriptor.Name))
				DeleteFromDescriptor = true
				AllowOverride = true
			}
		}
	}
	
	DeleteFromProjectDescriptor := false

	var projectDescriptor model.ProjectsDescriptor

	if ! existsInfrastructure {
		projectDescriptor, err = vmio.GetProjectDescriptor(ProjectName)
		if err != nil && projectDescriptor.InfraId != descriptor.InfraId {
			if ! Override {
				return Response{
					Message: fmt.Sprintf("Selected Project Name %s is already associated with another infrastructure %s, and no override clause specified", projectDescriptor.Name, projectDescriptor.InfraName),
					Status: false,
				},errors.New("Unable to execute task")
			} else {
				if ! Force {
					DeleteFromProjectDescriptor = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with delete of existing Project named '%s' and related Infrastructure named '%s'?", projectDescriptor.Name, projectDescriptor.InfraName))
					if !DeleteFromProjectDescriptor {
						return Response{
							Message: "User task interruption",
							Status: false,
						},errors.New("Unable to execute task")
					}
				} else {
					utils.PrintlnWarning(fmt.Sprintf("Project named %s and related Infrastructure name '%s'  will be removed ...", projectDescriptor.Name, projectDescriptor.InfraName))
					DeleteFromProjectDescriptor = true
				}
			}
		}
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Loading Infrastructure '%s' from file '%s'...", Name, File))
	
	var infrastructure model.Infrastructure
	err = infrastructure.Load(File)

	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Validating recovered Infrastructure '%s'...", Name))
	
	errorsList := infrastructure.Validate()
	
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete Recover for Infrastructure '%s'!!", Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error validating infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status: false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}
	
	if DeleteFromDescriptor {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project %s and Infrastructure '%s'...", descriptor.Name, descriptor.InfraName))
		if Backup {
			AllowBackup := Force
			if ! AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Project named '%s'?",descriptor.Name))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				ProjectBackup := fmt.Sprintf("%s.project-%s-%s.json", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name))
				project, err := vmio.LoadProject(descriptor.Id)
				if err == nil {
					vmio.ExportUserProject(project, ProjectBackup, "json", true)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Project '%s' backup at : %s", descriptor.Name, ProjectBackup))
				}
			}
			AllowBackup = Force
			if ! AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?",descriptor.InfraName))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				InfraBackup := fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name), utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
				infra, err := vmio.LoadInfrastructure(descriptor.Id)
				if err == nil {
					infra.Save(InfraBackup)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure '%s' backup at : %s", descriptor.InfraName, InfraBackup))
				}
			}
		}
		request.Type = DeleteConfig
		request.SubTypeStr = "delete-project"
		FoundName := false
		FoundForce := false
		for i := 0; i < len(request.Arguments.Options); i++ {
			if "force" == CorrectInput(request.Arguments.Options[i][0]) {
				request.Arguments.Options[i][1] = "true"
				FoundName = true
			} else if "name" == CorrectInput(request.Arguments.Options[i][0]) {
				request.Arguments.Options[i][1] = descriptor.Name
				FoundForce = true
			}
		}
		if ! FoundName {
			request.Arguments.Options = append(request.Arguments.Options, []string{"name", descriptor.Name})
		}
		if ! FoundForce {
			request.Arguments.Options = append(request.Arguments.Options, []string{"force", "true"})
		}
		request.DeleteProject()
	}
	
	if DeleteFromProjectDescriptor {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project %s and Infrastructure '%s'...", projectDescriptor.Name, projectDescriptor.InfraName))
		if Backup {
			AllowBackup := Force
			if ! AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Project named '%s'?",projectDescriptor.Name))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				ProjectBackup := fmt.Sprintf("%s.project-%s-%s.json", folder, utils.IdToFileFormat(projectDescriptor.Id), utils.NameToFileFormat(projectDescriptor.Name))
				project, err := vmio.LoadProject(projectDescriptor.Id)
				if err == nil {
					vmio.ExportUserProject(project, ProjectBackup, "json", true)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Project '%s' backup at : %s", projectDescriptor.Name, ProjectBackup))
				}
			}
			AllowBackup = Force
			if ! AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?",projectDescriptor.InfraName))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if ! strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				InfraBackup := fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(projectDescriptor.Id), utils.NameToFileFormat(projectDescriptor.Name), utils.IdToFileFormat(projectDescriptor.InfraId), utils.NameToFileFormat(projectDescriptor.InfraName))
				infra, err := vmio.LoadInfrastructure(projectDescriptor.Id)
				if err == nil {
					infra.Save(InfraBackup)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure '%s' backup at : %s", projectDescriptor.InfraName, InfraBackup))
				}
			}
		}
		
		request.Type = DeleteConfig
		request.SubTypeStr = "delete-project"
		FoundName := false
		FoundForce := false
		for i := 0; i < len(request.Arguments.Options); i++ {
			if "force" == CorrectInput(request.Arguments.Options[i][0]) {
				request.Arguments.Options[i][1] = "true"
				FoundName = true
			} else if "name" == CorrectInput(request.Arguments.Options[i][0]) {
				request.Arguments.Options[i][1] = projectDescriptor.Name
				FoundForce = true
			}
		}
		if ! FoundName {
			request.Arguments.Options = append(request.Arguments.Options, []string{"name", projectDescriptor.Name})
		}
		if ! FoundForce {
			request.Arguments.Options = append(request.Arguments.Options, []string{"force", "true"})
		}
		request.DeleteProject()
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Creating new Project %s from Infrastructure '%s'...", ProjectName, Name))
	
	newProject, err := InfrastructureToProject(infrastructure, ProjectName)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	newDescriptor := model.ProjectsDescriptor{
		Id: newProject.Id,
		Name: newProject.Name,
		Open: false,
		Active: false,
		Synced: true,
		InfraId: infrastructure.Id,
		InfraName: infrastructure.Name,
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Defining new Indexes from new Project %s and Infrastructure '%s'...", ProjectName, Name))

	err = UpdateIndexWithProjectsDescriptor(newDescriptor, true)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Saving new Project %s...", ProjectName))

	err = vmio.SaveProject(newProject)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Saving recovered Infrastructure '%s'...", Name))
	
	err = vmio.SaveInfrastructure(infrastructure)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Building new Project '%s'...", ProjectName))
	actionCouples, err := make([]operations.ActivityCouple, 0), errors.New("Unknown Error")
	if ! existsInfrastructure {
		creationCouples, err := operations.GetTaskActivities(newProject, infrastructure, operations.CreateMachine)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		
		inspectCouples, err := operations.GetTaskActivities(newProject, infrastructure, operations.MachineInspect)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		ipAddressCouples, err := operations.GetTaskActivities(newProject, infrastructure, operations.MachineIPAddress)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		stopCouples, err := operations.GetTaskActivities(newProject, infrastructure, operations.StopMachine)
		if err != nil {
			response := Response{
				Status: false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		extendsDiskCouples, err := operations.GetTaskActivities(newProject, infrastructure, operations.MachineExtendsDisk)
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
	}
	utils.PrintlnImportant("Now Proceding with machine creation ...!!")
	NumThreads := Threads
	if runtime.NumCPU() - 1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	
	var fixInfraValue int = len(actionCouples)
	errorsList = ExecuteInfrastructureActions(infrastructure, actionCouples, NumThreads,func(task scheduler.ScheduleTask){
		go func(task scheduler.ScheduleTask) {
			response := strings.Split(fmt.Sprintf("%s",task.Jobs[0].Runnable.Response()),"|")
			if len(response) > 3 {
				if len(response) > 4 {
					if response[0] == "ip" {
						instanceId := response[1]
						ipAddress := response[2]
						json := ""
						log := response[3] + response[4]
						FixInfrastructureElementValue(infrastructure, instanceId, ipAddress, json, log)
					} else if response[0] == "json" {
						instanceId := response[1]
						ipAddress := ""
						json := response[2]
						log := response[3] + response[4]
						FixInfrastructureElementValue(infrastructure, instanceId, ipAddress, json, log)
					}
				} else {
					instanceId := response[1]
					ipAddress := ""
					json := ""
					log := response[2] + response[3]
					FixInfrastructureElementValue(infrastructure, instanceId, ipAddress, json, log)
				}
			}
			//if len(response) > 2 {
			//	json := response[1]
			//	ipAddress := response[2]
			//	instanceId := response[0]
			//	FixInfrastructureElementValue(infrastructure, instanceId, ipAddress, json)
			//}
			fixInfraValue--
		}(task)
	})
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete Build of project '%s' : Errors building Infrastructure : '%s'!!", ProjectName, Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error building Infrastructure -> '%s' : ", Name), errorsList)
		response := Response{
			Status: false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnWarning(fmt.Sprintf("Waiting for Instance recovery information in Project '%s' Infrastrcucture '%s'", ProjectName, Name))
	
	for fixInfraValue > 0 {
		time.Sleep(1*time.Second)
	}
	
	
	utils.PrintlnWarning(fmt.Sprintf("Updating recovered Infrastructure '%s' with new Instance data...", Name))
	
	err = vmio.SaveInfrastructure(infrastructure)
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnSuccess(fmt.Sprintf("Recovery for Infrastructure named : %s completed successfully!!", Name))
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) StartInfra() (Response, error) {
	//TODO: Manage other error status ...
	
	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "overclock" == CorrectInput(option[0]) {
			Overclock = GetBoolean(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "threads" == CorrectInput(option[0]) {
			Threads = GetInteger(option[1], Threads)
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	if ! descriptor.Active {
		response := Response{
			Status: false,
			Message: fmt.Sprintf("Infrastructure %s not active, unable to start it. Please close project to perform this task", Name),
		}
		return response, errors.New("Unable to execute task")
	}
	if ! Force {
		AllowInfraStart := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with start machines process for Infrastructure named '%s'?", Name))
		if ! AllowInfraStart {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()
	
	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnImportant("Now Proceding with machines start ...!!")
	NumThreads := Threads
	if runtime.NumCPU() - 1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	
	startMachineActivities, err := operations.GetPostBuildTaskActivities(infrastructure, operations.StartMachine, []string{})
	
	startMachineActivities = operations.FilterByInstanceState(startMachineActivities, false)
	
	if len(startMachineActivities) == 0 {
		utils.PrintlnImportant("No machine to start ...")
		response := Response{
			Status: true,
			Message: "Success",
		}
		return  response, nil
	}
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, startMachineActivities, NumThreads,func(task scheduler.ScheduleTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete start machines process : Errors starting Infrastructure : %s!!",  Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error starting infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status: false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnSuccess(fmt.Sprintf("Start Infrastructure named : %s completed successfully!!", Name))
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) StopInfra() (Response, error) {
	//TODO: Manage other error status ...

	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "overclock" == CorrectInput(option[0]) {
			Overclock = GetBoolean(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "threads" == CorrectInput(option[0]) {
			Threads = GetInteger(option[1], Threads)
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	if ! Force {
		AllowInfraStop := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with stop machines process for Infrastructure named '%s'?", Name))
		if ! AllowInfraStop {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()
	
	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnImportant("Now Proceding with machines stop ...!!")
	NumThreads := Threads
	if runtime.NumCPU() - 1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	
	stopMachineActivities, err := operations.GetPostBuildTaskActivities(infrastructure, operations.StopMachine, []string{})
	
	stopMachineActivities = operations.FilterByInstanceState(stopMachineActivities, false)
	
	if len(stopMachineActivities) == 0 {
		utils.PrintlnImportant("No machine to stop ...")
		response := Response{
			Status: true,
			Message: "Success",
		}
		return  response, nil
	}
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, stopMachineActivities, NumThreads,func(task scheduler.ScheduleTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete stop machines process : Errors stopping Infrastructure : %s!!",  Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error stopping infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status: false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnSuccess(fmt.Sprintf("Stop Infrastructure named : %s completed successfully!!", Name))
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) RestartInfra() (Response, error) {
	//TODO: Manage other error status ...

	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "overclock" == CorrectInput(option[0]) {
			Overclock = GetBoolean(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "threads" == CorrectInput(option[0]) {
			Threads = GetInteger(option[1], Threads)
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	if ! descriptor.Active {
		response := Response{
			Status: false,
			Message: fmt.Sprintf("Infrastructure %s not active, unable to restart it. Please close project to perform this task", Name),
		}
		return response, errors.New("Unable to execute task")
	}
	if ! Force {
		AllowInfraStop := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with restart machines process for Infrastructure named '%s'?", Name))
		if ! AllowInfraStop {
			response := Response{
				Status: false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()
	
	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnImportant("Now Proceding with machines restart ...!!")
	NumThreads := Threads
	if runtime.NumCPU() - 1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))
	
	restartMachineActivities, err := operations.GetPostBuildTaskActivities(infrastructure, operations.RestartMachine, []string{})
	
	restartMachineActivities = operations.FilterByInstanceState(restartMachineActivities, false)
	
	if len(restartMachineActivities) == 0 {
		utils.PrintlnImportant("No machine to restart ...")
		response := Response{
			Status: true,
			Message: "Success",
		}
		return  response, nil
	}
	
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, restartMachineActivities, NumThreads,func(task scheduler.ScheduleTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete restart machines process : Errors restarting Infrastructure : %s!!",  Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error restarting infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status: false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}
	
	utils.PrintlnSuccess(fmt.Sprintf("Restart Infrastructure named : %s completed successfully!!", Name))
	
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) ListInfras() (Response, error) {
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
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
		return  response, errors.New("Unable to execute task")
	}
	if len(indexes.Projects) > 0 {
		utils.PrintlnImportant(fmt.Sprintf("%s  %s  %s  %s", utils.StrPad("Infrastructure Id", 40), utils.StrPad("Infrastructure Name", 40), utils.StrPad("Active", 6), utils.StrPad("Synced", 6)))
	} else {
		utils.PrintlnImportant("No Infrastructures found")
	}
	for _,index := range indexes.Projects {
		active := "no"
		synced := "no"
		if index.Active {
			active = "yes"
		}
		if index.Synced {
			synced = "yes"
		}
		fmt.Printf("%s  %s  %s  %s\n", utils.StrPad(index.InfraId, 40), utils.StrPad(index.InfraName, 40), utils.StrPad("  "+active, 6), utils.StrPad("  " + synced, 6))
		
	}
	response := Response{
		Status: true,
		Message: "Success",
	}
	return  response, nil
}

func (request *CmdRequest) StatusInfra() (Response, error) {
	Name := ""
	utils.NO_COLORS = false
	for _,option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status: false,},errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id: descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()
	
	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status: false,
			Message: err.Error(),
		}
		return  response, errors.New("Unable to execute task")
	}
	modified := "no"
	if infrastructure.Altered {
		modified = "yes"
	}
	errors := "no"
	if infrastructure.Errors {
		errors = "yes"
	}
	fmt.Printf("Id: %s\nInfrastructure: %s\nModified: %s\n", infrastructure.Id,infrastructure.Name, modified)
	created := "no"
	if infrastructure.Created {
		created = "yes"
	}
	fmt.Printf("Created [%s] : %d-%02d-%02d %02d:%02d:%02d\n", created,
		infrastructure.Creation.Year(), infrastructure.Creation.Month(), infrastructure.Creation.Day(),
		infrastructure.Creation.Hour(), infrastructure.Creation.Minute(), infrastructure.Creation.Second())
	fmt.Printf("Modified : %d-%02d-%02d %02d:%02d:%02d\n",
		infrastructure.Modified.Year(), infrastructure.Modified.Month(), infrastructure.Modified.Day(),
		infrastructure.Modified.Hour(), infrastructure.Modified.Minute(), infrastructure.Modified.Second())
	fmt.Printf("Errors: %s\nLast Message: %s\n", errors,infrastructure.LastMessage)
	fmt.Printf("Domains: %d\n", len(infrastructure.Domains))
	for _,domain := range infrastructure.Domains {
		num, options := vmio.StripOptions(domain.Options)
		fmt.Printf("Domain: %s (Id: %s) - Options [%d] :%s\n", domain.Name, domain.Id, num, options)
		fmt.Printf("Networks: %d\n", len(domain.Networks))
		for _,network := range domain.Networks {
			num, options := vmio.StripOptions(network.Options)
			fmt.Printf("   Network: %s (Id: %s) - Options [%d] :%s\n", network.Name, network.Id, num, options)
			fmt.Printf("   Local Instances: %d\n", len(network.LocalInstances))
			machinesMap := make(map[string]string)
			for _,machine := range network.LocalInstances {
				machinesMap[machine.Id] = machine.Name
				fmt.Printf("      Local Instance: %s (Id: %s) - Driver: %s - OS : %s:%s - IP Address: %s\n", machine.Name, machine.Id, machine.Driver, machine.OSType, machine.OSVersion, machine.IPAddress)
			}
			fmt.Printf("   Cloud Instances: %d\n", len(network.CloudInstances))
			for _,machine := range network.CloudInstances {
				machinesMap[machine.Id] = machine.Name
				num, options := vmio.StripOptions(machine.Options)
				fmt.Printf("      Cloud Instance: %s (Id: %s) - Driver: %s - IP Address: %s - Options [%d] :%s\n", machine.Name, machine.Id, machine.Driver, machine.IPAddress, num, options)
			}
			fmt.Printf("   Installation Plans: %d\n", len(network.Installations))
			for _,installation := range network.Installations {
				machineName,ok := machinesMap[installation.InstanceId]
				if !ok {
					machineName = "<invalid>"
				}
				cloud := "no"
				if installation.IsCloud {
					cloud = "yes"
				}
				success := "no"
				if installation.Success {
					success = "yes"
				}
				fmt.Printf("      Plan: Id: %s - Instance: %s [Id: %s] - Success: %s - Cloud: %s - Envoronment : %s  Role: %s  Type: %s\n", installation.Id, machineName, installation.InstanceId, success, cloud, model.InstanceEnvironmentToString(installation.Environment), model.InstanceRoleToString(installation.Role), model.InstanceInstallationToString(installation.Type))
			}
		}
	}
	return Response{
		Message: "Success",
		Status: true,}, nil
}
