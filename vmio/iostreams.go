package vmio

import (
	"io/ioutil"
	"vmkube/model"
)

func SaveIndex(index model.ProjectsIndex) error {
	info := ProjectIndexInfo{
		Format: "",
		Index: index,
	}
	return info.Write()
}

func LoadIndex() (model.ProjectsIndex, error) {
	index := model.ProjectsIndex{
		Projects: []model.ProjectsDescriptor{},
	}
	info := ProjectIndexInfo{
		Format: "",
		Index: index,
	}
	err := info.Read()
	return index, err
}

func ImportIndex(file string, format string) (model.ProjectsIndex, error) {
	index := model.ProjectsIndex{
		Projects: []model.ProjectsDescriptor{},
	}
	info := ProjectIndexInfo{
		Format: format,
		Index: index,
	}
	err := info.Import(file, format)
	return index, err
}

func ExportIndex(index model.ProjectsIndex, file string, format string, prettify bool) error {
	info := ProjectIndexInfo{
		Format: format,
		Index: index,
	}
	bytes, err := info.Export(prettify)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	return err
}

func SaveProject(project model.Project) error {
	info := ProjectInfo{
		Format: "",
		Project: project,
	}
	return info.Write()
}

func LoadProject(id string) (model.Project, error) {
	project := model.Project{
		Id: id,
	}
	info := ProjectInfo{
		Format: "",
		Project: project,
	}
	err := info.Read()
	return project, err
}

func ImportProject(file string, format string) (model.Project, error) {
	project := model.Project{
	}
	info := ProjectInfo{
		Format: format,
		Project: project,
	}
	err := info.Import(file, format)
	return project, err
}

func ImportUserProject(file string, format string) (model.Project, error) {
	imported := model.ProjectImport{
	}
	info := ProjectImportInfo{
		Format: format,
		ProjectImport: imported,
	}
	err := info.Import(file, format)
	if err != nil {
		return err
	}
	project := model.ProjectFromImport(imported)
	return project, err
}

func ExportProject(project model.Project, file string, format string, prettify bool) error {
	info := ProjectInfo{
		Format: format,
		Project: project,
	}
	bytes, err := info.Export(prettify)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	return err
}

func ExportUserProject(project model.Project, file string, format string, prettify bool) error {
	info := ProjectImportInfo{
		Format: format,
		ProjectImport: model.ProjectToImport(project),
	}
	bytes, err := info.Export(prettify)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	return err
}

func SaveInfrastructure(infrastructure model.Infrastructure) error {
	info := InfrastructureInfo{
		Format: "",
		Infra: infrastructure,
	}
	return info.Write()
}

func LoadInfrastructure(projectId string) (model.Infrastructure, error) {
	infrastructure := model.Infrastructure{
		ProjectId: projectId,
	}
	info := InfrastructureInfo{
		Format: "",
		Infra: infrastructure,
	}
	err := info.Read()
	return infrastructure, err
}

func ImportInfrastructure(file string, format string) (model.Infrastructure, error) {
	infrastructure := model.Infrastructure{
	}
	info := InfrastructureInfo{
		Format: format,
		Infra: infrastructure,
	}
	err := info.Import(file, format)
	return infrastructure, err
}

func ExportInfrastructure(infrastructure model.Infrastructure, file string, format string, prettify bool) error {
	info := InfrastructureInfo{
		Format: format,
		Infra: infrastructure,
	}
	bytes, err := info.Export(prettify)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	return err
}
