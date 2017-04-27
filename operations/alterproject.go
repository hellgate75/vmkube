package operations

import (
	"vmkube/model"
	"errors"
	"fmt"
	"vmkube/utils"
	"vmkube/vmio"
)

type CmdElementTypeDesc int

const (
	NoElementDesc					CmdElementTypeDesc = iota
	MachineDesc
	CMachineDesc
	PlanDesc
	NetworkDesc
	DomainDesc
	ProjectDesc
)

var ElementTypeDescriptors []string = []string{
	"No Element",
	"Local Machine",
	"Cloud Machine",
	"Installation Plan",
	"Network",
	"Domain",
	"Project",
}
var ElementTypeCodes []CmdElementTypeDesc = []CmdElementTypeDesc{
	NoElementDesc,
	MachineDesc,
	CMachineDesc,
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
					if code == MachineDesc {
						for k := 0; k < len(project.Domains[i].Networks[j].LocalMachines); k++ {
							if (id == "" &&  utils.CorrectInput(project.Domains[i].Networks[j].LocalMachines[k].Name) == utils.CorrectInput(name)) || project.Domains[i].Networks[j].LocalMachines[k].Id == id {
								indexes = append(indexes, []int{i, j, k})
								count ++
							}
						}
					} else if code == CMachineDesc {
						for k := 0; k < len(project.Domains[i].Networks[j].CloudMachines); k++ {
							if (id == "" &&  utils.CorrectInput(project.Domains[i].Networks[j].CloudMachines[k].Name) == utils.CorrectInput(name)) || project.Domains[i].Networks[j].CloudMachines[k].Id == id {
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

func AddElementToProject(project model.Project, typeElem int, name string, anchorTypeElem int, anchorName string, anchorId string, file string, format string) (string, error) {
	var ElementCode CmdElementTypeDesc = ElementTypeCodes[typeElem]
	var AnchorElementCode CmdElementTypeDesc = ElementTypeCodes[anchorTypeElem]
	
	if ElementCode == ProjectDesc {
		return "", errors.New("Uable to Add Entire Project in a project!!")
	}
	
	num, _ := lookupForDuplicates(project, ElementCode, name, "")
	if num > 1 {
		return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' has multiple occurances : %d!!", ElementTypeDescriptors[typeElem], name, num))
	}
	if num > 0 {
		return "", errors.New(fmt.Sprintf("Element Type '%s' Name %s already existing!!", ElementTypeDescriptors[typeElem], name))
	}
	var anchorIndexSet []int = make([]int, 0)
	
	if AnchorElementCode == ProjectDesc {
		if utils.CorrectInput(anchorName) != utils.CorrectInput(project.Name) && anchorId != project.Id {
			return "", errors.New(fmt.Sprintf("Project Name : '%s' or Id '%s' doesn't math with selected project!!", anchorName, anchorId))
		}
	} else {
		num2, indexList := lookupForDuplicates(project, AnchorElementCode, anchorName, anchorId)
		
		if num2 > 1 {
			return "", errors.New(fmt.Sprintf("Anchor Element Type '%s' Name '%s' (Id %s) has multiple occurances : %d!!", ElementTypeDescriptors[anchorTypeElem], anchorName, anchorId, num2))
		}
		if num2 == 0 {
			return "", errors.New(fmt.Sprintf("Anchor Element Type '%s' Name '%s' (Id %s) has not occurances!!", ElementTypeDescriptors[anchorTypeElem], anchorName, anchorId))
		}
		
		anchorIndexSet = indexList[0]
	}
	if int(ElementCode) >= int(AnchorElementCode) {
		return "", errors.New(fmt.Sprintf("Incompatible Anchor Type '%s' for Type '%s'!!", ElementTypeDescriptors[anchorTypeElem], ElementTypeDescriptors[typeElem]))
	}
	
	switch ElementCode {
	case DomainDesc:
		var domain model.MachineDomain
		err := domain.Import(file, format)
		if err != nil {
			return "", err
		}
		domain.Name = name
		errorsList :=  domain.Validate()
		if len(errorsList) == 0 {
			project.Domains = append(project.Domains, domain)
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Domain from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return domain.Id+"|"+domain.Name, nil
	case NetworkDesc:
		var network model.MachineNetwork
		err := network.Import(file, format)
		if err != nil {
			return "", err
		}
		network.Name = name
		errorsList :=  network.Validate()
		if len(errorsList) == 0 {
			project.Domains[anchorIndexSet[0]].Networks = append(project.Domains[anchorIndexSet[0]].Networks, network)
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Network from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return network.Id+"|"+network.Name, nil
	case MachineDesc:
		var machine model.LocalMachine
		err := machine.Import(file, format)
		if err != nil {
			return "", err
		}
		machine.Name = name
		errorsList :=  machine.Validate()
		if len(errorsList) == 0 {
			project.Domains[anchorIndexSet[0]].Networks[anchorIndexSet[1]].LocalMachines = append(project.Domains[anchorIndexSet[0]].Networks[anchorIndexSet[1]].LocalMachines, machine)
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Machine from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return machine.Id+"|"+machine.Name, nil
	case CMachineDesc:
		var machine model.CloudMachine
		err := machine.Import(file, format)
		if err != nil {
			return "", err
		}
		machine.Name = name
		errorsList :=  machine.Validate()
		if len(errorsList) == 0 {
			project.Domains[anchorIndexSet[0]].Networks[anchorIndexSet[1]].CloudMachines = append(project.Domains[anchorIndexSet[0]].Networks[anchorIndexSet[1]].CloudMachines, machine)
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Cloud Machine from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return machine.Id+"|"+machine.Name, nil
	default:
		//Plan
		var plan model.InstallationPlan
		err := plan.Import(file, format)
		if err != nil {
			return "", err
		}
		errorsList :=  plan.Validate()
		SubElementCode := MachineDesc
		if plan.IsCloud {
			SubElementCode = CMachineDesc
		}
		numOcc, _ := lookupForDuplicates(project, SubElementCode, "", plan.MachineId)
		if numOcc == 0 {
			cloud := "no"
			if plan.IsCloud {
				cloud = "yes"
			}
			numOcc, machineIndexes := lookupForDuplicates(project, SubElementCode, plan.MachineId, "")
			if numOcc == 1 {
				//First Occurrence of Machine or Cloud Machine is the Plan reference, we recover and replace the Id in the plan
				if plan.IsCloud {
					plan.MachineId = project.Domains[machineIndexes[0][0]].Networks[machineIndexes[0][1]].CloudMachines[machineIndexes[0][2]].Id
				} else {
					plan.MachineId = project.Domains[machineIndexes[0][0]].Networks[machineIndexes[0][1]].LocalMachines[machineIndexes[0][2]].Id
				}
				fmt.Printf("Plan Machine Id auto-discovery : %s\n", plan.MachineId)
			} else if numOcc > 1 {
				return "", errors.New(fmt.Sprintf("Machine Cloud: '%s' Id/Name : '%s' multiple occurrences : %d found in project", cloud, plan.MachineId, numOcc))
			} else {
				return "", errors.New(fmt.Sprintf("Machine Cloud: '%s' Id/Name : '%s' not found in project", cloud, plan.MachineId))
			}
		}
		if len(errorsList) == 0 {
			project.Domains[anchorIndexSet[0]].Networks[anchorIndexSet[1]].Installations = append(project.Domains[anchorIndexSet[0]].Networks[anchorIndexSet[1]].Installations, plan)
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Installation Plane from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return plan.Id+":"+plan.MachineId+"|"+fmt.Sprintf("Plan for machine id : %s", plan.MachineId), nil
	}
	return "", errors.New("Request Not Implemented!!")
}

func AlterElementInProject(project model.Project, typeElem int, name string, id string, file string, format string) (string, error) {
	var ElementCode CmdElementTypeDesc = ElementTypeCodes[typeElem]
	
	if ElementCode == ProjectDesc {
		return "", errors.New("Uable to Modify Entire Project in a project!!")
	}

	num, indexList := lookupForDuplicates(project, ElementCode, name, id)
	if num > 1 {
		return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has multiple occurances : %d!!", ElementTypeDescriptors[typeElem], name, id, num))
	}
	if num == 0 {
		return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has not occurances!!", ElementTypeDescriptors[typeElem], name, id))
	}
	
	targetIndexSet := indexList[0]
	
	switch ElementCode {
	case DomainDesc:
		if len(targetIndexSet) < 1 {
			return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
		}
		var domain model.MachineDomain
		err := domain.Import(file, format)
		if err != nil {
			return "", err
		}
		errorsList :=  domain.Validate()
		if len(errorsList) == 0 {
			project.Domains[targetIndexSet[0]] = domain
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Cloud Machine from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return project.Domains[targetIndexSet[0]].Id+"|"+project.Domains[targetIndexSet[0]].Name, nil
	case NetworkDesc:
		if len(targetIndexSet) < 2 {
			return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
		}
		var network model.MachineNetwork
		err := network.Import(file, format)
		if err != nil {
			return "", err
		}
		errorsList :=  network.Validate()
		if len(errorsList) == 0 {
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]] = network
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Cloud Machine from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Id+"|"+
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Name, nil
	case MachineDesc:
		if len(targetIndexSet) < 3 {
			return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
		}
		var machine model.LocalMachine
		err := machine.Import(file, format)
		if err != nil {
			return "", err
		}
		errorsList :=  machine.Validate()
		if len(errorsList) == 0 {
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].LocalMachines[targetIndexSet[2]] = machine
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Cloud Machine from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].LocalMachines[targetIndexSet[2]].Id+"|"+
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].LocalMachines[targetIndexSet[2]].Name, nil
	case CMachineDesc:
		if len(targetIndexSet) < 3 {
			return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
		}
		var machine model.CloudMachine
		err := machine.Import(file, format)
		if err != nil {
			return "",  err
		}
		errorsList :=  machine.Validate()
		if len(errorsList) == 0 {
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].CloudMachines[targetIndexSet[2]] = machine
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Cloud Machine from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].CloudMachines[targetIndexSet[2]].Id+"|"+
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].CloudMachines[targetIndexSet[2]].Name, nil
	default:
		//Plan
		if len(targetIndexSet) < 3 {
			return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
		}
		var plan model.InstallationPlan
		err := plan.Import(file, format)
		if err != nil {
			return "", err
		}
		errorsList :=  plan.Validate()
		if len(errorsList) == 0 {
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Installations[targetIndexSet[2]] = plan
		} else {
			_, desc := vmio.StripErrorMessages(fmt.Sprintf("Errors importing new Cloud Machine from file '%s' format '%s', stack trace : ", file, format), errorsList)
			return "", errors.New(desc)
		}
		return project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Installations[targetIndexSet[2]].Id+":"+
			project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Installations[targetIndexSet[2]].MachineId+"|"+
			fmt.Sprintf("Plan for machine id : %s", project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Installations[targetIndexSet[2]].MachineId), nil
	}
	return "", errors.New("Request Not Implemented!!")
}


func DeleteElementInProject(project model.Project, typeElem int, name string, id string) (string, error) {
	var ElementCode CmdElementTypeDesc = ElementTypeCodes[typeElem]
	
	if ElementCode == ProjectDesc {
		return "", errors.New("Uable to Delete Entire Project in a project!!")
	}

	num, indexList := lookupForDuplicates(project, ElementCode, name, id)
	if num > 1 {
		return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has multiple occurances : %d!!", ElementTypeDescriptors[typeElem], name, id, num))
	}
	if num == 0 {
		return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) has not occurances!!", ElementTypeDescriptors[typeElem], name, id))
	}
	
	targetIndexSet := indexList[0]
	
	switch ElementCode {
		case DomainDesc:
			if len(targetIndexSet) < 1 {
				return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
			}
			Id := project.Domains[targetIndexSet[0]].Id+"|"+project.Domains[targetIndexSet[0]].Name
			dSet := project.Domains[(targetIndexSet[0]+1):]
			project.Domains = project.Domains[:targetIndexSet[0]]
			project.Domains = append(project.Domains, dSet...)
			return Id, nil
		case NetworkDesc:
			if len(targetIndexSet) < 2 {
				return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
			}
			Id := project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Id+"|"+
				project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Name
			nSet := project.Domains[targetIndexSet[0]].Networks[(targetIndexSet[1]+1):]
			project.Domains[targetIndexSet[0]].Networks = project.Domains[targetIndexSet[0]].Networks[:targetIndexSet[1]]
			project.Domains[targetIndexSet[0]].Networks = append(project.Domains[targetIndexSet[0]].Networks, nSet...)
			return Id, nil
		case MachineDesc:
			if len(targetIndexSet) < 3 {
				return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
			}
			Id := project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].LocalMachines[targetIndexSet[2]].Id+"|"+
				project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].LocalMachines[targetIndexSet[2]].Name
			sSet := project.Domains[targetIndexSet[0]].Networks[1].LocalMachines[(targetIndexSet[2]+1):]
			project.Domains[targetIndexSet[0]].Networks[1].LocalMachines = project.Domains[targetIndexSet[0]].Networks[1].LocalMachines[:targetIndexSet[2]]
			project.Domains[targetIndexSet[0]].Networks[1].LocalMachines = append(project.Domains[targetIndexSet[0]].Networks[1].LocalMachines, sSet...)
			return Id, nil
		case CMachineDesc:
			if len(targetIndexSet) < 3 {
				return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
			}
			Id := project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].CloudMachines[targetIndexSet[2]].Id+"|"+
				project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].CloudMachines[targetIndexSet[2]].Name
			csSet := project.Domains[targetIndexSet[0]].Networks[1].CloudMachines[(targetIndexSet[2]+1):]
			project.Domains[targetIndexSet[0]].Networks[1].CloudMachines = project.Domains[targetIndexSet[0]].Networks[1].CloudMachines[:targetIndexSet[2]]
			project.Domains[targetIndexSet[0]].Networks[1].CloudMachines = append(project.Domains[targetIndexSet[0]].Networks[1].CloudMachines, csSet...)
			return Id, nil
		default:
			//Plan
			if len(targetIndexSet) < 3 {
				return "", errors.New(fmt.Sprintf("Element Type '%s' Name '%s' (Id %s) wrong position!!", ElementTypeDescriptors[typeElem], name, id))
			}
			Id := project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Installations[targetIndexSet[2]].Id+":"+
				project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Installations[targetIndexSet[2]].MachineId+"|"+
				fmt.Sprintf("Plan for machine id : %s", project.Domains[targetIndexSet[0]].Networks[targetIndexSet[1]].Installations[targetIndexSet[2]].MachineId)
			pSet := project.Domains[targetIndexSet[0]].Networks[1].Installations[(targetIndexSet[2]+1):]
			project.Domains[targetIndexSet[0]].Networks[1].Installations = project.Domains[targetIndexSet[0]].Networks[1].Installations[:targetIndexSet[2]]
			project.Domains[targetIndexSet[0]].Networks[1].Installations = append(project.Domains[targetIndexSet[0]].Networks[1].Installations, pSet...)
			return Id, nil
	}
	return "", errors.New("Request Not Implemented!!")
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