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
		fmt.Printf("Sucessfully Parser Command : %v", request)
	} else  {
		os.Exit(1)
	}
}
