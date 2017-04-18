package vmio

import (
	"vmkube/model"
	"errors"
)

type ProjectInfo struct {
	Project 			model.Project
	Format  			string
}

type ProjectStream interface {
	Read()		error
	Write() 	bool
	Export(prettify bool) 	([]byte, error)
	Import(file string, format string) 	error
}

func (info *ProjectInfo) Read() error {
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.Project.Id
	model.MakeFolderIfNotExists(folder)
	fileName := folder + "/project.ser"
	err := info.Project.Load(fileName)
	return err
}

func (info *ProjectInfo) Write() bool {
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.Project.Id
	model.MakeFolderIfNotExists(folder)
	fileName := folder + "/project.ser"
	err := info.Project.Save(fileName)
	return err == nil
}

func (info *ProjectInfo) Import(file string, format string) error {
	err := info.Project.Import(file, format)
	return  err
}

func (info *ProjectInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  GetJSONFromObj(info.Project, prettify)
	} else if "xml" == info.Format {
		return  GetXMLFromObj(info.Project, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not known ...")
	}
}
