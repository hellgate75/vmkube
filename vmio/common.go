package vmio

import (
	"errors"
	"fmt"
	"strings"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/utils"
)

type VMKubeElementsStream interface {
	Read() error
	Write() error
	Export(prettify bool) ([]byte, error)
	Import(file string, format string) error
}

func GetProjectDescriptor(name string) (model.ProjectsDescriptor, error) {
	descriptor := model.ProjectsDescriptor{}
	indexes, err := LoadIndex()
	if err != nil {
		return descriptor, err
	}
	for _, index := range indexes.Projects {
		if utils.CorrectInput(index.Name) == utils.CorrectInput(name) {
			return index, nil
		}
	}
	return descriptor, errors.New("Project '" + name + "' not found in project indexes")
}

func GetInfrastructureProjectDescriptor(name string) (model.ProjectsDescriptor, error) {
	descriptor := model.ProjectsDescriptor{}
	indexes, err := LoadIndex()
	if err != nil {
		return descriptor, err
	}
	for _, index := range indexes.Projects {
		if utils.CorrectInput(index.InfraName) == utils.CorrectInput(name) {
			return index, nil
		}
	}
	return descriptor, errors.New("Infrastructure '" + name + "' not found in project indexes")
}

func StripOptions(options [][]string) (int, string) {
	optionsStripped := ""
	counter := 0
	for _, option := range options {
		if counter > 0 {
			optionsStripped += ","
		} else {
			counter++
		}
		if len(option) == 1 {
			optionsStripped += option[0]
		} else {
			optionsStripped += strings.Join(option, "=")
		}
	}
	return len(options), optionsStripped
}
func StripErrorMessages(description string, errorList []error) (int, string) {
	errorsStripped := description
	for _, singleError := range errorList {
		if singleError != nil {
			errorsStripped += fmt.Sprintf("\n%s", singleError.Error())
		}
	}
	return len(errorList), errorsStripped
}
