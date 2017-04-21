package action

import (
	"strings"
	"fmt"
	"os"
	"errors"
	"vmkube/utils"
	"github.com/satori/go.uuid"
	"vmkube/model"
	"vmkube/vmio"
)

func RecoverCommandHelper(helpCommand string) CommandHelper {
	helperCommands := ParseArgumentHelper()
	for _, helper := range helperCommands {
		if helper.Command == strings.ToLower(helpCommand) {
			return  helper
		}
	}
	return  helperCommands[0]
}


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
		//request.CmdElementType = ??
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
				fmt.Fprintf(os.Stdout, "%s        %s\n",  utils.StrPad(option[0], 50), option[1])
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
			fmt.Fprintf(os.Stdout, "%s        %s\n",  utils.StrPad(option[0], 50), option[1])
		}
	}
	if found  {
		if len(helper.Options) > 0  {
			fmt.Fprintln(os.Stdout, "Options:")
		}
		for _,option := range helper.Options {
			validity := "optional"
			if "true" == option[3] {
				validity = "mandatory"
			}
			fmt.Fprintf(os.Stdout, "--%s  %s  %s  %s\n",  utils.StrPad(option[0],15),  utils.StrPad(option[1], 25), utils.StrPad(validity, 10), option[2])
		}
	} else  {
		fmt.Fprintln(os.Stdout, "Unable to complete help support ...")
	}
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
	active = true
	
	InfraId := ""
	InfraName := ""
	Found := false
	NewIndexes := make([]model.ProjectsDescriptor, 0)
	for _,prj := range indexes.Projects {
		if CorrectInput(prj.Id) != project.Id {
			NewIndexes = append(NewIndexes, )
		} else {
			synced = prj.Synced
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
		if CorrectInput(prj.Id) != project.Id {
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
		if CorrectInput(prj.Id) != project.Id {
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
		Active: CurrentIndex.Open,
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
