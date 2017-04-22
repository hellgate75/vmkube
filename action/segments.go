package action

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"vmkube/utils"
	"encoding/xml"
	"vmkube/model"
)

type ActionDescriptor struct {
	Id          string            `json:"Id" xml:"Id" mandatory:"yes" descr:"Project Unique Identifier" type:"text"`
	Request     CmdRequestType  	`json:"Request" xml:"Request" mandatory:"yes" descr:"Request Code" type:"int"`
	SubRequest  CmdSubRequestType `json:"SubRequest" xml:"SubRequest" mandatory:"yes" descr:"Sub-Request Code" type:"int"`
	ElementType CmdElementType  	`json:"ElementType" xml:"ElementType" mandatory:"yes" descr:"Elemrnt Type Code" type:"int"`
	ElementName string            `json:"ElementName" xml:"ElementName" mandatory:"yes" descr:"Element Name" type:"text"`
	JSONImage   string            `json:"JSONImage" xml:"JSONImage" mandatory:"yes" descr:"Element JSON image" type:"text"`
	FullProject bool  						`json:"FullProject" xml:"FullProject" mandatory:"yes" descr:"Describe if action infear on all project" type:"boolean"`
	DropAction  bool  						`json:"DropAction" xml:"DropAction" mandatory:"yes" descr:"Describe if action Drops Elements" type:"boolean"`
}

func (element *ActionDescriptor) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.Request == NoCommand {
		errorList = append(errorList, errors.New("Unassigned Command Request field"))
	}
	if element.SubRequest == NoSubCommand {
		errorList = append(errorList, errors.New("Unassigned Sub-Command Request field"))
	}
	if element.ElementType == NoElement {
		errorList = append(errorList, errors.New("Unassigned Element Type field"))
	}
	if element.ElementName == "" {
		errorList = append(errorList, errors.New("Unassigned Element Name field"))
	}
	if element.JSONImage == "" {
		errorList = append(errorList, errors.New("Unassigned JSON Image field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ActionDescriptor) Load(file string) error {
	if !model. ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(model.DecodeBytes(byteArray), &element)
	
}

func (element *ActionDescriptor) Import(file string, format string) error {
	if ! model.ExistsFile(file) {
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
	println(element)
	if err == nil {
		err := element.PostImport()
		if err != nil {
			return err
		}
	}
	println(element)
	return err
}

func (element *ActionDescriptor) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	return nil
}

func (element *ActionDescriptor) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	model.DeleteIfExists(file)
	return ioutil.WriteFile(file, model.EncodeBytes(byteArray) , 0777)
}

/*
Describe Projects Index, contains

  * Id          (string)                  Indexes Unique Identifier

  * Projects    ([]ProjectsDescriptor)    Projects indexed in VMKube
*/
type ProjectsActionIndex struct {
	Id          string                  `json:"Id" xml:"Id" mandatory:"yes" descr:"Project Unique Identifier" type:"text"`
	Actions		  []ActionDescriptor 		`json:"Projects" xml:"Projects" mandatory:"yes" descr:"Projects indexed in VMKube" type:"object Projects list"`
}

func (element *ProjectsActionIndex) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	//if len(element.Actions) == 0 {
	//	errorList = append(errorList, errors.New("Unassigned actions field"))
	//}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *ProjectsActionIndex) Load(file string) error {
	if ! model.ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(model.DecodeBytes(byteArray), &element)
}

func (element *ProjectsActionIndex) Import(file string, format string) error {
	if ! model.ExistsFile(file) {
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
	println(element)
	if err == nil {
		err := element.PostImport()
		if err != nil {
			return err
		}
	}
	println(element)
	return err
}

func (element *ProjectsActionIndex) PostImport() error {
	if element.Id == "" {
		element.Id = NewUUIDString()
	}
	for i := 0; i < len(element.Actions); i++ {
		element.Actions[i].PostImport()
	}
	return nil
}

func (element *ProjectsActionIndex) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	model.DeleteIfExists(file)
	return ioutil.WriteFile(file, model.EncodeBytes(byteArray) , 0777)
}
