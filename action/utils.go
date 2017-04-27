package action

import (
	"strings"
	"errors"
	"vmkube/model"
	"vmkube/vmio"
	"fmt"
	"os"
	"vmkube/utils"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
	"vmkube/term"
	"vmkube/scheduler"
	"vmkube/operations"
	"reflect"
)


func ParseCommandArguments(args	[]string) (*CmdArguments, error) {
	arguments := CmdArguments{}
	success := arguments.Parse(args[1:])
	if success  {
		return  &arguments, nil
	} else  {
		return  &arguments, errors.New("Unable to Parse Command Line")
	}
}

func ParseCommandLine(args []string) (CmdRequest, error) {
	request := CmdRequest{}
	arguments, error := ParseCommandArguments(args)
	if error == nil  {
		request.TypeStr = arguments.Cmd
		request.Type = arguments.CmdType
		request.SubTypeStr = arguments.SubCmd
		request.SubType = arguments.SubCmdType
		request.HelpType = arguments.SubCmdHelpType
		request.Arguments = arguments
	}
	return  request, error
}


func CmdParseElement(value string) (CmdElementType, error) {
	switch CorrectInput(value) {
	case "server":
		return  LServer, nil
	case "cloud-server":
		return  CLServer, nil
	case "network":
		return  SNetwork, nil
	case "domain":
		return  SDomain, nil
	case "project":
		return  SProject, nil
	case "plan":
		return  SPlan, nil
	default:
		return  NoElement, errors.New("Element '"+value+"' is not an infratructure element. Available ones : Server, Cloud-Server, Network, Domain, Plan, Project")
		
	}
}

func CorrectInput(input string) string {
	return  strings.TrimSpace(strings.ToLower(input))
}

func GetBoolean(input string) bool {
	return CorrectInput(input) == "true"
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
	for _,domain := range project.Domains {
		newDomain := model.Domain{
			Id: NewUUIDString(),
			Name: domain.Name,
			Options: domain.Options,
			Networks: []model.Network{},
		}
		for _,network := range domain.Networks {
			newNetwork := model.Network{
				Id: NewUUIDString(),
				Name: network.Name,
				Options: network.Options,
				Instances: []model.Instance{},
				CInstances: []model.CloudInstance{},
				Installations: []model.Installation{},
			}
			serverConvertionMap := make(map[string]string)
			for _,server := range network.Servers {
				var disks []model.Disk = make([]model.Disk,0)
				disks = append(disks, model.Disk{
					Id: NewUUIDString(),
					Name: "sda0",
					Size: server.DiskSize,
					Type: 0,
				})
				instanceId := NewUUIDString()
				instance := model.Instance{
					Id: instanceId,
					Name: server.Name,
					Options: server.Options,
					Cpus: server.Cpus,
					Memory: server.Memory,
					Disks: disks,
					Driver: server.Driver,
					Engine: model.ToInstanceEngineOpt(server.Engine),
					Swarm: model.ToInstanceSwarmOpt(server.Swarm),
					Hostname: server.Hostname,
					IPAddress: "",
					NoShare: server.NoShare,
					OSType: server.OSType,
					OSVersion: server.OSVersion,
					Roles: server.Roles,
					ServerId: server.Id,
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId: infrastructure.Id,
						ElementId: instanceId,
						LogLines: []string{},
					},
				}
				if _,ok := serverConvertionMap[server.Id]; ok {
					return infrastructure, errors.New("Duplicate Server Id in Project : " + server.Id)
				}
				serverConvertionMap[server.Id] = instance.Id
				newNetwork.Instances = append(newNetwork.Instances, instance)
			}
			for _,server := range network.CServers {
				instanceId := NewUUIDString()
				instance := model.CloudInstance{
					Id: instanceId,
					Name: server.Name,
					Driver: server.Driver,
					Hostname: server.Hostname,
					IPAddress: "",
					Options: server.Options,
					Roles: server.Roles,
					ServerId: server.Id,
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId: infrastructure.Id,
						ElementId: instanceId,
						LogLines: []string{},
					},
				}
				if _,ok := serverConvertionMap[server.Id]; ok {
					return infrastructure, errors.New("Duplicate Server Id in Project : " + server.Id)
				}
				serverConvertionMap[server.Id] = instance.Id
				newNetwork.CInstances = append(newNetwork.CInstances, instance)
			}
			for _,plan := range network.Installations {
				if _,ok := serverConvertionMap[plan.ServerId]; ! ok {
					return infrastructure, errors.New("Invalid server reference in plan : " + plan.ServerId)
				}
				instanceId, _  := serverConvertionMap[plan.ServerId]
				installationId := NewUUIDString()
				installation := model.Installation{
					Id: installationId,
					Environment: model.ToInstanceEnvironment(plan.Environment),
					Role: model.ToInstanceRole(plan.Role),
					Type: model.ToInstanceInstallation(plan.Type),
					Errors: false,
					InstanceId: instanceId,
					IsCloud: plan.IsCloud,
					Success: false,
					LastExecution: time.Now(),
					LastMessage: "",
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId: infrastructure.Id,
						ElementId: installationId,
						LogLines: []string{},
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
	Found := false
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _,prj := range indexes.Projects {
		if CorrectInput(prj.Name) != CorrectInput(project.Name) {
			NewIndexes = append(NewIndexes, )
		} else {
			synced = (prj.InfraId == "")
			active = prj.Active
			InfraId = prj.InfraId
			InfraName = prj.InfraName
			Found = true
		}
	}
	
	vmio.LockIndex(indexes)
	
	indexes.Projects = append(indexes.Projects, model.ProjectsDescriptor{
		Id: project.Id,
		Name: project.Name,
		Open: project.Open,
		Synced: synced,
		Active: active,
		InfraId: InfraId,
		InfraName: InfraName,
	})
	
	if Found {
		indexes.Projects = NewIndexes
	}
	
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
	
	Found := false
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _,prj := range indexes.Projects {
		if prj.Id != project.Id {
			NewIndexes = append(NewIndexes, )
		} else {
			Found = true
		}
	}
	
	vmio.LockIndex(indexes)
	
	if addDescriptor {
		NewIndexes = append(NewIndexes, project)
		indexes.Projects = NewIndexes
	} else if Found {
		indexes.Projects = NewIndexes
	}
	
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
	for _,prj := range indexes.Projects {
		if prj.Id == project.Id {
			prj.InfraId = infrastructure.Id
			prj.InfraName = infrastructure.Name
			prj.Open = false
			Found = true
		}
		NewIndexes = append(NewIndexes, prj)
	}
	
	if ! Found {
		return errors.New("Project Id: '"+infrastructure.ProjectId+"' for Infrastrcutre '"+infrastructure.Name+"' not found!!")
	}
	
	indexes.Projects = NewIndexes
	
	err = vmio.SaveIndex(indexes)
	
	vmio.UnlockIndex(indexes)
	
	return err
}

func CmdParseOption(key string, options []SubCommandHelper) (string, int, error) {
	if len(key) > 0 {
		if strings.Index(key, "--") == 0  {
			return key, -1, errors.New("Invalid Argument (wrong characters: --) : " + key)
		} else	if strings.Index(key, "-") == 0  {
			return key, -1, errors.New("Invalid Argument (wrong character: -) : " + key)
		} else  {
			for index,opts := range options {
				if CorrectInput(key) == opts.Command  {
					return  CorrectInput(key), index, nil
				}
			}
			return  key, -1, errors.New("Invalid Argument : " + key)
		}
	} else  {
		return  key, -1, errors.New("Unable to parse Agument : " + key)
	}
}

func RecoverCommandHelper(helpCommand string) CommandHelper {
	helperCommands := GetArgumentHelpers()
	for _, helper := range helperCommands {
		if helper.Command == strings.ToLower(helpCommand) {
			return  helper
		}
	}
	return  helperCommands[0]
}

func NewUUIDString()	string {
	return  uuid.NewV4().String()
}

func PrintCommandHelper(command	string, subCommand string) {
	helper := RecoverCommandHelper(command)
	fmt.Fprintln(os.Stdout, "Help: vmkube", helper.LineHelp)
	fmt.Fprintln(os.Stdout, "Action:", helper.Description)
	found := false
	if "" !=  strings.TrimSpace(strings.ToLower(subCommand)) && "help" !=  strings.TrimSpace(strings.ToLower(subCommand)) {
		fmt.Fprintln(os.Stdout, "Selected Sub-Command: " + subCommand)
		for _,option := range helper.SubCommands {
			fmt.Fprintf(os.Stdout, "%s        %s\n",  utils.StrPad(option.Command, 50), option.Description)
			found = true
		}
		if ! found  {
			fmt.Fprintln(os.Stdout, "Sub-Command Not found!!")
			if "help" !=  strings.TrimSpace(strings.ToLower(command)) {
				fmt.Fprintln(os.Stdout, "Please type: vmkube","help", command,"for full Sub-Command List")
			} else  {
				fmt.Fprintln(os.Stdout, "Please type: vmkube","help", "COMMAND","for full Sub-Command List")
			}
		}
	}  else {
		found = true
		if len(helper.SubCommands) > 0  {
			if len(helper.SubCmdTypes) > 0 {
				fmt.Fprintln(os.Stdout, "Sub-Commands:")
			} else {
				fmt.Fprintln(os.Stdout, "Commands:")
			}
		}
		for _,option := range helper.SubCommands {
			fmt.Fprintf(os.Stdout, "%s        %s\n",  utils.StrPad(option.Command, 55), option.Description)
		}
	}
	if found  {
		if len(helper.Options) > 0  {
			fmt.Fprintln(os.Stdout, "Options:")
		}
		for _,option := range helper.Options {
			validity := "optional"
			if option.Mandatory {
				validity = "mandatory"
			}
			fmt.Fprintf(os.Stdout, "--%s  %s  %s  %s\n",  utils.StrPad(option.Option,20),  utils.StrPad(option.Type, 25), utils.StrPad(validity, 10), option.Description)
		}
	} else  {
		fmt.Fprintln(os.Stdout, "Unable to complete help support ...")
	}
}

func ConvertActivityTaskInString(task operations.ActivityTask) string {
	operation := "Create"
	switch task {
	case operations.DestroyMachine:
		operation = "Destroy"
		break
	case operations.StartMachine:
		operation = "Start"
		break
	case operations.StopMachine:
		operation = "Stop"
		break
	case operations.RestartMachine:
		operation = "Restart"
		break
	case operations.MachineStatus:
		operation = "Get Status of"
		break
	case operations.MachineEnv:
		operation = "Get Environment for"
		break
	case operations.MachineInspect:
		operation = "Get Descriptor of"
		break
	case operations.MachineIPAddress:
		operation = "Get IP Address of"
		break
	}
	return operation
}

func ExecuteInfrastructureActions(infrastructure model.Infrastructure,infrastructureActionCouples []operations.ActivityCouple, NumThreads int, postTaskCallback func(task scheduler.ScheduleTask)) []error {
	var errorsList []error = make([]error, 0)
	var maxJobNameLen int = 0
	var ServerCreationAnswerChannel chan *operations.ServerOperationsJob = make(chan *operations.ServerOperationsJob)
	var jobsArrayLen int = 0
	var termElements []term.KeyValueElement = make([]term.KeyValueElement, 0)
	jobsArrayLen += len(infrastructureActionCouples)
	
	pool := scheduler.SchedulerPool{
		Id: NewUUIDString(),
		MaxParallel: NumThreads,
		KeepAlive: true,
		PostExecute: true,
		Callback: func(task scheduler.ScheduleTask) {
			//Any completed task come here ....
			postTaskCallback(task)
		},
	}
	pool.Init()
	var jobIds []string = make([]string, 0)
	for i := 0; i < jobsArrayLen; i++ {
		var jobId = NewUUIDString()
		jobIds = append(jobIds, jobId)
		var name string
		if infrastructureActionCouples[i].IsCloud {
			name = fmt.Sprintf("%s Cloud Instance Server %s", ConvertActivityTaskInString(infrastructureActionCouples[i].Task), infrastructureActionCouples[i].CInstance.Name)
		} else {
			name = fmt.Sprintf("%s Instance Server %s", ConvertActivityTaskInString(infrastructureActionCouples[i].Task), infrastructureActionCouples[i].Instance.Name)
		}
		
		if len(name) > maxJobNameLen {
			maxJobNameLen = len(name)
		}
		termElem := term.KeyValueElement{
			Id: NewUUIDString(),
			Name: name,
			State: term.StateColorWhite,
			Value: "waiting...",
			Ref: jobId,
		}
		termElements = append(termElements, termElem)
	}

	go pool.Start(func() {
		//println("Exit ...")
	})
	go func(){
		for i := 0; i < jobsArrayLen; i++ {
			if ! infrastructureActionCouples[i].IsCloud {
				pool.Tasks <- scheduler.ScheduleTask{
					Id: NewUUIDString(),
					Jobs: []scheduler.Job{
						{
							Id: jobIds[i],
							Name: fmt.Sprintf("Process Instance from Project Server Id: %s", infrastructureActionCouples[i].Instance.ServerId),
							Runnable: operations.RunnableStruct(&operations.ServerOperationsJob{
								Name: fmt.Sprintf("Process Instance from Project Server Id: %s", infrastructureActionCouples[i].Instance.ServerId),
								Infra: infrastructureActionCouples[i].Infra,
								Project:infrastructureActionCouples[i].Project,
								Activity: infrastructureActionCouples[i],
								InstanceId: infrastructureActionCouples[i].Instance.Id,
								OutChan: ServerCreationAnswerChannel,
								OwnState: termElements[i],
								SendStartMessage: true,
							}),
						},
					},
				}
			} else {
				pool.Tasks <- scheduler.ScheduleTask{
					Id: NewUUIDString(),
					Jobs: []scheduler.Job{
						{
							Id: jobIds[i],
							Name: fmt.Sprintf("Process Instance from Project Server Id: %s", infrastructureActionCouples[i].CInstance.ServerId),
							Runnable: operations.RunnableStruct(&operations.ServerOperationsJob{
								Name: fmt.Sprintf("Process Instance from Project Server Id: %s", infrastructureActionCouples[i].CInstance.ServerId),
								Infra: infrastructureActionCouples[i].Infra,
								Project:infrastructureActionCouples[i].Project,
								Activity: infrastructureActionCouples[i],
								InstanceId: infrastructureActionCouples[i].Instance.Id,
								OutChan: ServerCreationAnswerChannel,
								OwnState: termElements[i],
								SendStartMessage: true,
							}),
						},
					},
				}
			}
		}
	}()
	var answerCounter int = 0
	go func(){
		var resultsSeparator string = " status: "
		var screenManager term.KeyValueScreenManager
		if ! utils.NO_COLORS {
			screenManager = term.KeyValueScreenManager{
				Elements: termElements,
				MessageMaxLen: 25,
				Separator: resultsSeparator,
				OffsetCols: 0,
				OffsetRows: 0,
				TextLen: maxJobNameLen,
				BoldValue: false,
			}
			screenManager.Init()
			screenManager.Start()
		}
		var pending int = jobsArrayLen
		var answerScreenIds []string = make([]string, 0)
		var errorsInProgress bool = false
		for pending > 0 {
			serverOpsJob, ok := (<- ServerCreationAnswerChannel)
			if ok && serverOpsJob != nil {
				if ! serverOpsJob.State {
					answerCounter++
				}
				go func(serverOpsJob *operations.ServerOperationsJob) {
					machineMessage := serverOpsJob.MachineMessage
					activity := serverOpsJob.Activity
					
					if ! serverOpsJob.State {
						for _,domain := range infrastructure.Domains {
							for _,network := range domain.Networks {
								if activity.IsCloud {
									for _,instance := range network.CInstances {
										if instance.Id == activity.CInstance.Id {
											instance.InspectJSON = machineMessage.InspectJSON
											instance.IPAddress = machineMessage.IPAddress
											break
										}
									}
								} else {
									for _,instance := range network.Instances {
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
						if ! serverOpsJob.State {
							message := "success!!"
							if machineMessage.Error != nil {
								errorsList = append(errorsList, machineMessage.Error)
								message = "failed!!"
							}
							operation := ConvertActivityTaskInString(activity.Task)
							if activity.IsCloud {
								operation += " Cloud Instance Server "
								fmt.Println(fmt.Sprintf("%s%s%s", utils.StrPad(operation+activity.CInstance.Name,maxJobNameLen), resultsSeparator, message))
							} else {
								operation += " Instance Server "
								fmt.Println(fmt.Sprintf("%s%s%s", utils.StrPad(operation+activity.Instance.Name,maxJobNameLen), resultsSeparator, message))
							}
							if machineMessage.Error != nil {
								errorsList = append(errorsList, machineMessage.Error)
								if ! errorsInProgress {
									errorsInProgress = true
									fmt.Println(fmt.Sprintf(operation + "s interrupted, pending %d instance(s) will not be processed!!", (jobsArrayLen - answerCounter - pool.NumberOfWorkers() - 1)))
									if pool.IsRunning() {
										pool.Pause()
									}
									for pool.IsWorking() {
										time.Sleep(1*time.Second)
									}
								} else {
									pending --
								}
							} else {
								pending--
								if pending == 0 {
									close(ServerCreationAnswerChannel)
								}
							}
						}
					} else {
						//Interactive ...
						keyTerm := serverOpsJob.OwnState
						if serverOpsJob.State {
							keyTerm.State = term.StateColorYellow
							keyTerm.Value = "processing..."
						} else {
							answerScreenIds = append(answerScreenIds, keyTerm.Id)
							if machineMessage.Error != nil {
								keyTerm.State = term.StateColorRed
								keyTerm.Value = term.ScreenBold("failed!!")
							} else {
								keyTerm.State = term.StateColorGreen
								keyTerm.Value = term.ScreenBold("success!!")
							}
						}
						screenManager.CommChannel <- keyTerm
						if machineMessage.Error != nil {
							errorsList = append(errorsList, machineMessage.Error)
							if ! errorsInProgress {
								errorsInProgress = true
								if pool.IsRunning() {
									pool.Pause()
								}
								for pool.IsWorking() {
									time.Sleep(1*time.Second)
								}
								for _,signal := range screenManager.Elements {
									found := false
									for _,done := range answerScreenIds {
										if signal.Id == done {
											found = true
											break
										}
									}
									if ! found  && !pool.IsJobActive(fmt.Sprintf("%s", signal.Ref)) {
										signal.State = term.StateColorRed
										signal.Value = "interrupted!!"
										screenManager.CommChannel <- signal
									}
								}
							} else if ! serverOpsJob.State {
								pending --
							}
						} else {
							if ! serverOpsJob.State {
								pending--
							}
							if pending == 0 {
								close(ServerCreationAnswerChannel)
							}
						}
					}
				}(serverOpsJob)
			} else {
				pending = 0
				if utils.NO_COLORS {
					fmt.Println("Errors with legacy application ...")
				} else {
					for _,signal := range screenManager.Elements {
						found := false
						for _,done := range answerScreenIds {
							if signal.Id == done {
								found = true
								break
							}
						}
						if ! found {
							signal.State = term.StateColorRed
							signal.Value = "interrupted!!"
							screenManager.CommChannel <- signal
						}
					}
				}
				break
			}
		}
		if pool.IsRunning() {
			pool.Pause()
			for pool.IsWorking() || pending > 0 {
				time.Sleep(1*time.Second)
			}
			pool.Interrupt()
			pool.Stop()
		}
	}()
	pool.WG.Wait()
	time.Sleep(2*time.Second)
	utils.PrintlnImportant(fmt.Sprintf("Task executed:  %d", answerCounter))
	return errorsList
}

func FixInfrastructureElementValue(Infrastructure model.Infrastructure, instanceId string, ipAddress string, json string) bool {
	for i := 0; i < len(Infrastructure.Domains); i++ {
		for j := 0; j < len(Infrastructure.Domains[i].Networks); j++ {
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].Instances); k++ {
				if Infrastructure.Domains[i].Networks[j].Instances[k].Id == instanceId {
					Infrastructure.Domains[i].Networks[j].Instances[k].IPAddress = ipAddress
					Infrastructure.Domains[i].Networks[j].Instances[k].InspectJSON = json
					return true
				}
			}
			for k := 0; k < len(Infrastructure.Domains[i].Networks[j].CInstances); k++ {
				if Infrastructure.Domains[i].Networks[j].CInstances[k].Id == instanceId {
					Infrastructure.Domains[i].Networks[j].CInstances[k].IPAddress = ipAddress
					Infrastructure.Domains[i].Networks[j].CInstances[k].InspectJSON = json
					return true
				}
			}
		}
	}
	return false
}

func DefineDestroyActivityFromCreateOne(activity operations.ActivityCouple) operations.ActivityCouple {
	return operations.ActivityCouple{
		CInstance: activity.CInstance,
		CServer: activity.CServer,
		Infra: activity.Infra,
		Instance: activity.Instance,
		IsCloud: activity.IsCloud,
		IsInstance: activity.IsInstance,
		NewInfra: activity.NewInfra,
		Plans: activity.Plans,
		Project: activity.Project,
		Server: activity.Server,
		Task: operations.DestroyMachine,
	}
}

func DefineRebuildOfWholeInfrastructure(activities []operations.ActivityCouple) []operations.ActivityCouple {
	var outActivities []operations.ActivityCouple = make([]operations.ActivityCouple, 0)
	for _,activity := range activities {
		outActivities = append(outActivities, DefineDestroyActivityFromCreateOne(activity))
		outActivities = append(outActivities, activity)
	}
	return outActivities
}

func FindActivityById(activities []operations.ActivityCouple, id string) (operations.ActivityCouple, error) {
	for _, activity := range activities {
		if activity.IsCloud {
			if activity.CInstance.ServerId == id {
				return activity, nil
			}
		} else {
			if activity.CInstance.ServerId == id {
				return activity, nil
			}
		}
	}
	return operations.ActivityCouple{}, errors.New("Activity Not Found")
}

func FilterCreationBasedOnProjectActions(actions ProjectActionIndex, activities []operations.ActivityCouple) []operations.ActivityCouple {
	var outActivities []operations.ActivityCouple = make([]operations.ActivityCouple, 0)
	for _,action := range actions.Actions {
		if action.FullProject {
			return DefineRebuildOfWholeInfrastructure(activities)
		} else {
			var activity operations.ActivityCouple
			var err error
			if action.RelatedId == "" {
				activity, err = FindActivityById(activities, action.ElementId)
			} else {
				activity, err = FindActivityById(activities, action.RelatedId)
			}
			if err != nil {
				outActivities = append(outActivities, DefineDestroyActivityFromCreateOne(activity))
				if ! action.DropAction {
					outActivities = append(outActivities, activity)
				}
			}
		}
	}
	return outActivities
}

func MigrateProjectActionsToRollbackSegments(actions ProjectActionIndex) error {
	if actions.ProjectId != "" {
		return AddRollBackChangeActions(actions.ProjectId,actions.Actions...)
	}
	return errors.New("Invalid log descriptors for Project Id : " + actions.ProjectId)
}

func CopyStructure(origin interface{}, target interface{}) bool {
	valueOfOrigin := reflect.ValueOf(origin)
	if valueOfOrigin.Kind() == reflect.Struct {
		for i := 0; i < valueOfOrigin.NumField(); i++ {
			fieldValue :=valueOfOrigin.Field(i)
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