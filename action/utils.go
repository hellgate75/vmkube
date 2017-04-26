package action

import (
	"strings"
	"errors"
	"vmkube/model"
	"vmkube/vmio"
	"fmt"
	"os"
	"vmkube/utils"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
)


func ParseCommandArguments(args	[]string) (*CmdArguments, error) {
	arguments := CmdArguments{}
	success := arguments.Parse(args[1:])
	if success  {
		return  &arguments, nil
	} else  {
		return  &arguments, errors.New("Unable to Parse Command Line")
	}
}

func ParseCommandLine(args []string) (CmdRequest, error) {
	request := CmdRequest{}
	arguments, error := ParseCommandArguments(args)
	if error == nil  {
		request.TypeStr = arguments.Cmd
		request.Type = arguments.CmdType
		request.SubTypeStr = arguments.SubCmd
		request.SubType = arguments.SubCmdType
		request.HelpType = arguments.SubCmdHelpType
		request.Arguments = arguments
	}
	return  request, error
}


func CmdParseElement(value string) (CmdElementType, error) {
		switch CorrectInput(value) {
		case "server":
			return  LServer, nil
		case "cloud-server":
			return  CLServer, nil
		case "network":
			return  SNetwork, nil
		case "domain":
			return  SDomain, nil
		case "project":
			return  SProject, nil
		case "plan":
			return  SPlan, nil
		default:
			return  NoElement, errors.New("Element '"+value+"' is not an infratructure element. Available ones : Server, Cloud-Server, Network, Domain, Plan, Project")

		}
}

func CorrectInput(input string) string {
	return  strings.TrimSpace(strings.ToLower(input))
}

func GetBoolean(input string) bool {
	return CorrectInput(input) == "true"
}

func GetInteger(input string, defaultValue int) int {
	num, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return num
}

func ProjectToInfrastructure(project model.Project) (model.Infrastructure, error) {
	infrastructure := model.Infrastructure{}
	infrastructure.Id = NewUUIDString()
	infrastructure.ProjectId = project.Id
	infrastructure.Name = project.Name
	for _,domain := range project.Domains {
		newDomain := model.Domain{
			Id: NewUUIDString(),
			Name: domain.Name,
			Options: domain.Options,
			Networks: []model.Network{},
		}
		for _,network := range domain.Networks {
			newNetwork := model.Network{
				Id: NewUUIDString(),
				Name: network.Name,
				Options: network.Options,
				Instances: []model.Instance{},
				CInstances: []model.CloudInstance{},
				Installations: []model.Installation{},
			}
			serverConvertionMap := make(map[string]string)
			for _,server := range network.Servers {
				var disks []model.Disk = make([]model.Disk,0)
				disks = append(disks, model.Disk{
					Id: NewUUIDString(),
					Name: "sda0",
					Size: server.DiskSize,
					Type: 0,
				})
				instanceId := NewUUIDString()
				instance := model.Instance{
					Id: instanceId,
					Name: server.Name,
					Options: server.Options,
					Cpus: server.Cpus,
					Memory: server.Memory,
					Disks: disks,
					Driver: server.Driver,
					Engine: model.ToInstanceEngineOpt(server.Engine),
					Swarm: model.ToInstanceSwarmOpt(server.Swarm),
					Hostname: server.Hostname,
					IPAddress: "",
					NoShare: server.NoShare,
					OSType: server.OSType,
					OSVersion: server.OSVersion,
					Roles: server.Roles,
					ServerId: server.Id,
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId: infrastructure.Id,
						ElementId: instanceId,
						LogLines: []string{},
					},
				}
				if _,ok := serverConvertionMap[server.Id]; ok {
					return infrastructure, errors.New("Duplicate Server Id in Project : " + server.Id)
				}
				serverConvertionMap[server.Id] = instance.Id
				newNetwork.Instances = append(newNetwork.Instances, instance)
			}
			for _,server := range network.CServers {
				instanceId := NewUUIDString()
				instance := model.CloudInstance{
					Id: instanceId,
					Name: server.Name,
					Driver: server.Driver,
					Hostname: server.Hostname,
					IPAddress: "",
					Options: server.Options,
					Roles: server.Roles,
					ServerId: server.Id,
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId: infrastructure.Id,
						ElementId: instanceId,
						LogLines: []string{},
					},
				}
				if _,ok := serverConvertionMap[server.Id]; ok {
					return infrastructure, errors.New("Duplicate Server Id in Project : " + server.Id)
				}
				serverConvertionMap[server.Id] = instance.Id
				newNetwork.CInstances = append(newNetwork.CInstances, instance)
			}
			for _,plan := range network.Installations {
				if _,ok := serverConvertionMap[plan.ServerId]; ! ok {
					return infrastructure, errors.New("Invalid server reference in plan : " + plan.ServerId)
				}
				instanceId, _  := serverConvertionMap[plan.ServerId]
				installationId := NewUUIDString()
				installation := model.Installation{
					Id: installationId,
					Environment: model.ToInstanceEnvironment(plan.Environment),
					Role: model.ToInstanceRole(plan.Role),
					Type: model.ToInstanceInstallation(plan.Type),
					Errors: false,
					InstanceId: instanceId,
					IsCloud: plan.IsCloud,
					Success: false,
					LastExecution: time.Now(),
					LastMessage: "",
					Logs: model.LogStorage{
						ProjectId: project.Id,
						InfraId: infrastructure.Id,
						ElementId: installationId,
						LogLines: []string{},
					},
					Plan: plan,
				}
				newNetwork.Installations = append(newNetwork.Installations, installation)
			}
			newDomain.Networks = append(newDomain.Networks, newNetwork)
		}
		infrastructure.Domains = append(infrastructure.Domains, newDomain)
	}
	return infrastructure, nil
}

func UpdateIndexWithProject(project model.Project) error {
	indexes, err := vmio.LoadIndex()
	
	if err != nil {
		return err
	}
	
	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()
	
	indexes, err = vmio.LoadIndex()
	
	if err != nil {
		return err
	}
	
	var synced, active bool
	synced = true
	active = false
	
	InfraId := ""
	InfraName := ""
	Found := false
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _,prj := range indexes.Projects {
		if CorrectInput(prj.Name) != CorrectInput(project.Name) {
			NewIndexes = append(NewIndexes, )
		} else {
			synced = (prj.InfraId == "")
			active = prj.Active
			InfraId = prj.InfraId
			InfraName = prj.InfraName
			Found = true
		}
	}
	
	vmio.LockIndex(indexes)
	
	indexes.Projects = append(indexes.Projects, model.ProjectsDescriptor{
		Id: project.Id,
		Name: project.Name,
		Open: project.Open,
		Synced: synced,
		Active: active,
		InfraId: InfraId,
		InfraName: InfraName,
	})
	
	if Found {
		indexes.Projects = NewIndexes
	}
	
	err = vmio.SaveIndex(indexes)
	
	vmio.UnlockIndex(indexes)
	
	return err
}

func UpdateIndexWithProjectsDescriptor(project model.ProjectsDescriptor, addDescriptor bool) error {
	indexes, err := vmio.LoadIndex()
	
	if err != nil {
		return err
	}
	
	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()
	
	indexes, err = vmio.LoadIndex()
	
	if err != nil {
		return err
	}
	
	Found := false
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _,prj := range indexes.Projects {
		if prj.Id != project.Id {
			NewIndexes = append(NewIndexes, )
		} else {
			Found = true
		}
	}
	
	vmio.LockIndex(indexes)
	
	if addDescriptor {
		NewIndexes = append(NewIndexes, project)
		indexes.Projects = NewIndexes
	} else  if Found {
		indexes.Projects = NewIndexes
	}
	
	err = vmio.SaveIndex(indexes)
	
	vmio.UnlockIndex(indexes)
	
	return err
}

func UpdateIndexWithInfrastructure(infrastructure model.Infrastructure) error {
	indexes, err := vmio.LoadIndex()
	
	if err != nil {
		return err
	}
	
	iFaceIndex := vmio.IFaceIndex{
		Id: indexes.Id,
	}
	iFaceIndex.WaitForUnlock()
	
	indexes, err = vmio.LoadIndex()
	
	if err != nil {
		return err
	}
	
	vmio.LockIndex(indexes)
	
	iFaceProject := vmio.IFaceProject{
		Id: infrastructure.ProjectId,
	}
	iFaceProject.WaitForUnlock()
	
	
	project, err := vmio.LoadProject(infrastructure.ProjectId)
	
	if err != nil {
		return err
	}

	CurrentIndex := model.ProjectsDescriptor{}
	Found := false
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _,prj := range indexes.Projects {
		if prj.Id != project.Id {
			NewIndexes = append(NewIndexes, )
		} else {
			CurrentIndex = prj
			Found = true
		}
	}

	if ! Found {
		return errors.New("Project Id: '"+infrastructure.ProjectId+"' for Infrastrcutre '"+infrastructure.Name+"' not found!!")
	}
	
	indexes.Projects = append(indexes.Projects, model.ProjectsDescriptor{
		Id: project.Id,
		Name: project.Name,
		Open: project.Open,
		Synced: CurrentIndex.Synced,
		Active: CurrentIndex.Active,
		InfraId: infrastructure.Id,
		InfraName: infrastructure.Name,
	})

	UpdateIndex := len(indexes.Projects) > len(NewIndexes)
	
	if UpdateIndex {
		indexes.Projects = NewIndexes
	}
	
	err = vmio.SaveIndex(indexes)
	
	vmio.UnlockIndex(indexes)
	
	return err
}

func CmdParseOption(key string, options []SubCommandHelper) (string, int, error) {
	if len(key) > 0 {
		if strings.Index(key, "--") == 0  {
			return key, -1, errors.New("Invalid Argument (wrong characters: --) : " + key)
		} else	if strings.Index(key, "-") == 0  {
			return key, -1, errors.New("Invalid Argument (wrong character: -) : " + key)
		} else  {
			for index,opts := range options {
				if CorrectInput(key) == opts.Command  {
					return  CorrectInput(key), index, nil
				}
			}
			return  key, -1, errors.New("Invalid Argument : " + key)
		}
	} else  {
		return  key, -1, errors.New("Unable to parse Agument : " + key)
	}
}

func RecoverCommandHelper(helpCommand string) CommandHelper {
	helperCommands := GetArgumentHelpers()
	for _, helper := range helperCommands {
		if helper.Command == strings.ToLower(helpCommand) {
			return  helper
		}
	}
	return  helperCommands[0]
}

func NewUUIDString()	string {
	return  uuid.NewV4().String()
}

func PrintCommandHelper(command	string, subCommand string) {
	helper := RecoverCommandHelper(command)
	fmt.Fprintln(os.Stdout, "Help: vmkube", helper.LineHelp)
	fmt.Fprintln(os.Stdout, "Action:", helper.Description)
	found := false
	if "" !=  strings.TrimSpace(strings.ToLower(subCommand)) && "help" !=  strings.TrimSpace(strings.ToLower(subCommand)) {
		fmt.Fprintln(os.Stdout, "Selected Sub-Command: " + subCommand)
		for _,option := range helper.SubCommands {
			fmt.Fprintf(os.Stdout, "%s        %s\n",  utils.StrPad(option.Command, 50), option.Description)
			found = true
		}
		if ! found  {
			fmt.Fprintln(os.Stdout, "Sub-Command Not found!!")
			if "help" !=  strings.TrimSpace(strings.ToLower(command)) {
				fmt.Fprintln(os.Stdout, "Please type: vmkube","help", command,"for full Sub-Command List")
			} else  {
				fmt.Fprintln(os.Stdout, "Please type: vmkube","help", "COMMAND","for full Sub-Command List")
			}
		}
	}  else {
		found = true
		if len(helper.SubCommands) > 0  {
			if len(helper.SubCmdTypes) > 0 {
				fmt.Fprintln(os.Stdout, "Sub-Commands:")
			} else {
				fmt.Fprintln(os.Stdout, "Commands:")
			}
		}
		for _,option := range helper.SubCommands {
			fmt.Fprintf(os.Stdout, "%s        %s\n",  utils.StrPad(option.Command, 55), option.Description)
		}
	}
	if found  {
		if len(helper.Options) > 0  {
			fmt.Fprintln(os.Stdout, "Options:")
		}
		for _,option := range helper.Options {
			validity := "optional"
			if option.Mandatory {
				validity = "mandatory"
			}
			fmt.Fprintf(os.Stdout, "--%s  %s  %s  %s\n",  utils.StrPad(option.Option,20),  utils.StrPad(option.Type, 25), utils.StrPad(validity, 10), option.Description)
		}
	} else  {
		fmt.Fprintln(os.Stdout, "Unable to complete help support ...")
	}
}
