package operations

import (
	"vmkube/term"
	"vmkube/model"
	"vmkube/procedures"
	"errors"
	"os/exec"
	"fmt"
)


type RunnableStruct interface {
	Start()
	Stop()
	Status() bool
	Response() interface{}
}


type MachineOperationsJob struct {
	Name             string
	State            bool
	OutChan          chan *MachineOperationsJob
	OwnState         term.KeyValueElement
	Project          model.Project
	Infra            model.Infrastructure
	InstanceId       string
	Activity         ActivityCouple
	MachineMessage   procedures.MachineMessage
	SendStartMessage bool
	commandPipe      chan procedures.MachineMessage
	commandChannel   chan *exec.Cmd
	threadSafeCmd    bool
	control          procedures.MachineControlStructure
}

func (job *MachineOperationsJob) Response() interface{} {
	return fmt.Sprintf("%s|%s|%s",job.MachineMessage.InstanceId, job.MachineMessage.InspectJSON, job.MachineMessage.IPAddress)
}

func (job *MachineOperationsJob) Start() {
	if !job.State {
		job.State = true
		if job.SendStartMessage && ! job.control.Interrupt {
			job.OutChan <- job
		}
		job.commandPipe = make(chan procedures.MachineMessage)
		machineAdapter := procedures.GetCurrentMachineExecutor(job.Project, job.Infra, job.Activity.Machine, job.Activity.CMachine, job.Activity.Instance, job.Activity.CInstance, job.Activity.Instance.Id, job.Activity.IsCloud, job.Activity.NewInfra)
		job.threadSafeCmd = machineAdapter.IsThreadSafeCommand()
		machineAdapter.SetControlStructure(&job.control)
		job.commandChannel = make(chan *exec.Cmd)
		switch job.Activity.Task {
		case CreateMachine:
			if job.Activity.IsCloud {
				go machineAdapter.CreateCloudMachine(job.commandPipe, job.commandChannel)
			} else {
				go machineAdapter.CreateMachine(job.commandPipe, job.commandChannel)
			}
			break
			case DestroyMachine:
				go machineAdapter.RemoveMachine(job.commandPipe, job.commandChannel)
				break
			case StopMachine:
				go machineAdapter.StopMachine(job.commandPipe, job.commandChannel)
				break
			case StartMachine:
				go machineAdapter.StartMachine(job.commandPipe, job.commandChannel)
				break
			case RestartMachine:
				go machineAdapter.RestartMachine(job.commandPipe, job.commandChannel)
				break
			case MachineStatus:
				go machineAdapter.MachineStatus(job.commandPipe, job.commandChannel)
				break
			case MachineEnv:
				go machineAdapter.MachineEnv(job.commandPipe, job.commandChannel)
				break
			case MachineInspect:
				go machineAdapter.MachineInspect(job.commandPipe, job.commandChannel)
				break
			case MachineIPAddress:
				break
				go machineAdapter.MachineIPAddress(job.commandPipe, job.commandChannel)
			default:
				panic("No matching ActivityTask for Job")
		}
		go func(){
			for job.State {
				job.control.CurrentCommand = <- job.commandChannel
			}
		}()
		if ! job.control.Interrupt {
			message := <- job.commandPipe
			job.MachineMessage.IPAddress = message.IPAddress
			job.MachineMessage.InspectJSON = message.InspectJSON
			job.MachineMessage.Complete = message.Complete
			job.MachineMessage.Error = message.Error
			job.MachineMessage.Result = message.Result
			job.MachineMessage.State = message.State
			if job.Activity.IsCloud {
				job.MachineMessage.InstanceId = job.Activity.CInstance.Id
			} else {
				job.MachineMessage.InstanceId = job.Activity.Instance.Id
			}
			job.State = false
			job.OutChan <- job
		}
		job.State = false
	}
}

func (job *MachineOperationsJob) Stop() {
	if job.threadSafeCmd {
		job.control.Interrupt = true
		if job.control.CurrentCommand.Process.Pid > 0 {
			job.control.CurrentCommand.Process.Kill()
		}
		close(job.commandChannel)
		close(job.commandPipe)
	} else {
		job.control.Interrupt = true
		if job.control.CurrentCommand.Process.Pid > 0 {
			job.control.CurrentCommand.Wait()
		}
		close(job.commandChannel)
		close(job.commandPipe)
	}
	job.State = false
}

func (job *MachineOperationsJob) Status() bool {
	return job.State
}

type ActivityTask int

const(
	CreateMachine  ActivityTask = iota
	DestroyMachine
	StopMachine
	StartMachine
	RestartMachine
	MachineStatus
	MachineEnv
	MachineInspect
	MachineIPAddress
	
)

type ActivityCouple struct {
	Project     model.Project
	Infra       model.Infrastructure
	IsCloud     bool
	IsInstance  bool
	Machine      model.LocalMachine
	CMachine     model.CloudMachine
	Instance    model.LocalInstance
	CInstance   model.CloudInstance
	Plans       []model.InstallationPlan
	Task        ActivityTask
	NewInfra    bool
}

func filterPlansByMachine(id string, isCloud bool, network model.MachineNetwork) []model.InstallationPlan {
	var selectedPlans []model.InstallationPlan = make([]model.InstallationPlan, 0)
	for _,plan := range network.Installations {
		if plan.IsCloud == isCloud && plan.MachineId == id {
			selectedPlans = append(selectedPlans, plan)
		}
	}
	return selectedPlans
}

func filterInstanceByMachine(id string, infrastructure model.Infrastructure) (model.LocalInstance, error) {
	var Instance model.LocalInstance
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,machine := range network.LocalInstances {
				if machine.MachineId == id {
					return machine, nil
				}
			}
		}
	}
	return Instance, errors.New("Instance for Machine Id: "+id+" not found")
}

func filterCloudInstanceByMachine(id string, infrastructure model.Infrastructure) (model.CloudInstance, error) {
	var Instance model.CloudInstance
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,machine := range network.CloudInstances {
				if machine.MachineId == id {
					return machine, nil
				}
			}
		}
	}
	return Instance, errors.New("Instance for Machine Id: "+id+" not found")
}


func GetTaskActivities(project model.Project, infrastructure model.Infrastructure, task ActivityTask) ([]ActivityCouple, error) {
	var taskList []ActivityCouple = make([]ActivityCouple, 0)
	for _,domain := range project.Domains {
		for _,network := range domain.Networks {
			for _,machine := range network.LocalMachines {
				instance, err := filterInstanceByMachine(machine.Id, infrastructure)
				if err != nil {
					return taskList, err
				}
				taskList = append(taskList, ActivityCouple{
					IsCloud: false,
					Machine: machine,
					Instance: instance,
					Task: task,
					Plans: filterPlansByMachine(machine.Id, false, network),
					NewInfra: true,
					Project: project,
					Infra: infrastructure,
				})
			}
			for _,machine := range network.CloudMachines {
				instance, err := filterCloudInstanceByMachine(machine.Id, infrastructure)
				if err != nil {
					return taskList, err
				}
				taskList = append(taskList, ActivityCouple{
					IsCloud: true,
					CMachine: machine,
					CInstance: instance,
					Task: task,
					Plans: filterPlansByMachine(machine.Id, true, network),
					NewInfra: true,
					Project: project,
					Infra: infrastructure,
				})
			}
		}
	}
	return taskList, nil
}
func FilterByInstanceState(activities []ActivityCouple, isNew bool) []ActivityCouple {
	var outActivities []ActivityCouple = make([]ActivityCouple, 0)
	for _,activity := range activities {
		if activity.Task == StartMachine || activity.Task == StopMachine || activity.Task == RestartMachine {
			instanceId := ""
			if activity.IsCloud {
				instanceId = activity.CInstance.Id
			} else {
				instanceId = activity.Instance.Id
			}
			executor := procedures.GetCurrentMachineExecutor(
				activity.Project,
				activity.Infra,
				activity.Machine,
				activity.CMachine,
				activity.Instance,
				activity.CInstance,
				instanceId,
				activity.IsCloud,
				isNew,
			)
			var commandPipe chan procedures.MachineMessage = make(chan procedures.MachineMessage)
			var commandChannel chan *exec.Cmd = make(chan *exec.Cmd)
			executor.MachineStatus(commandPipe, commandChannel)
			message := <- commandPipe
			if message.State != procedures.Machine_State_None && (activity.Task == StartMachine && activity.Task == RestartMachine && message.State == procedures.Machine_State_Stopped) &&
				(activity.Task == StopMachine &&  message.State == procedures.Machine_State_Running){
				//TODO: Implement Filter by DefaultMachineExecutor.State
				outActivities = append(outActivities, activity)
			}
		} else {
			outActivities = append(outActivities, activity)
		}
	}
	return outActivities
}

func GetPostBuildTaskActivities(infrastructure model.Infrastructure, task ActivityTask) ([]ActivityCouple, error) {
	var taskList []ActivityCouple = make([]ActivityCouple, 0)
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,machine := range network.LocalInstances {
				taskList = append(taskList, ActivityCouple{
					IsCloud: false,
					Instance: machine,
					Task: task,
					Plans: []model.InstallationPlan{},
					Infra: infrastructure,
					NewInfra: false,
				})
			}
			for _,machine := range network.CloudInstances {
				taskList = append(taskList, ActivityCouple{
					IsCloud: true,
					CInstance: machine,
					Task: task,
					Plans: []model.InstallationPlan{},
					Infra: infrastructure,
					NewInfra: false,
				})
			}
		}
	}
	return taskList, nil
}
