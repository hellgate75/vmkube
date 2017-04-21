package model

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"vmkube/utils"
)

func (element *ProjectImport) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if len(element.Domains) == 0 {
		errorList = append(errorList, errors.New("Unassigned Domains List fields"))
	}
	for _,network := range element.Domains {
		errorList = append(errorList, network.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ProjectImport) Import(file string, format string) error {
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
		err := element.PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *ProjectImport) PostImport() error {
	element.Id=NewUUIDString()
	for _,domain := range element.Domains {
		domain.PostImport()
	}
	return nil
}

func (element *Project) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if len(element.Domains) == 0 {
		errorList = append(errorList, errors.New("Unassigned Domains List fields"))
	}
	for _,network := range element.Domains {
		errorList = append(errorList, network.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Project) Load(file string) error {
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

func (element *Project) Import(file string, format string) error {
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
		err := element.PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *Project) PostImport() error {
	element.Id=NewUUIDString()
	for _,domain := range element.Domains {
		domain.PostImport()
	}
	return nil
}

func (element *Project) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

func (element *Infrastructure) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if len(element.Domains) == 0 {
		errorList = append(errorList, errors.New("Unassigned Domains List fields"))
	}
	for _,network := range element.Domains {
		errorList = append(errorList, network.Validate()...)
	}
	errorList = append(errorList, element.State.Validate()...)
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Infrastructure) Load(file string) error {
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

func (element *Infrastructure) Import(file string, format string) error {
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
		err := element.PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *Infrastructure) PostImport() error {
	element.Id=NewUUIDString()
	for _,domain := range element.Domains {
		domain.PostImport()
	}
	return nil
}

func (element *Infrastructure) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

