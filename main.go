package main

import (
	"os"
	"vmkube/action"
	"vmkube/term"
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
			term.ScreenPrintln(term.ScreenColor("Errors During Command Execution!!", term.RED))
			term.ScreenFlush()
			os.Exit(1)
		}
	} else  {
		os.Exit(1)
	}
}
