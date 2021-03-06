package model

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/hellgate75/vmkube/utils"
)

func (element *Installation) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.InstanceId == "" {
		errorList = append(errorList, errors.New("Unassigned Instance Unique Identifier field"))
	}
	errorList = append(errorList, element.Plan.Validate()...)
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *Installation) Load(file string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *Installation) Import(file string, format string) error {
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
	} else if format == "yaml" {
		err = yaml.Unmarshal(byteArray, &element)
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

func (element *Installation) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	return nil
}

func (element *Installation) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray), 0777)
}

func (element *InstallationPlan) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.MachineId == "" {
		errorList = append(errorList, errors.New("Unassigned Machine Unique Identifier field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *InstallationPlan) Load(file string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *InstallationPlan) Import(file string, format string) error {
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
	} else if format == "yaml" {
		err = yaml.Unmarshal(byteArray, &element)
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

func (element *InstallationPlan) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	return nil
}

func (element *InstallationPlan) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray), 0777)
}
