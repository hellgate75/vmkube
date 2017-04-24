package action

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
	Command							string                `json:"Command" xml:"Command" mandatory:"yes" descr:"Required Command" type:"text"`
	Name      					string                `json:"Name" xml:"Name" mandatory:"yes" descr:"Command Name" type:"text"`
	Description					string                `json:"Description" xml:"Description" mandatory:"yes" descr:"Required Command" type:"text"`
	CmdType							CmdRequestType        `json:"CmdType" xml:"CmdType" mandatory:"yes" descr:"Required Command" type:"integer"`
	LineHelp						string                `json:"LineHelp" xml:"LineHelp" mandatory:"yes" descr:"Required Command" type:"text"`
	SubCommands					[]SubCommandHelper    `json:"SubCommands" xml:"SubCommands" mandatory:"yes" descr:"Required Command" type:"list of SubCommandHelper objects"`
	SubCmdTypes					[]CmdSubRequestType   `json:"SubCmdTypes" xml:"SubCmdTypes" mandatory:"yes" descr:"Required Command" type:"integer list"`
	SubCmdHelperTypes		[]CmdRequestType      `json:"SubCmdHelperTypes" xml:"SubCmdHelperTypes" mandatory:"yes" descr:"Required Command" type:"integer list"`
	Options							[]HelperOption        `json:"Options" xml:"Options" mandatory:"yes" descr:"Required Command" type:"tlist of HelperOption objects"`
}

type SubCommandHelper struct {
	Command             string                `json:"Command" xml:"Command" mandatory:"yes" descr:"Required Sub-Command" type:"text"`
	Description         string                `json:"Description" xml:"Description" mandatory:"yes" descr:"Sub-Command Description" type:"text"`
}

type HelperOption struct {
	Option							string                `json:"Option" xml:"Option" mandatory:"yes" descr:"Defined Option" type:"text"`
	Type							  string                `json:"Type" xml:"Type" mandatory:"yes" descr:"Defined Option Type Desription" type:"text"`
	Description					string                `json:"Description" xml:"Description" mandatory:"yes" descr:"Defined Option Desription" type:"text"`
	Mandatory					  bool                  `json:"Mandatory" xml:"Mandatory" mandatory:"yes" descr:"Describe a Mandatory option" type:"boolean"`
}

var(
	HelpCommand CommandHelper = CommandHelper{
		Command: "help",
		Name: "Help",
		Description: "Show help tips",
		CmdType: NoCommand,
		LineHelp: "help [COMMAND]",
		SubCommands: []SubCommandHelper{
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
		Options:	[]HelperOption{},
	}
	StartInfra CommandHelper = CommandHelper{
		Command: "start-infra",
		Name: "Start Infrastructure",
		Description: "Start infrastructre if stopped or nothing",
		CmdType: StartInfrastructure,
		LineHelp: "start-infra [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	StopInfra CommandHelper = CommandHelper{
		Command: "stop-infra",
		Name: "Stop Infrastructure",
		Description: "Stop infrastructre if running or nothing",
		CmdType: StopInfrastructure,
		LineHelp: "stop-infra [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	RestartInfra CommandHelper = CommandHelper{
		Command: "restart-infra",
		Name: "Restart Infrastructure",
		Description: "Restart infrastructre",
		CmdType: RestartInfrastructure,
		LineHelp: "restart-infra [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	DestroyInfra CommandHelper = CommandHelper{
		Command: "destroy-infra",
		Name: "Destroy Infrastructure",
		Description: "Destroy a desired infrastructre (No undo available)",
		CmdType: DestroyInfrastructure,
		LineHelp: "destroy-infra [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	BackupInfra CommandHelper = CommandHelper{
		Command: "backup-infra",
		Name: "Backup Infrastructure",
		Description: "Backup a specific Infrastructure to a backup file",
		CmdType: BackupInfrastructure,
		LineHelp: "backup-infra [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	RecoverInfra CommandHelper = CommandHelper{
		Command: "recover-infra",
		Name: "Recover Infrastructure",
		Description: "Recover a specific Infrastructure from a backup file",
		CmdType: RecoverInfrastructure,
		LineHelp: "recover-infra [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	ListInfra CommandHelper = CommandHelper{
		Command: "infra-status",
		Name: "Infrastructure Details",
		Description: "List information about a specific infrastructure",
		CmdType: ListInfrastructure,
		LineHelp: "infra-status [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	ListAllInfras CommandHelper = CommandHelper{
		Command: "status-all",
		Name: "List Infrastructures",
		Description: "List information about all existing infrastructures",
		CmdType: ListInfrastructures,
		LineHelp: "status-all",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	ListProjects CommandHelper = CommandHelper{
		Command: "list-projects",
		Name: "List Projects",
		Description: "List information about all existing projects",
		CmdType: ListConfigs,
		LineHelp: "list-projects",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	ListProject CommandHelper = CommandHelper{
		Command: "project-status",
		Name: "Project Details",
		Description: "List information about a specific project",
		CmdType: StatusConfig,
		LineHelp: "project-status [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	DefineProject CommandHelper = CommandHelper{
		Command: "define-project",
		Name: "Create Project",
		Description: "Define a new project",
		CmdType: DefineConfig,
		LineHelp: "define-project [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	InfoProject CommandHelper = CommandHelper{
		Command: "info-project",
		Name: "Require Project Schemas",
		Description: "Provides information about project elements definition",
		CmdType: InfoConfig,
		LineHelp: "info-project [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	AlterProject CommandHelper = CommandHelper{
		Command: "alter-project",
		Name: "Modify Project",
		Description: "Modify existing project, e.g.: open, close project or add, modify, delete items",
		CmdType: AlterConfig,
		LineHelp: "alter-project [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	DeleteProject CommandHelper = CommandHelper{
		Command: "delete-project",
		Name: "Delete Project",
		Description: "Delete an existing project",
		CmdType: DeleteConfig,
		LineHelp: "delete-project [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	BuildProject CommandHelper = CommandHelper{
		Command: "build-project",
		Name: "Build Project",
		Description: "Build and existing project and create/modify an infrstructure",
		CmdType: BuildConfig,
		LineHelp: "build-project [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	ImportProject CommandHelper = CommandHelper{
		Command: "import-project",
		Name: "Import Project",
		Description: "Import a new project from file",
		CmdType: ImportConfig,
		LineHelp: "import-project [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
	ExportProject CommandHelper = CommandHelper{
		Command: "export-project",
		Name: "Export Project",
		Description: "Export an existing project to file",
		CmdType: ExportConfig,
		LineHelp: "export-project [OPTIONS]",
		SubCommands: []SubCommandHelper{},
		SubCmdTypes: []CmdSubRequestType{},
		SubCmdHelperTypes: []CmdRequestType{},
		Options:	[]HelperOption{},
	}
)

func InitHelpers() {
	//Help
	HelpCommand.SubCommands = append(HelpCommand.SubCommands,
		SubCommandHelper{
			Command: "help",
			Description: "Show generic commands help",
		},
		SubCommandHelper{
			Command: "start-infra",
			Description: "Start an existing Infrastructure",
		},
		SubCommandHelper{
			Command: "stop-infra",
			Description: "Stop a Running Infrastructure",
		},
		SubCommandHelper{
			Command: "restart-infra",
			Description: "Restart a Running Infrastructure",
		},
		SubCommandHelper{
			Command: "destroy-infra",
			Description: "Destroy a specific Infrastructure",
		},
		SubCommandHelper{
			Command: "backup-infra",
			Description: "Backup a specific Infrastructure to a backup file",
		},
		SubCommandHelper{
			Command: "recover-infra",
			Description: "Recover a specific Infrastructure from a backup file",
		},
		SubCommandHelper{
			Command: "infra-status",
			Description: "Require information about a specific Infrastructure",
		},
		SubCommandHelper{
			Command: "status-all",
			Description: "Require list of all Infrastructures",
		},
		SubCommandHelper{
			Command: "list-projects",
			Description: "Require list of all available projects",
		},
		SubCommandHelper{
			Command: "project-status",
			Description: "Require information about a specific projects",
		},
		SubCommandHelper{
			Command: "define-project",
			Description: "Creates a new project",
		},
		SubCommandHelper{
			Command: "alter-project",
			Description: "Modify existing project, e.g.: open, close project or add, modify, delete items",
		},
		SubCommandHelper{
			Command: "info-project",
			Description: "Provides information about project elements definition",
		},
		SubCommandHelper{
			Command: "delete-project",
			Description: "Delete existing project",
		},
		SubCommandHelper{
			Command: "build-project",
			Description: "Build and existing project and create/modify an infrstructure",
		},
		SubCommandHelper{
			Command: "import-projec",
			Description: "Import project from existing configuration",
		},
		SubCommandHelper{
			Command: "export-project",
			Description: "Export existing project configuration",
		},
	)
	//Start Infrastructure
	StartInfra.Options = append(StartInfra.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Infrastructure name",
			Mandatory: true,
		},
	)

	//Stop Infrastructure
	StopInfra.Options = append(StopInfra.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Infrastructure name",
			Mandatory: true,
		},
	)

	//Restart Infrastructure
	RestartInfra.Options = append(RestartInfra.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Infrastructure name",
			Mandatory: true,
		},
	)

	//Destroy Infrastructure
	DestroyInfra.Options = append(DestroyInfra.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Infrastructure name",
			Mandatory: true,
		},
		HelperOption{
			Option: "force",
			Type: "<boolean>",
			Description: "Flag defining to force destroy, no confirmation will be prompted",
			Mandatory: false,
		},
	)
	
	//Backup Infrastructure
	BackupInfra.Options = append(BackupInfra.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Infrastructure name",
			Mandatory: true,
		},
		HelperOption{
			Option: "file",
			Type: "<text>",
			Description: "Full Backup file path, used to define the infrastructure (extension will be changed to .vmkube)",
			Mandatory: true,
		},
	)
	
	//Recover Infrastructure
	RecoverInfra.Options = append(RecoverInfra.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Infrastructure name",
			Mandatory: true,
		},
		HelperOption{
			Option: "file",
			Type: "<text>",
			Description: "Full Recovery file path, used to define the infrastructure (expected extension .vmkube)",
			Mandatory: true,
		},
		HelperOption{
			Option: "override",
			Type: "<boolean>",
			Description: "Flag defining to force override infrastructure if exists or elsewise fails in case of existing one (default: false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "project-name",
			Type: "<text>",
			Description: "Project Name used to assign a project to the recovered infrastructure",
			Mandatory: false,
		},
		HelperOption{
			Option: "force",
			Type: "<boolean>",
			Description: "Define to force a project creation, if it doesn't exist, using the recovered infrastructure",
			Mandatory: false,
		},
	)
	
	//Status Infrastructure
	ListInfra.Options = append(ListInfra.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Infrastructure name",
			Mandatory: true,
		},
	)

	//Status Project
	ListProject.Options = append(ListProject.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Project name",
			Mandatory: true,
		},
		HelperOption{
			Option: "show-full",
			Type: "<boolean>",
			Description: "Show full details of project on screen (default: false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "format",
			Type: "<text: json|xml>",
			Description: "Format used to show details on screen (default: json)",
			Mandatory: false,
		},
	)
	
	//New Project
	DefineProject.Options = append(DefineProject.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "New Project name",
			Mandatory: true,
		},
		HelperOption{
			Option: "file",
			Type: "<text>",
			Description: "Full Input file path, used to define the project",
			Mandatory: false,
		},
		HelperOption{
			Option: "format",
			Type: "<text: json|xml>",
			Description: "Format used to define the project (default: json)",
			Mandatory: false,
		},
		HelperOption{
			Option: "force",
			Type: "<boolean>",
			Description: "Flag used to force define project, overwriting existing and closed one, fails in case of built infrastructures (default: false), no confirmation will be prompted",
			Mandatory: false,
		},
		HelperOption{
			Option: "override-infra",
			Type: "<boolean>",
			Description: "Flag defining to force rebuild project and override infrastructure if exists or elsewise fails in case of built project (default: false)",
			Mandatory: false,
		},
	)
	
	//Build Project
	BuildProject.Options = append(BuildProject.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Project name",
			Mandatory: true,
		},
		HelperOption{
			Option: "override",
			Type: "<boolean>",
			Description: "Flag defining to override existing infrastructure (default: false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "force",
			Type: "<boolean>",
			Description: "Flag defining to force modify infrastructure, no confirmation will be prompted",
			Mandatory: false,
		},
	)
	//Information on Project Definition
	InfoProject.SubCommands = append(InfoProject.SubCommands,
		SubCommandHelper{
			Command: "list",
			Description: "List project elements, available for change commands",
		},
		SubCommandHelper{
			Command: "details",
			Description: "List of fields for a specific element, available for change commands",
		},
	)
	InfoProject.SubCmdTypes = append(InfoProject.SubCmdTypes,
		List,
		Detail,
	)
	
	InfoProject.Options = append(InfoProject.Options,
		HelperOption{
			Option: "elem-type",
			Type: "<text>",
			Description: "Type of entity to require field information (allowed: Server, Cloud-Server, Network, Domain,...)",
			Mandatory: false,
		},
		HelperOption{
			Option: "sample",
			Type: "<text: json|xml>",
			Description: "Print a sample schema for a specified element type",
			Mandatory: false,
		},
	)
	
	//Change Project
	AlterProject.SubCommands = append(AlterProject.SubCommands,
		SubCommandHelper{
			Command: "create",
			Description: "Create a new project item in a project",
		},
		SubCommandHelper{
			Command: "modify",
			Description: "Modify an existing project item in a project",
		},
		SubCommandHelper{
			Command: "delete",
			Description: "Delete an existing project item from a project",
		},
		SubCommandHelper{
			Command: "close",
			Description: "Close a project for deletion or build",
		},
		SubCommandHelper{
			Command: "open",
			Description: "Re-Open a closed project and eventually deactivate the related infrastructure",
		},
	)
	AlterProject.SubCmdTypes = append(AlterProject.SubCmdTypes,
		Create,
		Alter,
		Remove,
		Close,
		Open,
	)
	
	AlterProject.Options = append(AlterProject.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Project name",
			Mandatory: true,
		},
		HelperOption{
			Option: "file",
			Type: "<text>",
			Description: "Full Input file path, used to define the infrastructure element",
			Mandatory: true,
		},
		HelperOption{
			Option: "format",
			Type: "<text: json|xml>",
			Description: "Format used to define the infrastructure element (default: json)",
			Mandatory: false,
		},
		HelperOption{
			Option: "elem-type",
			Type: "<text>",
			Description: "Type of entity to create/modify/delete in the project (allowed: Server, Cloud-Server, Network, Domain,...)",
			Mandatory: true,
		},
		HelperOption{
			Option: "elem-name",
			Type: "<text>",
			Description: "Name of Entity to create/modify in the project",
			Mandatory: false,
		},
		HelperOption{
			Option: "elem-id",
			Type: "<text>",
			Description: "Id of Entity to modify/delete in the project (used in case of multiple elements with same name)",
			Mandatory: false,
		},
		HelperOption{
			Option: "anchor-elem-type",
			Type: "<text>",
			Description: "Type of anchor entity for new element to create in the project (allowed: Server, Cloud-Server, Network, Domain,...)",
			Mandatory: false,
		},
		HelperOption{
			Option: "anchor-elem-name",
			Type: "<text>",
			Description: "Name of anchor entity for new element to create in the project",
			Mandatory: false,
		},
		HelperOption{
			Option: "anchor-elem-id",
			Type: "<text>",
			Description: "Id of anchor entity for new element to create in the project (used in case of multiple elements with same name)",
			Mandatory: false,
		},
		HelperOption{
			Option: "override-infra",
			Type: "<boolean>",
			Description: "Flag defining to force rebuild project and override infrastructure element (default: false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "force",
			Type: "<boolean>",
			Description: "Flag defining to force modify infrastructure element, no confirmation will be prompted",
			Mandatory: false,
		},
		HelperOption{
			Option: "sample",
			Type: "<boolean>",
			Description: "Show a sample input format instead of alter project from file (default: false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "sample-format",
			Type: "<text: json|xml>",
			Description: "Output format for the required sample input format instead of alter project (default: json, valid if full-export = false)",
			Mandatory: false,
		},
	)

	//Delete Project
	DeleteProject.Options = append(DeleteProject.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "Project name",
			Mandatory: true,
		},
		HelperOption{
			Option: "force",
			Type: "<boolean>",
			Description: "Flag defining to force delete, no confirmation will be prompte",
			Mandatory: false,
		},
	)
	
	//Import Project
	ImportProject.Options = append(ImportProject.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "New Project name",
			Mandatory: true,
		},
		HelperOption{
			Option: "file",
			Type: "<text>",
			Description: "Full path for file used to import project",
			Mandatory: true,
		},
		HelperOption{
			Option: "format",
			Type: "<text: json|xml>",
			Description: "Format used to import project (default: json)",
			Mandatory: false,
		},
		HelperOption{
			Option: "full-import",
			Type: "<boolean>",
			Description: "Flag used to describe a full import (default: true), when true element list import will be ignored",
			Mandatory: false,
		},
		HelperOption{
			Option: "elem-type",
			Type: "<text>",
			Description: "Type of entity top level in the import (allowed: Server, Cloud-Server, Network, Domain,... valid if full-export = false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "sample",
			Type: "<boolean>",
			Description: "Show a sample input format instead of import from file (default: false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "sample-format",
			Type: "<text: json|xml>",
			Description: "Output format for the required sample import instead of import (default: json, valid if full-export = false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "override-infra",
			Type: "<boolean>",
			Description: "Flag defining to force rebuild project and override infrastructure element (default: false)",
			Mandatory: false,
		},
		HelperOption{
			Option: "force",
			Type: "<boolean>",
			Description: "Flag used to force import project, overwriting existing and closed one, project goes out of sync in case of built infrastructure (default: false), no confirmation will be prompted",
			Mandatory: false,
		},
	)
	
	//Export Project
	ExportProject.Options = append(ExportProject.Options,
		HelperOption{
			Option: "name",
			Type: "<text>",
			Description: "New Project name",
			Mandatory: true,
		},
		HelperOption{
			Option: "file",
			Type: "<text>",
			Description: "Full path for file to export project",
			Mandatory: true,
		},
		HelperOption{
			Option: "format",
			Type: "<text: json|xml>",
			Description: "Format used to export project (default: json)",
			Mandatory: false,
		},
		HelperOption{
			Option: "full-export",
			Type: "<boolean>",
			Description: "Flag used to describe a full export (default: true)",
			Mandatory: false,
		},
		HelperOption{
			Option: "elem-type",
			Type: "<text>",
			Description: "Type of entity top level in the export (allowed: Server, Cloud-Server, Network, Domain,... valid if full-export = false)",
			Mandatory: false,
		},
	)
	
	
}

func GetArgumentHelpers() []CommandHelper {
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