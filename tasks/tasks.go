package tasks

import (
	"time"
	"vmkube/state"
)

func (job *Job) Init(Sequence int, Global int) {
	job.Sequence = Sequence
	job.Of = Global
}

func (job *Job) Run() {
	job.State = true
	var exitChannel chan bool = make(chan bool, 1)
	job.Runnable.Start(exitChannel)
	go func() {
		select {
		case <-exitChannel:
			job.State = false
		case <-time.After(time.Second * MachineReadOperationTimeout):
			job.State = false
		}
	}()
}

func (job *Job) HasErrors() bool {
	return job.Runnable.IsError()
}

func (job *Job) IsAsync() bool {
	return job.Async
}

func (job *Job) IsRunning() bool {
	return job.State
}

func (job *Job) GetRunnable() Runnable {
	return job.Runnable
}

func (job *Job) WaitFor() {
	job.Runnable.WaitFor()
}
func (job *Job) Abort() {
	job.Runnable.Stop()
}

func (task *SchedulerTask) Init(context *state.StateContext) {
	task.Active = true
	task.Working = true
	task.LastIndex = 0
	task.State = context
}

func (task *SchedulerTask) Deactivate() {
	task.Active = false
	task.Working = false
	task.LastIndex = len(task.Jobs)
	time.Sleep(1 * time.Second)
	writeTaskContextState(*(task.State), task.Id, state.StateReferenceData{
		Id:     task.Id,
		Status: false,
	})
}

func (task *SchedulerTask) Execute() {
	writeTaskContextState(*(task.State), task.Id, state.StateReferenceData{
		Id:     task.Id,
		Status: true,
	})
	defer task.Abort()
	var index int = task.LastIndex
	for i := index; i < len(task.Jobs); i++ {
		task.Jobs[i].Init(i, len(task.Jobs))
		task.LastIndex = i
		if task.Jobs[i].IsAsync() {
			go task.Jobs[i].Run()
			task.Jobs[i].WaitFor()
		} else {
			task.Jobs[i].Run()
		}
		if !task.Active || task.Jobs[i].HasErrors() {
			task.Abort()
			task.Wait()
			break
		}
		if task.Jobs[i].IsRunning() {
			task.Jobs[i].Abort()
			task.Wait()
		}
		if i < len(task.Jobs)-1 {
			time.Sleep(1 * time.Second)
		}
		if i < len(task.Jobs)-1 {
			time.Sleep(1 * time.Second)
		}
	}
}

func (task *SchedulerTask) IsRunning() bool {
	return readTaskContextState(*(task.State), task.Id)
}

func (task *SchedulerTask) Abort() {
	time.Sleep(1 * time.Second)
	if task.LastIndex < len(task.Jobs) {
		for i := task.LastIndex; i < len(task.Jobs); i++ {
			task.Jobs[i].Abort()
		}
	}
	task.Deactivate()
}
func (task *SchedulerTask) Wait() {
	if task.LastIndex < len(task.Jobs) {
		for i := task.LastIndex; i < len(task.Jobs); i++ {
			for task.Jobs[i].IsRunning() {
				time.Sleep(1 * time.Second)
			}
		}
	}
}
