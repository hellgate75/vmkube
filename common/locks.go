package common

import (
	"strconv"
	"sync"
	"time"
	"github.com/hellgate75/vmkube/model"
)

type IFaceProjectActionIndex ProjectActionIndex

type IFaceRollBackIndex RollBackIndex

type IFaceRollBackSegment RollBackSegment

type IFaceLogStorage model.LogStorage

func (iFace *IFaceLogStorage) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		index := model.LogStorage{
			ProjectId: iFace.ProjectId,
			InfraId:   iFace.InfraId,
			ElementId: iFace.ElementId,
		}
		for {
			if IsLogLocked(index) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func (iFace *IFaceLogStorage) WaitForLogFileUnlock(logIndex int) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		index := model.LogStorage{
			ProjectId: iFace.ProjectId,
			InfraId:   iFace.InfraId,
			ElementId: iFace.ElementId,
		}
		for {
			if IsLogFileLocked(index, logIndex) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func (iFace *IFaceProjectActionIndex) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		index := ProjectActionIndex{
			Id:        iFace.Id,
			ProjectId: iFace.ProjectId,
		}
		for {
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
	go func() {
		index := RollBackIndex{
			Id:        iFace.Id,
			ProjectId: iFace.ProjectId,
		}
		for {
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
	go func() {
		index := RollBackSegment{
			Id:        iFace.Id,
			ProjectId: iFace.ProjectId,
			Index:     iFace.Index,
		}
		for {
			if IsRollBackSegmentLocked(index) {
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
	return model.WriteLock(index.ProjectId, index.Index.Index.Value)
}

func UnlockRollBackSegment(index RollBackSegment) bool {
	return model.RemoveLock(index.ProjectId, index.Index.Index.Value)
}

func LockRollBackSegmentById(projectId string, index RollBackSegmentIndex) bool {
	return model.WriteLock(projectId, index.Index.Value)
}

func UnlockRollBackSegmentById(projectId string, index RollBackSegmentIndex) bool {
	return model.RemoveLock(projectId, index.Index.Value)
}

func IsRollBackSegmentLockedById(projectId string, index RollBackSegmentIndex) bool {
	return model.HasLock(projectId, index.Index.Value)
}

func IsRollBackSegmentLocked(index RollBackSegment) bool {
	return model.HasLock(index.ProjectId, index.Index.Index.Value)
}

func IsLogLocked(index model.LogStorage) bool {
	return model.HasLock(index.ProjectId, index.ElementId)
}
func LockLog(index model.LogStorage) bool {
	return model.WriteLock(index.ProjectId, index.ElementId)
}

func UnlockLog(index model.LogStorage) bool {
	return model.RemoveLock(index.ProjectId, index.ElementId)
}

func IsLogFileLocked(index model.LogStorage, logIndex int) bool {
	return model.HasLock(index.ProjectId, index.ElementId+"-"+strconv.Itoa(logIndex))
}
func LockLogFile(index model.LogStorage, logIndex int) bool {
	return model.WriteLock(index.ProjectId, index.ElementId+"-"+strconv.Itoa(logIndex))
}

func UnlockLogFile(index model.LogStorage, logIndex int) bool {
	return model.RemoveLock(index.ProjectId, index.ElementId+"-"+strconv.Itoa(logIndex))
}
