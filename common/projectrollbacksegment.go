package common

import (
	"errors"
	"os"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/utils"
)

type ProjectRollbackSegmentInfo struct {
	Format string
	Index  RollBackSegment
}

func (info *ProjectRollbackSegmentInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Index.ProjectId) + "." + info.Index.Index.Index.Value + ".rollbacksegment"
	if _, err = os.Stat(fileName); err != nil {
		index := utils.Index{}
		index.New(SegmentIndexSize)
		info.Index = RollBackSegment{
			Id:        model.NewUUIDString(),
			ProjectId: info.Index.ProjectId,
			Size:      0,
			Index: RollBackSegmentIndex{
				Id:    model.NewUUIDString(),
				Index: index,
			},
			Storage: []ActionStorage{},
		}
		return nil
	}
	err = info.Index.Load(fileName)
	return err
}

func (info *ProjectRollbackSegmentInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Index.ProjectId) + "." + info.Index.Index.Index.Value + ".rollbacksegment"
	err := info.Index.Save(fileName)
	return err
}

func (info *ProjectRollbackSegmentInfo) Import(file string, format string) error {
	err := info.Index.Import(file, format)
	return err
}

func (info *ProjectRollbackSegmentInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + utils.IdToFileFormat(info.Index.ProjectId) + "." + info.Index.Index.Index.Value + ".rollbacksegment"
	_, err := os.Stat(fileName)
	if err == nil {
		return model.DeleteIfExists(fileName)
	}
	return nil
}

func (info *ProjectRollbackSegmentInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return utils.GetJSONFromElem(info.Index, prettify)
	} else if "yaml" == info.Format {
		return utils.GetYAMLFromElem(info.Index)
	} else if "xml" == info.Format {
		return utils.GetXMLFromElem(info.Index, prettify)
	} else {
		return []byte{}, errors.New("Format type : " + info.Format + " not provided ...")
	}
}
