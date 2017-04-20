package model

import (
	"github.com/satori/go.uuid"
	"time"
	"reflect"
	"vmkube/utils"
	"fmt"
	"os"
	"strings"
)

const TagMandatoryName = "mandatory"
const TagDescName = "descr"
const TagTypeName = "type"


func NewUUIDString()	string {
	return  uuid.NewV4().String()
}

func ProjectFromImport(imported ProjectImport) Project {
	return Project{
		Id: imported.Id,
		Name: imported.Name,
		Domains: imported.Domains,
		Open: true,
		LastMessage: "",
		Created: time.Now(),
		Modified: time.Now(),
		Errors: false,
	}
}

func ProjectToImport(project Project) ProjectImport {
	return ProjectImport{
		Id: project.Id,
		Name: project.Name,
		Domains: project.Domains,
	}
}

func EncodeBytes(decodedByteArray []byte) []byte {
	newBytes := make([]byte,0)
	for _,byteElem := range decodedByteArray {
		newBytes = append(newBytes, byteElem << 1)
	}
	return newBytes
	//return decodedByteArray
}

func DecodeBytes(encodedByteArray []byte) []byte {
	newBytes := make([]byte,0)
	for _,byteElem := range encodedByteArray {
		newBytes = append(newBytes, byteElem >> 1)
	}
	return newBytes
	//return encodedByteArray
}

type FieldData struct {
	JName       string
	XName       string
	Type        string
	Desc        string
	Mandatory   bool
	Fields      []FieldData
}

func recoverSubFields(sf reflect.StructField, vf reflect.Value) []FieldData {
	fields := make([]FieldData, 0)
	if sf.Type.Kind() == reflect.Struct {
		flds, err := DescribeStruct(vf.Interface())
		if err == nil {
			fields = append(fields, flds...)
		}
	} else if sf.Type.Kind() == reflect.Slice {
		if vf.Len() > 0 {
			sVal := vf.Index(0)
			if sVal.Kind() == reflect.Struct {
				flds, err := DescribeStruct(sVal.Interface())
				if err == nil {
					fields = append(fields, flds...)
				}
			}
		}
	}
	return fields
}

func DescribeStruct(s interface{}) ([]FieldData, error) {
	fields := make([]FieldData, 0)
	
	typeStruct := reflect.ValueOf(s)
	
	//typeElem := typeStruct.Elem()
	
	n := typeStruct.NumField()
	for i := 0; i < n; i++ {
		sField := typeStruct.Type().Field(i)
		vField := typeStruct.Field(i)
		fieldData := FieldData{}
		val := sField.Tag.Get("json")
		if val != "" {
			fieldData.JName = val
		}
		val = sField.Tag.Get("xml")
		if val != "" {
			fieldData.XName = val
		}
		val = sField.Tag.Get(TagDescName)
		if val != "" {
			fieldData.Desc = val
		}
		val = sField.Tag.Get(TagMandatoryName)
		if val != "" {
			fieldData.Mandatory = (utils.CorrectInput(val) == "yes")
		}
		val = sField.Tag.Get(TagTypeName)
		if val != "" {
			fieldData.Type = val
			fieldData.Fields = recoverSubFields(sField, vField)
		}
		fields = append(fields, fieldData)
	}
	return fields,nil
	
}

func PrintFieldsHeader(print bool) {
	if print {
		fmt.Fprintf(os.Stdout, "%s  %s  %s  %s  %s\n", utils.StrPad("JSON Field", 30), utils.StrPad("XML Tag", 20), utils.StrPad("Data Type Description", 30),utils.StrPad("Mandatory", 9), "Description")
		
	}
}

func PrintFieldsRecursively(fields []FieldData, index int) {
	for _,field := range fields {
		mandStr := "no"
		if field.Mandatory {
			mandStr = "yes"
		}
		tab := strings.Repeat("  ", index)
		fmt.Fprintf(os.Stdout, "%s  %s  %s  %s  %s\n", utils.StrPad(tab+field.JName, 30), utils.StrPad(field.XName, 20), utils.StrPad(field.Type, 30),utils.StrPad("   " + mandStr, 9), field.Desc)
		PrintFieldsRecursively(field.Fields, index + 1)
	}
}
