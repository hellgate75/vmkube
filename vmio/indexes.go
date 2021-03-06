package vmio

import (
	"errors"
	"os"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/utils"
)

type ProjectIndexInfo struct {
	Format string
	Index  model.ProjectsIndex
}

func (info *ProjectIndexInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	err := model.MakeFolderIfNotExists(baseFolder)
	if err != nil {
		return err
	}
	fileName := baseFolder + string(os.PathSeparator) + ".vmkubeindex"
	if _, err = os.Stat(fileName); err != nil {
		info.Index = model.ProjectsIndex{
			Id:       model.NewUUIDString(),
			Projects: []model.ProjectsDescriptor{},
		}
		return nil
	}
	err = info.Index.Load(fileName)
	return err
}

func (info *ProjectIndexInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".vmkubeindex"
	err := info.Index.Save(fileName)
	return err
}

func (info *ProjectIndexInfo) Import(file string, format string) error {
	err := info.Index.Import(file, format)
	return err
}

func (info *ProjectIndexInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".vmkubeindex"
	_, err := os.Stat(fileName)
	if err == nil {
		return model.DeleteIfExists(fileName)
	}
	return nil
}

func (info *ProjectIndexInfo) Export(prettify bool) ([]byte, error) {
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
