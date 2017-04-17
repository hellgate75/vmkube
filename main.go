package main

import (
	"os"
	"fmt"
	"vmkube/action"
	"strings"
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
