package vmio

import "vmkube/model"

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

func LockInfrastructureById(projectId string, uid string) bool {
	return model.WriteLock(projectId, uid)
}

func UnlockInfrastructureById(projectId string, uid string) bool {
	return model.RemoveLock(projectId, uid)
}

func IsInfrastructureLockedById(projectId string, uid string) bool {
	return model.HasLock(projectId, uid)
}
