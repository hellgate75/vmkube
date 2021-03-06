package scheduler

import (
	"runtime"
	"sync"
	"time"
	"github.com/hellgate75/vmkube/state"
	"github.com/hellgate75/vmkube/tasks"
)

type SchedulerPool struct {
	Id          string
	Tasks       chan tasks.SchedulerTask
	MaxParallel int
	WG          sync.WaitGroup
	KeepAlive   bool
	PostExecute bool
	Callback    func(task tasks.SchedulerTask)
	State       tasks.SchedulerState
}

//func (pool *SchedulerPool) RegisterWaitGroup(wg sync.WaitGroup) {
//	pool.WG = wg
//	pool.State.Active = false
//	pool.State.Pool = []tasks.ScheduleTask{}
//}

func (pool *SchedulerPool) Init() {
	pool.State.Paused = false
	pool.State.Active = false
	pool.Tasks = make(chan tasks.SchedulerTask)
	if pool.KeepAlive {
		pool.WG.Add(1)
	}
}

func (pool *SchedulerPool) Start(callback func()) {
	//TODO: remove cycles and activate : signal between objects to reduce CPU clock
	if !pool.State.Active {
		pool.State.Active = true
		var threads = pool.MaxParallel
		if threads == 0 {
			threads = runtime.NumCPU()
		}
		var Buffer []tasks.SchedulerTask = make([]tasks.SchedulerTask, 0)
		pumperExited := false
		var state state.StateContext = state.NewStateContext()

		// start Pool enqueue manager
		go func() {
			for pool.State.Active {
				if !pool.State.Paused {
					// Try add jobs from Buffer
					for threads > len(pool.State.Pool) && len(Buffer) > 0 {
						Task := Buffer[0]
						Buffer = Buffer[1:]
						pool.WG.Add(1)
						pool.State.Pool = append(pool.State.Pool, Task)
						pool.State.Pool[len(pool.State.Pool)-1].Init(&state)
						go pool.State.Pool[len(pool.State.Pool)-1].Execute()
						time.Sleep(1500 * time.Millisecond)
					}
					// Look for completed jobs to remove from Pool
					var i int = 0
					for len(pool.State.Pool) > 0 && i < len(pool.State.Pool) {
						if !pool.State.Pool[i].IsRunning() {
							if pool.PostExecute {
								go pool.Callback(pool.State.Pool[i])
							}
							if i > 0 && i < len(pool.State.Pool)-2 {
								pool.State.Pool = pool.State.Pool[:i]
								pool.State.Pool = append(pool.State.Pool, pool.State.Pool[i+1:]...)
							} else if i == 0 {
								pool.State.Pool = pool.State.Pool[i+1:]
							} else {
								pool.State.Pool = pool.State.Pool[:i]
							}
							pool.WG.Done()
							//time.Sleep(1 * time.Second)
						} else {
							i++
						}
					}
				} else {
					time.Sleep(3 * time.Second)
				}
			}
			pumperExited = true
		}()
		// start jobs
		go func() {
			//if runtime.NumCPU() < threads + 3 {
			//	runtime.GOMAXPROCS(threads + 5)
			//}
			for pool.State.Active {
				if !pool.State.Paused {
					Task := <-pool.Tasks
					if Task.Id == "<close>" {
						break
					} else if Task.Id != "" {
						Buffer = append(Buffer, Task)
					} else {
						time.Sleep(1500 * time.Millisecond)
					}
				} else {
					time.Sleep(3 * time.Second)
				}
			}
			for !pumperExited {
				time.Sleep(500 * time.Millisecond)
			}
			for i := 0; i < len(pool.State.Pool); i++ {
				if pool.State.Pool[i].IsRunning() {
					pool.State.Pool[i].Abort()
					for pool.State.Pool[i].IsRunning() {
						time.Sleep(1000 * time.Millisecond)
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
	if pool.State.Active && !pool.State.Paused {
		return true
	}
	for _, task := range pool.State.Pool {
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
	pool.Tasks <- tasks.SchedulerTask{
		Id:   "<close>",
		Jobs: []tasks.JobProcess{},
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
