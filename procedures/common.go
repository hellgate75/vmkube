package procedures

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"vmkube/model"
)

func DownloadISO(machineType string, version string) (string, string, error) {
	machineAction, error := model.GetMachineAction(machineType)
	var log string = ""
	if error == nil {
		if !machineAction.Check(version) {
			log += fmt.Sprintf("OS %s Version %s not present, downloading from internet...\n", strings.ToUpper(machineType), version)
			downloaded := machineAction.Download(version)
			log += fmt.Sprintf("OS %s Version %s dowanloaded: %t\n", strings.ToUpper(machineType), version, downloaded)
			if downloaded {
				return log, machineAction.Path(version), nil
			} else {
				return log, machineAction.Path(version), errors.New(fmt.Sprintf("Unable to download locally OS %s Version %s!!\n", strings.ToUpper(machineType), version))
			}
		} else {
			log += fmt.Sprintf("OS %s Version %s already dowanloaded...\n", strings.ToUpper(machineType), version)
			return log, machineAction.Path(version), nil
		}
	} else {
		log += fmt.Sprintf("OS %s v.%s not found!! - error: %v\n", strings.ToUpper(machineType), version, error)
		return log, "", errors.New(log)
	}
}

type MachineActions interface {
	CreateCloudMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	CreateMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	RemoveMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	StopMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	StartMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	RestartMachine(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	MachineStatus(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	MachineEnv(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	MachineInspect(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	MachineIPAddress(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	MachineExists(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	ExtendsDisk(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	SetControlStructure(Control *MachineControlStructure)
}

type MachineControlStructure struct {
	CurrentCommand *exec.Cmd
	Interrupt      bool
}

type DockerMachineExecutor struct {
	Project    model.Project
	Infra      model.Infrastructure
	IsCloud    bool
	InstanceId string
	Machine    model.LocalMachine
	CMachine   model.CloudMachine
	Instance   model.LocalInstance
	CInstance  model.CloudInstance
	NewInfra   bool
	Control    *MachineControlStructure
}

type MachineOperation int

const (
	CreateCloud MachineOperation = iota
	CreateMachine
	DestroyMachine
	StartMachine
	StopMachine
	RestartMachine
	StatusMachine
	MachineEnvironment
	MachineInspect
	MachineIPAddress
	MachineExists
	ExtendsDisk
)

type MachineMessage struct {
	Project     model.Project
	Infra       model.Infrastructure
	Operation   MachineOperation
	Cmd         []string
	InstanceId  string
	IPAddress   string
	InspectJSON string
	IsCloud     bool
	Complete    bool
	State       MachineState
	Result      string
	Supply      string
	Error       error
	OutReader   io.Reader
	ErrReader   io.Reader
}

func executeSyncCommand(command []string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	return cmd
}

func GetCurrentMachineExecutor(Project model.Project,
	Infra model.Infrastructure,
	Machine model.LocalMachine,
	CMachine model.CloudMachine,
	Instance model.LocalInstance,
	CInstance model.CloudInstance,
	InstanceId string,
	IsCloud bool,
	NewInfra bool) MachineActions {
	return MachineActions(&DockerMachineExecutor{
		Project:    Project,
		Infra:      Infra,
		Machine:    Machine,
		CMachine:   CMachine,
		Instance:   Instance,
		CInstance:  CInstance,
		InstanceId: InstanceId,
		IsCloud:    IsCloud,
		NewInfra:   NewInfra,
	})
}
