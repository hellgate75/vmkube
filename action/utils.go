package action

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"vmkube/model"
	"vmkube/procedures"
	"vmkube/scheduler"
	"vmkube/tasks"
	"vmkube/term"
	"vmkube/utils"
	"vmkube/vmio"
)

func ParseCommandArguments(args []string) (*CmdArguments, error) {
	arguments := CmdArguments{}
	success := arguments.Parse(args[1:])
	if success {
		return &arguments, nil
	} else {
		return &arguments, errors.New("Unable to Parse Command Line")
	}
}

func ParseCommandLine(args []string) (CmdRequest, error) {
	request := CmdRequest{}
	arguments, error := ParseCommandArguments(args)
	if error == nil {
		request.TypeStr = arguments.Cmd
		request.Type = arguments.CmdType
		request.SubTypeStr = arguments.SubCmd
		request.SubType = arguments.SubCmdType
		request.HelpType = arguments.SubCmdHelpType
		request.Arguments = arguments
	}
	return request, error
}

func CmdParseElement(value string) (CmdElementType, error) {
	switch CorrectInput(value) {
	case "machine":
		return LMachine, nil
	case "cloud-machine":
		return CLMachine, nil
	case "network":
		return SNetwork, nil
	case "domain":
		return SDomain, nil
	case "project":
		return SProject, nil
	case "plan":
		return SPlan, nil
	default:
		return NoElement, errors.New("Element '" + value + "' is not an infratructure element. Available ones : Machine, Cloud-Machine, Network, Domain, Plan, Project")

	}
}

func CorrectInput(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func GetBoolean(input string) bool {
	return CorrectInput(input) == "true"
}

func BoolToString(input bool) string {
	if input {
		return "yes"
	} else  {
		return "no"
	}
}

func GetInteger(input string, defaultValue int) int {
	num, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return num
}

func ProjectToInfrastructure(project model.Project) (model.Infrastructure, error) {
	infrastructure := model.Infrastructure{}
	infrastructure.Id = NewUUIDString()
	infrastructure.ProjectId = project.Id
	infrastructure.Name = project.Name
	for _, domain := range project.Domains {
		newDomain := model.Domain{
			Id:       NewUUIDString(),
			Name:     domain.Name,
			Options:  domain.Options,
			Networks: []model.Network{},
		}
		for _, network := range domain.Networks {
			newNetwork := model.Network{
				Id:             NewUUIDString(),
				Name:           network.Name,
				Options:        network.Options,
				LocalInstances: []model.LocalInstance{},
				CloudInstances: []model.CloudInstance{},
				Installations:  []model.Installation{},
			}
			machineRecogniseMap := make(map[string]string)
			for _, machine := range network.LocalMachines {
				var disks []model.Disk = make([]model.Disk, 0)
				disks = append(disks, model.Disk{
					Id:   NewUUIDString(),
					Name: "sda0",
					Size: machine.DiskSize,
					Type: 0,
				})
				instanceId := NewUUIDString()
				instance := model.LocalInstance{
					Id:        instanceId,
					Name:      machine.Name,
					Options:   machine.Options,
					Cpus:      machine.Cpus,
					Memory:    machine.Memory,
					Disks:     disks,
					Driver:    machine.Driver,
					Engine:    model.ToInstanceEngineOpt(machine.Engine),
					Swarm:     model.ToInstanceSwarmOpt(machine.Swarm),
					Hostname:  machine.Hostname,
					IPAddress: "",
					NoShare:   machine.NoShare,
					OSType:    machine.OSType,
					OSVersion: machine.OSVersion,
					Roles:     machine.Roles,
					MachineId: machine.Id,
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId:   infrastructure.Id,
						ElementId: instanceId,
						LogLines:  []string{},
					},
				}
				if _, ok := machineRecogniseMap[machine.Id]; ok {
					return infrastructure, errors.New("Duplicate Machine Id in Project : " + machine.Id)
				}
				machineRecogniseMap[machine.Id] = instance.Id
				newNetwork.LocalInstances = append(newNetwork.LocalInstances, instance)
			}
			for _, machine := range network.CloudMachines {
				instanceId := NewUUIDString()
				instance := model.CloudInstance{
					Id:        instanceId,
					Name:      machine.Name,
					Driver:    machine.Driver,
					Hostname:  machine.Hostname,
					IPAddress: "",
					Options:   machine.Options,
					Roles:     machine.Roles,
					MachineId: machine.Id,
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId:   infrastructure.Id,
						ElementId: instanceId,
						LogLines:  []string{},
					},
				}
				if _, ok := machineRecogniseMap[machine.Id]; ok {
					return infrastructure, errors.New("Duplicate Machine Id in Project : " + machine.Id)
				}
				machineRecogniseMap[machine.Id] = instance.Id
				newNetwork.CloudInstances = append(newNetwork.CloudInstances, instance)
			}
			for _, plan := range network.Installations {
				if _, ok := machineRecogniseMap[plan.MachineId]; !ok {
					return infrastructure, errors.New("Invalid machine reference in plan : " + plan.MachineId)
				}
				instanceId, _ := machineRecogniseMap[plan.MachineId]
				installationId := NewUUIDString()
				installation := model.Installation{
					Id:            installationId,
					Environment:   model.ToInstanceEnvironment(plan.Environment),
					Role:          model.ToInstanceRole(plan.Role),
					Type:          model.ToInstanceInstallation(plan.Type),
					Errors:        false,
					InstanceId:    instanceId,
					IsCloud:       plan.IsCloud,
					Success:       false,
					LastExecution: time.Now(),
					LastMessage:   "",
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId:   infrastructure.Id,
						ElementId: installationId,
						LogLines:  []string{},
					},
					Plan: plan,
				}
				newNetwork.Installations = append(newNetwork.Installations, installation)
			}
			newDomain.Networks = append(newDomain.Networks, newNetwork)
		}
		infrastructure.Domains = append(infrastructure.Domains, newDomain)
	}
	return infrastructure, nil
}

func InfrastructureToProject(infrastructure model.Infrastructure, projectName string) (model.Project, error) {
	newProject := model.Project{}
	newProject.Id = NewUUIDString()
	newProject.Name = projectName
	newProject.Created = time.Now()
	newProject.Modified = time.Now()
	newProject.Errors = false
	newProject.LastMessage = fmt.Sprintf("Imported from infrastructure %s", infrastructure.Name)
	for i := 0; i < len(infrastructure.Domains); i++ {
		newDomain := model.MachineDomain{
			Id:       NewUUIDString(),
			Name:     infrastructure.Domains[i].Name,
			Options:  infrastructure.Domains[i].Options,
			Networks: []model.MachineNetwork{},
		}
		for j := 0; j < len(infrastructure.Domains[i].Networks); j++ {
			newNetwork := model.MachineNetwork{
				Id:            NewUUIDString(),
				Name:          infrastructure.Domains[i].Networks[j].Name,
				Options:       infrastructure.Domains[i].Networks[j].Options,
				LocalMachines: []model.LocalMachine{},
				CloudMachines: []model.CloudMachine{},
				Installations: []model.InstallationPlan{},
			}
			machineRecogniseMap := make(map[string]string)
			for k := 0; k < len(infrastructure.Domains[i].Networks[j].LocalInstances); k++ {
				if infrastructure.Domains[i].Networks[j].LocalInstances[k].Disabled {
					infrastructure.Domains[i].Networks[j].LocalInstances[k].Disabled = false
				}
				machineId := NewUUIDString()
				infrastructure.Domains[i].Networks[j].LocalInstances[k].MachineId = machineId
				machine := model.LocalMachine{
					Id:        machineId,
					Name:      infrastructure.Domains[i].Networks[j].LocalInstances[k].Name,
					Options:   infrastructure.Domains[i].Networks[j].LocalInstances[k].Options,
					Cpus:      infrastructure.Domains[i].Networks[j].LocalInstances[k].Cpus,
					Memory:    infrastructure.Domains[i].Networks[j].LocalInstances[k].Memory,
					DiskSize:  infrastructure.Domains[i].Networks[j].LocalInstances[k].Disks[0].Size,
					Driver:    infrastructure.Domains[i].Networks[j].LocalInstances[k].Driver,
					Engine:    model.ToMachineEngineOpt(infrastructure.Domains[i].Networks[j].LocalInstances[k].Engine),
					Swarm:     model.ToMachineSwarmOpt(infrastructure.Domains[i].Networks[j].LocalInstances[k].Swarm),
					Hostname:  infrastructure.Domains[i].Networks[j].LocalInstances[k].Hostname,
					NoShare:   infrastructure.Domains[i].Networks[j].LocalInstances[k].NoShare,
					OSType:    infrastructure.Domains[i].Networks[j].LocalInstances[k].OSType,
					OSVersion: infrastructure.Domains[i].Networks[j].LocalInstances[k].OSVersion,
					Roles:     infrastructure.Domains[i].Networks[j].LocalInstances[k].Roles,
				}
				if _, ok := machineRecogniseMap[infrastructure.Domains[i].Networks[j].LocalInstances[k].Id]; ok {
					return newProject, errors.New("Duplicate Instance Id in Project : " + infrastructure.Domains[i].Networks[j].LocalInstances[k].Id)
				}
				machineRecogniseMap[infrastructure.Domains[i].Networks[j].LocalInstances[k].Id] = machineId
				newNetwork.LocalMachines = append(newNetwork.LocalMachines, machine)
			}
			for k := 0; k < len(infrastructure.Domains[i].Networks[j].CloudInstances); k++ {
				if infrastructure.Domains[i].Networks[j].CloudInstances[k].Disabled {
					infrastructure.Domains[i].Networks[j].CloudInstances[k].Disabled = false
				}
				machineId := NewUUIDString()
				infrastructure.Domains[i].Networks[j].CloudInstances[k].MachineId = machineId
				machine := model.CloudMachine{
					Id:       machineId,
					Name:     infrastructure.Domains[i].Networks[j].CloudInstances[k].Name,
					Driver:   infrastructure.Domains[i].Networks[j].CloudInstances[k].Driver,
					Hostname: infrastructure.Domains[i].Networks[j].CloudInstances[k].Hostname,
					Options:  infrastructure.Domains[i].Networks[j].CloudInstances[k].Options,
					Roles:    infrastructure.Domains[i].Networks[j].CloudInstances[k].Roles,
				}
				if _, ok := machineRecogniseMap[infrastructure.Domains[i].Networks[j].CloudInstances[k].Id]; ok {
					return newProject, errors.New("Duplicate Instance Id in Project : " + infrastructure.Domains[i].Networks[j].CloudInstances[k].Id)
				}
				machineRecogniseMap[infrastructure.Domains[i].Networks[j].CloudInstances[k].Id] = machineId
				newNetwork.CloudMachines = append(newNetwork.CloudMachines, machine)
			}
			for k := 0; k < len(infrastructure.Domains[i].Networks[j].Installations); k++ {
				if _, ok := machineRecogniseMap[infrastructure.Domains[i].Networks[j].Installations[k].InstanceId]; !ok {
					return newProject, errors.New("Invalid instance reference in plan : " + infrastructure.Domains[i].Networks[j].Installations[k].InstanceId)
				}
				machineId, _ := machineRecogniseMap[infrastructure.Domains[i].Networks[j].Installations[k].InstanceId]
				installation := infrastructure.Domains[i].Networks[j].Installations[k].Plan
				installation.Id = NewUUIDString()
				installation.MachineId = machineId
				installation.IsCloud = infrastructure.Domains[i].Networks[j].Installations[k].IsCloud
				newNetwork.Installations = append(newNetwork.Installations, installation)
			}
			newDomain.Networks = append(newDomain.Networks, newNetwork)
		}
		newProject.Domains = append(newProject.Domains, newDomain)
	}
	return newProject, nil
}

func UpdateIndexWithProject(project model.Project) error {
	indexes, err := vmio.LoadIndex()

	if err != nil {
		return err
	}

	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()

	indexes, err = vmio.LoadIndex()

	if err != nil {
		return err
	}

	var synced, active bool
	synced = true
	active = false

	InfraId := ""
	InfraName := ""
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _, prj := range indexes.Projects {
		if CorrectInput(prj.Name) != CorrectInput(project.Name) {
			NewIndexes = append(NewIndexes, prj)
		} else {
			synced = prj.Synced
			active = prj.InfraId != "" && !project.Open
			InfraId = prj.InfraId
			InfraName = prj.InfraName
		}
	}

	vmio.LockIndex(indexes)

	indexes.Projects = append(indexes.Projects, model.ProjectsDescriptor{
		Id:        project.Id,
		Name:      project.Name,
		Open:      project.Open,
		Synced:    synced,
		Active:    active,
		InfraId:   InfraId,
		InfraName: InfraName,
	})

	indexes.Projects = NewIndexes

	err = vmio.SaveIndex(indexes)

	vmio.UnlockIndex(indexes)

	return err
}

func UpdateIndexWithProjectStates(project model.Project, active bool, synced bool) error {
	indexes, err := vmio.LoadIndex()

	if err != nil {
		return err
	}

	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()

	indexes, err = vmio.LoadIndex()

	if err != nil {
		return err
	}

	InfraId := ""
	InfraName := ""
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _, prj := range indexes.Projects {
		if CorrectInput(prj.Name) != CorrectInput(project.Name) {
			NewIndexes = append(NewIndexes, prj)
		} else {
			InfraId = prj.InfraId
			InfraName = prj.InfraName
		}
	}

	vmio.LockIndex(indexes)

	indexes.Projects = append(indexes.Projects, model.ProjectsDescriptor{
		Id:        project.Id,
		Name:      project.Name,
		Open:      project.Open,
		Synced:    synced,
		Active:    active,
		InfraId:   InfraId,
		InfraName: InfraName,
	})

	indexes.Projects = NewIndexes

	err = vmio.SaveIndex(indexes)

	vmio.UnlockIndex(indexes)

	return err
}

func UpdateIndexWithProjectsDescriptor(project model.ProjectsDescriptor, addDescriptor bool) error {
	indexes, err := vmio.LoadIndex()

	if err != nil {
		return err
	}

	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()

	indexes, err = vmio.LoadIndex()

	if err != nil {
		return err
	}

	NewIndexes := make([]model.ProjectsDescriptor, 0)

	for _, prj := range indexes.Projects {
		if prj.Id != project.Id {
			NewIndexes = append(NewIndexes, prj)
		}
	}

	vmio.LockIndex(indexes)

	if addDescriptor {
		project.Open = project.InfraId == ""
		project.Active = !project.Open
		NewIndexes = append(NewIndexes, project)
	}

	indexes.Projects = NewIndexes

	err = vmio.SaveIndex(indexes)

	vmio.UnlockIndex(indexes)

	return err
}

func UpdateIndexWithInfrastructure(infrastructure model.Infrastructure) error {
	indexes, err := vmio.LoadIndex()

	if err != nil {
		return err
	}

	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()

	indexes, err = vmio.LoadIndex()

	if err != nil {
		return err
	}

	vmio.LockIndex(indexes)

	iFaceProject := vmio.IFaceProject{
		Id: infrastructure.ProjectId,
	}
	iFaceProject.WaitForUnlock()

	project, err := vmio.LoadProject(infrastructure.ProjectId)

	if err != nil {
		return err
	}

	vmio.LockProject(project)

	project.Open = false

	err = vmio.SaveProject(project)

	vmio.UnlockProject(project)

	if err != nil {
		return err
	}

	Found := false
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _, prj := range indexes.Projects {
		if prj.Id == project.Id {
			prj.InfraId = infrastructure.Id
			prj.InfraName = infrastructure.Name
			prj.Active = true
			prj.Open = false
			Found = true
		}
		NewIndexes = append(NewIndexes, prj)
	}

	if !Found {
		return errors.New("Project Id: '" + infrastructure.ProjectId + "' for Infrastrcutre '" + infrastructure.Name + "' not found!!")
	}

	indexes.Projects = NewIndexes

	err = vmio.SaveIndex(indexes)

	vmio.UnlockIndex(indexes)

	return err
}

func CmdParseOption(key string, options []SubCommandHelper) (string, int, error) {
	if len(key) > 0 {
		if strings.Index(key, "--") == 0 {
			return key, -1, errors.New("Invalid Argument (wrong characters: --) : " + key)
		} else if strings.Index(key, "-") == 0 {
			return key, -1, errors.New("Invalid Argument (wrong character: -) : " + key)
		} else {
			for index, opts := range options {
				if CorrectInput(key) == CorrectInput(opts.Command) {
					return CorrectInput(key), index, nil
				}
			}
			return key, -1, errors.New("Invalid Argument : " + key)
		}
	} else {
		return key, -1, errors.New("Unable to parse Agument : " + key)
	}
}

func RecoverCommandHelper(helpCommand string) CommandHelper {
	helperCommands := GetArgumentHelpers()
	for _, helper := range helperCommands {
		if strings.ToLower(helper.Command) == strings.ToLower(helpCommand) {
			return helper
		}
	}
	return helperCommands[0]
}

func NewUUIDString() string {
	return uuid.NewV4().String()
}

func PrintCommandHelper(command string, subCommand string) {
	helper := RecoverCommandHelper(command)
	fmt.Fprintln(os.Stdout, "Help: vmkube", helper.LineHelp)
	fmt.Fprintln(os.Stdout, "Action:", helper.Description)
	found := false
	if "" != strings.TrimSpace(strings.ToLower(subCommand)) && "help" != strings.TrimSpace(strings.ToLower(subCommand)) {
		fmt.Fprintln(os.Stdout, "Selected Sub-Command: "+subCommand)
		for _, option := range helper.SubCommands {
			fmt.Fprintf(os.Stdout, "%s        %s\n", utils.StrPad(option.Command, 50), option.Description)
			found = true
		}
		if !found {
			fmt.Fprintln(os.Stdout, "Sub-Command Not found!!")
			if "help" != strings.TrimSpace(strings.ToLower(command)) {
				fmt.Fprintln(os.Stdout, "Please type: vmkube", "help", command, "for full Sub-Command List")
			} else {
				fmt.Fprintln(os.Stdout, "Please type: vmkube", "help", "COMMAND", "for full Sub-Command List")
			}
		}
	} else {
		found = true
		if len(helper.SubCommands) > 0 {
			if len(helper.SubCmdTypes) > 0 {
				fmt.Fprintln(os.Stdout, "Sub-Commands:")
			} else {
				fmt.Fprintln(os.Stdout, "Commands:")
			}
		}
		for _, option := range helper.SubCommands {
			fmt.Fprintf(os.Stdout, "%s        %s\n", utils.StrPad(option.Command, 55), option.Description)
		}
	}
	if found {
		if len(helper.Options) > 0 {
			fmt.Fprintln(os.Stdout, "Options:")
		}
		for _, option := range helper.Options {
			validity := "optional"
			if option.Mandatory {
				validity = "mandatory"
			}
			fmt.Fprintf(os.Stdout, "--%s  %s  %s  %s\n", utils.StrPad(option.Option, 20), utils.StrPad(option.Type, 25), utils.StrPad(validity, 10), option.Description)
		}
	} else {
		fmt.Fprintln(os.Stdout, "Unable to complete help support ...")
	}
}

func ExecuteInfrastructureActions(infrastructure model.Infrastructure, infrastructureActionCouples []tasks.ActivityCouple, NumThreads int, postTaskCallback func(task tasks.ScheduleTask)) []error {
	return executeActions(infrastructure, tasks.GroupActivitiesBySubject(infrastructureActionCouples), NumThreads, postTaskCallback)

}

func executeActions(infrastructure model.Infrastructure, actionGroups []tasks.ActivityGroup, NumThreads int, postTaskCallback func(task tasks.ScheduleTask)) []error {

	var errorsList []error = make([]error, 0)
	var maxJobNameLen int = 0
	var MachineAlterationAnswerChannel chan *tasks.MachineOperationsJob = make(chan *tasks.MachineOperationsJob)
	var jobsArrayLen int = 0
	var termElements []term.KeyValueElement = make([]term.KeyValueElement, 0)
	jobsArrayLen += len(actionGroups)

	pool := scheduler.SchedulerPool{
		Id:          NewUUIDString(),
		MaxParallel: NumThreads,
		KeepAlive:   true,
		PostExecute: true,
		Callback: func(task tasks.ScheduleTask) {
			//Any completed task come here ....
			postTaskCallback(task)
		},
	}
	pool.Init()
	var jobIds []string = make([]string, 0)
	for i := 0; i < jobsArrayLen; i++ {
		var jobId = NewUUIDString()
		jobIds = append(jobIds, jobId)
		var Prefix string = ""
		if actionGroups[i].IsCloud {
			Prefix = "Cloud "
		}
		var name string = fmt.Sprintf("[%d/%d] %s %sMachine Instance: '%s'", 0, len(actionGroups[i].Activities), tasks.ConvertActivityTaskInString(actionGroups[i].Task), Prefix,
			actionGroups[i].Subject)

		if len(name) > maxJobNameLen {
			maxJobNameLen = len(name) + len(strconv.Itoa(len(actionGroups[i].Activities))) - 1
		}
		termElem := term.KeyValueElement{
			Id:      NewUUIDString(),
			Name:    name,
			State:   term.StateColorWhite,
			Value:   "waiting...",
			Ref:     jobId,
			Actions: len(actionGroups[i].Activities),
		}
		termElements = append(termElements, termElem)
	}

	go pool.Start(func() {
		//println("Exit ...")
	})
	go func() {
		for i := 0; i < jobsArrayLen; i++ {
			var jobs []tasks.JobProcess = make([]tasks.JobProcess, 0)
			for j := 0; j < len(actionGroups[i].Activities); j++ {
				jobs = append(jobs, tasks.JobProcess(&tasks.Job{
					Id:   jobIds[i],
					Name: fmt.Sprintf("Process Instance from Project, Machine Group Name: %s Task : %d", actionGroups[i].Name, j),
					Runnable: tasks.RunnableStruct(&tasks.MachineOperationsJob{
						Name:             fmt.Sprintf("Process Instance from Project, Machine Group Name: %s Task : %d", actionGroups[i].Name, j),
						Infra:            actionGroups[i].Activities[j].Infra,
						Project:          actionGroups[i].Activities[j].Project,
						Activity:         actionGroups[i].Activities[j],
						InstanceId:       actionGroups[i].Activities[j].Instance.Id,
						OutChan:          MachineAlterationAnswerChannel,
						OwnState:         termElements[i],
						SendStartMessage: true,
						Index:            j,
						PartOf:           len(actionGroups[i].Activities),
						Command:          tasks.ConvertActivityTaskInString(actionGroups[i].Activities[0].Task),
						ActivityGroup:    actionGroups[i],
					}),
				}))
			}
			pool.Tasks <- tasks.ScheduleTask{
				Id:   NewUUIDString(),
				Jobs: jobs,
			}
		}
	}()
	var answerCounter int = 0
	go func() {
		var mutex sync.Mutex
		var resultsSeparator string = " status: "
		var screenManager term.KeyValueScreenManager
		if !utils.NO_COLORS {
			screenManager = term.KeyValueScreenManager{
				Elements:      termElements,
				MessageMaxLen: 45,
				Separator:     resultsSeparator,
				OffsetCols:    0,
				OffsetRows:    0,
				TextLen:       maxJobNameLen,
				BoldValue:     false,
			}
			screenManager.Init()
			screenManager.Start()
			defer screenManager.Stop(false)
		}
		var pending int = jobsArrayLen
		var answerScreenIds []string = make([]string, 0)
		var errorsInProgress bool = false
		var channelOpen bool = true
		for pending > 0 {
			select {
			case machineOpsJob, ok := (<-MachineAlterationAnswerChannel):
				if ok && machineOpsJob != nil {
					if !machineOpsJob.State {
						answerCounter++
					}
					go func(machineOpsJob *tasks.MachineOperationsJob) {
						machineMessage := machineOpsJob.MachineMessage
						activity := machineOpsJob.Activity

						if !machineOpsJob.State {
							for _, domain := range infrastructure.Domains {
								for _, network := range domain.Networks {
									if activity.IsCloud {
										for _, instance := range network.CloudInstances {
											if instance.Id == activity.CInstance.Id {
												instance.InspectJSON = machineMessage.InspectJSON
												instance.IPAddress = machineMessage.IPAddress
												break
											}
										}
									} else {
										for _, instance := range network.LocalInstances {
											if instance.Id == activity.Instance.Id {
												instance.InspectJSON = machineMessage.InspectJSON
												instance.IPAddress = machineMessage.IPAddress
												break
											}
										}
									}
								}
							}

						}

						if utils.NO_COLORS {
							if !machineOpsJob.State {
								message := "success!!"
								if machineMessage.Error != nil {
									errorsList = append(errorsList, machineMessage.Error)
									message = "failed!!"
								}
								if machineOpsJob.Index == machineOpsJob.PartOf-1 || machineMessage.Error != nil {
									operation := machineOpsJob.Command
									if activity.IsCloud {
										operation += " Cloud Machine Instance "
										fmt.Println(fmt.Sprintf("%s%s%s", utils.StrPad(operation+"'"+machineOpsJob.Machine+"'", maxJobNameLen), resultsSeparator, message))
									} else {
										operation += " Machine Instance "
										fmt.Println(fmt.Sprintf("%s%s%s", utils.StrPad(operation+"'"+machineOpsJob.Machine+"'", maxJobNameLen), resultsSeparator, message))
									}
									if machineMessage.Error != nil {
										errorsList = append(errorsList, machineMessage.Error)
										if !errorsInProgress {
											errorsInProgress = true
											mutex.Lock()
											fmt.Println(fmt.Sprintf(operation+"s interrupted, pending %d instance(s) will not be processed!!", (jobsArrayLen - answerCounter - pool.NumberOfWorkers() - 1)))
											if pool.IsRunning() {
												pool.Pause()
												pool.Interrupt()
											}
											for pool.IsWorking() {
												time.Sleep(1 * time.Second)
											}
											pending = pool.NumberOfWorkers() + 1
											mutex.Unlock()
										}
										pending--
										if pending <= 0 {
											pending = 0
											if channelOpen {
												channelOpen = false
												close(MachineAlterationAnswerChannel)
											}
										}
									} else {
										if machineOpsJob.Index == machineOpsJob.PartOf-1 {
											pending--
										}
									}
									if pending <= 0 {
										pending = 0
										if channelOpen {
											channelOpen = false
											close(MachineAlterationAnswerChannel)
										}
									}
								}
							}
						} else {
							//Interactive ...
							keyTerm := machineOpsJob.OwnState
							mutex.Lock()
							var Prefix string = ""
							if machineOpsJob.ActivityGroup.IsCloud {
								Prefix = "Cloud "
							}
							var keyName string = fmt.Sprintf("[%d/%d] %s %sMachine Instance: '%s'", (machineOpsJob.Index + 1), machineOpsJob.PartOf, tasks.ConvertActivityTaskInString(machineOpsJob.ActivityGroup.Task), Prefix,
								machineOpsJob.ActivityGroup.Subject)
							keyTerm.Name = keyName
							if machineOpsJob.State {
								if machineMessage.Error != nil {
									answerScreenIds = append(answerScreenIds, keyTerm.Id)
									keyTerm.State = term.StateColorRed
									keyTerm.Value = utils.StrPad(tasks.ConvertSubActivityTaskInString(machineOpsJob.Activity.Task)+"..."+term.Screen.Bold("failed!!"), 35)
								} else {
									keyTerm.State = term.StateColorYellow
									keyTerm.Value = term.StrPad(tasks.ConvertSubActivityTaskInString(machineOpsJob.Activity.Task)+"...in progress", 35)
								}
							} else {
								if machineOpsJob.Index == machineOpsJob.PartOf-1 {
									answerScreenIds = append(answerScreenIds, keyTerm.Id)
									if machineMessage.Error != nil {
										keyTerm.State = term.StateColorRed
										keyTerm.Value = term.Screen.Bold(term.StrPad("process failed!!", 35))
									} else {
										keyTerm.State = term.StateColorGreen
										keyTerm.Value = term.Screen.Bold(term.StrPad("process success!!", 35))
									}
								} else {
									if machineMessage.Error != nil {
										answerScreenIds = append(answerScreenIds, keyTerm.Id)
										keyTerm.State = term.StateColorRed
										keyTerm.Value = term.StrPad(tasks.ConvertSubActivityTaskInString(machineOpsJob.Activity.Task)+"..."+term.Screen.Bold("failed!!"), 35)
									} else {
										keyTerm.State = term.StateColorYellow
										keyTerm.Value = term.StrPad(tasks.ConvertSubActivityTaskInString(machineOpsJob.Activity.Task)+"..."+term.Screen.Bold("completed!!"), 35)
									}
								}
							}
							screenManager.CommChannel <- keyTerm
							mutex.Unlock()
							if machineMessage.Error != nil {
								errorsList = append(errorsList, machineMessage.Error)
								if !errorsInProgress {
									errorsInProgress = true
									mutex.Lock()
									if pool.IsRunning() {
										pool.Pause()
										pool.Interrupt()
									}
									//pending=pool.NumberOfWorkers() + 1
									for pending > 1 && pool.IsWorking() {
										time.Sleep(1 * time.Second)
									}
									pending = 1
									go func() {
										for _, signal := range screenManager.Elements {
											found := false
											for _, done := range answerScreenIds {
												if signal.Id == done {
													found = true
													break
												}
											}
											if !found && !pool.IsJobActive(fmt.Sprintf("%s", signal.Ref)) {
												signal.State = term.StateColorRed
												signal.Value = "interrupted!!"
												screenManager.CommChannel <- signal
											}
										}
									}()
									pending = pool.NumberOfWorkers() + 1
									mutex.Unlock()
								}
								pending--
								if pending <= 0 {
									pending = 0
									if channelOpen {
										channelOpen = false
										close(MachineAlterationAnswerChannel)
									}
								}
							} else {
								if machineOpsJob.Index == machineOpsJob.PartOf-1 && !machineOpsJob.State {
									pending--
								}
							}
							if pending <= 0 {
								pending = 0
								if channelOpen {
									channelOpen = false
									close(MachineAlterationAnswerChannel)
								}
							}
						}
					}(machineOpsJob)
				} else {
					pending = 0
					if utils.NO_COLORS {
						fmt.Println("Errors with legacy application ...")
					} else {
						for _, signal := range screenManager.Elements {
							found := false
							for _, done := range answerScreenIds {
								if signal.Id == done {
									found = true
									break
								}
							}
							if !found {
								signal.State = term.StateColorRed
								signal.Value = "interrupted!!"
								screenManager.CommChannel <- signal
							}
						}
					}
					break
				}
			case <-time.After(time.Second * 60):
			}
		}
		if pool.IsRunning() {
			pool.Pause()
			pool.Interrupt()
			for pool.IsWorking() || pending > 0 {
				time.Sleep(1 * time.Second)
			}
		}
		pool.Stop()
	}()
	pool.WG.Wait()
	time.Sleep(2 * time.Second)
	term.Screen.MoveCursor(len(actionGroups)+1, 0)
	utils.PrintlnImportant(fmt.Sprintf("Number of executed processes :  %d", answerCounter))
	if term.Screen.HasCursorHidden() {
		term.Screen.ShowCursor()
	}
	return errorsList
}

var FixInfrastructureElementMutex sync.Mutex

func FixInfrastructureElementValue(Infrastructure *model.Infrastructure, instanceId string, ipAddress string, json string, log string) bool {
	defer FixInfrastructureElementMutex.Unlock()
	FixInfrastructureElementMutex.Lock()
	if instanceId != "" && (ipAddress != "" || json != "" || log != "") {
		for i := 0; i < len(Infrastructure.Domains); i++ {
			for j := 0; j < len(Infrastructure.Domains[i].Networks); j++ {
				for k := 0; k < len(Infrastructure.Domains[i].Networks[j].LocalInstances); k++ {
					if Infrastructure.Domains[i].Networks[j].LocalInstances[k].Id == instanceId {
						if strings.TrimSpace(ipAddress) != "" {
							Infrastructure.Domains[i].Networks[j].LocalInstances[k].IPAddress = strings.TrimSpace(ipAddress)
						}
						if strings.TrimSpace(json) != "" {
							Infrastructure.Domains[i].Networks[j].LocalInstances[k].InspectJSON = strings.TrimSpace(json)
						}
						//if log != "" {
						//	Infrastructure.Domains[i].Networks[j].LocalInstances[k]. = json
						//}
						return true
					}
				}
				for k := 0; k < len(Infrastructure.Domains[i].Networks[j].CloudInstances); k++ {
					if Infrastructure.Domains[i].Networks[j].CloudInstances[k].Id == instanceId {
						if ipAddress != "" {
							Infrastructure.Domains[i].Networks[j].LocalInstances[k].IPAddress = strings.TrimSpace(ipAddress)
						}
						if json != "" {
							Infrastructure.Domains[i].Networks[j].LocalInstances[k].InspectJSON = strings.TrimSpace(json)
						}
						//if log != "" {
						//	Infrastructure.Domains[i].Networks[j].LocalInstances[k]. = json
						//}
						return true
					}
				}
			}
		}
	}
	return false
}

func DefineDestroyActivityFromCreateOne(activity tasks.ActivityCouple) tasks.ActivityCouple {
	return tasks.ActivityCouple{
		CInstance:  activity.CInstance,
		CMachine:   activity.CMachine,
		Infra:      activity.Infra,
		Instance:   activity.Instance,
		IsCloud:    activity.IsCloud,
		IsInstance: activity.IsInstance,
		NewInfra:   activity.NewInfra,
		Plans:      activity.Plans,
		Project:    activity.Project,
		Machine:    activity.Machine,
		Task:       tasks.DestroyMachine,
	}
}

func IsActivitySelected(activity tasks.ActivityCouple, id string) bool {
	if activity.IsCloud {
		if activity.CInstance.MachineId == id {
			return true
		}
	} else {
		if activity.CInstance.MachineId == id {
			return true
		}
	}
	return false
}

func DefineRebuildOfWholeInfrastructure(activities []tasks.ActivityCouple, excludedIds []string) []tasks.ActivityCouple {
	var outActivities []tasks.ActivityCouple = make([]tasks.ActivityCouple, 0)
	for _, activity := range activities {
		outActivities = append(outActivities, DefineDestroyActivityFromCreateOne(activity))
		Excluded := false
		for _, id := range excludedIds {
			if IsActivitySelected(activity, id) {
				Excluded = true
				break
			}
		}
		if !Excluded {
			outActivities = append(outActivities, activity)
		}
	}
	return outActivities
}

const MachineReadOperationTimeout = 900

func ExistInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cInstance model.CloudInstance, isCloud bool, instanceId string) (procedures.MachineState, error) {
	procedureExecutor := procedures.GetCurrentMachineExecutor(model.Project{},
		infrastructure,
		model.LocalMachine{},
		model.CloudMachine{},
		instance,
		cInstance,
		instanceId,
		isCloud,
		false)
	var commandPipe chan procedures.MachineMessage = make(chan procedures.MachineMessage, 1)
	var commandChannel chan *exec.Cmd = make(chan *exec.Cmd, 1)
	procedureExecutor.MachineStatus(commandPipe, commandChannel)
	select {
	case answer, ok := <-commandPipe:
		defer func() {
			close(commandPipe)
			close(commandChannel)
		}()
		if ok {
			return answer.State, answer.Error

		} else {
			return procedures.Machine_State_None, errors.New("Legacy application not executed correctly!!")
		}
	case <-time.After(time.Second * MachineReadOperationTimeout):
		return procedures.Machine_State_None, errors.New("Legacy application timed out!!")
	}
}

func FilterForExistState(infrastructure model.Infrastructure) ([]string, error) {
	return filterForExistStateAndCondState(infrastructure, procedures.Machine_State_None, false, false, false)
}

func FilterForExistAndRunningState(infrastructure model.Infrastructure, state procedures.MachineState) ([]string, error) {
	return filterForExistStateAndCondState(infrastructure, state, true, false, true)
}

func FilterForExistAndNotRunningState(infrastructure model.Infrastructure, state procedures.MachineState) ([]string, error) {
	return filterForExistStateAndCondState(infrastructure, state, true, true, true)
}

func filterForExistStateAndCondState(infrastructure model.Infrastructure, state procedures.MachineState, checkState bool, notState bool, excludeDisabled bool) ([]string, error) {
	var inactiveMachines []string = make([]string, 0)
	var err error
	for _, domain := range infrastructure.Domains {
		for _, network := range domain.Networks {
			for _, instance := range network.LocalInstances {
				if excludeDisabled && instance.Disabled {
					inactiveMachines = append(inactiveMachines, instance.Id)
					continue
				}
				procedureExecutor := procedures.GetCurrentMachineExecutor(model.Project{},
					infrastructure,
					model.LocalMachine{},
					model.CloudMachine{},
					instance,
					model.CloudInstance{},
					instance.Id,
					false,
					false)
				var commandPipe chan procedures.MachineMessage = make(chan procedures.MachineMessage, 1)
				var commandChannel chan *exec.Cmd = make(chan *exec.Cmd, 1)
				procedureExecutor.MachineStatus(commandPipe, commandChannel)
				select {
				case answer, ok := <-commandPipe:
					if ok {
						if answer.Error == nil {
							if answer.State == procedures.Machine_State_None {
								inactiveMachines = append(inactiveMachines, instance.Id)
							} else if checkState {
								if (!notState && answer.State != state) || (notState && answer.State == state) {
									inactiveMachines = append(inactiveMachines, instance.Id)
								}
							}
						} else {
							err = answer.Error
							inactiveMachines = append(inactiveMachines, instance.Id)
						}
					} else {
						inactiveMachines = append(inactiveMachines, instance.Id)
					}
					close(commandPipe)
					close(commandChannel)
				case <-time.After(time.Second * MachineReadOperationTimeout):
					inactiveMachines = append(inactiveMachines, instance.Id)
					close(commandPipe)
					close(commandChannel)
				}

			}
			for _, instance := range network.CloudInstances {
				if excludeDisabled && instance.Disabled {
					inactiveMachines = append(inactiveMachines, instance.Id)
					continue
				}
				procedureExecutor := procedures.GetCurrentMachineExecutor(model.Project{},
					infrastructure,
					model.LocalMachine{},
					model.CloudMachine{},
					model.LocalInstance{},
					instance,
					instance.Id,
					true,
					false)
				var commandPipe chan procedures.MachineMessage = make(chan procedures.MachineMessage, 1)
				var commandChannel chan *exec.Cmd = make(chan *exec.Cmd, 1)
				procedureExecutor.MachineStatus(commandPipe, commandChannel)
				select {
				case answer, ok := <-commandPipe:
					if ok {
						if answer.Error == nil {
							if answer.State == procedures.Machine_State_None {
								inactiveMachines = append(inactiveMachines, instance.Id)
							} else if checkState {
								if (!notState && answer.State != state) || (notState && answer.State == state) {
									inactiveMachines = append(inactiveMachines, instance.Id)
								}
							}
						} else {
							err = answer.Error
							inactiveMachines = append(inactiveMachines, instance.Id)
						}
					} else {
						inactiveMachines = append(inactiveMachines, instance.Id)
					}
					close(commandPipe)
					close(commandChannel)
				case <-time.After(time.Second * MachineReadOperationTimeout):
					inactiveMachines = append(inactiveMachines, instance.Id)
					close(commandPipe)
					close(commandChannel)
				}
			}
		}
	}
	return inactiveMachines, err
}

func FilterCreationBasedOnProjectActions(actions ProjectActionIndex, activities []tasks.ActivityCouple) ([]string, []tasks.ActivityCouple) {
	var outActivities []tasks.ActivityCouple = make([]tasks.ActivityCouple, 0)
	var removedIds []string = make([]string, 0)
	var allActions bool = false
	for _, action := range actions.Actions {
		if action.FullProject {
			allActions = true
			//return DefineRebuildOfWholeInfrastructure(activities)
		}
	}
	for _, activity := range activities {
		for _, action := range actions.Actions {
			DeleteAction := false
			ElementId := ""
			if !action.FullProject {
				if action.RelatedId == "" {
					if IsActivitySelected(activity, action.ElementId) {
						ElementId = action.ElementId
						if action.DropAction {
							DeleteAction = true
						}
					}
				} else {
					if IsActivitySelected(activity, action.RelatedId) {
						ElementId = action.RelatedId
						if action.DropAction {
							DeleteAction = true
						}
					}
				}
			}
			outActivities = append(outActivities, DefineDestroyActivityFromCreateOne(activity))
			if !DeleteAction {
				outActivities = append(outActivities, activity)
			} else {
				removedIds = append(removedIds, ElementId)
			}
		}
	}
	if allActions {
		return removedIds, DefineRebuildOfWholeInfrastructure(outActivities, removedIds)
	}
	return removedIds, outActivities
}

func MigrateProjectActionsToRollbackSegments(actions ProjectActionIndex) error {
	if actions.ProjectId != "" {
		return AddRollBackChangeActions(actions.ProjectId, actions.Actions...)
	}
	return errors.New("Invalid log descriptors for Project Id : " + actions.ProjectId)
}

func CopyStructure(origin interface{}, target interface{}) bool {
	valueOfOrigin := reflect.ValueOf(origin)
	if valueOfOrigin.Kind() == reflect.Struct {
		for i := 0; i < valueOfOrigin.NumField(); i++ {
			fieldValue := valueOfOrigin.Field(i)
			reflect.ValueOf(origin).Field(i).Set(fieldValue)
		}
		return true
	}
	return false
}

func ExtractStructureValue(origin interface{}, field string) interface{} {
	valueOfOrigin := reflect.ValueOf(origin)
	return valueOfOrigin.FieldByName(field).Interface()
}

func GetDefault(value interface{}, nilVal interface{}, dafaultVal interface{}) interface{} {
	if value == nilVal {
		return dafaultVal
	}
	return value
}

func FindInfrastructureInstance(Infrastructure model.Infrastructure, InstanceId string, InstanceName string) (model.LocalInstance, error) {
	var instance model.LocalInstance
	for i := 0; i < len(Infrastructure.Domains); i++ {
		for j := 0; j < len(Infrastructure.Domains[i].Networks); j++ {
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].LocalInstances); k++ {
				if Infrastructure.Domains[i].Networks[j].LocalInstances[k].Id == InstanceId ||
					CorrectInput(Infrastructure.Domains[i].Networks[j].LocalInstances[k].Name) == CorrectInput(InstanceName) {
					return Infrastructure.Domains[i].Networks[j].LocalInstances[k], nil
				}
			}
		}
	}
	return instance, errors.New(fmt.Sprintf("Instance not Found by UID : %s or by Name : %s", InstanceId, InstanceName))
}

func FindInfrastructureCloudInstance(Infrastructure model.Infrastructure, InstanceId string, InstanceName string) (model.CloudInstance, error) {
	var instance model.CloudInstance
	for i := 0; i < len(Infrastructure.Domains); i++ {
		for j := 0; j < len(Infrastructure.Domains[i].Networks); j++ {
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].CloudInstances); k++ {
				if Infrastructure.Domains[i].Networks[j].CloudInstances[k].Id == InstanceId ||
					CorrectInput(Infrastructure.Domains[i].Networks[j].CloudInstances[k].Name) == CorrectInput(InstanceName) {
					return Infrastructure.Domains[i].Networks[j].CloudInstances[k], nil
				}
			}
		}
	}
	return instance, errors.New(fmt.Sprintf("Cloud Instance not Found by UID : %s or by Name : %s", InstanceId, InstanceName))
}

func RemoveInfrastructureInstanceById(Infrastructure *model.Infrastructure, InstanceId string) error {
	for i := 0; i < len(Infrastructure.Domains); i++ {
		for j := 0; j < len(Infrastructure.Domains[i].Networks); j++ {
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].LocalInstances); k++ {
				if Infrastructure.Domains[i].Networks[j].LocalInstances[k].Id == InstanceId {
					tmp := Infrastructure.Domains[i].Networks[j].LocalInstances[(k + 1):]
					Infrastructure.Domains[i].Networks[j].LocalInstances = Infrastructure.Domains[i].Networks[j].LocalInstances[:k]
					Infrastructure.Domains[i].Networks[j].LocalInstances = append(Infrastructure.Domains[i].Networks[j].LocalInstances, tmp...)
					return nil
				}
			}
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].CloudInstances); k++ {
				if Infrastructure.Domains[i].Networks[j].CloudInstances[k].Id == InstanceId {
					tmp := Infrastructure.Domains[i].Networks[j].CloudInstances[(k + 1):]
					Infrastructure.Domains[i].Networks[j].CloudInstances = Infrastructure.Domains[i].Networks[j].CloudInstances[:k]
					Infrastructure.Domains[i].Networks[j].CloudInstances = append(Infrastructure.Domains[i].Networks[j].CloudInstances, tmp...)
					return nil
				}
			}
		}
	}
	return errors.New(fmt.Sprintf("Instance not Found by UID : %s", InstanceId))
}

func RemoveProjectMachineById(Infrastructure *model.Project, InstanceId string) error {
	for i := 0; i < len(Infrastructure.Domains); i++ {
		for j := 0; j < len(Infrastructure.Domains[i].Networks); j++ {
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].LocalMachines); k++ {
				if Infrastructure.Domains[i].Networks[j].LocalMachines[k].Id == InstanceId {
					tmp := Infrastructure.Domains[i].Networks[j].LocalMachines[(k + 1):]
					Infrastructure.Domains[i].Networks[j].LocalMachines = Infrastructure.Domains[i].Networks[j].LocalMachines[:k]
					Infrastructure.Domains[i].Networks[j].LocalMachines = append(Infrastructure.Domains[i].Networks[j].LocalMachines, tmp...)
					return nil
				}
			}
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].CloudMachines); k++ {
				if Infrastructure.Domains[i].Networks[j].CloudMachines[k].Id == InstanceId {
					tmp := Infrastructure.Domains[i].Networks[j].CloudMachines[(k + 1):]
					Infrastructure.Domains[i].Networks[j].CloudMachines = Infrastructure.Domains[i].Networks[j].CloudMachines[:k]
					Infrastructure.Domains[i].Networks[j].CloudMachines = append(Infrastructure.Domains[i].Networks[j].CloudMachines, tmp...)
					return nil
				}
			}
		}
	}
	return errors.New(fmt.Sprintf("Machine not Found by UID : %s", InstanceId))
}
