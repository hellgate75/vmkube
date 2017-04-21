package model

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
	"vmkube/utils"
)

func (element *InstanceState) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Hostname == "" {
		errorList = append(errorList, errors.New("Unassigned host name field"))
	}
	if len(element.IPAddresses) == 0 {
		errorList = append(errorList, errors.New("Unassigned IP address list field"))
	}
	if element.InstanceId == "" {
		errorList = append(errorList, errors.New("Unassigned Instance Unique Identifier field"))
	}
	if element.NetworkId == "" {
		errorList = append(errorList, errors.New("Unassigned Network Unique Identifier field"))
	}
	if element.DomainId == "" {
		errorList = append(errorList, errors.New("Unassigned Domain Unique Identifier field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *InstanceState) Load(file string) error {
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

func (element *InstanceState) Import(file string, format string) error {
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

func (element *InstanceState) PostImport() error {
	element.Id = NewUUIDString()
	return nil
}

func (element *InstanceState) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

func (element *NetworkState) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.NetworkId == "" {
		errorList = append(errorList, errors.New("Unassigned Network Unique Identifier field"))
	}
	if element.DomainId == "" {
		errorList = append(errorList, errors.New("Unassigned Domain Unique Identifier field"))
	}
	if len(element.InstanceStates) == 0 {
		errorList = append(errorList, errors.New("Unassigned Instance States List fields"))
	}
	for _,state := range element.InstanceStates {
		errorList = append(errorList, state.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *NetworkState) Load(file string) error {
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

func (element *NetworkState) Import(file string, format string) error {
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

func (element *NetworkState) PostImport() error {
	element.Id = NewUUIDString()
	for _, states := range element.InstanceStates {
		states.PostImport()
	}
	return nil
}

func (element *NetworkState) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

func (element *DomainState) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.DomainId == "" {
		errorList = append(errorList, errors.New("Unassigned Domain Unique Identifier field"))
	}
	if len(element.NetworkStates) == 0 {
		errorList = append(errorList, errors.New("Unassigned Network States List fields"))
	}
	for _,state := range element.NetworkStates {
		errorList = append(errorList, state.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *DomainState) Load(file string) error {
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

func (element *DomainState) Import(file string, format string) error {
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

func (element *DomainState) PostImport() error {
	element.Id = NewUUIDString()
	for _, states := range element.NetworkStates {
		states.PostImport()
	}
	return nil
}

func (element *DomainState) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

func (element *State) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if len(element.DomainStates) == 0 {
		errorList = append(errorList, errors.New("Unassigned Domain States List fields"))
	}
	for _,state := range element.DomainStates {
		errorList = append(errorList, state.Validate()...)
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *State) Load(file string) error {
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

func (element *State) Import(file string, format string) error {
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

func (element *State) PostImport() error {
	element.Id = NewUUIDString()
	for _, states := range element.DomainStates {
		states.PostImport()
	}
	return nil
}

func (element *State) Save(file string) error {
	byteArray, err := json.MarshalIndent(element, "", "  ")
	if err != nil {
		return  err
	}
	DeleteIfExists(file)
	err = ioutil.WriteFile(file, EncodeBytes(byteArray) , 0777)
	return  err
}

