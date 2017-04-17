package model

import (
	"os"
)

func existsFile(file string) bool {
	_,err := os.Stat(file)
	return  err == nil
}

func deleteIfExists(file string) {
	_,err := os.Stat(file)
	if err != nil {
		os.Remove(file)
	}
}

type IONature interface {
	Load(file string) error
	Import(file string, format string) error
	Save(file string) error
}



