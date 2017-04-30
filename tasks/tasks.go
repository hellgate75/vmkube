package tasks

import (
	"fmt"
	"time"
)

type ScheduleTask struct {
	Id      string
	Jobs    []Job
	Active  bool
	Working bool
	LastIndex	int
}

type Job struct {
	Id       	string
	Name     	string
	Runnable 	MachineOperationsJob
	Async    	bool
	Sequence	int
	Of				int
}

func (job *Job) Init() {
}


func (job *Job) Run() {
	job.Runnable.Start()
}

func (job *Job) IsRunning() bool {
	return job.Runnable.Status()
}


func (job *Job) WaitFor()  {
	job.Runnable.WaitFor()
}
func (job *Job) Abort() {
	job.Runnable.Stop()
}

type JobProcess interface {
	Run()
	IsRunning() bool
	Abort()
	Init()
	WaitFor()
}

type TaskProcess interface {
	Run()
	IsRunning() bool
	Abort()
	Init()
	Wait()
	Deactivate()
}

func (task *ScheduleTask) Init() {
	task.Active = true
	task.Working = true
	task.LastIndex = 0
}


func (task *ScheduleTask) Deactivate() {
	DumpData("threads-request.txt", fmt.Sprintf("Request Deactivate for Task Id : %s [Sequece : %d] Job: %s Running : %t Last Index : %d Task Len : %d", task.Id, task.LastIndex, task.Jobs[task.LastIndex].Id, task.Jobs[task.LastIndex].IsRunning(), task.LastIndex,len(task.Jobs)), false)
	task.Active = false
	task.Working = false
	task.LastIndex = len(task.Jobs)
	task.Abort()
	DumpData("threads-request.txt", fmt.Sprintf("Post Deactivate for Task Id : %s [Sequece : %d] Last Index : %d Task Len : %d", task.Id, task.LastIndex, task.LastIndex,len(task.Jobs)), false)
}

func (task *ScheduleTask) Execute() {
	defer task.Deactivate()
	var index int = task.LastIndex
	for i := index; i < len(task.Jobs); i++ {
		//DumpData("threads-task.txt", fmt.Sprintf("Start Id : %s Job: %s Sequence: %d Async: %t", task.Id, task.Jobs[i].Id, i, task.Jobs[i].Async), false)

		task.Jobs[i].Init()
		task.Jobs[i].Sequence=i
		task.Jobs[i].Of=len(task.Jobs)
		task.LastIndex=i
		if task.Jobs[i].Async {
			go task.Jobs[i].Run()
			task.Jobs[i].WaitFor()
		} else {
			task.Jobs[i].Run()
		}
		//DumpData("threads-task.txt", fmt.Sprintf("Complete Id : %s Job: %s Sequence: %d Async: %t", task.Id, task.Jobs[i].Id, i, task.Jobs[i].Async), false)
		if ! task.Active || task.Jobs[i].Runnable.IsError() {
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
		//DumpData("threads-task.txt", fmt.Sprintf("Post Release Id : %s Job: %s Sequence: %d Count: %d", task.Id, task.Jobs[i].Id, i, task.Count), false)
	}
	//DumpData("threads-complete.txt", fmt.Sprintf("Id : %s Count: %d Active: %t", task.Id, task.Count, task.Active), false)
}

func (task *ScheduleTask) IsRunning() bool {

	//DumpData("threads-request.txt", fmt.Sprintf("Request IsRunning for Task Id : %s  Last Index : %d  Len : %d  Active : %t", task.Id,task.LastIndex, len(task.Jobs), task.Active), false)
	if task.Active || (task.Working && task.LastIndex < len(task.Jobs)) {
		return true
	}
	if task.LastIndex < len(task.Jobs)  {
		for i := task.LastIndex; i < len(task.Jobs); i++ {
			//DumpData("threads-request.txt", fmt.Sprintf("Request IsRunning for Task Id : %s [Sequece : %d] Job: %s Running : %t", task.Id, task.LastIndex, task.Jobs[task.LastIndex].Id, task.Jobs[task.LastIndex].IsRunning()), false)
			if task.Jobs[task.LastIndex].IsRunning() {
				return true
			}
		}
	}
	return false
}

func (task *ScheduleTask) Abort() {
	if task.LastIndex < len(task.Jobs)  {
		//DumpData("threads-request.txt", fmt.Sprintf("Request Abort for Task Id : %s [Sequece : %d] Job: %s Running : %t Last Index : %d Task Len : %d", task.Id, task.LastIndex, task.Jobs[task.LastIndex].Id, task.Jobs[task.LastIndex].IsRunning(), task.LastIndex,len(task.Jobs)), false)
		for i := task.LastIndex; i < len(task.Jobs); i++ {
			task.Jobs[i].Abort()
		}
	}
	task.Active = false
	task.LastIndex = len(task.Jobs)
	task.Working = false
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

type SchedulerState struct {
	Active        bool
	Paused        bool
	Pool          []ScheduleTask
}
