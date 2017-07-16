package vmio

import (
	"errors"
	"github.com/hellgate75/vmkube/model"
	"github.com/hellgate75/vmkube/utils"
)

type ProjectImportInfo struct {
	ProjectImport model.ProjectImport
	Format        string
}

func (info *ProjectImportInfo) Import(file string, format string) error {
	err := info.ProjectImport.Import(file, format)
	return err
}

func (info *ProjectImportInfo) Export(prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return utils.GetJSONFromElem(info.ProjectImport, prettify)
	} else if "yaml" == info.Format {
		return utils.GetYAMLFromElem(info.ProjectImport)
	} else if "xml" == info.Format {
		return utils.GetXMLFromElem(info.ProjectImport, prettify)
	} else {
		return []byte{}, errors.New("Format type : " + info.Format + " not provided ...")
	}
}
