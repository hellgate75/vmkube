package model

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"vmkube/utils"
)

func (element *ProjectsDescriptor) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Project Unique Identifier field"))
	}
	if element.InfraId == "" {
		errorList = append(errorList, errors.New("Unassigned Infrastructure Unique Identifier field"))
	}
	if element.InfraName == "" {
		errorList = append(errorList, errors.New("Unassigned Infrastructure Name field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Project Name field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ProjectsDescriptor) Load(file string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	err = json.Unmarshal(DecodeBytes(byteArray), &element)
	return err
}

func (element *ProjectsDescriptor) Import(file string, format string) error {
	if ! ExistsFile(file) {
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
	if err == nil && element.Id == "" {
		err = element.PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *ProjectsDescriptor) PostImport() error {
	return nil
}

func (element *ProjectsDescriptor) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

func (element *ProjectsIndex) Validate() []error {
	errorList := make([]error, 0)
	if len(element.Projects) == 0 {
		errorList = append(errorList, errors.New("Unassigned Project Descriptors List fields"))
	}
	for _,index := range element.Projects {
		errorList = append(errorList, index.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ProjectsIndex) Load(file string) error {
	if ! ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	err = json.Unmarshal(DecodeBytes(byteArray), element)
	return err
}

func (element *ProjectsIndex) Import(file string, format string) error {
	if ! ExistsFile(file) {
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
	if err == nil && element.Id == "" {
		err = element.PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *ProjectsIndex) PostImport() error {
	return nil
}

func (element *ProjectsIndex) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

