package vmio

import (
	"sync"
	"time"
	"github.com/hellgate75/vmkube/model"
)

type IFaceIndex model.ProjectsIndex
type IFaceProject model.Project
type IFaceInfra model.Infrastructure

func (iFace *IFaceIndex) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		index := model.ProjectsIndex{
			Id: iFace.Id,
		}
		for {
			if IsIndexLocked(index) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func (iFace *IFaceProject) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		project := model.Project{
			Id: iFace.Id,
		}
		for {
			if IsProjectLocked(project) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func (iFace *IFaceInfra) WaitForUnlock() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		infra := model.Infrastructure{
			Id:        iFace.Id,
			ProjectId: iFace.ProjectId,
		}
		for {
			if IsInfrastructureLocked(infra) {
				time.Sleep(500)
			} else {
				waitGroup.Done()
				break
			}
		}
	}()
	waitGroup.Wait()
}

func LockIndex(index model.ProjectsIndex) bool {
	return model.WriteLock(index.Id, index.Id)
}

func UnlockIndex(index model.ProjectsIndex) bool {
	return model.RemoveLock(index.Id, index.Id)
}

func IsIndexLocked(index model.ProjectsIndex) bool {
	return model.HasLock(index.Id, index.Id)
}

func LockProject(project model.Project) bool {
	return model.WriteLock(project.Id, project.Id)
}

func UnlockProject(project model.Project) bool {
	return model.RemoveLock(project.Id, project.Id)
}

func IsProjectLocked(project model.Project) bool {
	return model.HasLock(project.Id, project.Id)
}

func LockProjectById(uid string) bool {
	return model.WriteLock(uid, uid)
}

func UnlockProjectById(uid string) bool {
	return model.RemoveLock(uid, uid)
}

func IsProjectLockedById(uid string) bool {
	return model.HasLock(uid, uid)
}

func LockInfrastructure(infrastructure model.Infrastructure) bool {
	return model.WriteLock(infrastructure.ProjectId, infrastructure.Id)
}

func UnlockInfrastructure(infrastructure model.Infrastructure) bool {
	return model.RemoveLock(infrastructure.ProjectId, infrastructure.Id)
}

func IsInfrastructureLocked(infrastructure model.Infrastructure) bool {
	return model.HasLock(infrastructure.ProjectId, infrastructure.Id)
}

func LockInfrastructureById(projectId string, infraId string) bool {
	return model.WriteLock(projectId, infraId)
}

func UnlockInfrastructureById(projectId string, infraId string) bool {
	return model.RemoveLock(projectId, infraId)
}

func IsInfrastructureLockedById(projectId string, infraId string) bool {
	return model.HasLock(projectId, infraId)
}
