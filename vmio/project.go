package vmio

import (
	"vmkube/model"
	"os"
	"errors"
)

type ProjectInfo struct {
	ProjectId 	string
	ProjectName 	string
	Format  			string
}

type ProjectStream interface {
	Read()		(model.Project, error)
	Write(project model.Project) 	bool
	Export(project model.Project, prettify bool) 	([]byte, error)
	Import(file string, format string) 	(model.Project, error)
}

func (info *ProjectInfo) Read() (model.Project, error) {
	var project model.Project
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.ProjectId
	os.MkdirAll(folder, 0666)
	fileName := folder + "/project.ser"
	err := project.Load(fileName)
	return  project, err
}

func (info *ProjectInfo) Write(project model.Project) bool {
	home := model.VMBaseFolder()
	folder := home + "/metadata/" + info.ProjectId
	os.MkdirAll(folder, 0666)
	fileName := folder + "/project.ser"
	err := project.Save(fileName)
	return err == nil
}

func (info *ProjectInfo) Import(file string, format string) (model.Project, error) {
	var project model.Project
	err := project.Import(file, format)
	return  project, err
}

func (info *ProjectInfo) Export(project model.Project, prettify bool) ([]byte, error) {
	if "json" == info.Format {
		return  GetJSONFromObj(project, prettify)
	} else if "xml" == info.Format {
		return  GetXMLFromObj(project, prettify)
	} else {
		return  []byte{}, errors.New("Format type : "+info.Format+" not known ...")
	}
}
