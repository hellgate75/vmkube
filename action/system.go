package action

import (
	"vmkube/utils"
	"fmt"
	"os"
	"time"
	"strings"
)

type CmdRequest struct {
	TypeStr					string
	Type						CmdRequestType
	SubTypeStr			string
	SubType					CmdSubRequestType
	HelpType				CmdRequestType
	Element					CmdElementType
	Arguments				*CmdArguments
}

type CmdRequestType int

const (
	NoCommand								CmdRequestType = iota
	StartInfrastructure
	StopInfrastructure
	RestartInfrastructure
	DestroyInfrastructure
	ListInfrastructure
	ListInfrastructures
	ListConfigs
	StatusConfig
	ImportConfig
	ExportConfig
	DefineConfig
	BuildConfig
	DeleteConfig
	AlterConfig	
	InfoConfig	
)

type CmdSubRequestType int

const (
	NoSubCommand	CmdSubRequestType = iota
	Create
	Remove
	Alter
	Close
	Open
	List
	Detail
)

type CmdElementType int

const (
	NoElement					CmdElementType = iota
	LServer
	CLServer
	SNetwork
	SDomain
	SProject
	SPlan
)


type CmdArguments struct {
	Cmd							string
	CmdType					CmdRequestType
	SubCmd					string
	SubCmdType			CmdSubRequestType
	SubCmdHelpType	CmdRequestType
	Element					CmdElementType
	Options					[][]string
	Helper					CommandHelper
}

type CmdParser interface {
	Parse(args []string) bool
}

func (ArgCmd *CmdArguments) Parse(args []string) bool {
	if len(args) > 0 {
		command, error := utils.CmdParse(args[0])
		ArgCmd.Cmd = command
		if error == nil {
			helper := RecoverCommandHelper(command)
			ArgCmd.Helper = helper
			ArgCmd.Cmd = helper.Command
			ArgCmd.CmdType = helper.CmdType
			ArgCmd.SubCmd = ""
			ArgCmd.SubCmdType = NoSubCommand
			ArgCmd.SubCmdHelpType = NoCommand
			if len(args) > 1 && len(helper.SubCommands) > 0  {
				var SubCommand string
				var  index int
				SubCommand, index, error = utils.CmdParseOption(args[1], helper.SubCommands)
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
							if index % 2 == 0 && len(optsArgs) >= index + 1 {
								key, value, error := utils.OptionsParse(optsArgs[index], optsArgs[index+1])
								if error != nil {
									passed = false
									fmt.Println("Error: Unable to parse option", optsArgs[index],"for Command",command,"and Sub-Command",SubCommand)
									break
								} else {
									if "elem-type" == strings.ToLower(strings.TrimSpace(key)) {
										elementType, error := CmdParseElement(value)
										if error == nil && NoElement != elementType {
											ArgCmd.Element = elementType
										} else  {
											fmt.Println("Error: Invalid infrastructure element type", value,"for Command",command,"and Sub-Command",SubCommand)
											if error != nil {
												fmt.Fprintf(os.Stderr, "Details: %s \n", error.Error())
											}
											time.Sleep(100 * time.Millisecond)
											PrintCommandHelper(command, SubCommand)
											return  false
										}
									}
									options = append(options, []string{
										key,
										value,
									})
								}
							} else if len(optsArgs) < index + 1 {
								passed = false
								fmt.Println("Error: Uncomplete option", option[index],"for Command",command,"and Sub-Command",SubCommand)
								time.Sleep(100 * time.Millisecond)
								PrintCommandHelper(command, SubCommand)
							}
						}
						if passed  {
							ArgCmd.Options = options
							if "help" != command && "info-project" != command {
								fmt.Printf("Executing command %s ...\n", command)
							}
							return  true
						} else  {
							fmt.Println("One or more options parse failed!!")
							time.Sleep(100 * time.Millisecond)
							PrintCommandHelper(command, SubCommand)
							return  false
						}
					}
				} else {
					fmt.Println("Error:", error)
					time.Sleep(100 * time.Millisecond)
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
						if index % 2 == 0 && len(optsArgs) >= index + 1 {
							key, value, error := utils.OptionsParse(optsArgs[index], optsArgs[index+1])
							if error != nil {
								passed = false
								fmt.Println("Error: Unable to parse option", option[index],"for Command",command)
								break
							} else {
								if "elem-type" == strings.ToLower(strings.TrimSpace(key)) {
									elementType, error := CmdParseElement(value)
									if error == nil && NoElement != elementType {
										ArgCmd.Element = elementType
									} else  {
										fmt.Println("Error: Invalid infrastructure element type", value,"for Command",command)
										if error != nil {
											fmt.Fprintf(os.Stderr, "Details: %s \n", error.Error())
										}
										time.Sleep(100 * time.Millisecond)
										PrintCommandHelper(command, "")
										return  false
									}
								}
								options = append(options, []string{
									key,
									value,
								})
							}
						} else if len(optsArgs) < index + 1 {
							passed = false
							fmt.Println("Error: Uncomplete option", option[index],"for Command",command)
							time.Sleep(100 * time.Millisecond)
							PrintCommandHelper(command, "")
						}
					}
					if passed  {
						ArgCmd.Options = options
						if "help" != command && "info-project" != command {
							fmt.Printf("Executing command %s ...\n", command)
							time.Sleep(100 * time.Millisecond)
						}
						return  true
					} else  {
						fmt.Println("Error: One or more options parse failed!!")
						time.Sleep(100 * time.Millisecond)
						PrintCommandHelper(command, "")
						return  false
					}
				}
				if "help" != command && "info-project" != command {
					fmt.Printf("Executing command %s ...\n", command)
					time.Sleep(100 * time.Millisecond)
				}
				return  true
			} else if len(args) >= 1 {
				fmt.Println("Error: Unable to parse Sub-Command...")
				time.Sleep(100 * time.Millisecond)
				PrintCommandHelper(command, "")
			} else  {
				fmt.Println("Error: Unable to parse any parameter...")
				time.Sleep(100 * time.Millisecond)
				PrintCommandHelper(command, "")
			}
		} else {
			fmt.Printf("Error: Unable to parse command = %s\n",args[0])
			time.Sleep(100 * time.Millisecond)
			PrintCommandHelper("help", "")
		}
	} else {
		fmt.Printf("Error: Insufficient arguments = %d\n",len(args))
		time.Sleep(100 * time.Millisecond)
		PrintCommandHelper("help", "")
	}
	return  false
}