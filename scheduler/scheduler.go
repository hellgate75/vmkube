package scheduler


import (
	"runtime"
	"time"
	"strconv"
	"sync"
	"log"
	"reflect"
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
}

func (pool *SchedulerPool) RegisterWaitGroup(wg sync.WaitGroup) {
	pool.WG = wg
}

func (pool *SchedulerPool) Init() {
	pool.Tasks = make(chan ScheduleTask)
	
}

func (pool *SchedulerPool) Start() {
	if ! pool.Active {
		//start jobs
		go func(pool *SchedulerPool) {
			var threads = pool.MaxParallel
			if threads == 0 {
				threads = runtime.NumCPU()
			}
			log.Println("Max Parallel Processes = " + strconv.Itoa(threads))
			runtime.GOMAXPROCS(threads + 1)
			pool.Active = true
			pool.Pool = make([]ScheduleTask, 0)
			
			for pool.Active {
				log.Println("Pool active ...")
				if threads > len(pool.Pool) {
					log.Println("Wait for message ...")
					Task := <- pool.Tasks
					log.Println("Task : " + Task.Id)
					if Task.Id != "" {
						if Task.Id == "<close>" {
							break
						} else {
							log.Println("Pool Append : " + pool.Id)
							pool.Pool = append(pool.Pool, Task)
							go Task.Execute()
						}
					} else {
						time.Sleep(1000*time.Millisecond)
					}
				} else {
					//Thread Pool Full
					println("All full ... ")
					time.Sleep(1000*time.Millisecond)
					i := 0
					log.Println("Trimming completed task from : "+strconv.Itoa(len(pool.Pool)))
					for i < len(pool.Pool) {
						if ! pool.Pool[i].IsRunning() {
							if i > 0 && i < len(pool.Pool) - 1 {
								pool.Pool = pool.Pool[:i]
								pool.Pool = append(pool.Pool, pool.Pool[i+1:]...)
							} else if i == 0 {
								pool.Pool = pool.Pool[i+1:]
							}  else {
								pool.Pool = pool.Pool[:i]
							}
						} else {
							i++
						}
					}
				}
			}
			for _,task := range pool.Pool {
				if task.IsRunning() {
					task.Abort()
					for task.IsRunning() {
						time.Sleep(1000*time.Millisecond)
					}
				}
			}
			pool.WG.Done()
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
