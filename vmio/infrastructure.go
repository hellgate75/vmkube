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
	Read()		error
	Write() 	bool
	Export(prettify bool) 	([]byte, error)
	Import(file string, format string) 	error
}

func (info *InfrastructureInfo) Read() error {
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.Infra.ProjectId
	model.MakeFolderIfNotExists(folder)
	fileName := folder + "/infrastructure.ser"
	err := info.Infra.Load(fileName)
	return  err
}

func (info *InfrastructureInfo) Write() bool {
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.Infra.ProjectId
	model.MakeFolderIfNotExists(folder)
	fileName := folder + "/infrastructure.ser"
	err := info.Infra.Save(fileName)
	return err == nil
}

func (info *InfrastructureInfo) Import(file string, format string) error {
	err := info.Infra.Import(file, format)
	return  err
}

func (info *InfrastructureInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  GetJSONFromObj(info.Infra, prettify)
	} else if "xml" == info.Format {
		return  GetXMLFromObj(info.Infra, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not known ...")
	}
}
