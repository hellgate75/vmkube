package operations

import (
	"vmkube/model"
	"errors"
	"fmt"
	"vmkube/utils"
)

type CmdElementTypeDesc int

const (
	NoElementDesc					CmdElementTypeDesc = iota
	ServerDesc
	CServerDesc
	PlanDesc
	NetworkDesc
	DomainDesc
	ProjectDesc
)

var ElementTypeDescriptors []string = []string{
	"No Element",
	"Local Server",
	"Cloud Server",
	"Installation Plan",
	"Network",
	"Domain",
	"Project",
}
var ElementTypeCodes []CmdElementTypeDesc = []CmdElementTypeDesc{
	NoElementDesc,
	ServerDesc,
	CServerDesc,
	PlanDesc,
	NetworkDesc,
	DomainDesc,
	ProjectDesc,
}

func lookupForDuplicates(project model.Project, code CmdElementTypeDesc, name string, id string) (int, [][]int) {
	var count int = 0
	var indexes [][]int = make([][]int, 0)
	
	for i := 0; i < len(project.Domains); i++ {
		if code == DomainDesc {
			if (id == "" &&  utils.CorrectInput(project.Domains[i].Name) == utils.CorrectInput(name)) ||  project.Domains[i].Id == id {
				indexes = append(indexes, []int{i})
				count ++
			}
		} else {
			for j := 0; j < len(project.Domains[i].Networks); j++ {
				if code == NetworkDesc {
					if (id == "" &&  utils.CorrectInput(project.Domains[i].Networks[j].Name) == utils.CorrectInput(name)) || project.Domains[i].Networks[j].Id == id {
						indexes = append(indexes, []int{i, j})
						count ++
					}
				} else {
					if code == ServerDesc {
						for k := 0; k < len(project.Domains[i].Networks[j].Servers); k++ {
							if (id == "" &&  utils.CorrectInput(project.Domains[i].Networks[j].Servers[k].Name) == utils.CorrectInput(name)) || project.Domains[i].Networks[j].Servers[k].Id == id {
								indexes = append(indexes, []int{i, j, k})
								count ++
							}
						}
					} else if code == CServerDesc {
						for k := 0; k < len(project.Domains[i].Networks[j].CServers); k++ {
							if (id == "" &&  utils.CorrectInput(project.Domains[i].Networks[j].CServers[k].Name) == utils.CorrectInput(name)) || project.Domains[i].Networks[j].CServers[k].Id == id {
								indexes = append(indexes, []int{i, j, k})
								count ++
							}
						}
					} else if code == PlanDesc {
						for k := 0; k < len(project.Domains[i].Networks[j].Installations); k++ {
							if (id == "" &&  utils.CorrectInput(project.Domains[i].Networks[j].Installations[k].Id) == utils.CorrectInput(name)) || project.Domains[i].Networks[j].Installations[k].Id == id {
								indexes = append(indexes, []int{i, j, k})
								count ++
							}
						}
					}
				}
				
			}
		}
	}
	return count, indexes
}
func AddElementToProject(project model.Project, typeElem int, name string, anchorTypeElem int, anchorName string, anchorId string, file string, format string) error {
	var ElementCode CmdElementTypeDesc = ElementTypeCodes[typeElem]
	var AnchorElementCode CmdElementTypeDesc = ElementTypeCodes[anchorTypeElem]
	fmt.Printf("Selected Element : %s\n", ElementTypeDescriptors[typeElem])
	fmt.Printf("Selected Element Code : %d\n", ElementCode)
	
	if ElementCode == ProjectDesc {
		return errors.New("Uable to Add Entire Project in a project!!")
	}
	
	num, indexList := lookupForDuplicates(project, ElementCode, name, "")
	if num > 1 {
		return errors.New(fmt.Sprintf("Element Type '%s' Name '%s' has multiple occurances : %d!!", ElementTypeDescriptors[typeElem], name, num))
	}
	if num == 1 {
		return errors.New(fmt.Sprintf("Element Type '%s' Name %s already existing!!", ElementTypeDescriptors[typeElem], name))
	}
	
	
	num2, indexList2 := lookupForDuplicates(project, AnchorElementCode, anchorName, anchorId)
	
	if num2 > 1 {
		return errors.New(fmt.Sprintf("Anchor Element Type '%s' Name '%s' (Id %s) has multiple occurances : %d!!", ElementTypeDescriptors[anchorTypeElem], anchorName, anchorId, num2))
	}
	if num2 == 0 {
		return errors.New(fmt.Sprintf("Anchor Element Type '%s' Name '%s' (Id %s) has not occurances!!", ElementTypeDescriptors[anchorTypeElem], anchorName, anchorId))
	}
	
	anchorIndexSet := indexList2[0]
	
	fmt.Printf("Discovered Anchor Element(s) Locations : %v\n", anchorIndexSet)
	
	if ElementCode <= AnchorElementCode {
		return errors.New(fmt.Sprintf("Incompatible Anchor Type %s for Type %s!!", ElementTypeDescriptors[anchorTypeElem], ElementTypeDescriptors[typeElem]))
	}
	
	
	fmt.Printf("Discovered Element(s) Locations : %v\n", indexList)
	
	switch ElementCode {
	case DomainDesc:
		var domain model.ProjectDomain
		err := domain.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	case NetworkDesc:
		var network model.ProjectNetwork
		err := network.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	case ServerDesc:
		var server model.ProjectServer
		err := server.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	case CServerDesc:
		var server model.ProjectCloudServer
		err := server.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	default:
		//Plan
		var plan model.InstallationPlan
		err := plan.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
	}
	return errors.New("Request Not Implemented!!")
}

func AlterElementInProject(project model.Project, typeElem int, name string, id string, file string, format string) error {
	var ElementCode CmdElementTypeDesc = ElementTypeCodes[typeElem]
	fmt.Printf("Selected Element : %s\n", ElementTypeDescriptors[typeElem])
	fmt.Printf("Selected Element Code : %d\n", ElementCode)
	
	if ElementCode == ProjectDesc {
		return errors.New("Uable to Modify Entire Project in a project!!")
	}

	num, indexList := lookupForDuplicates(project, ElementCode, name, id)
	if num > 1 {
		return errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has multiple occurances : %d!!", ElementTypeDescriptors[typeElem], name, id, num))
	}
	if num == 0 {
		return errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has not occurances!!", ElementTypeDescriptors[typeElem], name, id))
	}
	
	targetIndexSet := indexList[0]
	
	fmt.Printf("Discovered Element(s) Locations : %v\n", targetIndexSet)
	
	switch ElementCode {
	case DomainDesc:
		var domain model.ProjectDomain
		err := domain.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	case NetworkDesc:
		var network model.ProjectNetwork
		err := network.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	case ServerDesc:
		var server model.ProjectServer
		err := server.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	case CServerDesc:
		var server model.ProjectCloudServer
		err := server.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
		break
	default:
	//Plan
		var plan model.InstallationPlan
		err := plan.Import(file, format)
		if err != nil {
			return err
		}
		return errors.New("Not Implemented!!")
	}
	return errors.New("Request Not Implemented!!")
}


func DeleteElementInProject(project model.Project, typeElem int, name string, id string) error {
	var ElementCode CmdElementTypeDesc = ElementTypeCodes[typeElem]
	fmt.Printf("Selected Element : %s\n", ElementTypeDescriptors[typeElem])
	fmt.Printf("Selected Element Code : %d\n", ElementCode)
	
	if ElementCode == ProjectDesc {
		return errors.New("Uable to Delete Entire Project in a project!!")
	}

	num, indexList := lookupForDuplicates(project, ElementCode, name, id)
	if num > 1 {
		return errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has multiple occurances : %d!!", ElementTypeDescriptors[typeElem], name, id, num))
	}
	if num == 0 {
		return errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has not occurances!!", ElementTypeDescriptors[typeElem], name, id))
	}
	
	targetIndexSet := indexList[0]
	
	fmt.Printf("Discovered Element(s) Locations : %v\n", targetIndexSet)
	
	switch ElementCode {
		case DomainDesc:
			return errors.New("Not Implemented!!")
			break
		case NetworkDesc:
			return errors.New("Not Implemented!!")
			break
		case ServerDesc:
			return errors.New("Not Implemented!!")
			break
		case CServerDesc:
			return errors.New("Not Implemented!!")
			break
		default:
			//Plan
			return errors.New("Not Implemented!!")
	}
	return errors.New("Request Not Implemented!!")
}

func OpenProject(project model.Project) (model.Project, error) {
	if project.Open {
		return project, errors.New(fmt.Sprintf("Project named %s is already opened", project.Name))
	}
	project.Open = true
	return project, nil
}

func CloseProject(project model.Project) (model.Project, error) {
	if ! project.Open {
		return project, errors.New(fmt.Sprintf("Project named %s is already closed", project.Name))
	}
	project.Open = false
	return project, nil
}