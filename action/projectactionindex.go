package action

import (
	"vmkube/model"
	"errors"
	"os"
	"vmkube/utils"
)

type ProjectActionIndexInfo struct {
	Format  	string
	Index			ProjectActionIndex
}


func (info *ProjectActionIndexInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Index.ProjectId + ".actionindex"
	if _,err = os.Stat(fileName); err!=nil {
		info.Index = ProjectActionIndex{
			Id: model.NewUUIDString(),
			ProjectId: info.Index.ProjectId,
			Actions: []ActionDescriptor{},
		}
		return nil
	}
	err = info.Index.Load(fileName)
	return  err
}

func (info *ProjectActionIndexInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Index.ProjectId  + ".actionindex"
	err := info.Index.Save(fileName)
	return err
}

func (info *ProjectActionIndexInfo) Import(file string, format string) error {
	err := info.Index.Import(file, format)
	return  err
}


func (info *ProjectActionIndexInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Index.ProjectId  + ".actionindex"
	return model.DeleteIfExists(fileName)
}


func (info *ProjectActionIndexInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  utils.GetJSONFromElem(info.Index, prettify)
	} else if "xml" == info.Format {
		return  utils.GetXMLFromElem(info.Index, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not provided ...")
	}
}

