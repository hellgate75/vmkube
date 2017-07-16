package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/procedures"
)

var Machine1 model.LocalMachine = model.LocalMachine{
	Id:        model.NewUUIDString(),
	Name:      "Machine1",
	Memory:    4096,
	DiskSize:  50,
	Cpus:      2,
	Driver:    "vmwarefusion",
	Hostname:  "machine1",
	OSType:    "rancheros",
	OSVersion: "1.0.0",
	NoShare:   true,
	Roles:     []string{"rancher", "server"},
	Options:   [][]string{},
}

var Instance1 model.LocalInstance = model.LocalInstance{
	Id:     model.NewUUIDString(),
	Name:   "Machine1",
	Memory: 4096,
	Disks: []model.Disk{
		{
			Id:   model.NewUUIDString(),
			Name: "dev0",
			Size: 50,
			Type: 0,
		},
	},
	Cpus:      2,
	Driver:    "virtualbox",
	Hostname:  "machine1",
	OSType:    "rancheros",
	OSVersion: "1.0.0",
	NoShare:   true,
	Roles:     []string{"rancher", "server"},
}

var Machine2 model.LocalMachine = model.LocalMachine{
	Id:        model.NewUUIDString(),
	Name:      "Machine2",
	Memory:    4096,
	DiskSize:  50,
	Cpus:      2,
	Driver:    "virtualbox",
	Hostname:  "machine2",
	OSType:    "rancheros",
	OSVersion: "1.0.0",
	NoShare:   false,
	Roles:     []string{"rancher", "server"},
	Options:   [][]string{},
	Engine: model.ProjectEngineOpt{
		StorageDriver: "overlay",
	},
}

var Instance2 model.LocalInstance = model.LocalInstance{
	Id:     model.NewUUIDString(),
	Name:   "Machine2",
	Memory: 4096,
	Disks: []model.Disk{
		{
			Id:   model.NewUUIDString(),
			Name: "dev0",
			Size: 50,
			Type: 0,
		},
	},
	Cpus:      2,
	Driver:    "virtualbox",
	Hostname:  "machine2",
	OSType:    "rancheros",
	OSVersion: "1.0.0",
	NoShare:   true,
	Roles:     []string{"rancher", "server"},
}

var myProject model.Project = model.Project{
	Id:          model.NewUUIDString(),
	Created:     time.Now(),
	Name:        "MyBuildProject",
	Errors:      false,
	LastMessage: "",
	Modified:    time.Now(),
	Open:        false,
	Domains: []model.MachineDomain{
		{
			Id:      model.NewUUIDString(),
			Name:    "Default Domain",
			Options: [][]string{},
			Networks: []model.MachineNetwork{
				{
					Id:            model.NewUUIDString(),
					Name:          "Default Network",
					CloudMachines: []model.CloudMachine{},
					Installations: []model.InstallationPlan{},
					LocalMachines: []model.LocalMachine{
						Machine1,
						Machine2,
					},
				},
			},
		},
	},
}

var myInfra model.Infrastructure = model.Infrastructure{
	Id:          model.NewUUIDString(),
	Created:     false,
	Name:        "MyBuildProject",
	Errors:      false,
	LastMessage: "",
	Modified:    time.Now(),
	Domains: []model.Domain{
		{
			Id:      model.NewUUIDString(),
			Name:    "Default Domain",
			Options: [][]string{},
			Networks: []model.Network{
				{
					Id:             model.NewUUIDString(),
					Name:           "Default Network",
					CloudInstances: []model.CloudInstance{},
					Installations:  []model.Installation{},
					LocalInstances: []model.LocalInstance{
						Instance1,
						Instance2,
					},
				},
			},
		},
	},
}

func TestVMFusionDockerMachineCreation() {
	var mySyncDockerMachine procedures.DockerMachineExecutor = procedures.DockerMachineExecutor{
		Infra:      myInfra,
		Project:    myProject,
		InstanceId: Instance1.Id,
		IsCloud:    false,
		Machine:    Machine1,
		Instance:   Instance1,
		NewInfra:   true,
	}
	commandChan := make(chan *exec.Cmd)
	responseChan := make(chan procedures.MachineMessage)
	go mySyncDockerMachine.CreateMachine(responseChan, commandChan)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go func(wGroup *sync.WaitGroup, responseChan chan procedures.MachineMessage) {
		myResponse := <-responseChan
		close(responseChan)
		fmt.Printf("SyncCommand : %s\nSuccess: %t\n", strings.Join(myResponse.Cmd, " "), (myResponse.Error == nil))
		fmt.Printf("Response :\n%s\nSupply :\n%s\nComplete: %t\nError: %v\n", myResponse.Result, myResponse.Supply, myResponse.Complete, myResponse.Error)
		fmt.Printf("Inspection JSON :\n%s\nIPAddress : %s\n", myResponse.InspectJSON, myResponse.IPAddress)
		wGroup.Done()
	}(&wGroup, responseChan)
	wGroup.Wait()
}

func TestVirtualBoxDockerMachineCreation() {
	var mySyncDockerMachine procedures.DockerMachineExecutor = procedures.DockerMachineExecutor{
		Infra:      myInfra,
		Project:    myProject,
		InstanceId: Instance1.Id,
		IsCloud:    false,
		Machine:    Machine2,
		Instance:   Instance2,
		NewInfra:   true,
	}
	commandChan := make(chan *exec.Cmd)
	responseChan := make(chan procedures.MachineMessage)
	go mySyncDockerMachine.CreateMachine(responseChan, commandChan)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go func(wGroup *sync.WaitGroup, responseChan chan procedures.MachineMessage) {
		myResponse := <-responseChan
		close(responseChan)
		fmt.Printf("ASyncCommand : %s\nSuccess: %t\n", strings.Join(myResponse.Cmd, " "), (myResponse.Error == nil))
		fmt.Printf("Response :\n%s\nSupply :\n%s\nComplete: %t\nError: %v\n", myResponse.Result, myResponse.Supply, myResponse.Complete, myResponse.Error)
		fmt.Printf("Inspection JSON :\n%s\nIPAddress : %s\n", myResponse.InspectJSON, myResponse.IPAddress)
		wGroup.Done()
	}(&wGroup, responseChan)
	wGroup.Wait()
}

func main() {
	//TestVMFusionDockerMachineCreation()
	TestVirtualBoxDockerMachineCreation()
}
