package tasks

import (
	"time"
	"vmkube/state"
	"sync"
)


type SchedulerState struct {
	Active        bool
	Paused        bool
	Pool          []ScheduleTask
}

type JobProcess interface {
	Run()
	IsRunning() bool
	IsAsync() bool
	HasErrors() bool
	Abort()
	Init(Sequence int, Global int)
	WaitFor()
	GetRunnable() RunnableStruct
}

type Job struct {
	Id       	string
	Name     	string
	Runnable 	RunnableStruct
	Async    	bool
	Sequence	int
	Of				int
}

func (job *Job) Init(Sequence int, Global int) {
	job.Sequence = Sequence
	job.Of = Global
}


func (job *Job) Run() {
	job.Runnable.Start()
}

func (job *Job) HasErrors() bool {
	return job.Runnable.IsError()
}

func (job *Job) IsAsync() bool {
	return job.Async
}

func (job *Job) IsRunning() bool {
	return job.Runnable.Status()
}

func (job *Job) GetRunnable() RunnableStruct {
	return job.Runnable
}

func (job *Job) WaitFor()  {
	job.Runnable.WaitFor()
}
func (job *Job) Abort() {
	job.Runnable.Stop()
}

type TaskProcess interface {
	Run()
	IsRunning() bool
	Abort()
	Init(context *state.StateContext)
	Wait()
	Deactivate()
}

type ScheduleTask struct {
	Id      	string
	Jobs    	[]JobProcess
	Active  	bool
	Working 	bool
	LastIndex	int
	State			*state.StateContext
}

func (task *ScheduleTask) Init(context *state.StateContext) {
	task.Active = true
	task.Working = true
	task.LastIndex = 0
	task.State = context
}


var TaskMutex sync.RWMutex

func readTaskContextState(state state.StateContext, taskId string) bool {
	TaskMutex.RLock()
	IsRegistered := state.HasValue(taskId)
	TaskMutex.RUnlock()
	if IsRegistered {
		defer  TaskMutex.RUnlock()
		TaskMutex.RLock()
		return state.State(taskId)
	}
	return true

}

func writeTaskContextState(state state.StateContext, taskId string, status state.StateReferenceData) {
	TaskMutex.Lock()
	defer  TaskMutex.Unlock()
	state.Collect(taskId) <- status
}


func (task *ScheduleTask) Deactivate() {
	task.Active = false
	task.Working = false
	task.LastIndex = len(task.Jobs)
	time.Sleep(1*time.Second)
	writeTaskContextState(*(task.State), task.Id, state.StateReferenceData{
		Id: task.Id,
		Status: false,
	})
}

func (task *ScheduleTask) Execute() {
	writeTaskContextState(*(task.State), task.Id, state.StateReferenceData{
		Id: task.Id,
		Status: true,
	})
	defer task.Abort()
	var index int = task.LastIndex
	for i := index; i < len(task.Jobs); i++ {
		task.Jobs[i].Init(i, len(task.Jobs))
		task.LastIndex=i
		if task.Jobs[i].IsAsync() {
			go task.Jobs[i].Run()
			task.Jobs[i].WaitFor()
		} else {
			task.Jobs[i].Run()
		}
		if ! task.Active || task.Jobs[i].HasErrors() {
			task.Abort()
			task.Wait()
			break
		}
		if task.Jobs[i].IsRunning() {
			task.Jobs[i].Abort()
			task.Wait()
		}
		if i < len(task.Jobs) - 1 {
			time.Sleep(1*time.Second)
		}
		if i < len(task.Jobs) - 1 {
			time.Sleep(1*time.Second)
		}
	}
}

func (task *ScheduleTask) IsRunning() bool {
	return readTaskContextState(*(task.State), task.Id)
}

func (task *ScheduleTask) Abort() {
	time.Sleep(1*time.Second)
	if task.LastIndex < len(task.Jobs)  {
		for i := task.LastIndex; i < len(task.Jobs); i++ {
			task.Jobs[i].Abort()
		}
	}
	task.Deactivate()
}
func (task *ScheduleTask) Wait() {
	if task.LastIndex < len(task.Jobs) {
		for i := task.LastIndex; i < len(task.Jobs); i++ {
			for task.Jobs[i].IsRunning() {
				time.Sleep(1*time.Second)
			}
		}
	}
}
