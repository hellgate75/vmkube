package action

import (
	"strings"
	"fmt"
	"os"
	"errors"
	"vmkube/utils"
	"github.com/satori/go.uuid"
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
		request.Arguments = arguments
	}
	return  request, error
}


func CmdParseElement(key string, value string) (CmdElementType, error) {
	if "elem-type" == strings.ToLower(strings.TrimSpace(key)) {
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "server":
			return  LServer, nil
		case "cloud-server":
			return  CLServer, nil
		case "network":
			return  SNetwork, nil
		case "domain":
			return  SDomain, nil
		case "project":
			return  SProject, nil
		case "plan":
			return  SPlan, nil
		default:
			return  NoElement, errors.New("Element '"+value+"' is not an infratructure element. Available ones : Server, Cloud-Server, Network, Domain, Plan, Project")

		}
	} else {
		return  NoElement, errors.New("Element '"+key+"' is not an infratructure element key. Available one : --elem-type")
	}
}

func CorrectInput(input string) string {
	return  strings.TrimSpace(strings.ToLower(input))
}

func NewUUIDString()	string {
	return  uuid.NewV4().String()
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
			//if option[0] == strings.TrimSpace(strings.ToLower(subCommand)) {
				fmt.Fprintf(os.Stdout, "%s\t%s\n",  utils.StrPad(option[0], 45), option[1])
				found = true
			//}
		}
		if ! found  {
			fmt.Fprintln(os.Stdout, "Sub-Command Not found!!")
			if "help" !=  strings.TrimSpace(strings.ToLower(command)) {
				fmt.Fprintln(os.Stdout, "Please type: vmkube","help", command,"for full Sub-Command List")
			} else  {
				fmt.Fprintln(os.Stdout, "Please type: vmkube","help", "COMMAND","for full Sub-Command List")
			}
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
			fmt.Fprintf(os.Stdout, "%s\t%s\n",  utils.StrPad(option[0], 45), option[1])
		}
	}
	if found  {
		if len(helper.Options) > 0  {
			fmt.Fprintln(os.Stdout, "Options:")
		}
		for _,option := range helper.Options {
			validity := "optional"
			if "true" == option[3] {
				validity = "mandatory"
			}
			fmt.Fprintf(os.Stdout, "--%s\t%s\t%s\n",  utils.StrPad(option[0]+option[1], 45), utils.StrPad(validity, 10), option[2])
		}
	} else  {
		fmt.Fprintln(os.Stdout, "Unable to complete help support ...")
	}
}

