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
}

type Job struct {
	Id       string
	Name     string
	Runnable operations.RunnableStruct
}

func (job *Job) Run() {
	job.Runnable.Start()
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
}

func (task *ScheduleTask) Execute() {
	task.Active = true
	for _,job := range task.Jobs {
		go job.Run()
		if ! task.Active {
			break
		}
	}
}

func (task *ScheduleTask) IsRunning() bool {
	if task.Active {
		running := false
		for _,job := range task.Jobs {
			if job.IsRunning() {
				running = true
			}
		}
		task.Active = running
		return task.Active
	}
	for _,job := range task.Jobs {
		if job.IsRunning() {
			return true
		}
	}
	return false
}

func (task *ScheduleTask) Abort() {
	for _,job := range task.Jobs {
		job.Abort()
	}
}

type SchedulerState struct {
	Active        bool
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
			//log.Println("Max Parallel Processes = " + strconv.Itoa(threads))
			if runtime.NumCPU() < threads + 1 {
				runtime.GOMAXPROCS(threads + 1)
			}
			//pool.State.Pool = make([]ScheduleTask, 0)
			
			for pool.State.Active {
				if threads > len(pool.State.Pool) {
					//log.Println("Waiting for New Task ...")
					Task := <- pool.Tasks
					if Task.Id != "" {
						if Task.Id == "<close>" {
							//log.Println("Pool manager exits on request ...")
							break
						} else {
							pool.WG.Add(1)
							//log.Println("Pool Append Task Id : " + pool.Id)
							mutex.Lock()
							pool.State.Pool = append(pool.State.Pool, Task)
							go pool.State.Pool[len(pool.State.Pool)-1].Execute()
							mutex.Unlock()
							//go Task.Execute()
						}
					} else {
						time.Sleep(1000*time.Millisecond)
					}
				} else {
					//Thread Pool Full
					time.Sleep(1000*time.Millisecond)
					i := 0
					//log.Println(fmt.Sprintf("Pool full - Removing completed task from : %s in pool ...", len(pool.Pool)))
					count := 0
					for i < len(pool.State.Pool) {
						if ! pool.State.Pool[i].IsRunning() {
							count ++
							if pool.PostExecute {
								pool.Callback(pool.State.Pool[i])
							}
							if i > 0 && i < len(pool.State.Pool) - 1 {
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
					//log.Println(fmt.Sprintf("Pool clean - Removed %d completed tasks!!", count))
				}
			}
			for _,task := range pool.State.Pool {
				if task.IsRunning() {
					task.Abort()
					for task.IsRunning() {
						time.Sleep(1000*time.Millisecond)
					}
				}
				if pool.PostExecute {
					pool.Callback(task)
				}
				pool.WG.Done()
			}
			callback()
			if pool.KeepAlive {
				pool.WG.Done()
			}
			
		}()
	}
}

func (pool *SchedulerPool) IsRunning() bool {
	if pool.State.Active {
		return true
	}
	for _,task := range pool.State.Pool {
		if task.IsRunning() {
			for task.IsRunning() {
				return true
			}
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
