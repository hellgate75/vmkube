package utils

import (
	"strings"
	"errors"
	"strconv"
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

