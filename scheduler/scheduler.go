package scheduler


import (
	"runtime"
	"time"
	"sync"
	"vmkube/operations"
)

type ScheduleTask struct {
	Id      string
	Jobs    []Job
	Active  bool
	Count   int
}

type Job struct {
	Id       string
	Name     string
	Runnable operations.RunnableStruct
	Async    bool
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
}

func (task *ScheduleTask) Init() {
	task.Count = 0
}


func (task *ScheduleTask) deactivate() {
	task.Active = false
	task.Count = len(task.Jobs)
}

func (task *ScheduleTask) Execute() {
	task.Active = true
	var index int = task.Count
	for i := index; i < len(task.Jobs); i++ {
		//DumpData("threads-task.txt", fmt.Sprintf("Start Id : %s Job: %s Sequence: %d Async: %t", task.Id, task.Jobs[i].Id, i, task.Jobs[i].Async), false)
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
		}
		if i < len(task.Jobs) - 1 {
			time.Sleep(1*time.Second)
		}
		task.Count=i+1
		if i < len(task.Jobs) - 1 {
			time.Sleep(1*time.Second)
		}
		//DumpData("threads-task.txt", fmt.Sprintf("Post Release Id : %s Job: %s Sequence: %d Count: %d", task.Id, task.Jobs[i].Id, i, task.Count), false)
	}
	task.Wait()
	task.deactivate()
	//DumpData("threads-complete.txt", fmt.Sprintf("Id : %s Count: %d Active: %t", task.Id, task.Count, task.Active), false)
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
	for i := task.Count; i < len(task.Jobs); i++ {
		task.Jobs[i].Abort()
	}
	task.Active = false
	task.Count = len(task.Jobs)
}
func (task *ScheduleTask) Wait() {
	for i := task.Count; i < len(task.Jobs); i++ {
		for task.Jobs[i].IsRunning() {
			time.Sleep(1*time.Second)
		}
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
		var threads = pool.MaxParallel
		if threads == 0 {
			threads = runtime.NumCPU()
		}
		var Buffer []ScheduleTask = make([]ScheduleTask, 0)
		pumperExited := false
		// start Pool enqueue manager
		go func() {
			for pool.State.Active {
				if  ! pool.State.Paused {
					// Try add jobs from Buffer
					for threads > len(pool.State.Pool) && len(Buffer) > 0 {
						Task := Buffer[0]
						Buffer = Buffer[1:]
						pool.WG.Add(1)
						pool.State.Pool = append(pool.State.Pool, Task)
						pool.State.Pool[len(pool.State.Pool)-1].Init()
						go pool.State.Pool[len(pool.State.Pool)-1].Execute()
					}
					// Look for completed jobs to remove from Pool
					var i int = 0
					for len(pool.State.Pool) > 0 && i < len(pool.State.Pool) {
						if ! pool.State.Pool[i].IsRunning() {
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
				}
			}
			pumperExited = true
		}()
		// start jobs
		go func() {
			//if runtime.NumCPU() < threads + 1 {
			//	runtime.GOMAXPROCS(threads + 1)
			//}
			for pool.State.Active {
				if ! pool.State.Paused {
						Task := <- pool.Tasks
						if Task.Id != "" {
							if Task.Id == "<close>" {
								break
							} else {
								Buffer = append(Buffer, Task)
							}
						}
				} else {
					time.Sleep(1500*time.Millisecond)
				}
			}
			for ! pumperExited {
				time.Sleep(500*time.Millisecond)
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
