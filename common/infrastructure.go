package common

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/procedures"
	"github.com/hellgate75/vmkube/tasks"
	"github.com/hellgate75/vmkube/utils"
	"github.com/hellgate75/vmkube/vmio"
)

type InfrastructureActions interface {
	CheckInfra(recoverHelpersFunc func() []CommandHelper) bool
	CreateInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	AlterInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	DeleteInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	StartInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	StopInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	RestartInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	ListInfras(recoverHelpersFunc func() []CommandHelper) (Response, error)
	StatusInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	BackupInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
	RecoverInfra(recoverHelpersFunc func() []CommandHelper) (Response, error)
}

func (request *CmdRequest) CheckInfra(recoverHelpersFunc func() []CommandHelper) bool {
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
					utils.PrintlnError(fmt.Sprintf("Option '--%s' is mandatory!!\n", option.Option))
				}
			}
		}
		if !correctness {
			return false
		}
	}
	return true
}

func (request *CmdRequest) CreateInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	response := Response{
		Status:  false,
		Message: "Not Implemented",
	}
	return response, errors.New("Unable to execute task")
}

func (request *CmdRequest) AlterInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	Name := ""
	InstanceId := ""
	InstanceName := ""
	Force := false
	IsCloud := false
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "force" == CorrectInput(option[0]) {
			Force = GetBoolean(option[1])
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "is-cloud" == CorrectInput(option[0]) {
			IsCloud = GetBoolean(option[1])
		} else if "instance-id" == CorrectInput(option[0]) {
			InstanceId = option[1]
		} else if "instance-name" == CorrectInput(option[0]) {
			InstanceName = option[1]
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id:        descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if request.SubType == AutoFix {
		err := AutoFixInfrastructureInstances(infrastructure)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		} else {
			utils.PrintlnImportant(fmt.Sprintf("Alter Infrastructure '%s' Command : '%s' executed successfully!!", descriptor.InfraName, CmdSubRequestDescriptors[int(request.SubType)]))
			response := Response{
				Status:  true,
				Message: "Success",
			}
			return response, nil
		}

	}

	if InstanceId == "" && InstanceName == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Instance Unique Identifier or Name not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	var instance model.LocalInstance
	var cloudInstance model.CloudInstance
	if !IsCloud {
		instance, err = FindInfrastructureInstance(infrastructure, InstanceId, InstanceName)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		if !Force && request.SubType != Status {
			AllowChange := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with changes with Instance '%s' (uid: %s) part of Infrastructure named '%s'?", instance.Name, instance.Id, Name))
			if !AllowChange {
				return Response{
					Message: "User task interruption",
					Status:  false,
				}, errors.New("Unable to execute task")
			}
		}
	} else {
		cloudInstance, err = FindInfrastructureCloudInstance(infrastructure, InstanceId, InstanceName)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		if !Force {
			AllowChange := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with changes with Instance '%s' (uid: %s) part of Infrastructure named '%s'?", cloudInstance.Name, cloudInstance.Id, Name))
			if !AllowChange {
				return Response{
					Message: "User task interruption",
					Status:  false,
				}, errors.New("Unable to execute task")
			}
		}
	}

	switch request.SubType {
	case Status:
		var instanceState procedures.MachineState
		instanceState, err = ExistInstance(infrastructure, instance, cloudInstance, IsCloud, instance.Id)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		err = DescribeInstance(infrastructure, instance, cloudInstance, IsCloud, instanceState)
		break
	case Start:
		var instanceState procedures.MachineState
		instanceState, err = ExistInstance(infrastructure, instance, cloudInstance, IsCloud, instance.Id)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		if instanceState != procedures.Machine_State_Stopped {
			response := Response{
				Status:  false,
				Message: fmt.Sprintf("Instance not stopped : %s", instanceState.String()),
			}
			return response, errors.New("Unable to execute task")
		}
		err = StartInstance(infrastructure, instance, cloudInstance, IsCloud, instanceState)
		break
	case Stop:
		var instanceState procedures.MachineState
		instanceState, err = ExistInstance(infrastructure, instance, cloudInstance, IsCloud, instance.Id)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		if instanceState != procedures.Machine_State_Running {
			response := Response{
				Status:  false,
				Message: fmt.Sprintf("Instance not running : %s", instanceState.String()),
			}
			return response, errors.New("Unable to execute task")
		}
		err = StopInstance(infrastructure, instance, cloudInstance, IsCloud, instanceState)
		break
	case Restart:
		var instanceState procedures.MachineState
		instanceState, err = ExistInstance(infrastructure, instance, cloudInstance, IsCloud, instance.Id)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		if instanceState != procedures.Machine_State_Running {
			response := Response{
				Status:  false,
				Message: fmt.Sprintf("Instance not running : %s", instanceState.String()),
			}
			return response, errors.New("Unable to execute task")
		}
		err = RestartInstance(infrastructure, instance, cloudInstance, IsCloud, instanceState)
		break
	case Disable:
		var instanceState procedures.MachineState
		instanceState, err = ExistInstance(infrastructure, instance, cloudInstance, IsCloud, instance.Id)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		err = DisableInstance(infrastructure, instance, cloudInstance, IsCloud, instanceState, descriptor)
		break
	case Enable:
		err = EnableInstance(infrastructure, instance, cloudInstance, IsCloud)
		break
	case Recreate:
		var instanceState procedures.MachineState
		instanceState, err = ExistInstance(infrastructure, instance, cloudInstance, IsCloud, instance.Id)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		err = RecreateInstance(infrastructure, instance, cloudInstance, IsCloud, instanceState, descriptor)
		break
	default:
		// Destroy
		var instanceState procedures.MachineState
		instanceState, err = ExistInstance(infrastructure, instance, cloudInstance, IsCloud, instance.Id)
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		err = DestroyInstance(infrastructure, instance, cloudInstance, IsCloud, instanceState, descriptor)
	}
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	utils.PrintlnImportant(fmt.Sprintf("Alter Infrastructure '%s' executed successfully!!", descriptor.InfraName))
	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) DeleteInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
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
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	if !Force {
		AllowInfraDeletion := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with deletion process for Infrastructure named '%s'?", Name))
		if !AllowInfraDeletion {
			response := Response{
				Status:  false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id:        descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	project, err := vmio.LoadProject(descriptor.Id)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnImportant("Now Proceding with machines destroy ...!!")
	NumThreads := Threads
	if runtime.NumCPU()-1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))

	exclusionList, _ := FilterForExistState(infrastructure)

	deleteActivities, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.DestroyMachine, exclusionList)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, deleteActivities, NumThreads, func(task tasks.SchedulerTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete Destroy machines process : Errors deleting Infrastructure : %s!!", Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error deleting infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status:  false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning("Removing Infrastructure logs ...")

	for _, domain := range infrastructure.Domains {
		for _, network := range domain.Networks {
			for _, instance := range network.LocalInstances {
				err = DeleteInfrastructureLogs(instance.Logs)
				if err != nil {
					response := Response{
						Status:  false,
						Message: err.Error(),
					}
					return response, errors.New("Unable to execute task")
				}
			}
			for _, instance := range network.CloudInstances {
				err = DeleteInfrastructureLogs(instance.Logs)
				if err != nil {
					response := Response{
						Status:  false,
						Message: err.Error(),
					}
					return response, errors.New("Unable to execute task")
				}
			}
			for _, installation := range network.Installations {
				err = DeleteInfrastructureLogs(installation.Logs)
				if err != nil {
					response := Response{
						Status:  false,
						Message: err.Error(),
					}
					return response, errors.New("Unable to execute task")
				}
			}
		}
	}

	utils.PrintlnWarning("Removing Rollback Index and Segments ...")

	rollbackIndexInfo := ProjectRollbackIndexInfo{
		Format: "",
		Index: RollBackIndex{
			Id:        "",
			ProjectId: descriptor.Id,
			IndexList: []RollBackSegmentIndex{},
		},
	}

	err = rollbackIndexInfo.Read()

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	iFaceRollbackIndex := IFaceRollBackIndex{
		Id:        rollbackIndexInfo.Index.Id,
		ProjectId: descriptor.Id,
	}

	iFaceRollbackIndex.WaitForUnlock()

	LockRollBackIndex(rollbackIndexInfo.Index)

	err = rollbackIndexInfo.Delete()

	UnlockRollBackIndex(rollbackIndexInfo.Index)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning("Removing Infrastructure ...")

	iFaceInfra.WaitForUnlock()

	vmio.LockInfrastructureById(descriptor.Id, descriptor.InfraId)
	var info vmio.InfrastructureInfo = vmio.InfrastructureInfo{
		Format: "",
		Infra: model.Infrastructure{
			ProjectId: descriptor.Id,
			Id:        descriptor.InfraId,
		},
	}
	err = info.Delete()
	vmio.UnlockInfrastructureById(descriptor.Id, descriptor.InfraId)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	descriptor.InfraId = ""
	descriptor.InfraName = ""

	utils.PrintlnWarning(fmt.Sprintf("Reopening Project '%s' ...", descriptor.Name))

	project.Open = true

	err = vmio.SaveProject(project)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning("Updating global indexes ...")

	descriptor.Open = true

	err = UpdateIndexWithProjectsDescriptor(descriptor, true)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnSuccess(fmt.Sprintf("Delete of Infrastructure named : %s completed successfully!!", Name))

	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) BackupInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	Name := ""
	File := ""
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "file" == CorrectInput(option[0]) {
			File = option[1]
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name field not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	if File == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "File Path field not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id:        descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	if !strings.HasSuffix(strings.ToLower(File), ".vmkube") {
		utils.PrintlnWarning("File extension changed, cause it was not standard...")
		File += ".vmkube"
	}

	err = infrastructure.Save(File)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnSuccess(fmt.Sprintf("Backup Infrastructure named : %s completed successfully!!", Name))

	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) RecoverInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
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
	for _, option := range request.Arguments.Options {
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
			ProjectName = option[1]
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
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name field not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	if File == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "File Path field not provided",
			Status:  false}, errors.New("Unable to execute task")
	}

	if !strings.HasSuffix(strings.ToLower(File), ".vmkube") {
		utils.PrintlnError("File extension wrong, it can cause panic in the system")
		response := Response{
			Status:  false,
			Message: fmt.Sprintf("Wrong file type for '%s', expected vmkube extension!!", File),
		}
		return response, errors.New("Unable to execute task")
	}

	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)

	existsInfrastructure := (err == nil) && descriptor.InfraId != ""

	if existsInfrastructure {
		utils.PrintlnWarning(fmt.Sprintf("Infrastructure '%s' already exists ...", Name))
	} else {
		utils.PrintlnInfo(fmt.Sprintf("Infrastructure '%s' doesn't exists ... Creating from scratch", Name))
	}

	DeleteFromDescriptor := false

	AllowOverride := false

	if existsInfrastructure && CorrectInput(ProjectName) != CorrectInput(descriptor.Name) {
		if !Override {
			return Response{
				Message: fmt.Sprintf("Project named %s already associated with infrastructure %s, and no override clause specified", descriptor.Name, descriptor.InfraName),
				Status:  false,
			}, errors.New("Unable to execute task")
		} else {
			if !Force {
				DeleteFromDescriptor = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with override of existing Project named '%s' and infrastructure '%s'?", descriptor.Name, descriptor.InfraName))
				if !DeleteFromDescriptor {
					return Response{
						Message: "User task interruption",
						Status:  false,
					}, errors.New("Unable to execute task")
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
				Status:  false,
			}, errors.New("Unable to execute task")
		} else {
			if !Force {
				AllowOverride = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with override of existing Infrastructure named '%s' associated with project '%s'?", descriptor.InfraName, descriptor.Name))
				if !AllowOverride {
					return Response{
						Message: "User task interruption",
						Status:  false,
					}, errors.New("Unable to execute task")
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
	projectDescriptor, err = vmio.GetProjectDescriptor(ProjectName)
	existsProject := (err == nil) && projectDescriptor.Id != ""

	if existsProject {
		if ProjectName == "" {
			PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
			return Response{
				Message: "Project Name field not provided, mandatory when we have an exiting one with same name",
				Status:  false,
			}, errors.New("Unable to execute task")
		}

		if projectDescriptor.InfraId != descriptor.InfraId {
			if !Override {
				return Response{
					Message: fmt.Sprintf("Selected Project Name %s is already associated with another infrastructure %s, and no override clause specified", projectDescriptor.Name, projectDescriptor.InfraName),
					Status:  false,
				}, errors.New("Unable to execute task")
			} else {
				if !Force {
					DeleteFromProjectDescriptor = utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with delete of existing Project named '%s' and related Infrastructure named '%s'?", projectDescriptor.Name, projectDescriptor.InfraName))
					if !DeleteFromProjectDescriptor {
						return Response{
							Message: "User task interruption",
							Status:  false,
						}, errors.New("Unable to execute task")
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
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	//utils.PrintlnInfo(fmt.Sprintf("Infrastructure : %s", utils.GetJSONFromObj(infrastructure, true)))
	utils.PrintlnWarning(fmt.Sprintf("Validating recovered Infrastructure '%s'...", Name))

	errorsList := infrastructure.Validate()

	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete Recover for Infrastructure '%s'!!", Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error validating infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status:  false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning(fmt.Sprintf("Creating new Project %s from Infrastructure '%s'...", ProjectName, Name))

	newProject, err := InfrastructureToProject(infrastructure, ProjectName)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	errorsList = newProject.Validate()

	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete Project Re-create '%s' from Infrastructure '%s'!!", ProjectName, Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error bulding project '%s' infrastructure -> '%s' : ", ProjectName, Name), errorsList)
		response := Response{
			Status:  false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}

	ProjectBackup := ""
	InfraBackup := ""
	if DeleteFromDescriptor {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project '%s' and Infrastructure '%s'...", descriptor.Name, descriptor.InfraName))
		if Backup {
			AllowBackup := Force
			if !AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Project named '%s'?", descriptor.Name))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if !strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				ProjectBackup = fmt.Sprintf("%s.project-%s-%s.json", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name))
				project, err := vmio.LoadProject(descriptor.Id)
				if err == nil {
					vmio.ExportUserProject(project, ProjectBackup, "json", true)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Project '%s' backup at : %s", descriptor.Name, ProjectBackup))
				}
			}
			AllowBackup = Force
			if !AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?", descriptor.InfraName))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if !strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				InfraBackup = fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(descriptor.Id), utils.NameToFileFormat(descriptor.Name), utils.IdToFileFormat(descriptor.InfraId), utils.NameToFileFormat(descriptor.InfraName))
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
				FoundForce = true
			} else if "name" == CorrectInput(request.Arguments.Options[i][0]) {
				request.Arguments.Options[i][1] = descriptor.Name
				FoundName = true
			}
		}
		if !FoundName {
			request.Arguments.Options = append(request.Arguments.Options, []string{"name", descriptor.Name})
		}
		if !FoundForce {
			request.Arguments.Options = append(request.Arguments.Options, []string{"force", "true"})
		}
		resp, err := request.DeleteProject(recoverHelpersFunc)
		if err != nil {
			return resp, err
		}
		time.Sleep(time.Second * 2)
	}
	ProjectProjectBackup := ""
	ProjectInfraBackup := ""
	if DeleteFromProjectDescriptor {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project '%s' and Infrastructure '%s'...", projectDescriptor.Name, projectDescriptor.InfraName))
		if Backup {
			AllowBackup := Force
			if !AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Project named '%s'?", projectDescriptor.Name))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if !strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				ProjectProjectBackup = fmt.Sprintf("%s.project-%s-%s.json", folder, utils.IdToFileFormat(projectDescriptor.Id), utils.NameToFileFormat(projectDescriptor.Name))
				project, err := vmio.LoadProject(projectDescriptor.Id)
				if err == nil {
					vmio.ExportUserProject(project, ProjectProjectBackup, "json", true)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Project '%s' backup at : %s", projectDescriptor.Name, ProjectProjectBackup))
				}
			}
			AllowBackup = Force
			if !AllowBackup {
				AllowBackup = utils.RequestConfirmation(fmt.Sprintf("Do you want backup Infrastructure named '%s'?", projectDescriptor.InfraName))
			}
			if AllowBackup {
				folder := model.GetEmergencyFolder()
				if BackupDir != "" {
					folder = BackupDir
				}
				if !strings.HasSuffix(folder, string(os.PathSeparator)) {
					folder += string(os.PathSeparator)
				}
				ProjectInfraBackup = fmt.Sprintf("%s.prj-%s-%s-infra-export-%s-%s.vmkube", folder, utils.IdToFileFormat(projectDescriptor.Id), utils.NameToFileFormat(projectDescriptor.Name), utils.IdToFileFormat(projectDescriptor.InfraId), utils.NameToFileFormat(projectDescriptor.InfraName))
				infra, err := vmio.LoadInfrastructure(projectDescriptor.Id)
				if err == nil {
					infra.Save(ProjectInfraBackup)
					utils.PrintlnImportant(fmt.Sprintf("Emergency Infrastructure '%s' backup at : %s", projectDescriptor.InfraName, ProjectInfraBackup))
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
				FoundForce = true
			} else if "name" == CorrectInput(request.Arguments.Options[i][0]) {
				request.Arguments.Options[i][1] = projectDescriptor.Name
				FoundName = true
			}
		}
		if !FoundName {
			request.Arguments.Options = append(request.Arguments.Options, []string{"name", projectDescriptor.Name})
		}
		if !FoundForce {
			request.Arguments.Options = append(request.Arguments.Options, []string{"force", "true"})
		}
		resp, err := request.DeleteProject(recoverHelpersFunc)
		if err != nil {
			return resp, err
		}
		time.Sleep(time.Second * 2)
	}
	utils.PrintlnWarning(fmt.Sprintf("Saving new Project '%s'...", ProjectName))

	err = vmio.SaveProject(newProject)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning(fmt.Sprintf("Saving recovered Infrastructure '%s'...", Name))

	err = vmio.SaveInfrastructure(infrastructure)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	newDescriptor := model.ProjectsDescriptor{
		Id:        newProject.Id,
		Name:      newProject.Name,
		Open:      false,
		Active:    false,
		Synced:    true,
		InfraId:   "",
		InfraName: "",
	}

	utils.PrintlnWarning(fmt.Sprintf("Defining new Indexes from new Project %s ...", ProjectName))

	err = UpdateIndexWithProjectsDescriptor(newDescriptor, true)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning(fmt.Sprintf("Building new Project '%s'...", ProjectName))
	actionCouples := make([]tasks.ActivityCouple, 0)
	creationCouples, err := tasks.GetTaskActivities(newProject, infrastructure, tasks.CreateMachine)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	inspectCouples, err := tasks.GetTaskActivities(newProject, infrastructure, tasks.MachineInspect)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	ipAddressCouples, err := tasks.GetTaskActivities(newProject, infrastructure, tasks.MachineIPAddress)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	stopCouples, err := tasks.GetTaskActivities(newProject, infrastructure, tasks.StopMachine)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	extendsDiskCouples, err := tasks.GetTaskActivities(newProject, infrastructure, tasks.MachineExtendsDisk)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	actionCouples = append(actionCouples, creationCouples...)
	actionCouples = append(actionCouples, inspectCouples...)
	actionCouples = append(actionCouples, ipAddressCouples...)
	actionCouples = append(actionCouples, stopCouples...)
	actionCouples = append(actionCouples, extendsDiskCouples...)

	utils.PrintlnImportant("Now Proceding with machine creation ...!!")
	NumThreads := Threads
	if runtime.NumCPU()-1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))

	var fixInfraValue int = len(actionCouples)
	errorsList = ExecuteInfrastructureActions(infrastructure, actionCouples, NumThreads, func(task tasks.SchedulerTask) {
		go func(task tasks.SchedulerTask) {
			for i := 0; i < len(task.Jobs); i++ {
				response := strings.Split(fmt.Sprintf("%s", task.Jobs[i].GetRunnable().Response()), "|")
				if len(response) > 3 {
					if len(response) > 4 {
						if response[0] == "ip" {
							instanceId := response[1]
							ipAddress := response[2]
							json := ""
							log := response[3] + response[4]
							FixInfrastructureElementValue(&infrastructure, instanceId, ipAddress, json, log)
						} else if response[0] == "json" {
							instanceId := response[1]
							ipAddress := ""
							json := response[2]
							log := response[3] + response[4]
							FixInfrastructureElementValue(&infrastructure, instanceId, ipAddress, json, log)
						}
					} else {
						instanceId := response[1]
						ipAddress := ""
						json := ""
						log := response[2] + response[3]
						FixInfrastructureElementValue(&infrastructure, instanceId, ipAddress, json, log)
					}
				} else if len(response) > 2 {
					planId := response[0]
					log := response[1] + "\n" + response[2]
					FixInfrastructureIntallationLogs(&infrastructure, planId, log)
				}
				fixInfraValue--
			}
		}(task)
	})

	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete Build of project '%s' : Errors building Infrastructure : '%s'!!", ProjectName, Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error building Infrastructure -> '%s' : ", Name), errorsList)
		response := Response{
			Status:  false,
			Message: message,
		}
		exclusionList, _ := FilterForExistState(infrastructure)
		rollbackActions, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.DestroyMachine, exclusionList)

		if err != nil {
			return response, errors.New("Unable to execute task")
		}
		utils.PrintlnWarning(fmt.Sprintf("Executing rollback for Project '%s' Infrastrcucture '%s'...", ProjectName, Name))
		if !utils.NO_COLORS {
			time.Sleep(4 * time.Second)
		}
		if existsInfrastructure || existsProject && (ProjectBackup != "" || ProjectProjectBackup != "" || InfraBackup != "" || ProjectInfraBackup != "") {
			utils.PrintlnImportant("Check logs for backup activities, you can use for recovery...")
		} else {
			utils.PrintlnImportant("No backup activities, you can not use recovery utils...")
		}
		ExecuteInfrastructureActions(infrastructure, rollbackActions, NumThreads, func(task tasks.SchedulerTask) {})
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnWarning(fmt.Sprintf("Waiting for Instance recovery information in Project '%s' Infrastrcucture '%s'", ProjectName, Name))

	for fixInfraValue > 0 {
		time.Sleep(1 * time.Second)
	}

	utils.PrintlnWarning(fmt.Sprintf("Updating recovered Infrastructure '%s' with new Instance data...", Name))

	err = vmio.SaveInfrastructure(infrastructure)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	newDescriptor = model.ProjectsDescriptor{
		Id:        newProject.Id,
		Name:      newProject.Name,
		Open:      false,
		Active:    false,
		Synced:    true,
		InfraId:   infrastructure.Id,
		InfraName: infrastructure.Name,
	}

	utils.PrintlnWarning(fmt.Sprintf("Defining new Indexes for new Project %s with recovered Infrastructure '%s'...", ProjectName, Name))

	err = UpdateIndexWithProjectsDescriptor(newDescriptor, true)

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnSuccess(fmt.Sprintf("Recovery for Infrastructure named : %s completed successfully!!", Name))

	if ProjectBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project backup file : %s", ProjectBackup))
		err = model.DeleteIfExists(ProjectBackup)
	}
	if ProjectProjectBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Project backup file : %s", ProjectProjectBackup))
		err = model.DeleteIfExists(ProjectProjectBackup)
	}
	if InfraBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Infrastructure backup file : %s", InfraBackup))
		err = model.DeleteIfExists(InfraBackup)
	}
	if ProjectInfraBackup != "" {
		utils.PrintlnWarning(fmt.Sprintf("Removing Infrastructure backup file : %s", ProjectInfraBackup))
		err = model.DeleteIfExists(ProjectInfraBackup)
	}
	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) StartInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
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
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if !descriptor.Active {
		response := Response{
			Status:  false,
			Message: fmt.Sprintf("Infrastructure %s not active, unable to start it. Please close project to perform this task", Name),
		}
		return response, errors.New("Unable to execute task")
	}
	if !Force {
		AllowInfraStart := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with start machines process for Infrastructure named '%s'?", Name))
		if !AllowInfraStart {
			response := Response{
				Status:  false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	iFaceInfra := vmio.IFaceInfra{
		Id:        descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnImportant("Now Proceding with machines start ...!!")
	NumThreads := Threads
	if runtime.NumCPU()-1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))

	exclusionList, _ := FilterForExistAndNotRunningState(infrastructure, procedures.Machine_State_Running)

	startMachineActivities, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StartMachine, exclusionList)

	if len(startMachineActivities) == 0 {
		utils.PrintlnImportant("No machine to start ...")
		response := Response{
			Status:  true,
			Message: "Success",
		}
		return response, nil
	}

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, startMachineActivities, NumThreads, func(task tasks.SchedulerTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete start machines process : Errors starting Infrastructure : %s!!", Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error starting infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status:  false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnSuccess(fmt.Sprintf("Start Infrastructure named : %s completed successfully!!", Name))

	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) StopInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
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
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	if !Force {
		AllowInfraStop := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with stop machines process for Infrastructure named '%s'?", Name))
		if !AllowInfraStop {
			response := Response{
				Status:  false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id:        descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnImportant("Now Proceding with machines stop ...!!")
	NumThreads := Threads
	if runtime.NumCPU()-1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))

	exclusionList, _ := FilterForExistAndRunningState(infrastructure, procedures.Machine_State_Running)

	stopMachineActivities, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.StopMachine, exclusionList)

	if len(stopMachineActivities) == 0 {
		utils.PrintlnImportant("No machine to stop ...")
		response := Response{
			Status:  true,
			Message: "Success",
		}
		return response, nil
	}

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, stopMachineActivities, NumThreads, func(task tasks.SchedulerTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete stop machines process : Errors stopping Infrastructure : %s!!", Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error stopping infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status:  false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnSuccess(fmt.Sprintf("Stop Infrastructure named : %s completed successfully!!", Name))

	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) RestartInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	Name := ""
	Force := false
	Threads := 1
	Overclock := false
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
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
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if !descriptor.Active {
		response := Response{
			Status:  false,
			Message: fmt.Sprintf("Infrastructure %s not active, unable to restart it. Please close project to perform this task", Name),
		}
		return response, errors.New("Unable to execute task")
	}
	if !Force {
		AllowInfraStop := utils.RequestConfirmation(fmt.Sprintf("Do you want proceed with restart machines process for Infrastructure named '%s'?", Name))
		if !AllowInfraStop {
			response := Response{
				Status:  false,
				Message: "User task interruption",
			}
			return response, errors.New("Unable to execute task")
		}
	}
	iFaceInfra := vmio.IFaceInfra{
		Id:        descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnImportant("Now Proceding with machines restart ...!!")
	NumThreads := Threads
	if runtime.NumCPU()-1 < Threads && !Overclock {
		NumThreads = runtime.NumCPU() - 1
		utils.PrintlnWarning(fmt.Sprintf("Number of threads in order to available processors : %d", NumThreads))
	}
	utils.PrintlnImportant(fmt.Sprintf("Number of threads assigned to scheduler : %d", NumThreads))

	exclusionList, _ := FilterForExistAndRunningState(infrastructure, procedures.Machine_State_Running)

	restartMachineActivities, err := tasks.GetPostBuildTaskActivities(infrastructure, tasks.RestartMachine, exclusionList)

	if len(restartMachineActivities) == 0 {
		utils.PrintlnImportant("No machine to restart ...")
		response := Response{
			Status:  true,
			Message: "Success",
		}
		return response, nil
	}

	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	errorsList := ExecuteInfrastructureActions(infrastructure, restartMachineActivities, NumThreads, func(task tasks.SchedulerTask) {
	})
	if len(errorsList) > 0 {
		utils.PrintlnError(fmt.Sprintf("Unable to complete restart machines process : Errors restarting Infrastructure : %s!!", Name))
		_, message := vmio.StripErrorMessages(fmt.Sprintf("Error restarting infrastructure -> %s : ", Name), errorsList)
		response := Response{
			Status:  false,
			Message: message,
		}
		return response, errors.New("Unable to execute task")
	}

	utils.PrintlnSuccess(fmt.Sprintf("Restart Infrastructure named : %s completed successfully!!", Name))

	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) ListInfras(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
		if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		}
	}
	indexes, error := vmio.LoadIndex()
	if error != nil {
		response := Response{
			Status:  false,
			Message: error.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if len(indexes.Projects) > 0 {
		utils.PrintlnImportant(fmt.Sprintf("%s  %s  %s  %s", utils.StrPad("Infrastructure Id", 40), utils.StrPad("Infrastructure Name", 40), utils.StrPad("Active", 6), utils.StrPad("Synced", 6)))
	} else {
		utils.PrintlnImportant("No Infrastructures found")
	}
	for _, index := range indexes.Projects {
		active := "no"
		synced := "no"
		if index.Active {
			active = "yes"
		}
		if index.Synced {
			synced = "yes"
		}
		fmt.Printf("%s  %s  %s  %s\n", utils.StrPad(index.InfraId, 40), utils.StrPad(index.InfraName, 40), utils.StrPad("  "+active, 6), utils.StrPad("  "+synced, 6))

	}
	response := Response{
		Status:  true,
		Message: "Success",
	}
	return response, nil
}

func (request *CmdRequest) StatusInfra(recoverHelpersFunc func() []CommandHelper) (Response, error) {
	Name := ""
	Format := "json"
	Details := false
	utils.NO_COLORS = false
	for _, option := range request.Arguments.Options {
		if "infra-name" == CorrectInput(option[0]) {
			Name = option[1]
		} else if "no-colors" == CorrectInput(option[0]) {
			utils.NO_COLORS = GetBoolean(option[1])
		} else if "show-full" == CorrectInput(option[0]) {
			Details = GetBoolean(option[1])
		} else if "format" == CorrectInput(option[0]) {
			Format = CorrectInput(option[1])
		}
	}
	if Name == "" {
		PrintCommandHelper(request.TypeStr, request.SubTypeStr, recoverHelpersFunc)
		return Response{
			Message: "Infrastrcuture Name not provided",
			Status:  false}, errors.New("Unable to execute task")
	}
	descriptor, err := vmio.GetInfrastructureProjectDescriptor(Name)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	iFaceInfra := vmio.IFaceInfra{
		Id:        descriptor.InfraId,
		ProjectId: descriptor.Id,
	}
	iFaceInfra.WaitForUnlock()

	infrastructure, err := vmio.LoadInfrastructure(descriptor.Id)
	if err != nil {
		response := Response{
			Status:  false,
			Message: err.Error(),
		}
		return response, errors.New("Unable to execute task")
	}
	if Details {
		var bytesArray []byte = make([]byte, 0)
		var err error
		if "json" == Format {
			bytesArray, err = utils.GetJSONFromElem(infrastructure, true)
		} else if "yaml" == Format {
			bytesArray, err = utils.GetYAMLFromElem(infrastructure)
		} else if "xml" == Format {
			bytesArray, err = utils.GetXMLFromElem(infrastructure, true)
		} else {
			response := Response{
				Status:  false,
				Message: "Sample Format '" + Format + "' not supported!!",
			}
			return response, errors.New("Unable to execute task")
		}
		if err != nil {
			response := Response{
				Status:  false,
				Message: err.Error(),
			}
			return response, errors.New("Unable to execute task")
		}
		fmt.Printf("%s\n", bytesArray)
	} else {
		fmt.Printf("Id: %s\nInfrastructure: %s\nModified: %s\n", infrastructure.Id, infrastructure.Name, BoolToString(infrastructure.Altered))
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
		fmt.Printf("Errors: %s\nLast Message: %s\n", BoolToString(infrastructure.Errors), infrastructure.LastMessage)
		fmt.Printf("Domains: %d\n", len(infrastructure.Domains))
		for _, domain := range infrastructure.Domains {
			num, options := vmio.StripOptions(domain.Options)
			fmt.Printf("Domain: %s (Id: %s) - Options [%d] :%s\n", domain.Name, domain.Id, num, options)
			fmt.Printf("Networks: %d\n", len(domain.Networks))
			for _, network := range domain.Networks {
				num, options := vmio.StripOptions(network.Options)
				fmt.Printf("   Network: %s (Id: %s) - Options [%d] :%s\n", network.Name, network.Id, num, options)
				fmt.Printf("   Local Instances: %d\n", len(network.LocalInstances))
				instancesMap := make(map[string]string)
				for _, instance := range network.LocalInstances {
					instancesMap[instance.Id] = instance.Name
					var instanceState procedures.MachineState
					instanceState, _ = ExistInstance(infrastructure, instance, model.CloudInstance{}, false, instance.Id)
					fmt.Printf("      Local Instance: %s (Id: %s) - Driver: %s - OS : %s:%s - IP Address: %s State: %s\n", instance.Name, instance.Id, instance.Driver, instance.OSType, instance.OSVersion, strings.TrimSpace(instance.IPAddress), instanceState.String())
				}
				fmt.Printf("   Cloud Instances: %d\n", len(network.CloudInstances))
				for _, instance := range network.CloudInstances {
					instancesMap[instance.Id] = instance.Name
					var instanceState procedures.MachineState
					instanceState, _ = ExistInstance(infrastructure, model.LocalInstance{}, instance, true, instance.Id)
					num, options := vmio.StripOptions(instance.Options)
					fmt.Printf("      Cloud Instance: %s (Id: %s) - Driver: %s - IP Address: %s - Options [%d] :%s State: %s\n", instance.Name, instance.Id, instance.Driver, strings.TrimSpace(instance.IPAddress), num, options, instanceState.String())
				}
				fmt.Printf("   Installation Plans: %d\n", len(network.Installations))
				for _, installation := range network.Installations {
					machineName, ok := instancesMap[installation.InstanceId]
					if !ok {
						machineName = "<invalid>"
					}
					fmt.Printf("      Plan: Id: %s - Instance: %s [Id: %s] - Success: %s - Cloud: %s - Envoronment : %s  Role: %s  Type: %s\n", installation.Id, machineName, installation.InstanceId, BoolToString(installation.Success), BoolToString(installation.IsCloud), model.InstanceEnvironmentToString(installation.Environment), model.InstanceRoleToString(installation.Role), model.InstanceInstallationToString(installation.Type))
				}
			}
		}
	}
	return Response{
		Message: "Success",
		Status:  true}, nil
}
