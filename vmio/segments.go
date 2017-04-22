package vmio

import (
	"vmkube/model"
	"errors"
	"vmkube/utils"
	"os"
)

type ActionIndexInfo struct {
	Format  	string
	Index			model.ProjectsIndex
}


func (info *ActionIndexInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + ".vmkubeaction"
	if _,err = os.Stat(fileName); err!=nil {
		info.Index = model.ProjectsIndex{
			Id: model.NewUUIDString(),
			Projects: []model.ProjectsDescriptor{},
		}
		return nil
	}
	err = info.Index.Load(fileName)
	return  err
}

func (info *ActionIndexInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".vmkubeaction"
	err := info.Index.Save(fileName)
	return err
}

func (info *ActionIndexInfo) Import(file string, format string) error {
	err := info.Index.Import(file, format)
	return  err
}


func (info *ActionIndexInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".vmkubeaction"
	return model.DeleteIfExists(fileName)
}


func (info *ActionIndexInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  utils.GetJSONFromElem(info.Index, prettify)
	} else if "xml" == info.Format {
		return  utils.GetXMLFromElem(info.Index, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not known ...")
	}
}
