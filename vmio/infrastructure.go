package vmio

import (
	"vmkube/model"
	"errors"
	"vmkube/utils"
	"os"
)

type InfrastructureInfo struct {
	Format  	string
	Infra			model.Infrastructure
}

func (info *InfrastructureInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".infra-" + utils.IdToFileFormat(info.Infra.ProjectId) + ".vmkube"
	if _,err := os.Stat(fileName); err!=nil {
		return err
	}
	err := info.Infra.Load(fileName)
	return  err
}

func (info *InfrastructureInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".infra-" + utils.IdToFileFormat(info.Infra.ProjectId) + ".vmkube"
	err := info.Infra.Save(fileName)
	return err
}

func (info *InfrastructureInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".infra-" + utils.IdToFileFormat(info.Infra.ProjectId) + ".vmkube"
	_,err := os.Stat(fileName)
	if err == nil {
		return model.DeleteIfExists(fileName)
	}
	return nil
}

func (info *InfrastructureInfo) Evacuate() error {
	//TODO: Implement infrastructure evacuate
	return nil
}

func (info *InfrastructureInfo) Import(file string, format string) error {
	err := info.Infra.Import(file, format)
	return  err
}

func (info *InfrastructureInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  utils.GetJSONFromElem(info.Infra, prettify)
	} else if "xml" == info.Format {
		return  utils.GetXMLFromElem(info.Infra, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not provided ...")
	}
}
