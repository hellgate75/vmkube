package vmio

import (
	"errors"
	"os"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/utils"
)

type ProjectInfo struct {
	Project model.Project
	Format  string
}

func (info *ProjectInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".prj-" + utils.IdToFileFormat(info.Project.Id) + ".vmkube"
	if _, err := os.Stat(fileName); err != nil {
		return err
	}
	err := info.Project.Load(fileName)
	return err
}

func (info *ProjectInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".prj-" + utils.IdToFileFormat(info.Project.Id) + ".vmkube"
	err := info.Project.Save(fileName)
	return err
}

func (info *ProjectInfo) Delete() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) + ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + ".prj-" + utils.IdToFileFormat(info.Project.Id) + ".vmkube"
	_, err := os.Stat(fileName)
	if err == nil {
		return model.DeleteIfExists(fileName)
	}
	return nil
}

func (info *ProjectInfo) Import(file string, format string) error {
	err := info.Project.Import(file, format)
	return err
}

func (info *ProjectInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return utils.GetJSONFromElem(info.Project, prettify)
	} else if "yaml" == info.Format {
		return utils.GetYAMLFromElem(info.Project)
	} else if "xml" == info.Format {
		return utils.GetXMLFromElem(info.Project, prettify)
	} else {
		return []byte{}, errors.New("Format type : " + info.Format + " not provided ...")
	}
}
