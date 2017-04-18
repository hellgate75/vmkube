package model

import (
	"errors"
	"io/ioutil"
	"encoding/base64"
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
	if ! existsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	by, err := base64.RawStdEncoding.DecodeString(string(byteArray))
	if err != nil {
		return  err
	}
	err = json.Unmarshal(by, &element)
	return err
}

func (element *Network) Import(file string, format string) error {
	if ! existsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	if format != "json" && format != "xml" {
		return  errors.New("Format "+format+" nor reknown!!")
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
	return err
}

func (element *Network) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	deleteIfExists(file)
	value := base64.RawStdEncoding.EncodeToString(byteArray)
	newBytes := []byte(value)
	err = ioutil.WriteFile(file, newBytes , 0666)
	return  err
}

func (element *ProjectNetwork) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if len(element.CServers) == 0 && len(element.Servers) == 0 {
		errorList = append(errorList, errors.New("Unassigned Cloud Servers or Servers List fields"))
	}
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
	if ! existsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	by, err := base64.RawStdEncoding.DecodeString(string(byteArray))
	if err != nil {
		return  err
	}
	err = json.Unmarshal(by, &element)
	return err
}

func (element *ProjectNetwork) Import(file string, format string) error {
	if ! existsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	if format != "json" && format != "xml" {
		return  errors.New("Format "+format+" nor reknown!!")
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
		element.Id=NewUUIDString()
		serverMap := make(map[string]string, 0)
		for _,server := range element.Servers {
			id := server.Id
			if id == "" {
				id = server.Name
			}
			if id != "" {
				if _,ok := serverMap[id]; ok {
					bytes := []byte(`Duplicate server Id/Name reference in json : `)
					bytes = append(bytes,utils.GetJSONFromObj(server, true)...)
					return errors.New(string(bytes))
				}
			}
			server.Id = NewUUIDString()
			if id != "" {
				serverMap[id] = server.Id
			}
		}
		for _,server := range element.CServers {
			id := server.Id
			if id == "" {
				id = server.Name
			}
			if id != "" {
				if _,ok := serverMap[id]; ok {
					bytes := []byte(`Duplicate cloud server or server Id/Name reference in json : `)
					bytes = append(bytes,utils.GetJSONFromObj(server, true)...)
					return errors.New(string(bytes))
				}
			}
			server.Id = NewUUIDString()
			if id != "" {
				serverMap[id] = server.Id
			}
		}
		for _,installPlan := range element.Installations {
			installPlan.Id = NewUUIDString()
			oldId := installPlan.ServerId
			if _,ok := serverMap[oldId]; ! ok || oldId == "" {
				bytes := []byte(`Unable to locate cloud server or server Id/Name in installation plan reference in json : `)
				bytes = append(bytes,utils.GetJSONFromObj(installPlan, true)...)
				return errors.New(string(bytes))
			}
			value, _ := serverMap[oldId]
			installPlan.ServerId = value
		}
	}
	return err
}

func (element *ProjectNetwork) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	deleteIfExists(file)
	value := base64.RawStdEncoding.EncodeToString(byteArray)
	newBytes := []byte(value)
	err = ioutil.WriteFile(file, newBytes , 0666)
	return  err
}

