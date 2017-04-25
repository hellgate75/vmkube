package utils

import (
	"strings"
	"errors"
	"strconv"
	"encoding/json"
	"encoding/xml"
	"bufio"
	"os"
	"fmt"
	"io/ioutil"
	"syscall"
	"reflect"
	"vmkube/term"
)

func StrPad(instr string, capping int) string {
	strlen := len(instr)
	if strlen == capping  {
		return  instr
	} else  {
		if strlen < capping {
			padding := strings.Repeat(" ", (capping-strlen))
			return  instr + padding
		} else {
			val := instr[0:(capping-2)]
			val += ".."
			return  val
		}
	}
}


func RequestConfirmation(reason string) bool {
	text := ""
	reader := bufio.NewReader(os.Stdin)
	options := "y/n/yes/no"
	options = term.ScreenBold(options)
	allText := fmt.Sprintf("%s Confirm operation [%s] : ", reason, options)
	term.ScreenPrint(allText)
	term.ScreenFlush()
	term.ScreenMoveCursorRelative(-1, len(allText) + 2)
	text, _ = reader.ReadString('\n')
	fmt.Println("")
	if CorrectInput(text) != "y" && CorrectInput(text) != "yes" && CorrectInput(text) != "n" && CorrectInput(text) != "no" {
		text := "y/n/yes/no"
		text = term.ScreenBold(text)
		answer := "Current text is not allowed :"
		answer = term.ScreenColor(answer, term.RED)
		term.ScreenPrintln(fmt.Sprintf("%s : %s\n", answer, text))
		term.ScreenFlush()
		return  RequestConfirmation(reason)
	}
	return (CorrectInput(text) == "y" || CorrectInput(text) == "yes")
}



func CmdParse(key string) (string, error) {
	if len(key) > 0 {
		if strings.Index(key, "--") == 0  {
			return key, errors.New("Invalid Argument : " + key)
		} else	if strings.Index(key, "-") == 0  {
			return key, errors.New("Invalid Argument : " + key)
		} else  {
			return  CorrectInput(key), nil
		}
	} else  {
		return  key, errors.New("Unable to parse Agument : " + key)
	}
}

func OptionsParse(key string, val string) (string, string, error) {
	if strings.Index(key, "--") == 0  {
		return  CorrectInput(key[2:]), val, nil
	} else	if strings.Index(key, "-") == 0  {
		return  CorrectInput(key[1:]), val, nil
	} else  {
		return  key, val, errors.New("Unable to parse Agument : " + key)
	}
}

func CorrectInput(input string) string {
	return  strings.TrimSpace(strings.ToLower(input))
}


func StringToInt(s string) (int,error) {
	return strconv.Atoi(s)
}

func IntToString(n int) string {
	return strconv.Itoa(n)
}


// go binary decoder
func GetJSONFromObj(m interface{}, prettify bool) []byte {
	if prettify {
		bytes,err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return []byte{}
		}
		return bytes
	}
	bytes,err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return bytes
}


func GetXMLFromObj(m interface{}, prettify bool) []byte {
	if prettify {
		bytes,err := xml.MarshalIndent(m, "", "  ")
		if err != nil {
			return []byte{}
		}
		return bytes
	}
	bytes,err := xml.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return bytes
}


func ExportStructureToFile(File string, Format string, structure interface{}) error {
	var bytesArray []byte = make([]byte, 0)
	var err error
	if CorrectInput(Format) == "json" {
		bytesArray, err = GetJSONFromElem(structure, true)
		if err != nil {
			return  err
		}
	} else if CorrectInput(Format) == "xml" {
		bytesArray, err = GetXMLFromElem(structure, true)
		if err != nil {
			return  err
		}
	} else {
			return errors.New("File Format '"+Format+"' not available")
	}
	return ioutil.WriteFile(File, bytesArray, 0777)
}

// go binary decoder
func GetJSONFromElem(m interface{}, prettify bool) ([]byte, error) {
	if prettify {
		return  json.MarshalIndent(m, "", "  ")
	}
	return json.Marshal(m)
}

func GetXMLFromElem(m interface{}, prettify bool) ([]byte, error) {
	if prettify {
		return xml.MarshalIndent(m, "", "  ")
	}
	return xml.Marshal(m)
}

func ToMap(m interface{}) map[string]interface{} {
	var inInterface interface{}
	inrec, _ := json.Marshal(&m)
	json.Unmarshal(inrec, &inInterface)
	return  inInterface.(map[string]interface{})
}

func CreateNewEmptyFile(file string) error {
	return ioutil.WriteFile(file, []byte{}, 0777)
}

func ReverseString(str string) string {
	bytesArray, err := syscall.ByteSliceFromString(str)
	if err != nil {
		return str
	}
	size := len(bytesArray);
	cycle := len(bytesArray)/2;
	for i := 0; i < cycle; i++ {
		var b byte =  bytesArray[i]
		bytesArray[i] = bytesArray[size -1 - i]
		bytesArray[size -1 - i] = b
	}
	return string(bytesArray);
}

func ReverseBytesArray(bytesArray []byte) []byte {
	size := len(bytesArray);
	cycle := len(bytesArray)/2;
	for i := 0; i < cycle; i++ {
		var b byte =  bytesArray[i]
		bytesArray[i] = bytesArray[size -1 - i]
		bytesArray[size -1 - i] = b
	}
	return bytesArray;
}

func IdToFileFormat(id string) string {
	return id//strings.Replace(id, "-", "_", len(id))
}

func NameToFileFormat(id string) string {
	return strings.Replace(id, " ", "_", len(id))
}

func ReducedToStringsSlice(reduced []interface{}) []string {
	arrayOfStrings := make([]string, 0)
	for _,interfaceX := range reduced {
		arrayOfStrings = append(arrayOfStrings, reflect.ValueOf(interfaceX).String())
	}
	return arrayOfStrings
}

func ReducedToIntsSlice(reduced []interface{}) []int {
	arrayOfStrings := make([]int, 0)
	for _,interfaceX := range reduced {
		arrayOfStrings = append(arrayOfStrings, int(reflect.ValueOf(interfaceX).Int()))
	}
	return arrayOfStrings
}

func ReduceStruct(field string, list []interface{}) ([]interface{}, error) {
	sample  := make([]interface{}, 0)
	fieldIndex := -1
	if len(list) == 0 {
		return sample, nil
	}
	structure := list[0]
	if reflect.TypeOf(structure).Kind() != reflect.Struct {
		return sample, errors.New("Type of interface is not a Struct")
	}
	var typeOfField reflect.Type
	typeStruct := reflect.ValueOf(structure)
	found := false
	extraction := false
	for i := 0; i < typeStruct.NumField(); i++ {
		typeField := typeStruct.Type().Field(i)
		typeOfField = typeField.Type
		if typeField.Name == field {
			fieldIndex = i
			found = true
			valueField := typeStruct.Field(i)
			if valueField.IsValid(){
				extraction = true
				
			}
			break
		}
	}
	if ! found {
		return sample, errors.New(fmt.Sprintf("Field %s not found in structure", field))
	}
	if ! extraction {
		return sample, errors.New(fmt.Sprintf("Field %s not valid in structure", field))
	}
//	arrType := typeStruct.Type()
//	arrElemType := arrType.Elem()
	resultSliceType := reflect.SliceOf(typeOfField)
	reduced := reflect.MakeSlice(resultSliceType, 0, len(list))
	
	for _,structure := range list {
		typeStruct := reflect.ValueOf(structure)
//		typeField := typeStruct.Type().Field(fieldIndex)
		valueField := typeStruct.Field(fieldIndex)
		reduced = reflect.Append(reduced, reflect.ValueOf(valueField.Interface()))
	}
	ret := make([]interface{}, 0)
	
	for i:=0; i<reduced.Len(); i++ {
		ret = append(ret,reduced.Index(i).Interface())
	}
	
	return ret, nil
}
