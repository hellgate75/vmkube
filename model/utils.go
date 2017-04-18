package model

import "github.com/satori/go.uuid"

func NewUUIDString()	string {
	return  uuid.NewV4().String()
}

