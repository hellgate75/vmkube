package scheduler


import (
	"runtime"
	"time"
	"sync"
	//"log"
	//"strconv"
	//"fmt"
	"vmkube/operations"
)

type ScheduleTask struct {
	Id      string
	Jobs    []Job
	Active  bool
	Count    int
}

type Job struct {
	Id       string
	Name     string
	Runnable operations.RunnableStruct
	Async    bool
}

func (job *Job) Run() {
	job.Runnable.Start()
	if job.Async {
		job.Runnable.WaitFor()
	}
}

func (job *Job) IsRunning() bool {
	return job.Runnable.Status()
}
func (job *Job) Abort() {
	job.Runnable.Stop()
}

type JobProcess interface {
	Run()
	IsRunning() bool
	Abort()
	Init()
}


func (task *ScheduleTask) Init() {
	task.Count = 0
}


func (task *ScheduleTask) deactivate() {
	task.Active = false
}

func (task *ScheduleTask) Execute() {
	task.Active = true
	defer task.deactivate()
	for i := task.Count; i < len(task.Jobs); i++ {
		task.Jobs[i].Run()
		if ! task.Active || task.Jobs[i].Runnable.IsError() {
			task.Abort()
			break
		}
		task.Count++
		if i < len(task.Jobs) - 1 {
			time.Sleep(3*time.Second)
		}
	}
}

func (task *ScheduleTask) IsRunning() bool {
	if task.Active || task.Count < len(task.Jobs) {
		return true
	}
	for i := task.Count; i < len(task.Jobs); i++ {
		if task.Jobs[i].IsRunning() {
			return true
		}
	}
	return false
}

func (task *ScheduleTask) Abort() {
	task.Active = false
	task.Count = len(task.Jobs)
	for i := task.Count; i < len(task.Jobs); i++ {
		task.Jobs[i].Abort()
	}
}

type SchedulerState struct {
	Active        bool
	Paused        bool
	Pool          []ScheduleTask
}

type SchedulerPool struct {
	Id            string
	Tasks         chan ScheduleTask
	MaxParallel   int
	WG            sync.WaitGroup
	KeepAlive     bool
	PostExecute   bool
	Callback      func(task ScheduleTask)
	State         SchedulerState
}

func (pool *SchedulerPool) RegisterWaitGroup(wg sync.WaitGroup) {
	pool.WG = wg
	pool.State.Active = false
	pool.State.Pool = []ScheduleTask {}
}

func (pool *SchedulerPool) Init() {
	pool.State.Paused = false
	pool.State.Active = false
	pool.Tasks = make(chan ScheduleTask)
	if pool.KeepAlive {
		pool.WG.Add(1)
	}
}

func (pool *SchedulerPool) Start(callback func()) {
	if ! pool.State.Active {
		pool.State.Active = true
		//start jobs
		go func() {
			var mutex sync.Mutex
			var threads = pool.MaxParallel
			if threads == 0 {
				threads = runtime.NumCPU() - 1
			}
			if runtime.NumCPU() < threads + 1 {
				runtime.GOMAXPROCS(threads + 1)
			}
			for pool.State.Active {
				if ! pool.State.Paused {
					if threads > len(pool.State.Pool) {
						Task := <- pool.Tasks
						if Task.Id != "" {
							if Task.Id == "<close>" {
								break
							} else {
								mutex.Lock()
								pool.WG.Add(1)
								pool.State.Pool = append(pool.State.Pool, Task)
								go pool.State.Pool[len(pool.State.Pool)-1].Execute()
								mutex.Unlock()
							}
						} else {
							time.Sleep(500*time.Millisecond)
						}
					} else {
						//Thread Pool Full
						time.Sleep(1000*time.Millisecond)
						i := 0
						count := 0
						mutex.Lock()
						for i < len(pool.State.Pool) {
							if ! pool.State.Pool[i].IsRunning() {
								count ++
								if pool.PostExecute {
									go pool.Callback(pool.State.Pool[i])
								}
								if i > 0 && i < len(pool.State.Pool) - 2 {
									pool.State.Pool = pool.State.Pool[:i]
									pool.State.Pool = append(pool.State.Pool, pool.State.Pool[i+1:]...)
								} else if i == 0 {
									pool.State.Pool = pool.State.Pool[i+1:]
								}  else {
									pool.State.Pool = pool.State.Pool[:i]
								}
								pool.WG.Done()
							} else {
								i++
							}
						}
						mutex.Unlock()
					}
				} else {
					time.Sleep(1500*time.Millisecond)
				}
			}
			for i :=0; i<len(pool.State.Pool); i++ {
				if pool.State.Pool[i].IsRunning() {
					pool.State.Pool[i].Abort()
					for pool.State.Pool[i].IsRunning() {
						time.Sleep(1000*time.Millisecond)
					}
				}
				pool.WG.Done()
				if pool.PostExecute {
					go pool.Callback(pool.State.Pool[i])
				}
			}
			go callback()
			if pool.KeepAlive {
				pool.WG.Done()
			}
			
		}()
	}
}

func (pool *SchedulerPool) IsRunning() bool {
	if pool.State.Active && ! pool.State.Paused {
		return true
	}
	for _,task := range pool.State.Pool {
		if task.IsRunning() {
				return true
		}
	}
	return false
}

func (pool *SchedulerPool) IsWorking() bool {
	for i := 0; i < len(pool.State.Pool); i++ {
		if pool.State.Pool[i].IsRunning() {
			return true
		}
	}
	return false
}

func (pool *SchedulerPool) NumberOfWorkers() int {
	var counter int = 0
	for i := 0; i < len(pool.State.Pool); i++ {
		if pool.State.Pool[i].IsRunning() {
			counter++
		}
	}
	return counter
}

func (pool *SchedulerPool) IsJobActive(id string) bool {
	for i := 0; i < len(pool.State.Pool); i++ {
		if pool.State.Pool[i].Id == id {
			return pool.State.Pool[i].IsRunning()
		}
	}
	return false
}

func (pool *SchedulerPool) Stop() {
	pool.State.Active = false
	pool.Tasks <- ScheduleTask{
		Id: "<close>",
		Jobs: []Job{},
	}
	close(pool.Tasks)
	
}

func (pool *SchedulerPool) Pause() {
	pool.State.Paused = true
}

func (pool *SchedulerPool) IsPaused() bool {
	return pool.State.Paused
}

func (pool *SchedulerPool) Resume() {
	pool.State.Paused = false
}

func (pool *SchedulerPool) Interrupt() {
	for i := 0; i < len(pool.State.Pool); i++ {
			pool.State.Pool[i].Abort()
	}
}
