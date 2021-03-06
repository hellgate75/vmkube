package model

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/hellgate75/vmkube/utils"
)

func (element *CloudInstance) Validate() []error {
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
	if len(element.Options) == 0 {
		errorList = append(errorList, errors.New("Unassigned Vendor specific Options field"))
	}
	if element.Hostname == "" {
		errorList = append(errorList, errors.New("Unassigned host name field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *CloudInstance) Load(file string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *CloudInstance) Import(file string, format string) error {
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

func (element *CloudInstance) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	if element.Disabled {
		element.Disabled = false
	}
	return nil
}

func (element *CloudInstance) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray), 0777)
}

func (element *CloudMachine) Validate() []error {
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
	if len(element.Options) == 0 {
		errorList = append(errorList, errors.New("Unassigned Vendor specific Options field"))
	}
	if element.Hostname == "" {
		errorList = append(errorList, errors.New("Unassigned host name field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *CloudMachine) Load(file string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(DecodeBytes(byteArray), &element)
}

func (element *CloudMachine) Import(file string, format string) error {
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

func (element *CloudMachine) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	return nil
}

func (element *CloudMachine) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray), 0777)
}
