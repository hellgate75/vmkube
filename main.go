package main

import (
	"os"
	"fmt"
	"vmkube/action"
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
	request, error := action.ParseCommandLine(os.Args)
	if error == nil  {
		response := action.ExecuteRequest(request)
		if response  {
			os.Exit(0)
		} else  {
			fmt.Fprintln(os.Stderr, "Errors During Command Execution!!")
			os.Exit(1)
		}
	} else  {
		os.Exit(1)
	}
}
