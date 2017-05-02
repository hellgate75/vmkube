package action

import (
	"vmkube/model"
)

func SaveProjectActionIndex(index ProjectActionIndex) error {
	info := ProjectActionIndexInfo{
		Format: "",
		Index:  index,
	}
	return info.Write()
}

func LoadProjectActionIndex(projectId string) (ProjectActionIndex, error) {
	index := ProjectActionIndex{
		ProjectId: projectId,
	}
	info := ProjectActionIndexInfo{
		Format: "",
		Index:  index,
	}
	err := info.Read()
	return info.Index, err
}

func DeleteActionIndex(projectId string) error {
	index := ProjectActionIndex{
		ProjectId: projectId,
	}
	info := ProjectActionIndexInfo{
		Format: "",
		Index:  index,
	}
	err := info.Delete()
	return err
}

func SaveRollbackIndex(index RollBackIndex) error {
	info := ProjectRollbackIndexInfo{
		Format: "",
		Index:  index,
	}
	return info.Write()
}

func LoadRollbackIndex(projectId string) (RollBackIndex, error) {
	index := RollBackIndex{
		ProjectId: projectId,
	}
	info := ProjectRollbackIndexInfo{
		Format: "",
		Index:  index,
	}
	err := info.Read()
	return info.Index, err
}

func DeleteRollbackIndex(projectId string) error {
	index, err := LoadRollbackIndex(projectId)
	if err != nil {
		return err
	}
	for _, segmentIndex := range index.IndexList {
		err = DeleteRollbackSegment(projectId, segmentIndex)
		if err != nil {
			return err
		}
	}
	info := ProjectRollbackIndexInfo{
		Format: "",
		Index:  index,
	}

	err = info.Delete()
	return err
}

func SaveRollbackSegment(index RollBackSegment) error {
	info := ProjectRollbackSegmentInfo{
		Format: "",
		Index:  index,
	}
	return info.Write()
}

func LoadRollbackSegment(projectId string, rollbackIndex RollBackSegmentIndex) (RollBackSegment, error) {
	index := RollBackSegment{
		Index:     rollbackIndex,
		ProjectId: projectId,
	}
	info := ProjectRollbackSegmentInfo{
		Format: "",
		Index:  index,
	}
	err := info.Read()
	return info.Index, err
}

func DeleteRollbackSegment(projectId string, rollbackIndex RollBackSegmentIndex) error {
	index := RollBackSegment{
		Index:     rollbackIndex,
		ProjectId: projectId,
	}
	info := ProjectRollbackSegmentInfo{
		Format: "",
		Index:  index,
	}
	err := info.Delete()
	return err
}

func SaveInfrastructureLogs(log model.LogStorage) error {
	info := InfrastructureLogsInfo{
		Format: "",
		Logs:   log,
	}
	return info.SaveLogFile()
}

func LoadInfrastructureLogs(projectId string, infraId string, elementId string) (model.LogStorage, error) {
	logs := model.LogStorage{
		InfraId:   infraId,
		ProjectId: projectId,
		ElementId: elementId,
		LogLines:  []string{},
	}
	info := InfrastructureLogsInfo{
		Format: "",
		Logs:   logs,
	}
	err := info.Read()
	return info.Logs, err
}

func DeleteInfrastructureLogs(log model.LogStorage) error {
	info := InfrastructureLogsInfo{
		Format: "",
		Logs:   log,
	}
	return info.Delete()
}

func DeleteInfrastructureLogsById(projectId string, infraId string, elementId string) error {
	logs := model.LogStorage{
		InfraId:   infraId,
		ProjectId: projectId,
		ElementId: elementId,
		LogLines:  []string{},
	}
	info := InfrastructureLogsInfo{
		Format: "",
		Logs:   logs,
	}
	return info.Delete()
}

func LoadInfrastructureLogFiles(log model.LogStorage) error {
	info := InfrastructureLogsInfo{
		Format: "",
		Logs:   log,
	}
	return info.ReadLogFiles()
}

func DeleteInfrastructureLogFiles(projectId string, infraId string, elementId string) error {
	info := InfrastructureLogsInfo{
		Format: "",
		Logs: model.LogStorage{
			InfraId:   infraId,
			ProjectId: projectId,
			ElementId: elementId,
			LogLines:  []string{},
		},
	}
	return info.ReadLogFiles()
}
