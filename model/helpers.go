package model

type CommandHelper struct {
	Command							string
	Description					string
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
		Description: "Show help tips",
		CmdType: NoCommand,
		LineHelp: "help [COMMAND]",
		SubCommands: [][]string{
		},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{
			NoCommand,
			StartInfrastructure,
			StopInfrastructure,
			RestartInfrastructure,
			DestroyInfrastructure,
			ListInfrastructure,
			ListInfrastructures,
			ListConfigs,
			StatusConfig,
			AlterConfig,
			DefineConfig,
			DeleteConfig,
			ImportConfig,
			ExportConfig,
		},
		Options:	[][]string{},
	}
	StartInfra CommandHelper = CommandHelper{
		Command: "start-infra",
		Description: "Start infrastructre if stopped or nothing",
		CmdType: StartInfrastructure,
		LineHelp: "start-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	StopInfra CommandHelper = CommandHelper{
		Command: "stop-infra",
		Description: "Stop infrastructre if running or nothing",
		CmdType: StopInfrastructure,
		LineHelp: "stop-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	RestartInfra CommandHelper = CommandHelper{
		Command: "restart-infra",
		Description: "Restart infrastructre",
		CmdType: RestartInfrastructure,
		LineHelp: "restart-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	DestroyInfra CommandHelper = CommandHelper{
		Command: "destroy-infra",
		Description: "Destroy a desired infrastructre (No undo available)",
		CmdType: DestroyInfrastructure,
		LineHelp: "destroy-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	ListInfra CommandHelper = CommandHelper{
		Command: "status-infra",
		Description: "List information about a specific infrastructure",
		CmdType: ListInfrastructure,
		LineHelp: "status-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	ListAllInfras CommandHelper = CommandHelper{
		Command: "status-all",
		Description: "List information about all existing infrastructures",
		CmdType: ListInfrastructures,
		LineHelp: "status-all",
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
		[]string{"destroy-infra", "Destroy a specific Infrastructure"},
		[]string{"status-infra", "Require information about a specific Infrastructure"},
		[]string{"status-all", "Require list and information about all Infrastructures"},
		[]string{"list-projects", "Require list of all released and not released projects"},
		[]string{"status-project", "Require information about a specific projects"},
		[]string{"alter-project", "Alter existing project, e.g.: open, close, finalize, or simply change nodes"},
		[]string{"define-project", "Define new project"},
		[]string{"delete-project", "Delete existing closed project"},
		[]string{"import-project", "Import project from existing configuration"},
		[]string{"export-project", "Export existing project configuration"},
	)
	StartInfra.Options = append(StartInfra.Options,
		[]string{"name", " <project name>", "Infrastructure project name"},
	)

	StopInfra.Options = append(StopInfra.Options,
		[]string{"name", " <project name>", "Infrastructure project name"},
	)

	RestartInfra.Options = append(RestartInfra.Options,
		[]string{"name", " <project name>", "Infrastructure project name"},
	)

	DestroyInfra.Options = append(DestroyInfra.Options,
		[]string{"name", " <project name>", "Infrastructure project name"},
	)

	ListInfra.Options = append(ListInfra.Options,
		[]string{"name", " <project name>", "Infrastructure project name"},
	)

}

func ParseArgumentHelper() []CommandHelper {
	/*
	Order :
			NoCommand,
			StartInfrastructure,
			StopInfrastructure,
			RestartInfrastructure,
			DestroyInfrastructure,
			ListInfrastructure,
			ListInfrastructures,
			ListConfigs,
			StatusConfig,
			AlterConfig,
			DefineConfig,
			DeleteConfig,
			ImportConfig,
			ExportConfig,
	*/
	return  []CommandHelper{
		HelpCommand,
		StartInfra,
		StopInfra,
		RestartInfra,
		DestroyInfra,
		ListInfra,
		ListAllInfras,
	}
}