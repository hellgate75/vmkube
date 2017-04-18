package model

import (
	"github.com/satori/go.uuid"
	"time"
)

const TagMandatoryName = "mandatory"
const TagScriptionName = "descr"


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