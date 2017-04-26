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
	CreateCloudServer(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	CreateServer(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	RemoveServer(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	StopServer(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	StartServer(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	RestartServer(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	ServerStatus(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	ServerEnv(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	ServerInspect(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
	ServerIPAddress(commandPipe chan MachineMessage, commandChannel chan *exec.Cmd)
}

type DockerMachine struct {
	Project     model.Project
	Infra       model.Infrastructure
	IsCloud     bool
	InstanceId  string
	Server      model.ProjectServer
	CServer     model.ProjectCloudServer
	Instance    model.Instance
	CInstance   model.CloudInstance
	NewInfra    bool
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

func GetCurrentServerMachine( Project     model.Project,
												Infra       model.Infrastructure,
												Server      model.ProjectServer,
												CServer     model.ProjectCloudServer,
												Instance    model.Instance,
												CInstance   model.CloudInstance,
												InstanceId  string,
												IsCloud     bool,
												NewInfra    bool) MachineActions {
	return MachineActions(&DockerMachine{
		Project: Project,
		Infra: Infra,
		Server: Server,
		CServer: CServer,
		Instance:Instance,
		CInstance:CInstance,
		InstanceId: InstanceId,
		IsCloud: IsCloud,
		NewInfra: NewInfra,
	})
}