package main

import (
	"os"
	"fmt"
	"vmkube/model"
)

func init() {
	if len(os.Args) == 0  {
		println("Error: No arguments for command")
		model.PrintCommandHelper("", "")
		os.Exit(1)
	}
	model.InitHelpers()
}

func main() {
	request, error := model.ParseCommandLine(os.Args)
	if error == nil  {
		//fmt.Fprintf(os.Stdout, "Successfully Parser Command : %v\n", request)
		response := model.ExecuteRequest(request)
		if response  {
			fmt.Fprintln(os.Stdout, "Successfully Executed Command!!")
			os.Exit(0)
		} else  {
			fmt.Fprintln(os.Stderr, "Errors During Command Execution!!")
			os.Exit(1)
		}
	} else  {
		os.Exit(1)
	}
}
