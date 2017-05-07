package common

import "errors"

/*
Describe Command Helper options, contains

  * Command           (string)              Execution Command Text

	* Name              (string)              Logical Command Name

	* Description       (string)              Command Full Description

  * CmdType           (CmdRequestType)      Dimension of disk root (in GB)

	* LineHelp          (string)              Command Line Helper Text

  * SubCommands       ([][]string)          Sub Commands Descriptors

  * SubCmdTypes       ([]CmdSubRequestType) Sub Commands Types list

  * SubCmdHelperTypes ([]CmdRequestType)    Main Helper Available Commands List

  * Options           ([][]string)          Command Line Options
*/
type CommandHelper struct {
	Command           string              `json:"Command" xml:"Command" mandatory:"yes" descr:"Required Command" type:"text"`
	Name              string              `json:"Name" xml:"Name" mandatory:"yes" descr:"Command Name" type:"text"`
	Description       string              `json:"Description" xml:"Description" mandatory:"yes" descr:"Required Command" type:"text"`
	CmdType           CmdRequestType      `json:"CmdType" xml:"CmdType" mandatory:"yes" descr:"Required Command" type:"integer"`
	LineHelp          string              `json:"LineHelp" xml:"LineHelp" mandatory:"yes" descr:"Required Command" type:"text"`
	SubCommands       []SubCommandHelper  `json:"SubCommands" xml:"SubCommands" mandatory:"yes" descr:"Required Command" type:"list of SubCommandHelper objects"`
	SubCmdTypes       []CmdSubRequestType `json:"SubCmdTypes" xml:"SubCmdTypes" mandatory:"yes" descr:"Required Command" type:"integer list"`
	SubCmdHelperTypes []CmdRequestType    `json:"SubCmdHelperTypes" xml:"SubCmdHelperTypes" mandatory:"yes" descr:"Required Command" type:"integer list"`
	Options           []HelperOption      `json:"Options" xml:"Options" mandatory:"yes" descr:"Required Command" type:"tlist of HelperOption objects"`
}

func (Helper *CommandHelper) HasOption(value string) bool {
	for _, option := range Helper.Options {
		if CorrectInput(option.Option) == CorrectInput(value) ||
			CorrectInput(option.Short) == CorrectInput(value) {
			return true
		}
	}
	return false
}

type SubCommandHelper struct {
	Command     string `json:"Command" xml:"Command" mandatory:"yes" descr:"Required Sub-Command" type:"text"`
	Description string `json:"Description" xml:"Description" mandatory:"yes" descr:"Sub-Command Description" type:"text"`
}

type HelperOption struct {
	Option      string `json:"Option" xml:"Option" mandatory:"yes" descr:"Defined Option" type:"text"`
	Short       string `json:"Short" xml:"Short" mandatory:"no" descr:"Defined Short Option" type:"text"`
	Type        string `json:"Type" xml:"Type" mandatory:"yes" descr:"Defined Option Type Desription" type:"text"`
	Description string `json:"Description" xml:"Description" mandatory:"yes" descr:"Defined Option Desription" type:"text"`
	Mandatory   bool   `json:"Mandatory" xml:"Mandatory" mandatory:"yes" descr:"Describe a Mandatory option" type:"boolean"`
}

func (Option *HelperOption) Match(value string) bool {
	if CorrectInput(Option.Option) == CorrectInput(value) ||
		CorrectInput(Option.Short) == CorrectInput(value) {
		return true
	}
	return false
}

func (Option *HelperOption) Equals(option HelperOption) bool {
	if CorrectInput(Option.Option) == CorrectInput(Option.Option) ||
		CorrectInput(Option.Short) == CorrectInput(Option.Short) {
		return true
	}
	return false
}

func ParseCommandArguments(args []string, recoverHelpersFunc func() []CommandHelper) (*CmdArguments, error) {
	arguments := CmdArguments{}
	success := arguments.Parse(args[1:], recoverHelpersFunc)
	if success {
		return &arguments, nil
	} else {
		return &arguments, errors.New("Unable to Parse Command Line")
	}
}

func ParseCommandLine(args []string, recoverHelpersFunc func() []CommandHelper) (CmdRequest, error) {
	request := CmdRequest{}
	var arguments *CmdArguments
	var err error
	arguments, err = ParseCommandArguments(args, recoverHelpersFunc)
	if err == nil {
		request.TypeStr = arguments.Cmd
		request.Type = arguments.CmdType
		request.SubTypeStr = arguments.SubCmd
		request.SubType = arguments.SubCmdType
		request.HelpType = arguments.SubCmdHelpType
		request.Arguments = arguments
	}
	return request, err
}

func CmdParseElement(value string) (CmdElementType, error) {
	switch CorrectInput(value) {
	case "local-machine":
		return LMachine, nil
	case "cloud-machine":
		return CLMachine, nil
	case "network":
		return SNetwork, nil
	case "domain":
		return SDomain, nil
	case "project":
		return SProject, nil
	case "plan":
		return SPlan, nil
	default:
		return NoElement, errors.New("Element '" + value + "' is not an infratructure element. Available ones : Local-Machine, Cloud-Machine, Network, Domain, Plan, Project")

	}
}
