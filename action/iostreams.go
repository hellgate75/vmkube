package action

func SaveProjectActionIndex(index ProjectActionIndex) error {
	info := ProjectActionIndexInfo{
		Format: "",
		Index: index,
	}
	return info.Write()
}

func LoadProjectActionIndex(projectId string) (ProjectActionIndex, error) {
	index := ProjectActionIndex{
		ProjectId: projectId,
	}
	info := ProjectActionIndexInfo{
		Format: "",
		Index: index,
	}
	err := info.Read()
	return info.Index, err
}


func SaveRollbackIndex(index RollBackIndex) error {
	info := ProjectRollbackIndexInfo{
		Format: "",
		Index: index,
	}
	return info.Write()
}

func LoadRollbackIndex(projectId string) (RollBackIndex, error) {
	index := RollBackIndex{
		ProjectId: projectId,
	}
	info := ProjectRollbackIndexInfo{
		Format: "",
		Index: index,
	}
	err := info.Read()
	return info.Index, err
}

func SaveRollbackSegment(index RollBackSegment) error {
	info := ProjectRollbackSegmentInfo{
		Format: "",
		Index: index,
	}
	return info.Write()
}

func LoadRollbackSegment(projectId string, rollbackIndex RollBackSegmentIndex) (RollBackSegment, error) {
	index := RollBackSegment{
		Index: rollbackIndex,
		ProjectId: projectId,
	}
	info := ProjectRollbackSegmentInfo {
		Format: "",
		Index: index,
	}
	err := info.Read()
	return info.Index, err
}
