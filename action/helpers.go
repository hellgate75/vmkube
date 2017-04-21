package action

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
			BackupInfrastructure,
			RecoverInfrastructure,
			ListInfrastructure,
			ListInfrastructures,
			ListConfigs,
			StatusConfig,
			DefineConfig,
			AlterConfig,
			InfoConfig,
			DeleteConfig,
			BuildConfig,
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
	BackupInfra CommandHelper = CommandHelper{
		Command: "backup-infra",
		Description: "Backup a specific Infrastructure to a backup file",
		CmdType: BackupInfrastructure,
		LineHelp: "backup-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	RecoverInfra CommandHelper = CommandHelper{
		Command: "recover-infra",
		Description: "Recover a specific Infrastructure from a backup file",
		CmdType: RecoverInfrastructure,
		LineHelp: "recover-infra [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	ListInfra CommandHelper = CommandHelper{
		Command: "infra-status",
		Description: "List information about a specific infrastructure",
		CmdType: ListInfrastructure,
		LineHelp: "infra-status [OPTIONS]",
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
	ListProjects CommandHelper = CommandHelper{
		Command: "list-projects",
		Description: "List information about all existing projects",
		CmdType: ListConfigs,
		LineHelp: "list-projects",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	ListProject CommandHelper = CommandHelper{
		Command: "project-status",
		Description: "List information about a specific project",
		CmdType: StatusConfig,
		LineHelp: "project-status [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	DefineProject CommandHelper = CommandHelper{
		Command: "define-project",
		Description: "Define a new project",
		CmdType: DefineConfig,
		LineHelp: "define-project [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	InfoProject CommandHelper = CommandHelper{
		Command: "info-project",
		Description: "Provides information about project elements definition",
		CmdType: InfoConfig,
		LineHelp: "info-project [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	AlterProject CommandHelper = CommandHelper{
		Command: "alter-project",
		Description: "Modify existing project, e.g.: open, close project or add, modify, delete items",
		CmdType: AlterConfig,
		LineHelp: "alter-project [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	DeleteProject CommandHelper = CommandHelper{
		Command: "delete-project",
		Description: "Delete an existing project",
		CmdType: DeleteConfig,
		LineHelp: "delete-project [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	BuildProject CommandHelper = CommandHelper{
		Command: "build-project",
		Description: "Build and existing project and create/modify an infrstructure",
		CmdType: BuildConfig,
		LineHelp: "build-project [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	ImportProject CommandHelper = CommandHelper{
		Command: "import-project",
		Description: "Import a new project from file",
		CmdType: ImportConfig,
		LineHelp: "import-project [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
	ExportProject CommandHelper = CommandHelper{
		Command: "export-project",
		Description: "Export an existing project to file",
		CmdType: ExportConfig,
		LineHelp: "export-project [OPTIONS]",
		SubCommands: [][]string{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[][]string{},
	}
)

func InitHelpers() {
	//Help
	HelpCommand.SubCommands = append(HelpCommand.SubCommands,
		[]string{"help", "Show generic commands help"},
		[]string{"start-infra", "Start Ready Infrastructure"},
		[]string{"stop-infra", "Stop Running Infrastructure"},
		[]string{"restart-infra", "Restart Running Infrastructure"},
		[]string{"destroy-infra", "Destroy a specific Infrastructure"},
		[]string{"backup-infra", "Backup a specific Infrastructure to a backup file"},
		[]string{"recover-infra", "Recover a specific Infrastructure from a backup file"},
		[]string{"infra-status", "Require information about a specific Infrastructure"},
		[]string{"status-all", "Require list of all Infrastructures"},
		[]string{"list-projects", "Require list of all available projects"},
		[]string{"project-status", "Require information about a specific projects"},
		[]string{"define-project", "Define new project"},
		[]string{"alter-project", "Modify existing project, e.g.: open, close project or add, modify, delete items"},
		[]string{"info-project", "Provides information about project elements definition"},
		[]string{"delete-project", "Delete existing project"},
		[]string{"build-project", "Build and existing project and create/modify an infrstructure"},
		[]string{"import-project", "Import project from existing configuration"},
		[]string{"export-project", "Export existing project configuration"},
	)
	//Start Infrastructure
	StartInfra.Options = append(StartInfra.Options,
		[]string{"name", " <infrastructure name>", "Infrastructure name", "true"},
	)

	//Stop Infrastructure
	StopInfra.Options = append(StopInfra.Options,
		[]string{"name", " <infrastructure name>", "Infrastructure name", "true"},
	)

	//Restart Infrastructure
	RestartInfra.Options = append(RestartInfra.Options,
		[]string{"name", " <infrastructure name>", "Infrastructure name", "true"},
	)

	//Destroy Infrastructure
	DestroyInfra.Options = append(DestroyInfra.Options,
		[]string{"name", " <infrastructure name>", "Infrastructure name", "true"},
	)

	DestroyInfra.Options = append(DestroyInfra.Options,
		[]string{"force", " <boolean>", "Flag defining to force destroy, no confirmation will be prompted", "false"},
	)
	
	//Backup Infrastructure
	BackupInfra.Options = append(BackupInfra.Options,
		[]string{"name", " <infrastructure name>", "Infrastructure name", "true"},
	)
	
	BackupInfra.Options = append(BackupInfra.Options,
		[]string{"file", " <file path>", "Full Backup file path, used to define the infrastructure (extension will be changed to .vmkube)", "true"},
	)
	
	//Recover Infrastructure
	RecoverInfra.Options = append(RecoverInfra.Options,
		[]string{"name", " <infrastructure name>", "Infrastructure name", "true"},
	)
	
	RecoverInfra.Options = append(RecoverInfra.Options,
		[]string{"file", " <file path>", "Full Recovery file path, used to define the infrastructure (expected extension .vmkube)", "true"},
	)
	
	RecoverInfra.Options = append(RecoverInfra.Options,
		[]string{"override", " <boolean>", "Flag defining to force override infrastructure if exists or elsewise fails in case of existing one (default: false)", "false"},
	)
	
	RecoverInfra.Options = append(RecoverInfra.Options,
		[]string{"project-name", " <project name>", "Project Name used to assign a project to the recovered infrastructure", "false"},
	)
	
	RecoverInfra.Options = append(RecoverInfra.Options,
		[]string{"force", " <boolean>", "Define to force a project creation, if it doesn't exist, using the recovered infrastructure", "false"},
	)
	
	//Status Infrastructure
	ListInfra.Options = append(ListInfra.Options,
		[]string{"name", " <infrastructure name>", "Infrastructure name", "true"},
	)

	//Status Project
	ListProject.Options = append(ListProject.Options,
		[]string{"name", " <project name>", "Project name", "true"},
	)
	
	ListProject.Options = append(ListProject.Options,
		[]string{"show-full", " <boolean>", "Show full details of project on screen (default: false)", "false"},
	)
	
	ListProject.Options = append(ListProject.Options,
		[]string{"format", " <json|xml>", "Format used to show details on screen (default: json)", "false"},
	)
	
	//New Project
	DefineProject.Options = append(DefineProject.Options,
		[]string{"name", " <project name>", "Project name", "true"},
	)

	DefineProject.Options = append(DefineProject.Options,
		[]string{"file", " <file path>", "Full Input file path, used to define the project", "false"},
	)

	DefineProject.Options = append(DefineProject.Options,
		[]string{"format", " <json|xml>", "Format used to define the project (default: json)", "false"},
	)
	
	DefineProject.Options = append(DefineProject.Options,
		[]string{"force", " <boolean>", "Flag used to force define project, overwriting existing and closed one, fails in case of built infrastructures (default: false), no confirmation will be prompted", "false"},
	)
	
	DefineProject.Options = append(DefineProject.Options,
		[]string{"destroy-infra", " <boolean>", "Flag defining to force destroy infrastructure if exists or elsewise fails in case of built project (default: false)", "false"},
	)
	
	//Build Project
	BuildProject.Options = append(BuildProject.Options,
		[]string{"name", " <project name>", "Project name", "true"},
	)

	BuildProject.Options = append(BuildProject.Options,
		[]string{"override", " <boolean>", "Flag defining to override existing infrastructure (default: false)", "false"},
	)

	BuildProject.Options = append(BuildProject.Options,
		[]string{"force", " <boolean>", "Flag defining to force modify infrastructure, no confirmation will be prompted", "false"},
	)
	//Information on Project Definition
	InfoProject.SubCommands = append(InfoProject.SubCommands,
		[]string{"list", "List project elements, available for change commands"},
		[]string{"details", "List of fields for a specific element, available for change commands"},
	)
	InfoProject.SubCmdTypes = append(InfoProject.SubCmdTypes,
		List,
		Detail,
	)
	
	InfoProject.Options = append(InfoProject.Options,
		[]string{"elem-type", " <infra element type>", "Type of entity to require field information (allowed: Server, Cloud-Server, Network, Domain,...)", "false"},
	)
	
	InfoProject.Options = append(InfoProject.Options,
		[]string{"sample", " <json|xml>", "Print a sample schema for a specified element type", "false"},
	)
	
	//Change Project
	AlterProject.SubCommands = append(AlterProject.SubCommands,
		[]string{"create", "Create a project item"},
		[]string{"modify", "Alter a project item"},
		[]string{"delete", "Delete a project item"},
		[]string{"close", "Close a project for deletion or build"},
		[]string{"open", "Re-Open a closed project and eventually deactivate the infrastructure"},
	)
	AlterProject.SubCmdTypes = append(AlterProject.SubCmdTypes,
		Create,
		Alter,
		Remove,
		Close,
		Open,
	)
	
	AlterProject.Options = append(AlterProject.Options,
		[]string{"name", " <project name>", "Project name", "true"},
	)

	AlterProject.Options = append(AlterProject.Options,
		[]string{"file", " <file path>", "Full Input file path, used to define the infrastructure element", "true"},
	)

	AlterProject.Options = append(AlterProject.Options,
		[]string{"format", " <json|xml>", "Format used to define the infrastructure element (default: json)", "true"},
	)

	AlterProject.Options = append(AlterProject.Options,
		[]string{"elem-type", " <infra element type>", "Type of entity to create/modify/delete in the project (allowed: Server, Cloud-Server, Network, Domain,...)", "true"},
	)

	AlterProject.Options = append(AlterProject.Options,
		[]string{"elem-name", " <name pf entity>", "Entity to create/modify in the project", "true"},
	)

	AlterProject.Options = append(AlterProject.Options,
		[]string{"override", " <boolean>", "Flag defining to override existing infrastructure element (default: false)", "false"},
	)

	AlterProject.Options = append(AlterProject.Options,
		[]string{"force", " <boolean>", "Flag defining to force modify infrastructure element, no confirmation will be prompted", "false"},
	)

	//Delete Project
	DeleteProject.Options = append(DeleteProject.Options,
		[]string{"name", " <project name>", "Project name", "true"},
	)

	DeleteProject.Options = append(DeleteProject.Options,
		[]string{"force", " <boolean>", "Flag defining to force delete, no confirmation will be prompted", "false"},
	)
	
	//Import Project
	ImportProject.Options = append(ImportProject.Options,
		[]string{"name", " <project name>", "New project name", "true"},
	)

	ImportProject.Options = append(ImportProject.Options,
		[]string{"file", " <file path>", "Full path for file used to import project", "true"},
	)

	ImportProject.Options = append(ImportProject.Options,
		[]string{"format", " <json|xml>", "Format used to import project (default: json)", "false"},
	)
	
	ImportProject.Options = append(ImportProject.Options,
		[]string{"full-import", " <boolean>", "Flag used to describe a full import (default: true), when true element list import will be ignored", "false"},
	)
	
	ImportProject.Options = append(ImportProject.Options,
		[]string{"elem-type", " <infra element type>", "Type of entity top level in the import (allowed: Server, Cloud-Server, Network, Domain,... valid if full-export = false)", "false"},
	)
	ImportProject.Options = append(ImportProject.Options,
		[]string{"sample", " <boolean>", "Show a sample import output instead of import from file (default: false, valid if full-export = false)", "false"},
	)
	
	ImportProject.Options = append(ImportProject.Options,
		[]string{"sample-format", " <json|xml>", "Output format for the required sample import instead of import (default: json, valid if full-export = false)", "false"},
	)
	
	ImportProject.Options = append(ImportProject.Options,
		[]string{"force", " <boolean>", "Flag used to force import project, overwriting existing and closed one, project goes out of sync in case of built infrastructure (default: false), no confirmation will be prompted", "false"},
	)
	
	
	//Export Project
	ExportProject.Options = append(ExportProject.Options,
		[]string{"name", " <project name>", "Project name", "true"},
	)

	ExportProject.Options = append(ExportProject.Options,
		[]string{"file", " <file path>", "Full path for file to export project", "true"},
	)

	ExportProject.Options = append(ExportProject.Options,
		[]string{"format", " <json|xml>", "Format used to export project (default: json)", "false"},
	)

	ExportProject.Options = append(ExportProject.Options,
		[]string{"full-export", " <boolean>", "Flag used to describe a full export (default: true)", "false"},
	)
	
	ExportProject.Options = append(ExportProject.Options,
		[]string{"elem-type", " <infra element type>", "Type of entity top level in the export (allowed: Server, Cloud-Server, Network, Domain,... valid if full-export = false)", "false"},
	)
	
	
}

func ParseArgumentHelper() []CommandHelper {
	return  []CommandHelper{
		HelpCommand,
		StartInfra,
		StopInfra,
		RestartInfra,
		DestroyInfra,
		BackupInfra,
		RecoverInfra,
		ListInfra,
		ListAllInfras,
		ListProjects,
		ListProject,
		DefineProject,
		AlterProject,
		InfoProject,
		DeleteProject,
		BuildProject,
		ImportProject,
		ExportProject,
	}
}