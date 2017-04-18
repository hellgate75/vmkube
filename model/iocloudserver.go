package model

import (
	"errors"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"vmkube/vmio"
)

func (element *CloudInstance) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Name == "" {
		errorList = append(errorList, errors.New("Unassigned Name field"))
	}
	if element.ServerId == "" {
		errorList = append(errorList, errors.New("Unassigned Project Server Id field"))
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
		bytes = append(bytes,vmio.GetJSONFromObj(element, true))
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *CloudInstance) Load(file string) error {
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

func (element *CloudInstance) Import(file string, format string) error {
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

func (element *CloudInstance) Save(file string) error {
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

func (element *ProjectCloudServer) Validate() []error {
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
		bytes = append(bytes,vmio.GetJSONFromObj(element, true))
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ProjectCloudServer) Load(file string) error {
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

func (element *ProjectCloudServer) Import(file string, format string) error {
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

func (element *ProjectCloudServer) Save(file string) error {
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

