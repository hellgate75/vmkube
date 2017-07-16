package common

import (
	"fmt"
	"strings"
	"time"
	"github.com/hellgate75/vmkube/utils"
)

type Response struct {
	Status  bool
	Message string
}

type CmdRequest struct {
	TypeStr    string
	Type       CmdRequestType
	SubTypeStr string
	SubType    CmdSubRequestType
	HelpType   CmdRequestType
	Element    CmdElementType
	Arguments  *CmdArguments
}

type CmdRequestType int

const (
	NoCommand CmdRequestType = iota
	StartInfrastructure
	StopInfrastructure
	RestartInfrastructure
	DestroyInfrastructure
	AlterInfrastructure
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
	NoSubCommand CmdSubRequestType = iota
	Create
	Remove
	Alter
	Close
	Open
	List
	Detail
	Start
	Stop
	Restart
	Disable
	Enable
	Status
	Recreate
	Destroy
	AutoFix
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
	"Start",
	"Stop",
	"Disable",
	"Enable",
	"Status",
	"Recreate",
	"Destroy",
	"Autofix",
}

type CmdElementType int

const (
	NoElement CmdElementType = iota
	LMachine
	CLMachine
	SPlan
	SNetwork
	SDomain
	SProject
)

var CmdElementTypeDescriptors []string = []string{
	"No Element",
	"Local Machine",
	"Cloud Machine",
	"Installation Plan",
	"Network",
	"Domain",
	"Project",
}

type CmdArguments struct {
	Cmd            string
	CmdType        CmdRequestType
	SubCmd         string
	SubCmdType     CmdSubRequestType
	SubCmdHelpType CmdRequestType
	Element        CmdElementType
	Options        [][]string
	Helper         CommandHelper
}

type CmdParser interface {
	Parse(args []string) bool
}

func (ArgCmd *CmdArguments) Parse(args []string, recoverHelpersFunc func() []CommandHelper) bool {
	if len(args) > 0 {
		command, error := utils.CmdParse(args[0])
		if error == nil {
			for strings.Index(command, "  ") >= 0 {
				command = strings.Replace(command, "  ", "", len(command)/2)
			}
		}
		ArgCmd.Cmd = command
		if error == nil {
			helper := RecoverCommandHelper(command, recoverHelpersFunc)
			ArgCmd.Helper = helper
			ArgCmd.Cmd = helper.Command
			ArgCmd.CmdType = helper.CmdType
			ArgCmd.SubCmd = ""
			ArgCmd.SubCmdType = NoSubCommand
			ArgCmd.SubCmdHelpType = NoCommand
			if len(args) > 1 && len(helper.SubCommands) > 0 {
				var SubCommand string
				var index int
				SubCommand, index, error = CmdParseOption(args[1], helper.SubCommands)
				if error == nil {
					ArgCmd.SubCmd = SubCommand
					if ArgCmd.CmdType != NoCommand {
						ArgCmd.SubCmdType = helper.SubCmdTypes[index]
					} else {
						ArgCmd.SubCmdHelpType = helper.SubCmdHelperTypes[index]
					}
					if ArgCmd.CmdType == NoCommand {
						return true
					}
					if len(args) > 2 {
						optsArgs := args[2:]
						options := make([][]string, 0)
						passed := true
						for index, option := range optsArgs {
							if index%2 == 0 {
								if len(optsArgs) > index+1 {
									key, value, error := utils.OptionsParse(optsArgs[index], optsArgs[index+1])
									if error != nil {
										passed = false
										utils.PrintlnBoldError(fmt.Sprintf("Error: Unable to parse option %s for Command %s and Sub-Command %s", option, command, SubCommand))
										break
									} else {
										if "elem-type" == strings.ToLower(strings.TrimSpace(key)) {
											elementType, error := CmdParseElement(value)
											if error == nil && NoElement != elementType {
												ArgCmd.Element = elementType
											} else {
												utils.PrintlnBoldError(fmt.Sprintf("Error: Invalid infrastructure element type %s for Command %s and Sub-Command %s", value, command, SubCommand))
												if error != nil {
													utils.PrintlnBoldError(fmt.Sprintf("Details: %s", error.Error()))
												}
												time.Sleep(100 * time.Millisecond)
												PrintCommandHelper(command, SubCommand, recoverHelpersFunc)
												return false
											}
										}
										options = append(options, []string{
											key,
											value,
										})
									}
								} else {
									passed = false
									utils.PrintlnBoldError(fmt.Sprintf("Error: Uncomplete option %s for Command %s and Sub-Command %s", option, command, SubCommand))
									time.Sleep(100 * time.Millisecond)
									PrintCommandHelper(command, SubCommand, recoverHelpersFunc)
									return false
								}
							}
						}
						if passed {
							ArgCmd.Options = options
							if "help" != command && "info-project" != command {
								utils.PrintInfo("Executing command ")
								utils.PrintImportant(command)
								utils.PrintlnInfo(" ...")
							}
							return true
						} else {
							utils.PrintlnBoldError("Error: One or more options parse failed!!")
							time.Sleep(100 * time.Millisecond)
							PrintCommandHelper(command, SubCommand, recoverHelpersFunc)
							return false
						}
					}
				} else {
					utils.PrintlnBoldError(fmt.Sprintln("Error:", error))
					time.Sleep(100 * time.Millisecond)
					PrintCommandHelper(command, SubCommand, recoverHelpersFunc)
					return false
				}
				return true
			} else if len(args) >= 1 && len(helper.SubCommands) == 0 {
				if len(args) > 1 {
					optsArgs := args[1:]
					options := make([][]string, 0)
					passed := true
					for index, option := range optsArgs {
						if index%2 == 0 && len(optsArgs) > index+1 {
							key, value, error := utils.OptionsParse(optsArgs[index], optsArgs[index+1])
							if error != nil {
								passed = false
								utils.PrintlnBoldError(fmt.Sprintf("Error: Unable to parse option %s for Command %s", optsArgs[index], command))
								break
							} else {
								if "elem-type" == strings.ToLower(strings.TrimSpace(key)) {
									elementType, error := CmdParseElement(value)
									if error == nil && NoElement != elementType {
										ArgCmd.Element = elementType
									} else {
										utils.PrintlnBoldError(fmt.Sprintf("Error: Invalid infrastructure element type %s for Command %s", value, command))
										if error != nil {
											utils.PrintlnBoldError(fmt.Sprintf("Details: %s", error.Error()))
										}
										time.Sleep(100 * time.Millisecond)
										PrintCommandHelper(command, "", recoverHelpersFunc)
										return false
									}
								}
								options = append(options, []string{
									key,
									value,
								})
							}
						} else if len(optsArgs) < index+1 {
							passed = false
							utils.PrintlnBoldError(fmt.Sprintln("Error: Uncomplete option", option[index], "for Command", command))
							time.Sleep(100 * time.Millisecond)
							PrintCommandHelper(command, "", recoverHelpersFunc)
						}
					}
					if passed {
						ArgCmd.Options = options
						if "help" != command && "info-project" != command {
							utils.PrintInfo("Executing command ")
							utils.PrintImportant(command)
							utils.PrintlnInfo(" ...")
							time.Sleep(100 * time.Millisecond)
						}
						return true
					} else {
						utils.PrintlnBoldError("Error: One or more options parse failed!!")
						time.Sleep(100 * time.Millisecond)
						PrintCommandHelper(command, "", recoverHelpersFunc)
						return false
					}
				}
				if "help" != command && "info-project" != command {
					utils.PrintInfo("Executing command ")
					utils.PrintImportant(command)
					utils.PrintlnInfo(" ...")
					time.Sleep(100 * time.Millisecond)
				}
				return true
			} else if len(args) >= 1 {
				utils.PrintlnBoldError("Error: Unable to parse Sub-Command...")
				time.Sleep(100 * time.Millisecond)
				PrintCommandHelper(command, "", recoverHelpersFunc)
			} else {
				utils.PrintlnBoldError("Error: Unable to parse any parameter...")
				time.Sleep(100 * time.Millisecond)
				PrintCommandHelper(command, "", recoverHelpersFunc)
			}
		} else {
			utils.PrintlnBoldError(fmt.Sprintf("Error: Unable to parse command = %s\n", args[0]))
			time.Sleep(100 * time.Millisecond)
			PrintCommandHelper("help", "", recoverHelpersFunc)
		}
	} else {
		utils.PrintlnBoldError(fmt.Sprintf("Error: Insufficient arguments = %d\n", len(args)))
		time.Sleep(100 * time.Millisecond)
		PrintCommandHelper("help", "", recoverHelpersFunc)
	}
	return false
}
