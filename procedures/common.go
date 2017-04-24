package procedures

import (
	"fmt"
	"log"
	"vmkube/model"
	"os/exec"
	"io"
	"strings"
)

func DownloadISO(machineType string, version string) (string, bool) {
	machineAction, error := model.GetMachineAction(machineType)
	if error == nil {
		if ! machineAction.Check(version) {
			fmt.Printf("Machine %s Version %s not present, downloading from internet...\n",strings.ToUpper(machineType),version)
			downloaded := machineAction.Download(version)
			fmt.Printf("Machine %s Version %s dowanloaded: %t\n",strings.ToUpper(machineType),version,downloaded)
			return machineAction.Path(version), downloaded
		} else {
			fmt.Printf("Machine %s Version %s already dowanloaded...\n",strings.ToUpper(machineType),version)
			return machineAction.Path(version), true
		}
	} else {
		log.Fatal(fmt.Sprintf("Machine %s v.%s not found!! - error: \n", strings.ToUpper(machineType),version), error)
		return  "", false
	}
}

type MachineActions interface {
	CreateCloudServer(server model.ProjectCloudServer, commandPipe chan MachineMessage)
	CreateServer(server model.ProjectServer, commandPipe chan MachineMessage)
	RemoveServer(name string, id string, commandPipe chan MachineMessage)
	StopServer(name string, id string, commandPipe chan MachineMessage)
	StartServer(name string, id string, commandPipe chan MachineMessage)
	RestartServer(name string, id string, commandPipe chan MachineMessage)
	ServerStatus(name string, id string, commandPipe chan MachineMessage)
	ServerEnv(name string, id string, commandPipe chan MachineMessage)
	ServerInspect(name string, id string, commandPipe chan MachineMessage)
	ServerIPAddress(name string, id string, commandPipe chan MachineMessage)
}

type SyncDockerMachine struct {
	Project     model.Project
	Infra       model.Infrastructure
	InstanceId  string
	IsCloud     bool
}

type AsyncDockerMachine struct {
	Project     model.Project
	Infra       model.Infrastructure
	InstanceId  string
	IsCloud     bool
}

type MachineOperation int

const(
	CreateCloud MachineOperation = iota
	CreateServer
	DestroyServer
	StartServer
	StopServer
	RestartServer
	StatusServer
	ServerEnvironment
	ServerInspect
	ServerIPAddress
)

type MachineMessage struct {
	Project     model.Project
	Infra       model.Infrastructure
	Operation   MachineOperation
	Cmd         []string
	InstanceId  string
	IsCloud     bool
	Complete    bool
	State       MachineState
	Result      string
	Error       error
	OutReader   io.ReadCloser
	ErrReader   io.ReadCloser
}

func executeSyncCommand(command []string) ([]byte, error) {
	cmd := exec.Command(command[0], command[1:]...)
	return cmd.CombinedOutput()
}

func executeAsyncCommand(message MachineMessage, commandPipe chan MachineMessage) {
	go func(message MachineMessage, commandPipe chan MachineMessage){
		var outBuff io.ReadCloser
		var errBuff io.ReadCloser
		cmd := exec.Command(message.Cmd[0], message.Cmd[1:]...)
		message.Complete = false
		//reader, err := cmd.StdoutPipe()
		//if err == nil {
		//	outBuff = reader
		//}
		//errReader, err := cmd.StderrPipe()
		//if err == nil {
		//	errBuff = errReader
		//}
		message.OutReader = outBuff
		message.ErrReader = errBuff
		commandPipe <- message
		bytes, err := cmd.CombinedOutput()
		message.Result = string(bytes)
		message.Error = err
		message.Complete = true
		commandPipe <- message
		
	}(message, commandPipe)
}