package scheduler


import (
	"runtime"
	"time"
	"strconv"
	"sync"
	"log"
	"reflect"
	"fmt"
)

type ScheduleTask struct {
	Id      string
	Jobs    []Job
	Active  bool
}

type Job struct {
	Id      string
	Name    string
	Struct  interface{}
}

func (job *Job) Run() {
	log.Println("Job ["+job.Id+"->"+job.Name+"] - Calling Start on interface implementation ...")
	valueType := reflect.ValueOf(job.Struct)
	method := valueType.MethodByName("Start")
	method.Call(nil)
}

func (job *Job) IsRunning() bool {
	log.Println("Job ["+job.Id+"->"+job.Name+"] - Calling Status on interface implementation ...")
	valueType := reflect.ValueOf(job.Struct)
	method := valueType.MethodByName("Status")
	ret := method.Call(nil)
	return ret[0].Bool()
}
func (job *Job) Abort() {
	log.Println("Job ["+job.Id+"->"+job.Name+"] - Calling Stop on interface implementation ...")
	valueType := reflect.ValueOf(job.Struct)
	method := valueType.MethodByName("Stop")
	method.Call(nil)
}

type JobProcess interface {
	Run()
	IsRunning() bool
	Abort()
}

func (task *ScheduleTask) Execute() {
	task.Active = true
	for _,job := range task.Jobs {
		job.Run()
		if ! task.Active {
			break
		}
	}
}

func (task *ScheduleTask) IsRunning() bool {
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


type SchedulerPool struct {
	Id            string
	Tasks         chan ScheduleTask
	MaxParallel   int
	Active        bool
	Pool          []ScheduleTask
	WG            sync.WaitGroup
	KeepAlive     bool
	PostExecute   bool
	Callback      func(task ScheduleTask)
}

func (pool *SchedulerPool) RegisterWaitGroup(wg sync.WaitGroup) {
	pool.WG = wg
}

func (pool *SchedulerPool) Init() {
	pool.Tasks = make(chan ScheduleTask)
	
}

func (pool *SchedulerPool) Start(callback func()) {
	if ! pool.Active {
		//start jobs
		go func(pool *SchedulerPool) {
			if pool.KeepAlive {
				pool.WG.Add(1)
			}
			var threads = pool.MaxParallel
			if threads == 0 {
				threads = runtime.NumCPU() - 1
			}
			log.Println("Max Parallel Processes = " + strconv.Itoa(threads))
			runtime.GOMAXPROCS(threads + 1)
			pool.Active = true
			pool.Pool = make([]ScheduleTask, 0)
			

			for pool.Active {
				if threads > len(pool.Pool) {
					log.Println("Waiting for New Task ...")
					Task := <- pool.Tasks
					if Task.Id != "" {
						if Task.Id == "<close>" {
							log.Println("Pool manager exits on request ...")
							break
						} else {
							pool.WG.Add(1)
							log.Println("Pool Append Task Id : " + pool.Id)
							pool.Pool = append(pool.Pool, Task)
							go Task.Execute()
						}
					} else {
						time.Sleep(1000*time.Millisecond)
					}
				} else {
					//Thread Pool Full
					time.Sleep(1000*time.Millisecond)
					i := 0
					log.Println(fmt.Sprintf("Pool full - Removing completed task from : %s in pool ...", len(pool.Pool)))
					count := 0
					for i < len(pool.Pool) {
						if ! pool.Pool[i].IsRunning() {
							count ++
							if pool.PostExecute {
								pool.Callback(pool.Pool[i])
							}
							if i > 0 && i < len(pool.Pool) - 1 {
								pool.Pool = pool.Pool[:i]
								pool.Pool = append(pool.Pool, pool.Pool[i+1:]...)
							} else if i == 0 {
								pool.Pool = pool.Pool[i+1:]
							}  else {
								pool.Pool = pool.Pool[:i]
							}
							pool.WG.Done()
						} else {
							i++
						}
					}
					log.Println(fmt.Sprintf("Pool clean - Removed %d completed tasks!!", count))
				}
			}
			for _,task := range pool.Pool {
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
			
		}(pool)
	}
}

func (pool *SchedulerPool) IsRunning() bool {
	if pool.Active {
		return true
	}
	for _,task := range pool.Pool {
		if task.IsRunning() {
			for task.IsRunning() {
				return true
			}
		}
	}
	return false
}

func (pool *SchedulerPool) Stop() {
	pool.Active = false
	pool.Tasks <- ScheduleTask{
		Id: "<close>",
		Jobs: []Job{},
	}
	close(pool.Tasks)
	
}

type RunnableStruct interface {
	Start()
	Stop()
	Status() bool
}
