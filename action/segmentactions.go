package action

import "time"

const MAX_SEGMENT_SIZE int = 1000

func AddProjectChangeActions(projectId string,actions ...ActionDescriptor) error {
	actionIndex, err := LoadProjectActionChanges(projectId)
	if err != nil {
		return err
	}

	iFaceIndex := IFaceProjectActionIndex{
		ProjectId: projectId,
		Id: actionIndex.Id,
	}
	iFaceIndex.WaitForUnlock()

	LockActionIndex(actionIndex)

	actionIndex.Actions = append(actionIndex.Actions, actions...)

	err = SaveProjectActionIndex(actionIndex)

	UnlockActionIndex(actionIndex)

	return err
}

func AddRollBackChangeActions(projectId string,actions ...ActionDescriptor) error {
	rollbackIndex, err := LoadProjectRollBackIndex(projectId)

	if err != nil {
		return err
	}
	
	rollbackIFace := IFaceRollBackIndex {
		Id: rollbackIndex.Id,
		ProjectId: rollbackIndex.ProjectId,
	}
	
	latestIndex := RollBackSegmentIndex{}
	latestIndex.New()
	isNewIndex := false
	
	var segment RollBackSegment

	if len(rollbackIndex.IndexList) == 0 {
		firstSegment := RollBackSegment{}
		//Creating new segment index
		FirstSegmentIndex := RollBackSegmentIndex{}
		FirstSegmentIndex.New()
		//Creating new segment
		firstSegment.Storage = []ActionStorage{}
		firstSegment.Size = 0
		firstSegment.Id = NewUUIDString()
		firstSegment.Index = FirstSegmentIndex
		rollbackIndex.IndexList = append(rollbackIndex.IndexList, FirstSegmentIndex)
		rollbackIFace.WaitForUnlock()
		LockRollBackIndex(rollbackIndex)
		err = SaveRollbackIndex(rollbackIndex)
		UnlockRollBackIndex(rollbackIndex)
		if err != nil {
			return err
		}
		LockRollBackSegment(firstSegment)
		err = SaveRollbackSegment(firstSegment)
		UnlockRollBackSegment(firstSegment)
		if err != nil {
			return err
		}
		latestIndex = FirstSegmentIndex
		segment = firstSegment
	} else {
		latestIndex = rollbackIndex.IndexList[len(rollbackIndex.IndexList)-1]
	}

	if ! isNewIndex {
		segment, err = LoadProjectRollBackSegment(projectId, latestIndex)
		if err != nil {
			return err
		}
	}
	
	
	newSegmentRequired := false
	actionsInRollback1 := actions
	var actionsInRollback2 []ActionDescriptor = make([]ActionDescriptor, 0)
	var newSegment RollBackSegment = RollBackSegment{}
	
	if segment.Size + len(actions) > MAX_SEGMENT_SIZE {
		newSegmentRequired = true
		if segment.Size < MAX_SEGMENT_SIZE {
			maxSegmentDiff := MAX_SEGMENT_SIZE - segment.Size
			actionsInRollback1 = actions[:maxSegmentDiff]
			actionsInRollback2 = actions[maxSegmentDiff:]
		} else {
			actionsInRollback1 = make([]ActionDescriptor, 0)
			actionsInRollback2 = actions
		}
	}
	
	if len(actionsInRollback1) > 0 {
		for _,action := range actionsInRollback1 {
			segment.Storage = append(segment.Storage, ActionStorage{
				Action: action,
				Date: time.Now(),
				Id: NewUUIDString(),
			})
		}
		segment.Size += len(actionsInRollback1)
		//Save updated rollback segment
		LockRollBackSegment(segment)
		err = SaveRollbackSegment(segment)
		UnlockRollBackSegment(segment)
		if err != nil {
			return err
		}
		
	}
	
	if newSegmentRequired {
		//Creating new segment index
		NewSegmentIndex := RollBackSegmentIndex{}
		NewSegmentIndex.NewNextFrom(segment.Index.Index)
		//Creating new segment
		for _,action := range actionsInRollback2 {
			newSegment.Storage = append(newSegment.Storage, ActionStorage{
				Action: action,
				Date: time.Now(),
				Id: NewUUIDString(),
			})
		}
		newSegment.Size = len(actionsInRollback2)
		newSegment.Id = NewUUIDString()
		newSegment.Index = NewSegmentIndex
		// Updating main index
		rollbackIndex.IndexList = append(rollbackIndex.IndexList, NewSegmentIndex)
		rollbackIFace.WaitForUnlock()
		LockRollBackIndex(rollbackIndex)
		err = SaveRollbackIndex(rollbackIndex)
		UnlockRollBackIndex(rollbackIndex)
		if err != nil {
			return err
		}
		// Saving new rollback segment
		LockRollBackSegment(newSegment)
		err = SaveRollbackSegment(newSegment)
		UnlockRollBackSegment(newSegment)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadProjectActionChanges(projectId string) (ProjectActionIndex, error) {
	return LoadProjectActionIndex(projectId)
}

func LoadProjectRollBackIndex(projectId string) (RollBackIndex, error) {
	return LoadRollbackIndex(projectId)
}

func LoadProjectRollBackSegment(projectId string, index RollBackSegmentIndex) (RollBackSegment, error) {
	return LoadRollbackSegment(projectId, index)
}
