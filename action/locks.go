package action

import (
	"vmkube/model"
	"sync"
	"time"
)

type IFaceProjectActionIndex ProjectActionIndex

type IFaceRollBackIndex RollBackIndex


type IFaceRollBackSegment RollBackSegment

func (iFace *IFaceProjectActionIndex) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func(){
		index := ProjectActionIndex{
			Id: iFace.Id,
			ProjectId: iFace.ProjectId,
		}
		for{
			if IsActionIndexLocked(index) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func (iFace *IFaceRollBackIndex) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func(){
		index := RollBackIndex{
			Id: iFace.Id,
			ProjectId: iFace.ProjectId,
		}
		for{
			if IsRollBackIndexLocked(index) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func (iFace *IFaceRollBackSegment) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func(){
		index := RollBackSegment{
			Id: iFace.Id,
			ProjectId: iFace.ProjectId,
		}
		for{
			if UnlockRollBackSegment(index) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func LockActionIndex(index ProjectActionIndex) bool {
	return model.WriteLock(index.ProjectId, index.Id)
}

func UnlockActionIndex(index ProjectActionIndex) bool {
	return model.RemoveLock(index.ProjectId, index.Id)
}

func IsActionIndexLocked(index ProjectActionIndex) bool {
	return model.HasLock(index.ProjectId, index.Id)
}

func LockRollBackIndex(index RollBackIndex) bool {
	return model.WriteLock(index.ProjectId, index.Id)
}

func UnlockRollBackIndex(index RollBackIndex) bool {
	return model.RemoveLock(index.ProjectId, index.Id)
}

func IsRollBackIndexLocked(index RollBackIndex) bool {
	return model.HasLock(index.ProjectId, index.Id)
}
func LockRollBackSegment(index RollBackSegment) bool {
	return model.WriteLock(index.ProjectId, index.Id)
}

func UnlockRollBackSegment(index RollBackSegment) bool {
	return model.RemoveLock(index.ProjectId, index.Id)
}

func IsRollBackSegmentLocked(index RollBackSegment) bool {
	return model.HasLock(index.ProjectId, index.Id)
}
