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
	if len(element.CInstances) == 0 && len(element.Instances) == 0 {
		errorList = append(errorList, errors.New("Unassigned Cloud Instances or Instances List fields"))
	}
	for _,instance := range element.Instances {
		errorList = append(errorList, instance.Validate()...)
	}
	for _,instance := range element.CInstances {
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
	serverMap := make(map[string]string, 0)
	for i := 0; i < len(element.Instances); i++ {
		id := element.Instances[i].Id
		if id == "" {
			id = element.Instances[i].Name
		}
		if id != "" {
			if _,ok := serverMap[id]; ok {
				bytes := []byte(`Duplicate instance Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.Instances[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.Instances[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			serverMap[id] = element.Instances[i].Id
		}
	}
	for i := 0; i < len(element.CInstances); i++ {
		id := element.CInstances[i].Id
		if id == "" {
			id = element.CInstances[i].Name
		}
		if id != "" {
			if _,ok := serverMap[id]; ok {
				bytes := []byte(`Duplicate cloud instance or server Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.CInstances[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.CInstances[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			serverMap[id] = element.CInstances[i].Id
		}
	}
	for i := 0; i < len(element.Installations); i++ {
		err := element.Installations[i].PostImport()
		if err != nil {
			return err
		}
		oldId := element.Installations[i].InstanceId
		if _,ok := serverMap[oldId]; ! ok || oldId == "" {
			bytes := []byte(`Unable to locate cloud server or server Id/Name in installation plan reference in json : `)
			bytes = append(bytes,utils.GetJSONFromObj(element.Installations[i], true)...)
			return errors.New(string(bytes))
		}
		value, _ := serverMap[oldId]
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

func (element *ProjectNetwork) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	// Permissive approach for empty projects in development
	//if len(element.CServers) == 0 && len(element.Servers) == 0 {
	//	errorList = append(errorList, errors.New("Unassigned Cloud Servers or Servers List fields"))
	//}
	for _,server := range element.Servers {
		errorList = append(errorList, server.Validate()...)
	}
	for _,server := range element.CServers {
		errorList = append(errorList, server.Validate()...)
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

func (element *ProjectNetwork) Load(file string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *ProjectNetwork) Import(file string, format string) error {
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

func (element *ProjectNetwork) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	serverMap := make(map[string]string, 0)
	for i := 0; i < len(element.Servers); i++ {
		id := element.Servers[i].Id
		if id == "" {
			id = element.Servers[i].Name
		}
		if id != "" {
			if _,ok := serverMap[id]; ok {
				bytes := []byte(`Duplicate server Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.Servers[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.Servers[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			serverMap[id] = element.Servers[i].Id
		}
	}
	for i := 0; i < len(element.CServers); i++ {
		id := element.CServers[i].Id
		if id == "" {
			id = element.CServers[i].Name
		}
		if id != "" {
			if _,ok := serverMap[id]; ok {
				bytes := []byte(`Duplicate cloud server or server Id/Name reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(element.CServers[i], true)...)
				return errors.New(string(bytes))
			}
		}
		err := element.CServers[i].PostImport()
		if err != nil {
			return err
		}
		if id != "" {
			serverMap[id] = element.CServers[i].Id
		}
	}
	for i := 0; i < len(element.Installations); i++ {
		err := element.Installations[i].PostImport()
		if err != nil {
			return err
		}
		oldId := element.Installations[i].ServerId
		if _,ok := serverMap[oldId]; ! ok || oldId == "" {
			bytes := []byte(`Unable to locate cloud server or server Id/Name in installation plan reference in json : `)
			bytes = append(bytes,utils.GetJSONFromObj(element.Installations[i], true)...)
			return errors.New(string(bytes))
		}
		value, _ := serverMap[oldId]
		element.Installations[i].ServerId = value
	}
	return nil
}

func (element *ProjectNetwork) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
}

