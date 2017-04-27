package model

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"vmkube/utils"
)

func (element *Instance) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if element.MachineId == "" {
		errorList = append(errorList, errors.New("Unassigned Project Machine Id field"))
	}
	if element.Driver == "" {
		errorList = append(errorList, errors.New("Unassigned Driver field"))
	}
	if element.Cpus == 0 {
		errorList = append(errorList, errors.New("Unassigned Cpu Count field"))
	}
	if element.Memory == 0 {
		errorList = append(errorList, errors.New("Unassigned Memory Size field"))
	}
	if len(element.Disks) == 0 {
		errorList = append(errorList, errors.New("Unassigned Disks field"))
	}
	if element.Hostname == "" {
		errorList = append(errorList, errors.New("Unassigned host name field"))
	}
	if element.OSType == "" {
		errorList = append(errorList, errors.New("Unassigned OS type name field"))
	}
	if element.OSVersion == "" {
		errorList = append(errorList, errors.New("Unassigned OS version field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Instance) Load(file string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *Instance) Import(file string, format string) error {
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

func (element *Instance) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	return nil
}

func (element *Instance) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
}

func (element *LocalMachine) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if element.Driver == "" {
		errorList = append(errorList, errors.New("Unassigned Driver field"))
	}
	// Not mandatory, we can inspect one the assigned one in the machine inspection
	//if element.Cpus == 0 {
	//	errorList = append(errorList, errors.New("Unassigned Cpu Count field"))
	//}
	//if element.Memory == 0 {
	//	errorList = append(errorList, errors.New("Unassigned Memory Size field"))
	//}
	//if element.DiskSize == 0 {
	//	errorList = append(errorList, errors.New("Unassigned Disk Size field"))
	//}
	if element.Hostname == "" {
		errorList = append(errorList, errors.New("Unassigned host name field"))
	}
	if element.OSType == "" {
		errorList = append(errorList, errors.New("Unassigned OS type name field"))
	}
	if element.OSVersion == "" {
		errorList = append(errorList, errors.New("Unassigned OS version field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *LocalMachine) Load(file string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *LocalMachine) Import(file string, format string) error {
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

func (element *LocalMachine) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	return nil
}

func (element *LocalMachine) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
}

