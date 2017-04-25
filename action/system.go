package action

import (
	"vmkube/utils"
	"fmt"
	"time"
	"strings"
)

type Response struct {
	Status	bool
	Message	string
}


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
	BackupInfrastructure
	RecoverInfrastructure
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

var CmdRequestDescriptors []string = []string{
	"No Command",
	"Start Infrastructure",
	"Stop Infrastructure",
	"Restart Infrastructure",
	"Destroy Infrastructure",
	"Backup Infrastructure",
	"Recover Infrastructure",
	"Infrastructure Status",
	"List Infrastructures",
	"List Projects",
	"Project Status",
	"Import Project",
	"Export Project",
	"Define Project",
	"BuildC Project",
	"Delete Project",
	"Alter Project",
	"Describe Project",
}

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

var CmdSubRequestDescriptors []string = []string{
	"No Sub-Command",
	"Create",
	"Remove",
	"Modify",
	"Close",
	"Open",
	"List",
	"Detail",
}



type CmdElementType int

const (
	NoElement					CmdElementType = iota
	LServer
	CLServer
	SPlan
	SNetwork
	SDomain
	SProject
)

var CmdElementTypeDescriptors []string = []string{
	"No Element",
	"Local Server",
	"Cloud Server",
	"Installation Plan",
	"Network",
	"Domain",
	"Project",
}


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
				SubCommand, index, error = CmdParseOption(args[1], helper.SubCommands)
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
									utils.PrintlnError(fmt.Sprintf("Error: Unable to parse option %s for Command %s and Sub-Command %s", optsArgs[index], command,SubCommand))
									break
								} else {
									if "elem-type" == strings.ToLower(strings.TrimSpace(key)) {
										elementType, error := CmdParseElement(value)
										if error == nil && NoElement != elementType {
											ArgCmd.Element = elementType
										} else  {
											utils.PrintlnError(fmt.Sprintf("Error: Invalid infrastructure element type %s for Command %s and Sub-Command %s", value, command,SubCommand))
											if error != nil {
												utils.PrintlnError(fmt.Sprintf("Details: %s", error.Error()))
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
								utils.PrintlnError(fmt.Sprintf("Error: Uncomplete option %s for Command %s and Sub-Command %s",option[index],command,SubCommand))
								time.Sleep(100 * time.Millisecond)
								PrintCommandHelper(command, SubCommand)
							}
						}
						if passed  {
							ArgCmd.Options = options
							if "help" != command && "info-project" != command {
								utils.PrintInfo("Executing command ")
								utils.PrintImportant(command)
								utils.PrintlnInfo(" ...")
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
								utils.PrintlnError(fmt.Sprintf("Error: Unable to parse option %s for Command %s", option[index], command))
								break
							} else {
								if "elem-type" == strings.ToLower(strings.TrimSpace(key)) {
									elementType, error := CmdParseElement(value)
									if error == nil && NoElement != elementType {
										ArgCmd.Element = elementType
									} else  {
										utils.PrintlnError(fmt.Sprintf("Error: Invalid infrastructure element type %s for Command %s", value, command))
										if error != nil {
											utils.PrintlnError(fmt.Sprintf("Details: %s", error.Error()))
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
							utils.PrintInfo("Executing command ")
							utils.PrintImportant(command)
							utils.PrintlnInfo(" ...")
							time.Sleep(100 * time.Millisecond)
						}
						return  true
					} else  {
						utils.PrintlnError("Error: One or more options parse failed!!")
						time.Sleep(100 * time.Millisecond)
						PrintCommandHelper(command, "")
						return  false
					}
				}
				if "help" != command && "info-project" != command {
					utils.PrintInfo("Executing command ")
					utils.PrintImportant(command)
					utils.PrintlnInfo(" ...")
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