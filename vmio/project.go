package vmio

import (
	"vmkube/model"
	"errors"
	"vmkube/utils"
	"os"
)

type ProjectInfo struct {
	Project 			model.Project
	Format  			string
}

func (info *ProjectInfo) Read() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Project.Id + ".project"
	if _,err := os.Stat(fileName); err!=nil {
		info.Write()
	}
	err := info.Project.Load(fileName)
	return err
}

func (info *ProjectInfo) Write() error {
	baseFolder := model.VMBaseFolder() + string(os.PathSeparator) +  ".data"
	model.MakeFolderIfNotExists(baseFolder)
	fileName := baseFolder + string(os.PathSeparator) + "." + info.Project.Id + ".project"
	model.DeleteIfExists(fileName)
	err := info.Project.Save(fileName)
	return err
}

func (info *ProjectInfo) Import(file string, format string) error {
	err := info.Project.Import(file, format)
	return  err
}

func (info *ProjectInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  utils.GetJSONFromElem(info.Project, prettify)
	} else if "xml" == info.Format {
		return  utils.GetXMLFromElem(info.Project, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not known ...")
	}
}
