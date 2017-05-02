package main

import (
	"os"
	"vmkube/action"
	"vmkube/term"
	"vmkube/utils"
)

func init() {
	if len(os.Args) == 0 {
		println("Error: No arguments for command")
		action.PrintCommandHelper("help", "help")
		os.Exit(1)
	}
	action.InitHelpers()
}

func main() {
	request, error := action.ParseCommandLine(os.Args)
	if error == nil {
		response := action.ExecuteRequest(request)
		if response {
			os.Exit(0)
		} else {
			utils.PrintlnBoldError("Errors During Command Execution!!")
			term.Screen.Flush()
			os.Exit(1)
		}
	} else {
		os.Exit(1)
	}
}
