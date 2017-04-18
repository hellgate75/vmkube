package vmio

import (
	"vmkube/model"
	"errors"
)

type InfrastructureInfo struct {
	Format  	string
	Infra			model.Infrastructure
}

type InfrastructureStream interface {
	Read()		(model.Infrastructure, error)
	Write(project model.Infrastructure) 	bool
	Export(project model.Infrastructure, prettify bool) 	([]byte, error)
	Import(file string, format string) 	(model.Infrastructure, error)
}

func (info *InfrastructureInfo) Read() (model.Infrastructure, error) {
	var project model.Infrastructure
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.Infra.ProjectId
	model.MakeFolderIfNotExists(folder)
	fileName := folder + "/infrastructure.ser"
	err := project.Load(fileName)
	return  project, err
}

func (info *InfrastructureInfo) Write(project model.Infrastructure) bool {
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.Infra.ProjectId
	model.MakeFolderIfNotExists(folder)
	fileName := folder + "/infrastructure.ser"
	err := project.Save(fileName)
	return err == nil
}

func (info *InfrastructureInfo) Import(file string, format string) (model.Infrastructure, error) {
	var project model.Infrastructure
	err := project.Import(file, format)
	return  project, err
}

func (info *InfrastructureInfo) Export(project model.Infrastructure, prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  GetJSONFromObj(project, prettify)
	} else if "xml" == info.Format {
		return  GetXMLFromObj(project, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not known ...")
	}
}
