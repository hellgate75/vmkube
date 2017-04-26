package operations

import (
	"vmkube/term"
	"vmkube/model"
	"vmkube/procedures"
	"errors"
)


type RunnableStruct interface {
	Start()
	Stop()
	Status() bool
}


type ServerOperationsJob struct {
	Name             string
	State            bool
	OutChan          chan *ServerOperationsJob
	OwnState         term.KeyValueElement
	Project          model.Project
	Infra            model.Infrastructure
	InstanceId       string
	Activity         ActivityCouple
	CommandPipe      chan procedures.MachineMessage
	MachineMessage   procedures.MachineMessage
	SendStartMessage bool
}

func (job *ServerOperationsJob) Start() {
	if !job.State {
		job.State = true
		if job.SendStartMessage {
			job.OutChan <- job
		}
		job.CommandPipe = make(chan procedures.MachineMessage)
		var message procedures.MachineMessage
		machineAdapter := procedures.GetCurrentServerMachine(job.Project, job.Infra, job.Activity.Server, job.Activity.CServer, job.Activity.Instance, job.Activity.CInstance, job.Activity.Instance.Id, job.Activity.IsCloud, job.Activity.NewInfra)
		switch job.Activity.Task {
		case CreateMachine:
			if job.Activity.IsCloud {
				go machineAdapter.CreateCloudServer(job.CommandPipe)
			} else {
				go machineAdapter.CreateServer(job.CommandPipe)
			}
			break
			case DestroyMachine:
				go machineAdapter.RemoveServer(job.CommandPipe)
				break
			case StopMachine:
				go machineAdapter.StopServer(job.CommandPipe)
				break
			case StartMachine:
				go machineAdapter.StartServer(job.CommandPipe)
				break
			case RestartMachine:
				go machineAdapter.RestartServer(job.CommandPipe)
				break
			case MachineStatus:
				go machineAdapter.ServerStatus(job.CommandPipe)
				break
			case MachineEnv:
				go machineAdapter.ServerEnv(job.CommandPipe)
				break
			case MachineInspect:
				go machineAdapter.ServerInspect(job.CommandPipe)
				break
			case MachineIPAddress:
				break
				go machineAdapter.ServerIPAddress(job.CommandPipe)
			default:
				panic("No matching ActivityTask for Job")
		}
		message = <- job.CommandPipe
		close(job.CommandPipe)
		job.MachineMessage = message
		job.State = false
		job.OutChan <- job
	}
}

func (job *ServerOperationsJob) Stop() {
	close(job.CommandPipe)
	job.State = false
}

func (job *ServerOperationsJob) Status() bool {
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
	Server      model.ProjectServer
	CServer     model.ProjectCloudServer
	Instance    model.Instance
	CInstance   model.CloudInstance
	Plans       []model.InstallationPlan
	Task        ActivityTask
	NewInfra    bool
}

func filterPlansByServer(id string, isCloud bool, network model.ProjectNetwork) []model.InstallationPlan {
	var selectedPlans []model.InstallationPlan = make([]model.InstallationPlan, 0)
	for _,plan := range network.Installations {
		if plan.IsCloud == isCloud && plan.ServerId == id {
			selectedPlans = append(selectedPlans, plan)
		}
	}
	return selectedPlans
}

func filterInstanceByServer(id string, infrastructure model.Infrastructure) (model.Instance, error) {
	var Instance model.Instance
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,server := range network.Instances {
				if server.ServerId == id {
					return server, nil
				}
			}
		}
	}
	return Instance, errors.New("Instance for Server Id: "+id+" not found")
}

func filterCloudInstanceByServer(id string, infrastructure model.Infrastructure) (model.CloudInstance, error) {
	var Instance model.CloudInstance
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,server := range network.CInstances {
				if server.ServerId == id {
					return server, nil
				}
			}
		}
	}
	return Instance, errors.New("Instance for Server Id: "+id+" not found")
}

func GetTaskActivities(project model.Project, infrastructure model.Infrastructure, task ActivityTask) ([]ActivityCouple, error) {
	var taskList []ActivityCouple = make([]ActivityCouple, 0)
	for _,domain := range project.Domains {
		for _,network := range domain.Networks {
			for _,server := range network.Servers {
				instance, err := filterInstanceByServer(server.Id, infrastructure)
				if err != nil {
					return taskList, err
				}
				taskList = append(taskList, ActivityCouple{
					IsCloud: false,
					Server: server,
					Instance: instance,
					Task: task,
					Plans: filterPlansByServer(server.Id, false, network),
					NewInfra: true,
					Project: project,
					Infra: infrastructure,
				})
			}
			for _,server := range network.CServers {
				instance, err := filterCloudInstanceByServer(server.Id, infrastructure)
				if err != nil {
					return taskList, err
				}
				taskList = append(taskList, ActivityCouple{
					IsCloud: true,
					CServer: server,
					CInstance: instance,
					Task: task,
					Plans: filterPlansByServer(server.Id, true, network),
					NewInfra: true,
					Project: project,
					Infra: infrastructure,
				})
			}
		}
	}
	return taskList, nil
}

func GetPostBuildTaskActivities(infrastructure model.Infrastructure, task ActivityTask) ([]ActivityCouple, error) {
	var taskList []ActivityCouple = make([]ActivityCouple, 0)
	for _,domain := range infrastructure.Domains {
		for _,network := range domain.Networks {
			for _,server := range network.Instances {
				taskList = append(taskList, ActivityCouple{
					IsCloud: false,
					Instance: server,
					Task: task,
					Plans: []model.InstallationPlan{},
					Infra: infrastructure,
					NewInfra: false,
				})
			}
			for _,server := range network.CInstances {
				taskList = append(taskList, ActivityCouple{
					IsCloud: true,
					CInstance: server,
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
