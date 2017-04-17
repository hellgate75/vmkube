package main

import (
	"os"
	"fmt"
	"vmkube/action"
	"strings"
	"vmkube/model"
	"vmkube/procedures"
)

func init() {
	if len(os.Args) == 0  {
		println("Error: No arguments for command")
		action.PrintCommandHelper("", "")
		os.Exit(1)
	}
	action.InitHelpers()
}

func main() {
	//server := model.ProjectServer{
	//	Memory: 2048,
	//	Cpus: 2,
	//	Roles: []string{"server","rancher"},
	//	Name: "MyServer",
	//	Driver: "virtualbox",
	//	Hostname: "myserver",
	//	OSType: "rancheros",
	//	OSVersion: "1.0.0",
	//	NoShare: false,
	//	Id: action.NewUUIDString(),
	//	Options: [][]string{},
	//	DiskSize: 60000,
	//	Engine: model.ProjectEngineOpt{
	//		Options: []string{},
	//		Environment: []string{},
	//		InsecureRegistry: []string{},
	//		InstallURL: "",
	//		Labels: []string{},
	//		StorageDriver: "",
	//	},
	//	Swarm: model.ProjectSwarmOpt{
	//		Host: "",
	//		Image: "",
	//		IsMaster: false,
	//		JoinOpts: []string{},
	//		Strategy: "",
	//		TLSSan: []string{},
	//		UseAddress: false,
	//		DiscoveryToken: "",
	//		UseExperimental: false,
	//	},
	//}
	//message, err := procedures.CreateServer(server)
	//fmt.Fprintf(os.Stdout, "Create Server Message: %s\n", message)
	//fmt.Fprintf(os.Stdout, "Create Server Error: %s\n", err)
	//
	//message0, err0 := procedures.ServerStatus(server.Name, server.Id)
	//fmt.Fprintf(os.Stdout, "Status Server Message: %s\n", message0)
	//fmt.Fprintf(os.Stdout, "Status Server Error: %s\n", err0)
	//
	//message1, err1 := procedures.ServerEnv(server.Name, server.Id)
	//fmt.Fprintf(os.Stdout, "Status Server Message: %s\n", message1)
	//fmt.Fprintf(os.Stdout, "Status Server Error: %s\n", err1)
	//
	//message2, err2 := procedures.ServerInspect(server.Name, server.Id)
	//fmt.Fprintf(os.Stdout, "Status Server Message: %s\n", message2)
	//fmt.Fprintf(os.Stdout, "Status Server Error: %s\n", err2)
	//
	//message3, err3 :=  procedures.RemoveServer(server.Name, server.Id)
	//fmt.Fprintf(os.Stdout, "Delete Server Message: %s\n", message3)
	//fmt.Fprintf(os.Stdout, "Delete Server Error: %s\n", err3)
	//if err != nil {
	//	os.Exit(0)
	//}
	//server := model.Server{
	//	Memory: 2048,
	//	Cpus: 2,
	//	Roles: []string{},
	//	Name: "MyServer",
	//	Driver: "virtualbox",
	//	Hostname: "myserver",
	//	OSType: "rancheros",
	//	OSVersion: "0.9.0",
	//	NoShare: false,
	//	Id: action.NewUUIDString(),
	//	Options: [][]string{},
	//	Disks: []model.Disk{},
	//	Engine: model.EngineOpt{
	//		Options: []string{},
	//		Environment: []string{},
	//		InsecureRegistry: []string{},
	//		InstallURL: "",
	//		Labels: []string{},
	//		StorageDriver: "",
	//	},
	//	Swarm: model.SwarmOpt{
	//		Host: "",
	//		Image: "",
	//		IsMaster: false,
	//		JoinOpts: []string{},
	//		Strategy: "",
	//		TLSSan: []string{},
	//		UseAddress: false,
	//		UseDiscovery: false,
	//		UseExperimental: false,
	//	},
	//}
	//fmt.Fprintf(os.Stdout, "Original %v\n", server)
	//err := server.Save("./test.ser")
	//fmt.Fprintf(os.Stdout, "Error %v\n", err)
	//var intent model.Server
	//err = intent.Load("./test.ser")
	//fmt.Fprintf(os.Stdout, "Error %v\n", err)
	//fmt.Fprintf(os.Stdout, "Decoded Map %v\n", intent)
	request, error := action.ParseCommandLine(os.Args)
	if error == nil  {
		//fmt.Fprintf(os.Stdout, "Successfully Parser Command : %v\n", request)
		response := action.ExecuteRequest(request)
		if response  {
			if len(os.Args) <= 1 || "help" != strings.TrimSpace(strings.ToLower(os.Args[1])) {
				fmt.Fprintln(os.Stdout, "Successfully Executed Command!!")
			}
			os.Exit(0)
		} else  {
			fmt.Fprintln(os.Stderr, "Errors During Command Execution!!")
			os.Exit(1)
		}
	} else  {
		os.Exit(1)
	}
}
