package model

import (
	"errors"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
)

func (element *ProjectsDescriptor) Load(file string) error {
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

func (element *ProjectsDescriptor) Import(file string, format string) error {
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

func (element *ProjectsDescriptor) Save(file string) error {
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

func (element *ProjectsIndex) Load(file string) error {
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

func (element *ProjectsIndex) Import(file string, format string) error {
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

func (element *ProjectsIndex) Save(file string) error {
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
