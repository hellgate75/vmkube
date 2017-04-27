package operations

import (
	"vmkube/term"
	"vmkube/model"
	"vmkube/procedures"
	"errors"
	"os/exec"
)


type RunnableStruct interface {
	Start()
	Stop()
	Status() bool
	Response() interface{}
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
	MachineMessage   procedures.MachineMessage
	SendStartMessage bool
	commandPipe      chan procedures.MachineMessage
	commandChannel   chan *exec.Cmd
	threadSafeCmd    bool
	control          procedures.ControlStructure
}

func (job *ServerOperationsJob) Response() interface{} {
	return []string{job.MachineMessage.InstanceId, job.MachineMessage.InspectJSON, job.MachineMessage.IPAddress}
}

func (job *ServerOperationsJob) Start() {
	if !job.State {
		job.State = true
		if job.SendStartMessage && ! job.control.Interrupt {
			job.OutChan <- job
		}
		job.commandPipe = make(chan procedures.MachineMessage)
		var message procedures.MachineMessage
		machineAdapter := procedures.GetCurrentServerMachine(job.Project, job.Infra, job.Activity.Server, job.Activity.CServer, job.Activity.Instance, job.Activity.CInstance, job.Activity.Instance.Id, job.Activity.IsCloud, job.Activity.NewInfra)
		job.threadSafeCmd = machineAdapter.IsThreadSafeCommand()
		machineAdapter.SetControlStructure(&job.control)
		job.commandChannel = make(chan *exec.Cmd)
		switch job.Activity.Task {
		case CreateMachine:
			if job.Activity.IsCloud {
				go machineAdapter.CreateCloudServer(job.commandPipe, job.commandChannel)
			} else {
				go machineAdapter.CreateServer(job.commandPipe, job.commandChannel)
			}
			break
			case DestroyMachine:
				go machineAdapter.RemoveServer(job.commandPipe, job.commandChannel)
				break
			case StopMachine:
				go machineAdapter.StopServer(job.commandPipe, job.commandChannel)
				break
			case StartMachine:
				go machineAdapter.StartServer(job.commandPipe, job.commandChannel)
				break
			case RestartMachine:
				go machineAdapter.RestartServer(job.commandPipe, job.commandChannel)
				break
			case MachineStatus:
				go machineAdapter.ServerStatus(job.commandPipe, job.commandChannel)
				break
			case MachineEnv:
				go machineAdapter.ServerEnv(job.commandPipe, job.commandChannel)
				break
			case MachineInspect:
				go machineAdapter.ServerInspect(job.commandPipe, job.commandChannel)
				break
			case MachineIPAddress:
				break
				go machineAdapter.ServerIPAddress(job.commandPipe, job.commandChannel)
			default:
				panic("No matching ActivityTask for Job")
		}
		go func(){
			for job.State {
				job.control.CurrentCommand = <- job.commandChannel
			}
		}()
		if ! job.control.Interrupt {
			message = <- job.commandPipe
			job.MachineMessage = message
			job.State = false
			job.OutChan <- job
		}
		job.State = false
	}
}

func (job *ServerOperationsJob) Stop() {
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
