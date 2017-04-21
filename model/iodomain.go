package model

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"vmkube/utils"
)

func (element *Domain) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if len(element.Networks) == 0 {
		errorList = append(errorList, errors.New("Unassigned Networks List fields"))
	}
	for _,network := range element.Networks {
		errorList = append(errorList, network.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Domain) Load(file string) error {
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

func (element *Domain) Import(file string, format string) error {
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
	for i := 0; i < len(element.Networks); i++ {
		err := element.Networks[i].PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *Domain) PostImport() error {
	element.Id=NewUUIDString()
	for _,network := range element.Networks {
		err := network.PostImport()
		if err != nil {
			return err
		}
	}
	return nil
}

func (element *Domain) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

func (element *ProjectDomain) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if len(element.Networks) == 0 {
		errorList = append(errorList, errors.New("Unassigned Networks List fields"))
	}
	for _,network := range element.Networks {
		errorList = append(errorList, network.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ProjectDomain) Load(file string) error {
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

func (element *ProjectDomain) Import(file string, format string) error {
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
	if err == nil && element.Id == "" {
		err := element.PostImport()
		if err != nil {
			return err
		}
	}
	return err
}

func (element *ProjectDomain) PostImport() error {
	element.Id=NewUUIDString()
	for i := 0; i < len(element.Networks); i++ {
		err := element.Networks[i].PostImport()
		if err != nil {
			return err
		}
	}
	return nil
}

func (element *ProjectDomain) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

