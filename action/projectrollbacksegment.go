package action

import (
	"vmkube/model"
	"errors"
	"os"
	"vmkube/utils"
)

type ProjectRollbackRollBackSegmentInfo struct {
	Format  	string
	Index			RollBackSegment
}


func (info *ProjectRollbackRollBackSegmentInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + info.Index.ProjectId + ".rollbacksegment"
	if _,err = os.Stat(fileName); err!=nil {
		index := RollBackSegmentIndex{}
		index.New()
		info.Index = RollBackSegment{
			Id: model.NewUUIDString(),
			ProjectId: info.Index.ProjectId,
			Size: 0,
			Index: RollBackSegmentIndex{
				Id: model.NewUUIDString(),
				Index: index,
			},
			Storage: []ActionStorage{},
		}
		return nil
	}
	err = info.Index.Load(fileName)
	return  err
}

func (info *ProjectRollbackRollBackSegmentInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + info.Index.ProjectId  + ".rollbacksegment"
	err := info.Index.Save(fileName)
	return err
}

func (info *ProjectRollbackRollBackSegmentInfo) Import(file string, format string) error {
	err := info.Index.Import(file, format)
	return  err
}


func (info *ProjectRollbackRollBackSegmentInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + info.Index.ProjectId  + ".rollbacksegment"
	return model.DeleteIfExists(fileName)
}


func (info *ProjectRollbackRollBackSegmentInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  utils.GetJSONFromElem(info.Index, prettify)
	} else if "xml" == info.Format {
		return  utils.GetXMLFromElem(info.Index, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not provided ...")
	}
}

