package action

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"vmkube/utils"
	"encoding/xml"
	"vmkube/model"
	"time"
)

var SegmentIndexSize int = 30

type SegmentIndexNature interface {
	New()
	NewNextFrom(previousIndex utils.Index)
	NewPreviousFrom(nextIndex utils.Index)
	Before(segmentIndex RollBackSegmentIndex)
}

/*
Describe Action Storage, contains

  * Id          (string)    Indexes Unique Identifier

  * Action      (Action)    Specific Action

  * Date        (time.Time) Action Store/Update Date
*/
type ActionStorage struct {
	Id          string            `json:"Id" xml:"Id" mandatory:"yes" descr:"Rollback Descriptor Unique Identifier" type:"text"`
	Action      ActionDescriptor  `json:"Action" xml:"Action" mandatory:"yes" descr:"Specific Action" type:"text"`
	Date        time.Time         `json:"Date" xml:"Date" mandatory:"yes" descr:"Specific Action Rollback registration Date" type:"datetime"`
}

/*
Describe RollBack Segment, contains

  * Id          (string)                  Indexes Unique Identifier

  * Projects    ([]ActionStorage)    Projects indexed in VMKube
*/
type RollBackSegment struct {
	Id                string                  `json:"Id" xml:"Id" mandatory:"yes" descr:"Action Index Unique Identifier" type:"text"`
	ProjectId         string                  `json:"ProjectId" xml:"ProjectId" mandatory:"yes" descr:"Project Unique Identifier" type:"text"`
	Storage   []ActionStorage 		            `json:"Storage" xml:"Storage" mandatory:"yes" descr:"Project Storage list" type:"object ActionStorage list"`
	Index     RollBackSegmentIndex 	          `json:"Index" xml:"Index" mandatory:"yes" descr:"Rollback Segment Index" type:"object RollBackSegmentIndex list"`
	Size              int 	                  `json:"Size" xml:"Size" mandatory:"yes" descr:"Rollback Segment Size in Action Storage Elements" type:"object RollBackSegmentIndex list"`
}

func (element *RollBackSegment) Load(file string) error {
	if !model. ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(model.DecodeBytes(byteArray), &element)
	
}

func (element *RollBackSegment) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	model.DeleteIfExists(file)
	return ioutil.WriteFile(file, model.EncodeBytes(byteArray) , 0777)
}

func (element *RollBackSegment) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.ProjectId == "" {
		errorList = append(errorList, errors.New("Unassigned Project Unique Identifier field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}

func (element *RollBackSegment) Import(file string, format string) error {
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

func (element *RollBackSegment) PostImport() error {
	if element.Id == "" {
		element.Id = model.NewUUIDString()
	}
	return nil
}



/*
Describe Projects Index, contains

  * Id          (string)                  Indexes Unique Identifier

  * Index        (utils.Index)            Projects indexed in VMKube
*/
type RollBackSegmentIndex struct {
	Id                string                `json:"Id" xml:"Id" mandatory:"yes" descr:"Action Index Unique Identifier" type:"text"`
	Index             utils.Index 		      `json:"Index" xml:"Index" mandatory:"yes" descr:"Project Rollback Segment Index" type:"string"`
}

func (element *RollBackSegmentIndex) New() {
	element.Id = model.NewUUIDString()
	element.Index.New(SegmentIndexSize)
	element.Index = element.Index.Next()
}

func (element *RollBackSegmentIndex) NewNextFrom(previousIndex utils.Index) {
	element.Id = model.NewUUIDString()
	element.Index = previousIndex.Next()
}

func (element *RollBackSegmentIndex) NewPreviousFrom(nextIndex utils.Index) {
	element.Id = model.NewUUIDString()
	element.Index = nextIndex.Previous()
}

func (element *RollBackSegmentIndex) Before(segmentIndex RollBackSegmentIndex) bool {
	return element.Index.Compare(segmentIndex.Index) < 0
}

/*
Describe Projects Index, contains

  * Id          (string)                  Indexes Unique Identifier

  * Projects    ([]ProjectsDescriptor)    Projects indexed in VMKube
*/
type RollBackIndex struct {
	Id                string                  `json:"Id" xml:"Id" mandatory:"yes" descr:"Action Index Unique Identifier" type:"text"`
	ProjectId         string                  `json:"ProjectId" xml:"ProjectId" mandatory:"yes" descr:"Project Unique Identifier" type:"text"`
	IndexList         []RollBackSegmentIndex 	`json:"IndexList" xml:"IndexList" mandatory:"yes" descr:"Rollback Segments Index list" type:"object RollBackSegmentIndex list"`
}

func (element *RollBackIndex) Validate() []error {
	errorList := make([]error, 0)
	if element.Id == "" {
		errorList = append(errorList, errors.New("Unassigned Unique Identifier field"))
	}
	if element.ProjectId == "" {
		errorList = append(errorList, errors.New("Unassigned Project Unique Identifier field"))
	}
	if len(errorList) > 0 {
		bytes := []byte(`Errors reported in json : `)
		bytes = append(bytes,utils.GetJSONFromObj(element, true)...)
		errorList = append(errorList, errors.New(string(bytes)))
	}
	return errorList
}


func (element *RollBackIndex) Load(file string) error {
	if !model. ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(model.DecodeBytes(byteArray), &element)
	
}

func (element *RollBackIndex) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	model.DeleteIfExists(file)
	return ioutil.WriteFile(file, model.EncodeBytes(byteArray) , 0777)
}

func (element *RollBackIndex) Import(file string, format string) error {
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

func (element *RollBackIndex) PostImport() error {
	if element.Id == "" {
		element.Id = model.NewUUIDString()
	}
	return nil
}


type ActionDescriptor struct {
	Id          string            `json:"Id" xml:"Id" mandatory:"yes" descr:"Action Unique Identifier" type:"text"`
	Request     CmdRequestType  	`json:"Request" xml:"Request" mandatory:"yes" descr:"Request Code" type:"int"`
	SubRequest  CmdSubRequestType `json:"SubRequest" xml:"SubRequest" mandatory:"yes" descr:"Sub-Request Code" type:"int"`
	ElementType CmdElementType  	`json:"ElementType" xml:"ElementType" mandatory:"yes" descr:"Elemrnt Type Code" type:"int"`
	ElementName string            `json:"ElementName" xml:"ElementName" mandatory:"yes" descr:"Element Name" type:"text"`
	JSONImage   string            `json:"JSONImage" xml:"JSONImage" mandatory:"yes" descr:"Element JSON image" type:"text"`
	FullProject bool  						`json:"FullProject" xml:"FullProject" mandatory:"yes" descr:"Describe if action infear on all project" type:"boolean"`
	DropAction  bool  						`json:"DropAction" xml:"DropAction" mandatory:"yes" descr:"Describe if action Drops Elements" type:"boolean"`
	Date        time.Time         `json:"Date" xml:"Date" mandatory:"yes" descr:"Specific Action Date" type:"text"`
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
		element.Id = model.NewUUIDString()
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
Describe Project Action Index, contains

  * Id          (string)               Indexes Unique Identifier

  * Actions    ([]ActionDescriptor)    Projects ActionDescriptor in VMKube Project
*/
type ProjectActionIndex struct {
	Id          string                  `json:"Id" xml:"Id" mandatory:"yes" descr:"Action Index Unique Identifier" type:"text"`
	ProjectId   string                  `json:"ProjectId" xml:"ProjectId" mandatory:"yes" descr:"Project Unique Identifier" type:"text"`
	Actions		  []ActionDescriptor 		  `json:"Actions" xml:"Actions" mandatory:"yes" descr:"Project Actions, indexed in VMKube" type:"object ActionDescriptor list"`
}

func (element *ProjectActionIndex) Validate() []error {
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

func (element *ProjectActionIndex) Load(file string) error {
	if ! model.ExistsFile(file) {
		return  errors.New("File "+file+" doesn't exist!!")
	}
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		return  err
	}
	return json.Unmarshal(model.DecodeBytes(byteArray), &element)
}

func (element *ProjectActionIndex) Import(file string, format string) error {
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

func (element *ProjectActionIndex) PostImport() error {
	if element.Id == "" {
		element.Id = model.NewUUIDString()
	}
	for i := 0; i < len(element.Actions); i++ {
		element.Actions[i].PostImport()
	}
	return nil
}

func (element *ProjectActionIndex) Save(file string) error {
	byteArray, err := json.Marshal(element)
	if err != nil {
		return  err
	}
	model.DeleteIfExists(file)
	return ioutil.WriteFile(file, model.EncodeBytes(byteArray) , 0777)
}
