package tasks

import (
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
	"vmkube/model"
	"vmkube/procedures"
	"vmkube/utils"
)

func (job *MachineOperationsJob) Start(exitChannel chan bool) {
	if !job.State {
		defer func() {
			job.State = false
			exitChannel <- true
		}()
		job.State = true
		name := ""
		if job.Activity.NewInfra {
			if job.Activity.IsCloud {
				name = job.Activity.CMachine.Name
			} else {
				name = job.Activity.Machine.Name
			}
		} else {
			if job.Activity.IsCloud {
				name = job.Activity.CInstance.Name
			} else {
				name = job.Activity.Instance.Name
			}
		}
		if job.SendStartMessage {
			if job.control.Interrupt {
				job.MachineMessage.Complete = true
				job.MachineMessage.Error = errors.New(fmt.Sprintf("Interrupted Machine %s Command %s", name, ConvertActivityTaskInString(job.Activity.Task)))
			}
			defer func() {
				// recover from panic caused by writing to a closed channel
				if r := recover(); r != nil {
				}
			}()
			job.OutChan <- job
		}
		if job.control.Interrupt {
			return
		}
		job.commandPipe = make(chan procedures.MachineMessage)
		machineAdapter := procedures.GetCurrentMachineExecutor(job.Project, job.Infra, job.Activity.Machine, job.Activity.CMachine, job.Activity.Instance, job.Activity.CInstance, job.Activity.Instance.Id, job.Activity.IsCloud, job.Activity.NewInfra)
		machineAdapter.SetControlStructure(&job.control)
		job.commandChannel = make(chan *exec.Cmd)
		go func() {
			for job.State {
				job.control.CurrentCommand = <-job.commandChannel
			}
		}()
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
			go machineAdapter.MachineIPAddress(job.commandPipe, job.commandChannel)
			break
		case MachineExtendsDisk:
			go machineAdapter.ExtendsDisk(job.commandPipe, job.commandChannel)
			break
		default:
			panic("No matching ActivityTask for Job")
		}
		if job.control.Interrupt {
			job.MachineMessage.Complete = true
			job.MachineMessage.Error = errors.New(fmt.Sprintf("Interrupted Machine %s Command %s", name, ConvertActivityTaskInString(job.Activity.Task)))
		} else {
			var message procedures.MachineMessage
			select {
			case message = <-job.commandPipe:
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
				job.MachineMessage.Complete = true
			case <-time.After(time.Second * MachineReadOperationTimeout):
				job.MachineMessage.Complete = true
				job.MachineMessage.Error = errors.New(fmt.Sprintf("Timeout for Machine %s Command %s reached", name, ConvertActivityTaskInString(job.Activity.Task)))
			}
		}
		if !job.control.Interrupt {
			defer func() {
				// recover from panic caused by writing to a closed channel
				if r := recover(); r != nil {
				}
			}()
			job.OutChan <- job
		}
	}
}

func (job *MachineOperationsJob) getInstanceId() string {
	if job.Activity.IsCloud {
		return job.Activity.CInstance.Id
	}
	return job.Activity.Instance.Id
}

func (job *MachineOperationsJob) Response() interface{} {
	if job.MachineMessage.InspectJSON != "" {
		return fmt.Sprintf("%s|%s|%s|%s|%s", "json", job.getInstanceId(), job.MachineMessage.InspectJSON, job.MachineMessage.Supply, job.MachineMessage.Result)
	} else if job.MachineMessage.IPAddress != "" {
		return fmt.Sprintf("%s|%s|%s|%s|%s", "ip", job.getInstanceId(), job.MachineMessage.IPAddress, job.MachineMessage.Supply, job.MachineMessage.Result)
	}
	return fmt.Sprintf("%s|%s|%s|%s", "message", job.getInstanceId(), job.MachineMessage.Supply, job.MachineMessage.Result)
}

func (job *MachineOperationsJob) WaitFor() {
	for job.State {
		time.Sleep(1 * time.Second)
	}
}

func (job *MachineOperationsJob) IsError() bool {
	return job.MachineMessage.Error != nil
}

func (job *MachineOperationsJob) Stop() {
	if job.State {
		job.control.Interrupt = true
		if job.control.CurrentCommand != nil {
			if job.control.CurrentCommand != nil && job.control.CurrentCommand.Process != nil && job.control.CurrentCommand.Process.Pid > 0 {
				job.control.CurrentCommand.Process.Kill()
			}
		}
		if job.commandPipe != nil {
			close(job.commandPipe)
			job.commandPipe = nil
		}
		if job.commandChannel != nil {
			close(job.commandChannel)
			job.commandChannel = nil
		}
		job.State = false
	}
}

func (job *MachineOperationsJob) IsInterrupted() bool {
	return job.control.Interrupt
}

func (job *MachineOperationsJob) Status() bool {
	return job.State
}

type ActivityTask int

const (
	CreateMachine ActivityTask = iota
	DestroyMachine
	StopMachine
	StartMachine
	RestartMachine
	MachineStatus
	MachineEnv
	MachineInspect
	MachineIPAddress
	MachineExtendsDisk
)

type ActivityGroup struct {
	Name       string
	Subject    string
	Activities []ActivityCouple
	NewInfra   bool
	Task       ActivityTask
	IsCloud    bool
}

type ActivityCouple struct {
	Project    model.Project
	Infra      model.Infrastructure
	IsCloud    bool
	IsInstance bool
	Machine    model.LocalMachine
	CMachine   model.CloudMachine
	Instance   model.LocalInstance
	CInstance  model.CloudInstance
	Plans      []model.InstallationPlan
	Task       ActivityTask
	NewInfra   bool
}

func filterPlansByMachine(id string, isCloud bool, network model.MachineNetwork) []model.InstallationPlan {
	var selectedPlans []model.InstallationPlan = make([]model.InstallationPlan, 0)
	for _, plan := range network.Installations {
		if plan.IsCloud == isCloud && plan.MachineId == id {
			selectedPlans = append(selectedPlans, plan)
		}
	}
	return selectedPlans
}

func filterPlansByInstance(id string, isCloud bool, network model.Network) []model.InstallationPlan {
	var selectedPlans []model.InstallationPlan = make([]model.InstallationPlan, 0)
	for _, plan := range network.Installations {
		if plan.IsCloud == isCloud && plan.InstanceId == id {
			if !plan.Success {
				selectedPlans = append(selectedPlans, plan.Plan)
			}
		}
	}
	return selectedPlans
}

func filterInstanceByMachine(id string, infrastructure model.Infrastructure) (model.LocalInstance, error) {
	var Instance model.LocalInstance
	for _, domain := range infrastructure.Domains {
		for _, network := range domain.Networks {
			for _, machine := range network.LocalInstances {
				if machine.MachineId == id {
					return machine, nil
				}
			}
		}
	}
	return Instance, errors.New("Instance for Machine Id: " + id + " not found")
}

func filterCloudInstanceByMachine(id string, infrastructure model.Infrastructure) (model.CloudInstance, error) {
	var Instance model.CloudInstance
	for _, domain := range infrastructure.Domains {
		for _, network := range domain.Networks {
			for _, machine := range network.CloudInstances {
				if machine.MachineId == id {
					return machine, nil
				}
			}
		}
	}
	return Instance, errors.New("Instance for Machine Id: " + id + " not found")
}

func GetTaskActivities(project model.Project, infrastructure model.Infrastructure, task ActivityTask) ([]ActivityCouple, error) {
	return GetTaskActivitiesExclusion(project, infrastructure, task, []string{})
}

func GetTaskActivitiesExclusion(project model.Project, infrastructure model.Infrastructure, task ActivityTask, exclusionList []string) ([]ActivityCouple, error) {
	var taskList []ActivityCouple = make([]ActivityCouple, 0)
	for _, domain := range project.Domains {
		for _, network := range domain.Networks {
			for _, machine := range network.LocalMachines {
				instance, err := filterInstanceByMachine(machine.Id, infrastructure)
				if err != nil {
					return taskList, err
				}
				if task == MachineExtendsDisk && utils.CorrectInput(machine.Driver) != "virtualbox" && utils.CorrectInput(machine.Driver) != "vmwarefusion" && utils.CorrectInput(machine.Driver) != "vmwarevsphere" {
					continue
				} else if machine.DiskSize == 0 {
					continue
				}
				if containsString(exclusionList, instance.Id) {
					continue
				}
				taskList = append(taskList, ActivityCouple{
					IsCloud:  false,
					Machine:  machine,
					Instance: instance,
					Task:     task,
					Plans:    filterPlansByMachine(machine.Id, false, network),
					NewInfra: true,
					Project:  project,
					Infra:    infrastructure,
				})
			}
			for _, machine := range network.CloudMachines {
				if task == MachineExtendsDisk {
					continue
				}
				instance, err := filterCloudInstanceByMachine(machine.Id, infrastructure)
				if err != nil {
					return taskList, err
				}
				if containsString(exclusionList, instance.Id) {
					continue
				}
				taskList = append(taskList, ActivityCouple{
					IsCloud:   true,
					CMachine:  machine,
					CInstance: instance,
					Task:      task,
					Plans:     filterPlansByMachine(machine.Id, true, network),
					NewInfra:  true,
					Project:   project,
					Infra:     infrastructure,
				})
			}
		}
	}
	return taskList, nil
}

func containsString(intSlice []string, searchInt string) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}

func GroupActivitiesBySubject(activities []ActivityCouple) []ActivityGroup {
	var groups []ActivityGroup = make([]ActivityGroup, 0)
	var collector map[string][]ActivityCouple = make(map[string][]ActivityCouple)
	for i := 0; i < len(activities); i++ {
		var name string = ""
		if activities[i].NewInfra {
			if !activities[i].IsCloud {
				name = activities[i].Machine.Name + "-" + activities[i].Machine.Id
			} else {
				name = activities[i].CMachine.Name + "-" + activities[i].CMachine.Id
			}
		} else {
			if !activities[i].IsCloud {
				name = activities[i].Instance.Name + "-" + activities[i].Instance.Id
			} else {
				name = activities[i].CInstance.Name + "-" + activities[i].CInstance.Id
			}
		}
		if _, ok := collector[name]; !ok {
			collector[name] = []ActivityCouple{}
		}
		collector[name] = append(collector[name], activities[i])
	}
	for name := range collector {
		value, _ := collector[name]
		groups = append(groups, ActivityGroup{
			Name:       name,
			Subject:    strings.Split(name, "-")[0],
			Activities: value,
			NewInfra:   value[0].NewInfra,
			Task:       value[0].Task,
			IsCloud:    value[0].IsCloud,
		})
	}
	SortGroups(groups)
	return groups
}

type SortGroupType []ActivityGroup

func (a SortGroupType) Len() int           { return len(a) }
func (a SortGroupType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGroupType) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) < 0 }

func SortGroups(groups []ActivityGroup) {
	sort.Sort(SortGroupType(groups))
}

func GetExclusionListExceptInstanceList(infrastructure model.Infrastructure, instanceIds []string) []string {
	var exclusionList []string = make([]string, 0)
	for _, domain := range infrastructure.Domains {
		for _, network := range domain.Networks {
			for _, instance := range network.LocalInstances {
				if !containsString(instanceIds, instance.Id) {
					exclusionList = append(exclusionList, instance.Id)
				}
			}
			for _, instance := range network.CloudInstances {
				if !containsString(instanceIds, instance.Id) {
					exclusionList = append(exclusionList, instance.Id)
				}
			}
		}
	}
	return exclusionList
}

func GetPostBuildTaskActivities(infrastructure model.Infrastructure, task ActivityTask, exclusionList []string) ([]ActivityCouple, error) {
	var taskList []ActivityCouple = make([]ActivityCouple, 0)
	for _, domain := range infrastructure.Domains {
		for _, network := range domain.Networks {
			for _, instance := range network.LocalInstances {
				if !containsString(exclusionList, instance.Id) {
					if task == MachineExtendsDisk && utils.CorrectInput(instance.Driver) != "virtualbox" && utils.CorrectInput(instance.Driver) != "vmwarefusion" && utils.CorrectInput(instance.Driver) != "vmwarevsphere" {
						continue
					}
					taskList = append(taskList, ActivityCouple{
						IsCloud:  false,
						Instance: instance,
						Task:     task,
						Plans:    filterPlansByInstance(instance.Id, false, network),
						Infra:    infrastructure,
						NewInfra: false,
					})
				}
			}
			for _, instance := range network.CloudInstances {
				if !containsString(exclusionList, instance.Id) {
					if task == MachineExtendsDisk {
						continue
					}
					taskList = append(taskList, ActivityCouple{
						IsCloud:   true,
						CInstance: instance,
						Task:      task,
						Plans:     filterPlansByInstance(instance.Id, true, network),
						Infra:     infrastructure,
						NewInfra:  false,
					})
				}
			}
		}
	}
	return taskList, nil
}
