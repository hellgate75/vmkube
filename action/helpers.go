package action

import "github.com/hellgate75/vmkube/common"

var (
	HelpCommand common.CommandHelper = common.CommandHelper{
		Command:     "help",
		Name:        "Help",
		Description: "Show help tips",
		CmdType:     common.NoCommand,
		LineHelp:    "help [COMMAND]",
		SubCommands: []common.SubCommandHelper{},
		SubCmdTypes: []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{
			common.NoCommand,
			common.StartInfrastructure,
			common.StopInfrastructure,
			common.RestartInfrastructure,
			common.DestroyInfrastructure,
			common.AlterInfrastructure,
			common.BackupInfrastructure,
			common.RecoverInfrastructure,
			common.ListInfrastructure,
			common.ListInfrastructures,
			common.ListConfigs,
			common.StatusConfig,
			common.DefineConfig,
			common.AlterConfig,
			common.InfoConfig,
			common.DeleteConfig,
			common.BuildConfig,
			common.ImportConfig,
			common.ExportConfig,
		},
		Options: []common.HelperOption{},
	}
	StartInfra common.CommandHelper = common.CommandHelper{
		Command:           "start-infra",
		Name:              "Start Infrastructure",
		Description:       "Start infrastructre if stopped or nothing",
		CmdType:           common.StartInfrastructure,
		LineHelp:          "start-infra [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	StopInfra common.CommandHelper = common.CommandHelper{
		Command:           "stop-infra",
		Name:              "Stop Infrastructure",
		Description:       "Stop infrastructre if running or nothing",
		CmdType:           common.StopInfrastructure,
		LineHelp:          "stop-infra [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	RestartInfra common.CommandHelper = common.CommandHelper{
		Command:           "restart-infra",
		Name:              "Restart Infrastructure",
		Description:       "Restart infrastructre",
		CmdType:           common.RestartInfrastructure,
		LineHelp:          "restart-infra [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	DestroyInfra common.CommandHelper = common.CommandHelper{
		Command:           "destroy-infra",
		Name:              "Destroy Infrastructure",
		Description:       "Destroy a desired infrastructre (No undo available)",
		CmdType:           common.DestroyInfrastructure,
		LineHelp:          "destroy-infra [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	AlterInfra common.CommandHelper = common.CommandHelper{
		Command:           "alter-infra",
		Name:              "Alter Infrastructure",
		Description:       "Alter a desired infrastructre (instance start,stop,status,recreate,remove,...)",
		CmdType:           common.AlterInfrastructure,
		LineHelp:          "alter-infra [COMMAND] [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	BackupInfra common.CommandHelper = common.CommandHelper{
		Command:           "backup-infra",
		Name:              "Backup Infrastructure",
		Description:       "Backup a specific Infrastructure to a backup file",
		CmdType:           common.BackupInfrastructure,
		LineHelp:          "backup-infra [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	RecoverInfra common.CommandHelper = common.CommandHelper{
		Command:           "recover-infra",
		Name:              "Recover Infrastructure",
		Description:       "Recover a specific Infrastructure from a backup file",
		CmdType:           common.RecoverInfrastructure,
		LineHelp:          "recover-infra [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	ListInfra common.CommandHelper = common.CommandHelper{
		Command:           "infra-status",
		Name:              "Infrastructure Details",
		Description:       "List information about a specific infrastructure",
		CmdType:           common.ListInfrastructure,
		LineHelp:          "infra-status [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	ListAllInfras common.CommandHelper = common.CommandHelper{
		Command:           "list-all-infra",
		Name:              "List Infrastructures",
		Description:       "List information about all existing infrastructures",
		CmdType:           common.ListInfrastructures,
		LineHelp:          "list-all-infra",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	ListProjects common.CommandHelper = common.CommandHelper{
		Command:           "list-projects",
		Name:              "List Projects",
		Description:       "List information about all existing projects",
		CmdType:           common.ListConfigs,
		LineHelp:          "list-projects",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	ListProject common.CommandHelper = common.CommandHelper{
		Command:           "project-status",
		Name:              "Project Details",
		Description:       "List information about a specific project",
		CmdType:           common.StatusConfig,
		LineHelp:          "project-status [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	DefineProject common.CommandHelper = common.CommandHelper{
		Command:           "define-project",
		Name:              "Create Project",
		Description:       "Define a new project",
		CmdType:           common.DefineConfig,
		LineHelp:          "define-project [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	InfoProject common.CommandHelper = common.CommandHelper{
		Command:           "info-project",
		Name:              "Require Project Schemas",
		Description:       "Provides information about project elements definition",
		CmdType:           common.InfoConfig,
		LineHelp:          "info-project [COMMAND] [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	AlterProject common.CommandHelper = common.CommandHelper{
		Command:           "alter-project",
		Name:              "Modify Project",
		Description:       "Modify existing project, e.g.: open, close project or add, modify, delete items",
		CmdType:           common.AlterConfig,
		LineHelp:          "alter-project [COMMAND] [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	DeleteProject common.CommandHelper = common.CommandHelper{
		Command:           "delete-project",
		Name:              "Delete Project",
		Description:       "Delete an existing project",
		CmdType:           common.DeleteConfig,
		LineHelp:          "delete-project [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	BuildProject common.CommandHelper = common.CommandHelper{
		Command:           "build-project",
		Name:              "Build Project",
		Description:       "Build and existing project and create/modify an infrastructure",
		CmdType:           common.BuildConfig,
		LineHelp:          "build-project [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	ImportProject common.CommandHelper = common.CommandHelper{
		Command:           "import-project",
		Name:              "Import Project",
		Description:       "Import a new project from file",
		CmdType:           common.ImportConfig,
		LineHelp:          "import-project [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
	ExportProject common.CommandHelper = common.CommandHelper{
		Command:           "export-project",
		Name:              "Export Project",
		Description:       "Export an existing project to file",
		CmdType:           common.ExportConfig,
		LineHelp:          "export-project [OPTIONS]",
		SubCommands:       []common.SubCommandHelper{},
		SubCmdTypes:       []common.CmdSubRequestType{},
		SubCmdHelperTypes: []common.CmdRequestType{},
		Options:           []common.HelperOption{},
	}
)

func InitHelpers() {
	//Help
	HelpCommand.SubCommands = append(HelpCommand.SubCommands,
		common.SubCommandHelper{
			Command:     "help",
			Description: "Show generic commands help",
		},
		common.SubCommandHelper{
			Command:     "start-infra",
			Description: "Start an existing Infrastructure",
		},
		common.SubCommandHelper{
			Command:     "stop-infra",
			Description: "Stop a Running Infrastructure",
		},
		common.SubCommandHelper{
			Command:     "restart-infra",
			Description: "Restart a Running Infrastructure",
		},
		common.SubCommandHelper{
			Command:     "destroy-infra",
			Description: "Destroy a specific Infrastructure",
		},
		common.SubCommandHelper{
			Command:     "alter-infra",
			Description: "Alter a desired infrastructre (instance start,stop,status,recreate,remove,...)",
		},
		common.SubCommandHelper{
			Command:     "backup-infra",
			Description: "Backup a specific Infrastructure to a backup file",
		},
		common.SubCommandHelper{
			Command:     "recover-infra",
			Description: "Recover a specific Infrastructure from a backup file",
		},
		common.SubCommandHelper{
			Command:     "infra-status",
			Description: "Get information about a specific Infrastructure",
		},
		common.SubCommandHelper{
			Command:     "list-all-infra",
			Description: "Get list of all Infrastructures",
		},
		common.SubCommandHelper{
			Command:     "list-projects",
			Description: "Get list of all available projects",
		},
		common.SubCommandHelper{
			Command:     "project-status",
			Description: "Get information about a specific projects",
		},
		common.SubCommandHelper{
			Command:     "define-project",
			Description: "Creates a new project",
		},
		common.SubCommandHelper{
			Command:     "alter-project",
			Description: "Modify existing project, e.g.: open, close project or add, modify, delete items",
		},
		common.SubCommandHelper{
			Command:     "info-project",
			Description: "Provides information about project elements definition",
		},
		common.SubCommandHelper{
			Command:     "delete-project",
			Description: "Delete existing project",
		},
		common.SubCommandHelper{
			Command:     "build-project",
			Description: "Build and existing project and create/modify an infrastructure",
		},
		common.SubCommandHelper{
			Command:     "import-project",
			Description: "Import project from existing configuration",
		},
		common.SubCommandHelper{
			Command:     "export-project",
			Description: "Export existing project configuration",
		},
	)
	//Start Infrastructure
	StartInfra.Options = append(StartInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force the command, no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "threads",
			Type:        "<integer>",
			Description: "Number of parallel threads used in Infrastructure command (default: 1)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "overclock",
			Type:        "<boolean>",
			Description: "Ignore the capping form the maximum available processors (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Stop Infrastructure
	StopInfra.Options = append(StopInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force the command, no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "threads",
			Type:        "<integer>",
			Description: "Number of parallel threads used in Infrastructure command (default: 1)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "overclock",
			Type:        "<boolean>",
			Description: "Ignore the capping form the maximum available processors (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Restart Infrastructure
	RestartInfra.Options = append(RestartInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force the command, no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "threads",
			Type:        "<integer>",
			Description: "Number of parallel threads used in Infrastructure command (default: 1)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "overclock",
			Type:        "<boolean>",
			Description: "Ignore the capping form the maximum available processors (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Destroy Infrastructure
	DestroyInfra.Options = append(DestroyInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force delete, no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "threads",
			Type:        "<integer>",
			Description: "Number of parallel threads used in Infrastructure delete (default: 1)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "overclock",
			Type:        "<boolean>",
			Description: "Ignore the capping form the maximum available processors (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Alter Infrastructure
	AlterInfra.SubCommands = append(AlterInfra.SubCommands,
		common.SubCommandHelper{
			Command:     "status",
			Description: "Display informations about an instance part of an infrastructure",
		},
		common.SubCommandHelper{
			Command:     "start",
			Description: "Start an instance part of an infrastructure",
		},
		common.SubCommandHelper{
			Command:     "stop",
			Description: "Stop an instance part of an infrastructure",
		},
		common.SubCommandHelper{
			Command:     "restart",
			Description: "Restart an instance part of an infrastructure",
		},
		common.SubCommandHelper{
			Command:     "disable",
			Description: "Disable an instance part of an infrastructure and no actions available as group",
		},
		common.SubCommandHelper{
			Command:     "enable",
			Description: "Enable a disabled instance part of an infrastructure and no actions available as group",
		},
		common.SubCommandHelper{
			Command:     "recreate",
			Description: "Recreate an instance part of an infrastructure",
		},
		common.SubCommandHelper{
			Command:     "remove",
			Description: "Destory and remove an instance from own infrastructure and the original project",
		},
		common.SubCommandHelper{
			Command:     "autofix",
			Description: "Start fixing issues for instances part of an infrastructure",
		},
	)

	AlterInfra.SubCmdTypes = append(AlterInfra.SubCmdTypes,
		common.Status,
		common.Start,
		common.Stop,
		common.Restart,
		common.Disable,
		common.Enable,
		common.Recreate,
		common.Destroy,
		common.AutoFix,
	)

	AlterInfra.Options = append(AlterInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "instance-id",
			Type:        "<text>",
			Description: "Instance Unique identifier information, used to recover the instance to alter (allowed: Instance Id, Cloud Instance Id)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "instance-name",
			Type:        "<text>",
			Description: "Instance Name information, valid only only if it is unique in the whole infrastruture, used to recover the instance to alter (allowed: Name, Cloud Name)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "is-cloud",
			Type:        "<boolean>",
			Description: "Flag defining the instance is local or is on the cloud, useful to find and instance by name (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force alter the instance or whole instances, no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Backup Infrastructure
	BackupInfra.Options = append(BackupInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "file",
			Type:        "<text>",
			Description: "Full Backup file path, used to define the infrastructure (extension will be changed to .vmkube)",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Recover Infrastructure
	RecoverInfra.Options = append(RecoverInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "file",
			Type:        "<text>",
			Description: "Full Recovery file path, used to define the infrastructure (expected extension .vmkube)",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "override",
			Type:        "<boolean>",
			Description: "Flag defining to force override infrastructure if exists or elsewise fails in case of existing one (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "project-name",
			Type:        "<text>",
			Description: "Project Name used to assign a project to the recovered infrastructure",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Define to force a project creation / infrastructure assignment, removing any previous build,index,segment",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "threads",
			Type:        "<integer>",
			Description: "Number of parallel threads used in Infrastructure build (default: 1)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "overclock",
			Type:        "<boolean>",
			Description: "Ignore the capping form the maximum available processors (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup",
			Type:        "<boolean>",
			Description: "Allow to check need for backup of being removed elements (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup-dir",
			Type:        "<string>",
			Description: "Folder used to store backup files (default: '')",
			Mandatory:   false,
		},
	)

	ListAllInfras.Options = append(ListAllInfras.Options,
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Status Infrastructure
	ListInfra.Options = append(ListInfra.Options,
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Infrastructure name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "show-full",
			Type:        "<boolean>",
			Description: "Show full details of infrastructures on screen (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "format",
			Type:        "<text: json|xml|yaml>",
			Description: "Format used to show details on screen (default: json)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Status Project
	ListProject.Options = append(ListProject.Options,
		common.HelperOption{
			Option:      "name",
			Type:        "<text>",
			Description: "Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "show-full",
			Type:        "<boolean>",
			Description: "Show full details of project on screen (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "format",
			Type:        "<text: json|xml|yaml>",
			Description: "Format used to show details on screen (default: json)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//New Project
	DefineProject.Options = append(DefineProject.Options,
		common.HelperOption{
			Option:      "name",
			Type:        "<text>",
			Description: "New Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "file",
			Type:        "<text>",
			Description: "Full Input file path, used to define the project",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "format",
			Type:        "<text: json|xml|yaml>",
			Description: "Format used to define the project (default: json)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag used to force define project, overwriting existing and closed one, fails in case of built infrastructures (default: false), no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "override-infra",
			Type:        "<boolean>",
			Description: "Flag defining to force rebuild project and override infrastructure if exists or elsewise fails in case of built project (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup",
			Type:        "<boolean>",
			Description: "Allow to check need for backup of being removed elements (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup-dir",
			Type:        "<string>",
			Description: "Folder used to store backup files (default: '')",
			Mandatory:   false,
		},
	)

	//Build Project
	BuildProject.Options = append(BuildProject.Options,
		common.HelperOption{
			Option:      "name",
			Type:        "<text>",
			Description: "Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "infra-name",
			Type:        "<text>",
			Description: "Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force modify infrastructure (default: false), no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "rebuild",
			Type:        "<boolean>",
			Description: "Flag defining to rebuild an existing infrastructure (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "threads",
			Type:        "<integer>",
			Description: "Number of parallel threads used in Infrastructure build (default: 1)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "overclock",
			Type:        "<boolean>",
			Description: "Ignore the capping form the maximum available processors (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup",
			Type:        "<boolean>",
			Description: "Allow to check need for backup of being removed elements (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup-dir",
			Type:        "<string>",
			Description: "Folder used to store backup files (default: '')",
			Mandatory:   false,
		},
	)
	//Information on Project Definition
	InfoProject.SubCommands = append(InfoProject.SubCommands,
		common.SubCommandHelper{
			Command:     "list",
			Description: "List project elements, available for change commands",
		},
		common.SubCommandHelper{
			Command:     "details",
			Description: "List of fields for a specific element, available for change commands",
		},
	)
	InfoProject.SubCmdTypes = append(InfoProject.SubCmdTypes,
		common.List,
		common.Detail,
	)

	InfoProject.Options = append(InfoProject.Options,
		common.HelperOption{
			Option:      "elem-type",
			Type:        "<text>",
			Description: "Type of entity to require field information (allowed: Machine, Local-Machine, Cloud-Machine, Network, Domain,...)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "sample",
			Type:        "<text: json|xml|yaml>",
			Description: "Print a sample schema for a specified element type",
			Mandatory:   false,
		},
	)

	//Change Project
	AlterProject.SubCommands = append(AlterProject.SubCommands,
		common.SubCommandHelper{
			Command:     "create",
			Description: "Create a new project item in a project",
		},
		common.SubCommandHelper{
			Command:     "modify",
			Description: "Modify an existing project item in a project",
		},
		common.SubCommandHelper{
			Command:     "delete",
			Description: "Delete an existing project item from a project",
		},
		common.SubCommandHelper{
			Command:     "close",
			Description: "Close a project for deletion or build",
		},
		common.SubCommandHelper{
			Command:     "open",
			Description: "Re-Open a closed project and eventually deactivate the related infrastructure",
		},
	)
	AlterProject.SubCmdTypes = append(AlterProject.SubCmdTypes,
		common.Create,
		common.Alter,
		common.Remove,
		common.Close,
		common.Open,
	)

	AlterProject.Options = append(AlterProject.Options,
		common.HelperOption{
			Option:      "name",
			Type:        "<text>",
			Description: "Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "file",
			Type:        "<text>",
			Description: "Full Input file path, used to define the infrastructure element",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "format",
			Type:        "<text: json|xml|yaml>",
			Description: "Format used to define the infrastructure element (default: json)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "elem-type",
			Type:        "<text>",
			Description: "Type of entity to create/modify/delete in the project (allowed: cs, Cloud-Machine, Network, Domain,...)",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "elem-name",
			Type:        "<text>",
			Description: "Name of Entity to create/modify in the project",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "elem-id",
			Type:        "<text>",
			Description: "Id of Entity to modify/delete in the project (used in case of multiple elements with same name)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "anchor-elem-type",
			Type:        "<text>",
			Description: "Type of anchor entity for new element to create in the project (allowed: Local-Machine, Cloud-Machine, Network, Domain,...)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "anchor-elem-name",
			Type:        "<text>",
			Description: "Name of anchor entity for new element to create in the project",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "anchor-elem-id",
			Type:        "<text>",
			Description: "Id of anchor entity for new element to create in the project (used in case of multiple elements with same name)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "override-infra",
			Type:        "<boolean>",
			Description: "Flag defining to force rebuild project and override infrastructure element (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force modify infrastructure element, no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "sample",
			Type:        "<boolean>",
			Description: "Show a sample input format instead of alter project from file (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "sample-format",
			Type:        "<text: json|xml|yaml>",
			Description: "Output format for the required sample input format instead of alter project (default: json, valid if full-export = false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup",
			Type:        "<boolean>",
			Description: "Allow to check need for backup of being removed elements (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup-dir",
			Type:        "<string>",
			Description: "Folder used to store backup files (default: '')",
			Mandatory:   false,
		},
	)

	//Delete Project
	DeleteProject.Options = append(DeleteProject.Options,
		common.HelperOption{
			Option:      "name",
			Type:        "<text>",
			Description: "Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag defining to force delete, no confirmation will be prompte",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

	//Import Project
	ImportProject.Options = append(ImportProject.Options,
		common.HelperOption{
			Option:      "name",
			Type:        "<text>",
			Description: "New Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "file",
			Type:        "<text>",
			Description: "Full path for file used to import project",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "format",
			Type:        "<text: json|xml|yaml>",
			Description: "Format used to import project (default: json)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "full-import",
			Type:        "<boolean>",
			Description: "Flag used to describe a full import (default: true), when true element list import will be ignored",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "elem-type",
			Type:        "<text>",
			Description: "Type of entity top level in the import (allowed: Local-Machine, Cloud-Machine, Network, Domain,... valid if full-export = false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "sample",
			Type:        "<boolean>",
			Description: "Show a sample input format instead of import from file (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "sample-format",
			Type:        "<text: json|xml|yaml>",
			Description: "Output format for the required sample import instead of import (default: json, valid if full-export = false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "override-infra",
			Type:        "<boolean>",
			Description: "Flag defining to force rebuild project and override infrastructure element (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "force",
			Type:        "<boolean>",
			Description: "Flag used to force import project, overwriting existing and closed one, project goes out of sync in case of built infrastructure (default: false), no confirmation will be prompted",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup",
			Type:        "<boolean>",
			Description: "Allow to check need for backup of being removed elements (default: false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "backup-dir",
			Type:        "<string>",
			Description: "Folder used to store backup files (default: '')",
			Mandatory:   false,
		},
	)

	//Export Project
	ExportProject.Options = append(ExportProject.Options,
		common.HelperOption{
			Option:      "name",
			Type:        "<text>",
			Description: "New Project name",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "file",
			Type:        "<text>",
			Description: "Full path for file to export project",
			Mandatory:   true,
		},
		common.HelperOption{
			Option:      "format",
			Type:        "<text: json|xml|yaml>",
			Description: "Format used to export project (default: json)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "full-export",
			Type:        "<boolean>",
			Description: "Flag used to describe a full export (default: true)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "elem-type",
			Type:        "<text>",
			Description: "Type of entity top level in the export (allowed: Local-Machine, Cloud-Machine, Network, Domain,... valid if full-export = false)",
			Mandatory:   false,
		},
		common.HelperOption{
			Option:      "no-colors",
			Type:        "<boolean>",
			Description: "Prevent to print a colorful output, useful for piping results to a file (default: false)",
			Mandatory:   false,
		},
	)

}

func GetArgumentHelpers() []common.CommandHelper {
	return []common.CommandHelper{
		HelpCommand,
		StartInfra,
		StopInfra,
		RestartInfra,
		DestroyInfra,
		AlterInfra,
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
