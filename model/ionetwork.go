package model

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"vmkube/utils"
)

func (element *Network) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if len(element.CloudInstances) == 0 && len(element.LocalInstances) == 0 {
		errorList = append(errorList, errors.New("Unassigned Cloud Instances or Instances List fields"))
	}
	for _,instance := range element.LocalInstances {
		errorList = append(errorList, instance.Validate()...)
	}
	for _,instance := range element.CloudInstances {
		errorList = append(errorList, instance.Validate()...)
	}
	for _,installation := range element.Installations {
		errorList = append(errorList, installation.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Network) Load(file string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *Network) Import(file string, format string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	if format != "json" && format != "xml" {
		return  errors.New("Format "+format+" not supported!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	if format == "json" {
		err = json.Unmarshal(byteArray, &element)
	} else  {
		err = xml.Unmarshal(byteArray, &element)
	}
	if err == nil {
		err := element.PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *Network) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	machineMap := make(map[string]string, 0)
	for i := 0; i < len(element.LocalInstances); i++ {
		id := element.LocalInstances[i].Id
		if id == "" {
			id = element.LocalInstances[i].Name
		}
		if id != "" {
			if _,ok := machineMap[id]; ok {
				bytes := []byte(`Duplicate instance Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.LocalInstances[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.LocalInstances[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			machineMap[id] = element.LocalInstances[i].Id
		}
	}
	for i := 0; i < len(element.CloudInstances); i++ {
		id := element.CloudInstances[i].Id
		if id == "" {
			id = element.CloudInstances[i].Name
		}
		if id != "" {
			if _,ok := machineMap[id]; ok {
				bytes := []byte(`Duplicate cloud instance or machine Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.CloudInstances[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.CloudInstances[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			machineMap[id] = element.CloudInstances[i].Id
		}
	}
	for i := 0; i < len(element.Installations); i++ {
		err := element.Installations[i].PostImport()
		if err != nil {
			return err
		}
		oldId := element.Installations[i].InstanceId
		if _,ok := machineMap[oldId]; ! ok || oldId == "" {
			bytes := []byte(`Unable to locate cloud machine or machine Id/Name in installation plan reference in json : `)
			bytes = append(bytes,utils.GetJSONFromObj(element.Installations[i], true)...)
			return errors.New(string(bytes))
		}
		value, _ := machineMap[oldId]
		element.Installations[i].InstanceId = value
	}
	return nil
}

func (element *Network) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
}

func (element *MachineNetwork) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	// Permissive approach for empty projects in development
	//if len(element.CMachines) == 0 && len(element.Machines) == 0 {
	//	errorList = append(errorList, errors.New("Unassigned Cloud Machines or Machines List fields"))
	//}
	for _,machine := range element.LocalMachines {
		errorList = append(errorList, machine.Validate()...)
	}
	for _,machine := range element.CloudMachines {
		errorList = append(errorList, machine.Validate()...)
	}
	for _,plan := range element.Installations {
		errorList = append(errorList, plan.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *MachineNetwork) Load(file string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *MachineNetwork) Import(file string, format string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	if format != "json" && format != "xml" {
		return  errors.New("Format "+format+" not supported!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	if format == "json" {
		err = json.Unmarshal(byteArray, &element)
	} else  {
		err = xml.Unmarshal(byteArray, &element)
	}
	if err == nil {
		if err == nil {
			err := element.PostImport()
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (element *MachineNetwork) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	machineMap := make(map[string]string, 0)
	for i := 0; i < len(element.LocalMachines); i++ {
		id := element.LocalMachines[i].Id
		if id == "" {
			id = element.LocalMachines[i].Name
		}
		if id != "" {
			if _,ok := machineMap[id]; ok {
				bytes := []byte(`Duplicate machine Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.LocalMachines[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.LocalMachines[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			machineMap[id] = element.LocalMachines[i].Id
		}
	}
	for i := 0; i < len(element.CloudMachines); i++ {
		id := element.CloudMachines[i].Id
		if id == "" {
			id = element.CloudMachines[i].Name
		}
		if id != "" {
			if _,ok := machineMap[id]; ok {
				bytes := []byte(`Duplicate cloud machine or machine Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.CloudMachines[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.CloudMachines[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			machineMap[id] = element.CloudMachines[i].Id
		}
	}
	for i := 0; i < len(element.Installations); i++ {
		err := element.Installations[i].PostImport()
		if err != nil {
			return err
		}
		oldId := element.Installations[i].MachineId
		if _,ok := machineMap[oldId]; ! ok || oldId == "" {
			bytes := []byte(`Unable to locate cloud machine or machine Id/Name in installation plan reference in json : `)
			bytes = append(bytes,utils.GetJSONFromObj(element.Installations[i], true)...)
			return errors.New(string(bytes))
		}
		value, _ := machineMap[oldId]
		element.Installations[i].MachineId = value
	}
	return nil
}

func (element *MachineNetwork) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
}

