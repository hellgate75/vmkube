package model

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
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
	for _, network := range element.Domains {
		errorList = append(errorList, network.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ProjectImport) Import(file string, format string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	if format != "json" && format != "xml" {
		return errors.New("Format " + format + " not supported!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if format == "json" {
		err = json.Unmarshal(byteArray, &element)
	} else {
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

func (element *ProjectImport) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	for i := 0; i < len(element.Domains); i++ {
		err := element.Domains[i].PostImport()
		if err != nil {
			return err
		}
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
	for _, network := range element.Domains {
		errorList = append(errorList, network.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Project) Load(file string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *Project) Import(file string, format string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	if format != "json" && format != "xml" {
		return errors.New("Format " + format + " not supported!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if format == "json" {
		err = json.Unmarshal(byteArray, &element)
	} else {
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

func (element *Project) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	for i := 0; i < len(element.Domains); i++ {
		err := element.Domains[i].PostImport()
		if err != nil {
			return err
		}
	}
	return nil
}

func (element *Project) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray), 0777)
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
	for _, domain := range element.Domains {
		errorList = append(errorList, domain.Validate()...)
	}
	//errorList = append(errorList, element.State.Validate()...)
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Infrastructure) Load(file string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *Infrastructure) Import(file string, format string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	if format != "json" && format != "xml" {
		return errors.New("Format " + format + " not supported!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if format == "json" {
		err = json.Unmarshal(byteArray, &element)
	} else {
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

func (element *Infrastructure) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	for i := 0; i < len(element.Domains); i++ {
		err := element.Domains[i].PostImport()
		if err != nil {
			return err
		}
	}
	return nil
}

func (element *Infrastructure) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray), 0777)

}
