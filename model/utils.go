package model

import (
	"strings"
	"fmt"
	"os"
	"errors"
	"vmkube/utils"
)

func RecoverCommandHelper(helpCommand string) CommandHelper {
	helperCommands := ParseArgumentHelper()
	for _, helper := range helperCommands {
		if helper.Command == strings.ToLower(helpCommand) {
			return  helper
		}
	}
	return  helperCommands[0]
}


func ParseCommandArguments(args	[]string) (*CmdArguments, error) {
	arguments := CmdArguments{}
	success := arguments.Parse(args[1:])
	if success  {
		return  &arguments, nil
	} else  {
		return  &arguments, errors.New("Unable to Parse Command Line")
	}
}

func ParseCommandLine(args []string) (CmdRequest, error) {
	request := CmdRequest{}
	arguments, error := ParseCommandArguments(args)
	if error == nil  {
		request.TypeStr = arguments.Cmd
		request.Type = arguments.CmdType
		request.SubTypeStr = arguments.SubCmd
		request.SubType = arguments.SubCmdType
		request.HelpType = arguments.SubCmdHelpType
		//request.CmdElementType = ??
		request.Arguments = arguments.Options
	}
	return  request, error
}


func PrintCommandHelper(command	string, subCommand string) {
	helper := RecoverCommandHelper(command)
	//fmt.Fprintln(os.Stdout, "vmkube Command [SubCommand] [OPTIONS]")
	//fmt.Fprintln(os.Stdout, "Required Command:", command)
	//fmt.Fprintln(os.Stdout, "Parsed Command:", helper.Command)
	fmt.Fprintln(os.Stdout, "Help: vmkube", helper.LineHelp)
	fmt.Fprintln(os.Stdout, "Action:", helper.Description)
	found := false
	if "" !=  strings.TrimSpace(strings.ToLower(subCommand)) && "help" !=  strings.TrimSpace(strings.ToLower(subCommand)) {
		fmt.Fprintln(os.Stdout, "Selected Sub-Command: " + subCommand)
		for _,option := range helper.SubCommands {
			if option[0] == strings.TrimSpace(strings.ToLower(subCommand)) {
				fmt.Fprintln(os.Stdout, "%s\t%s",  utils.StrPad(option[0], 15), option[1])
				found = true
			}
		}
		if ! found  {
			fmt.Fprintln(os.Stdout, "Sub-Command Not found!!")
			fmt.Fprintln(os.Stdout, "Please type: vmkube","help", command,"for full Sub-Command List")
		}
	}  else {
		found = true
		if len(helper.SubCommands) > 0  {
			if len(helper.SubCmdTypes) > 0 {
				fmt.Fprintln(os.Stdout, "Sub-Commands:")
			} else {
				fmt.Fprintln(os.Stdout, "Commands:")
			}
		}
		for _,option := range helper.SubCommands {
			fmt.Fprintf(os.Stdout, "%s\t%s\n",  utils.StrPad(option[0], 15), option[1])
		}
	}
	if found  {
		if len(helper.Options) > 0  {
			fmt.Fprintln(os.Stdout, "Options:")
		}
		for _,option := range helper.Options {
			fmt.Fprintf(os.Stdout, "--%s\t%s\n",  utils.StrPad(option[0]+option[1], 30), option[2])
		}
	} else  {
		fmt.Fprintln(os.Stdout, "Unable to complete help support ...")
	}
}

