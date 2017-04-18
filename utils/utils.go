package utils

import (
	"strings"
	"errors"
	"strconv"
	"encoding/json"
	"encoding/xml"
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

func CmdParseOption(key string, options [][]string) (string, int, error) {
	if len(key) > 0 {
		if strings.Index(key, "--") == 0  {
			return key, -1, errors.New("Invalid Argument (wrong characters: --) : " + key)
		} else	if strings.Index(key, "-") == 0  {
			return key, -1, errors.New("Invalid Argument (wrong character: -) : " + key)
		} else  {
			for index,opts := range options {
				if CorrectInput(key) == opts[0]  {
					return  CorrectInput(key), index, nil
				}
			}
			return  key, -1, errors.New("Invalid Argument : " + key)
		}
	} else  {
		return  key, -1, errors.New("Unable to parse Agument : " + key)
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


