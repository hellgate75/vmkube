package model

import (
	"vmkube/utils"
	"fmt"
	"os"
)

type CmdRequestType int

const (
	NoCommand								CmdRequestType = iota;
	StartInfrastructure			CmdRequestType = iota + 1;
	StopInfrastructure			CmdRequestType = iota + 1;
	RestartInfrastructure		CmdRequestType = iota + 1;
	DestroyInfrastructure		CmdRequestType = iota + 1;
	ListInfrastructure			CmdRequestType = iota + 1;
	ListInfrastructures			CmdRequestType = iota + 1;
	ListConfigs							CmdRequestType = iota + 1;
	StatusConfig							CmdRequestType = iota + 1;
	ImportConfig						CmdRequestType = iota + 1;
	ExportConfig						CmdRequestType = iota + 1;
	DefineConfig						CmdRequestType = iota + 1;
	DeleteConfig						CmdRequestType = iota + 1;
	AlterConfig							CmdRequestType = iota + 1;
)

type CmdSubRequestType int

const (
	NoSubCommand	CmdSubRequestType = iota;
	Create 		CmdSubRequestType = iota + 1;
	Remove 		CmdSubRequestType = iota + 1;
	List 			CmdSubRequestType = iota + 1;
	Add  			CmdSubRequestType= iota + 1;
	Finalize  CmdSubRequestType= iota + 1;
	Export  	CmdSubRequestType= iota + 1;
)

type CmdElementType int

const (
	NoElement					CmdElementType = iota;
	LServer					CmdElementType = iota + 1;
	CLServer					CmdElementType = iota + 1;
	SNetwork						CmdElementType = iota + 1;
	SDomain						CmdElementType = iota + 1;
	SInfrastructure		CmdElementType = iota + 1;
)


type CmdRequest struct {
	TypeStr					string
	Type						CmdRequestType
	SubTypeStr			string
	SubType					CmdSubRequestType
	HelpType				CmdRequestType
	CmdElementType	CmdElementType
	Arguments				[][]string
}

type CmdArguments struct {
	Cmd							string
	CmdType					CmdRequestType
	SubCmd					string
	SubCmdType			CmdSubRequestType
	SubCmdHelpType	CmdRequestType
	Options					[][]string
}

type CmdParser interface {
	Parse(args []string) bool
}

func (ArgCmd *CmdArguments) Parse(args []string) bool {
	if len(args) > 0 {
		command, error := utils.CmdParse(args[0])
		ArgCmd.Cmd = command
		if error == nil {
			//fmt.Fprintf(os.Stdout, "Arguments: %v\n", args)
			helper := RecoverCommandHelper(command)
			//fmt.Fprintf(os.Stdout, "Helper: %v\n", helper)
			ArgCmd.Cmd = helper.Command
			ArgCmd.CmdType = helper.CmdType
			ArgCmd.SubCmd = ""
			ArgCmd.SubCmdType = NoSubCommand
			ArgCmd.SubCmdHelpType = NoCommand
			if len(args) > 1 && len(helper.SubCommands) > 0  {
				var SubCommand string
				var  index int
				SubCommand, index, error = utils.CmdParseOption(args[1], helper.SubCommands)
				//fmt.Fprintf(os.Stdout, "Index: %d\n", index)
				if error == nil  {
					ArgCmd.SubCmd = SubCommand
					if ArgCmd.CmdType != NoCommand {
						ArgCmd.SubCmdType = helper.SubCmdTypes[index]
					} else  {
						ArgCmd.SubCmdHelpType = helper.SubCmdHelperTypes[index]
					}
					if len(args) > 2 {
						optsArgs  := args[2:]
						options  := make([][]string, 0)
						passed := true
						for index, option := range optsArgs {
							if index % 2 == 0 && len(optsArgs) > index + 1 {
								key, value, error := utils.OptionsParse(optsArgs[index], optsArgs[index+1])
								if error != nil {
									passed = false
									fmt.Fprintln(os.Stdout, "Unable to parse option", option[index],"for Command",command,"and Sub-Command",SubCommand)
								} else {
									options = append(options, []string{
										key,
										value,
									})
								}
							} else if len(optsArgs) <= index + 1 {
								passed = false
								fmt.Fprintln(os.Stdout, "Uncompleted option", option[index],"for Command",command,"and Sub-Command",SubCommand)
							}
						}
						if passed  {
							ArgCmd.Options = options
							fmt.Fprintf(os.Stdout, "Executing command %s ...\n", command)
							return  true
						} else  {
							fmt.Fprintln(os.Stderr, "One or more options parse failed!!")
							PrintCommandHelper(command, SubCommand)
							return  false
						}
					}
				} else {
					fmt.Fprintln(os.Stderr, "Error:", error)
					PrintCommandHelper(command, SubCommand)
					return  false
				}
				return  true
			} else if len(args) >= 1 && len(helper.SubCommands) == 0 {
				if len(args) > 1 {
					optsArgs  := args[1:]
					options  := make([][]string, 0)
					passed := true
					for index, option := range optsArgs {
						if index % 2 == 0 && len(optsArgs) > index + 1 {
							key, value, error := utils.OptionsParse(optsArgs[index], optsArgs[index+1])
							if error != nil {
								passed = false
								fmt.Fprintln(os.Stdout, "Unable to parse option", option[index],"for Command",command)
							} else {
								options = append(options, []string{
									key,
									value,
								})
							}
						} else if len(optsArgs) <= index + 1 {
							passed = false
							fmt.Fprintln(os.Stdout, "Uncompleted option", option[index],"for Command",command)
						}
					}
					if passed  {
						ArgCmd.Options = options
						fmt.Fprintf(os.Stdout, "Executing command %s ...\n", command)
						return  true
					} else  {
						fmt.Fprintln(os.Stderr, "One or more options parse failed!!")
						PrintCommandHelper(command, "")
						return  false
					}
				}
				fmt.Fprintf(os.Stdout, "Executing command %s ...\n", command)
				return  true
			} else if len(args) >= 1 {
				fmt.Fprintln(os.Stderr, "Error: Unable to parse Sub-Command...")
				PrintCommandHelper(command, "")
			} else  {
				fmt.Fprintln(os.Stderr, "Error: Unable to parse any parameter...")
				PrintCommandHelper(command, "")
			}
		} else {
			fmt.Fprintln(os.Stderr, "Error: Unable to parse command =",args[0])
			PrintCommandHelper("help", "")
		}
	} else {
		fmt.Fprintln(os.Stderr, "Error: Insufficient arguments =",len(args))
		PrintCommandHelper("help", "")
	}
	return  false
}