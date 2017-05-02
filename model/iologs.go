package model

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"vmkube/utils"
)

func (element *LogStorage) Validate() []error {
	errorList := make([]error, 0)
	if element.InfraId == "" {
		errorList = append(errorList, errors.New("Unassigned Infrastructure Unique Identifier field"))
	}
	if element.InfraId == "" {
		errorList = append(errorList, errors.New("Unassigned Project Unique Identifier field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes, utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *LogStorage) Load(file string) error {
	if !ExistsFile(file) {
		return errors.New("File " + file + " doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(DecodeBytes(byteArray), &element)
	if err == nil {
		element.LogLines = make([]string, 0)
	}
	return err
}

func (element *LogStorage) Import(file string, format string) error {
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

func (element *LogStorage) PostImport() error {
	return nil
}

func (element *LogStorage) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return err
	}
	DeleteIfExists(file)
	return ioutil.WriteFile(file, EncodeBytes(byteArray), 0777)
}
