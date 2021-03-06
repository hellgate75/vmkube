package main

import (
	"os"
	"github.com/hellgate75/vmkube/action"
	"github.com/hellgate75/vmkube/common"
	"github.com/hellgate75/vmkube/term"
	"github.com/hellgate75/vmkube/utils"
)

func init() {
	if len(os.Args) == 0 {
		println("Error: No arguments for command")
		common.PrintCommandHelper("help", "help", action.GetArgumentHelpers)
		os.Exit(1)
	}
	action.InitHelpers()
}

func main() {
	request, error := common.ParseCommandLine(os.Args, action.GetArgumentHelpers)
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
