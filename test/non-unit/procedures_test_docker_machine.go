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
	DiskSize: 0,
	Cpus: 2,
	Driver: "virtualbox",
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
	DiskSize: 25,
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

func TestSyncDockerMachineCreation() {
	var mySyncDockerMachine procedures.SyncDockerMachine = procedures.SyncDockerMachine{
		Infra: myInfra,
		Project: myProject,
		InstanceId: Instance1.Id,
		IsCloud: false,
	}
	responseChan := make(chan procedures.MachineMessage)
	go mySyncDockerMachine.CreateServer(Server1, responseChan)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go func(wGroup *sync.WaitGroup, responseChan chan procedures.MachineMessage) {
		myResponse := <- responseChan
		close(responseChan)
		fmt.Printf("SyncCommand : %s\nSuccess: %t\n", strings.Join(myResponse.Cmd, " "), (myResponse.Error == nil))
		fmt.Printf("Response : %s\nComplete: %t\nError: %v\n", myResponse.Result, myResponse.Complete, myResponse.Error)
		wGroup.Done()
	}(&wGroup, responseChan)
	wGroup.Wait()
}

func TestASyncDockerMachineCreation() {
	var mySyncDockerMachine procedures.AsyncDockerMachine = procedures.AsyncDockerMachine{
		Infra: myInfra,
		Project: myProject,
		InstanceId: Instance1.Id,
		IsCloud: false,
	}
	responseChan := make(chan procedures.MachineMessage)
	mySyncDockerMachine.CreateServer(Server2, responseChan)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go func(wGroup *sync.WaitGroup, responseChan chan procedures.MachineMessage) {
		workOn := true
		for workOn {
			myResponse := <- responseChan
			//fmt.Printf("Task message : %v \n", myResponse)
			if myResponse.Complete {
				fmt.Println("Task completed ...")
				close(responseChan)
				fmt.Printf("ASyncCommand : %s\nSuccess: %t\n", strings.Join(myResponse.Cmd, " "), (myResponse.Error == nil))
				fmt.Printf("Response : %s\nComplete: %t\nError: %v\n", myResponse.Result, myResponse.Complete, myResponse.Error)
				break
			} else {
				fmt.Println("Opening stream ...")
				//go func(){
				//	for workOn {
				//		var buffRead []byte = make([]byte,1024)
				//		intOut,err := myResponse.OutReader.Read(buffRead)
				//		if err == nil && intOut > 0 {
				//			fmt.Println(string(buffRead))
				//		}
				//		buffRead = make([]byte,1024)
				//		intErr,err := myResponse.ErrReader.Read(buffRead)
				//		if err == nil && intErr > 0 {
				//			fmt.Println(string(buffRead))
				//		}
				//	}
				//}()
			}
		}
		wGroup.Done()
	}(&wGroup, responseChan)
	wGroup.Wait()
}

func main() {
	//TestSyncDockerMachineCreation()
	TestASyncDockerMachineCreation()
}