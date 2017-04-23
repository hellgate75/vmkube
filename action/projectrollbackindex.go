package action

import (
	"vmkube/model"
	"errors"
	"os"
	"vmkube/utils"
)

type ProjectRollbackIndexInfo struct {
	Format  	string
	Index			RollBackIndex
}


func (info *ProjectRollbackIndexInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Index.ProjectId + ".rollbacksegmentindex"
	if _,err = os.Stat(fileName); err!=nil {
		info.Index = RollBackIndex{
			Id: model.NewUUIDString(),
			ProjectId: info.Index.ProjectId,
			IndexList: []RollBackSegmentIndex{},
		}
		return nil
	}
	err = info.Index.Load(fileName)
	return  err
}

func (info *ProjectRollbackIndexInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Index.ProjectId  + ".rollbacksegmentindex"
	err := info.Index.Save(fileName)
	return err
}

func (info *ProjectRollbackIndexInfo) Import(file string, format string) error {
	err := info.Index.Import(file, format)
	return  err
}


func (info *ProjectRollbackIndexInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Index.ProjectId  + ".rollbacksegmentindex"
	return model.DeleteIfExists(fileName)
}


func (info *ProjectRollbackIndexInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  utils.GetJSONFromElem(info.Index, prettify)
	} else if "xml" == info.Format {
		return  utils.GetXMLFromElem(info.Index, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not provided ...")
	}
}

