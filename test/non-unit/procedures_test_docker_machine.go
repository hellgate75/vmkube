package main

import (
	"vmkube/model"
	"time"
	"vmkube/procedures"
	"fmt"
	"strings"
	"sync"
)

var Server1 model.ProjectServer = model.ProjectServer{
	Id: model.NewUUIDString(),
	Name: "Server1",
	Memory: 4096,
	DiskSize: 50,
	Cpus: 2,
	Driver: "vmwarefusion",
	Hostname: "server1",
	OSType: "rancheros",
	OSVersion: "1.0.0",
	NoShare: true,
	Roles: []string{"rancher", "server"},
	Options: [][]string{},
}

var Instance1 model.Instance = model.Instance{
	Id: model.NewUUIDString(),
	Name: "Server1",
	Memory: 4096,
	Disks: []model.Disk{
		{
			Id: model.NewUUIDString(),
			Name: "dev0",
			Size: 50,
			Type: 0,
		},
	},
	Cpus: 2,
	Driver: "virtualbox",
	Hostname: "server1",
	OSType: "rancheros",
	OSVersion: "1.0.0",
	NoShare: true,
	Roles: []string{"rancher", "server"},
}

var Server2 model.ProjectServer = model.ProjectServer{
	Id: model.NewUUIDString(),
	Name: "Server2",
	Memory: 4096,
	DiskSize: 50,
	Cpus: 2,
	Driver: "virtualbox",
	Hostname: "server2",
	OSType: "rancheros",
	OSVersion: "1.0.0",
	NoShare: false,
	Roles: []string{"rancher", "server"},
	Options: [][]string{},
	Engine: model.ProjectEngineOpt{
		StorageDriver: "overlay",
	},
}

var Instance2 model.Instance = model.Instance{
	Id: model.NewUUIDString(),
	Name: "Server2",
	Memory: 4096,
	Disks: []model.Disk{
		{
			Id: model.NewUUIDString(),
			Name: "dev0",
			Size: 50,
			Type: 0,
		},
	},
	Cpus: 2,
	Driver: "virtualbox",
	Hostname: "server2",
	OSType: "rancheros",
	OSVersion: "1.0.0",
	NoShare: true,
	Roles: []string{"rancher", "server"},
}


var myProject model.Project = model.Project{
	Id: model.NewUUIDString(),
	Created: time.Now(),
	Name: "MyBuildProject",
	Errors: false,
	LastMessage: "",
	Modified: time.Now(),
	Open: false,
	Domains: []model.ProjectDomain{
		{
			Id: model.NewUUIDString(),
			Name: "Default Domain",
			Options: [][]string{},
			Networks: []model.ProjectNetwork{
				{
					Id: model.NewUUIDString(),
					Name: "Default Network",
					CServers: []model.ProjectCloudServer{},
					Installations: []model.InstallationPlan{},
					Servers: []model.ProjectServer{
						Server1,
						Server2,
					},
				},
			},
		},
	},
}

var myInfra model.Infrastructure = model.Infrastructure{
	Id: model.NewUUIDString(),
	Created: false,
	Name: "MyBuildProject",
	Errors: false,
	LastMessage: "",
	Modified: time.Now(),
	Domains: []model.Domain{
		{
			Id: model.NewUUIDString(),
			Name: "Default Domain",
			Options: [][]string{},
			Networks: []model.Network{
				{
					Id: model.NewUUIDString(),
					Name: "Default Network",
					CInstances: []model.CloudInstance{},
					Installations: []model.Installation{},
					Instances: []model.Instance{
						Instance1,
						Instance2,
					},
				},
			},
		},
	},
}

func TestVMFusionDockerMachineCreation() {
	var mySyncDockerMachine procedures.DockerMachine = procedures.DockerMachine{
		Infra: myInfra,
		Project: myProject,
		InstanceId: Instance1.Id,
		IsCloud: false,
		Server: Server1,
		Instance: Instance1,
		NewInfra: true,
	}
	responseChan := make(chan procedures.MachineMessage)
	go mySyncDockerMachine.CreateServer(responseChan)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go func(wGroup *sync.WaitGroup, responseChan chan procedures.MachineMessage) {
		myResponse := <- responseChan
		close(responseChan)
		fmt.Printf("SyncCommand : %s\nSuccess: %t\n", strings.Join(myResponse.Cmd, " "), (myResponse.Error == nil))
		fmt.Printf("Response :\n%s\nSupply :\n%s\nComplete: %t\nError: %v\n", myResponse.Result, myResponse.Supply, myResponse.Complete, myResponse.Error)
		fmt.Printf("Inspection JSON :\n%s\nIPAddress : %s\n", myResponse.InspectJSON, myResponse.IPAddress)
		wGroup.Done()
	}(&wGroup, responseChan)
	wGroup.Wait()
}

func TestVirtualBoxDockerMachineCreation() {
	var mySyncDockerMachine procedures.DockerMachine = procedures.DockerMachine{
		Infra: myInfra,
		Project: myProject,
		InstanceId: Instance1.Id,
		IsCloud: false,
		Server: Server2,
		Instance: Instance2,
		NewInfra: true,
	}
	responseChan := make(chan procedures.MachineMessage)
	go mySyncDockerMachine.CreateServer(responseChan)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go func(wGroup *sync.WaitGroup, responseChan chan procedures.MachineMessage) {
		myResponse := <- responseChan
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