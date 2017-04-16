package model

type CommandHelper struct {
	Command							string
	CmdType							CmdRequestType
	LineHelp						string
	SubCommands					[][]string
	SubCmdTypes					[]CmdSubRequestType
	SubCmdHelperTypes		[]CmdRequestType
	Options							[][]string
}

var(
	HelpCommand CommandHelper = CommandHelper{
		Command: "help",
		CmdType: NoCommand,
		LineHelp: "help [COMMAND]",
		SubCommands: [][]string{
		},
		SubCmdHelperTypes: []CmdRequestType{
			StartInfrastructure,
		},
		Options:	[][]string{},
	}
	StartInfra CommandHelper = CommandHelper{
		Command: "start-infra",
		CmdType: StartInfrastructure,
		LineHelp: "start-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	StopInfra CommandHelper = CommandHelper{
		Command: "stop-infra",
		CmdType: StopInfrastructure,
		LineHelp: "stop-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
)

func InitHelpers() {
	HelpCommand.SubCommands = append(HelpCommand.SubCommands,
		[]string{"help", "Show generic commands help"},
		[]string{"start-infra", "Start Ready Infrastructure"},
		[]string{"stop-infra", "Stop Running Infrastructure"},
		[]string{"restart-infra", "Restart Running Infrastructure"},
		[]string{"status-infra", "Require information about an Infrastructure"},
		[]string{"status-all", "Require list and information about all Infrastructures"},
		[]string{"list-projects", "Require list of all released and not released projects"},
		[]string{"alter-project", "Alter existing project, e.g.: open, close, finalize, or simply change nodes"},
		[]string{"define-project", "Define new project"},
		[]string{"delete-project", "Delete existing closed project"},
		[]string{"import-project", "Import project from existing configuration"},
		[]string{"export-project", "Export existing project configuration"},
	)
	StartInfra.Options = append(StartInfra.Options,
		[]string{"name", " <project name>", "Infrastrucure project defined name"},
	)

	StopInfra.Options = append(StopInfra.Options,
		[]string{"name", " <project name>", "Infrastrucure project defined name"},
	)

}

func ParseArgumentHelper() []CommandHelper {
	return  []CommandHelper{
		HelpCommand,
		StartInfra,
		StopInfra,
	}
}