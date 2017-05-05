package test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
	"vmkube/action"
	"vmkube/scheduler"
	"vmkube/tasks"
)

/*
* Test Job implements vmkube/scheduler.RunnableStruct
 */
type TestJob struct {
	Name    string
	Count   int
	State   bool
	OutChan chan string
}

func (job TestJob) Start(channel chan bool) {
	if !job.State {
		job.State = true
		for i := 0; i < job.Count; i++ {
			job.OutChan <- strconv.Itoa(i)
			if !job.State {
				break
			}
		}
		job.State = false
		channel <- true
	}
}

func (job TestJob) Stop() {
	job.State = false
}

func (job TestJob) Status() bool {
	return job.State
}

func (job TestJob) IsInterrupted() bool {
	return false
}
func (job TestJob) IsError() bool {
	return false
}
func (job TestJob) Response() interface{} {
	return ""
}
func (job TestJob) WaitFor() {

}

func testJobs(chan1 chan string, chan2 chan string, chan3 chan string) {
	pool := scheduler.SchedulerPool{
		Id:          action.NewUUIDString(),
		MaxParallel: 2,
		KeepAlive:   false,
	}
	pool.Init()
	go pool.Start(func() {
		pool.WG.Done()
	})
	task1 := tasks.SchedulerTask{
		Id: action.NewUUIDString(),
		Jobs: []tasks.JobProcess{
			tasks.JobProcess(&tasks.Job{
				Id:   action.NewUUIDString(),
				Name: "TestJob1",
				Runnable: TestJob{
					Name:    "Job 1",
					Count:   10,
					OutChan: chan1,
				},
			}),
			tasks.JobProcess(&tasks.Job{
				Id:   action.NewUUIDString(),
				Name: "TestJob1",
				Runnable: TestJob{
					Name:    "Job 2",
					Count:   5,
					OutChan: chan1,
				},
			}),
		},
	}
	task2 := tasks.SchedulerTask{
		Id: action.NewUUIDString(),
		Jobs: []tasks.JobProcess{
			tasks.JobProcess(&tasks.Job{
				Id:   action.NewUUIDString(),
				Name: "TestJob1",
				Runnable: TestJob{
					Name:    "Job 3",
					Count:   5,
					OutChan: chan2,
				},
			}),
			tasks.JobProcess(&tasks.Job{
				Id:   action.NewUUIDString(),
				Name: "TestJob1",
				Runnable: TestJob{
					Name:    "Job 4",
					Count:   10,
					OutChan: chan2,
				},
			}),
		},
	}
	task3 := tasks.SchedulerTask{
		Id: action.NewUUIDString(),
		Jobs: []tasks.JobProcess{
			tasks.JobProcess(&tasks.Job{
				Id:   action.NewUUIDString(),
				Name: "TestJob1",
				Runnable: TestJob{
					Name:    "Job 5",
					Count:   5,
					OutChan: chan3,
				},
			}),
			tasks.JobProcess(&tasks.Job{
				Id:   action.NewUUIDString(),
				Name: "TestJob1",
				Runnable: TestJob{
					Name:    "Job 6",
					Count:   5,
					OutChan: chan3,
				},
			}),
		},
	}
	pool.WG.Add(1)
	go func() {
		pool.Tasks <- task1
		pool.Tasks <- task2
		pool.Tasks <- task3
	}()
	go func(pool *scheduler.SchedulerPool) {
		time.Sleep(3 * time.Second)
		pool.Stop()
	}(&pool)
	pool.WG.Wait()
}

func collectChanValues(chanX chan string, buffer *bytes.Buffer) {
	for {
		val := <-chanX
		if val == "" {
			break
		}
		buffer.Write([]byte(val))
		buffer.Write([]byte(" "))
	}
}

func TestSchedulerJobExecution(t *testing.T) {
	chan1 := make(chan string)
	chan2 := make(chan string)
	chan3 := make(chan string)
	buffer1 := bytes.NewBuffer([]byte{})
	buffer2 := bytes.NewBuffer([]byte{})
	buffer3 := bytes.NewBuffer([]byte{})
	go collectChanValues(chan1, buffer1)
	go collectChanValues(chan2, buffer2)
	go collectChanValues(chan3, buffer3)

	testJobs(chan1, chan2, chan3)
	chan1 <- ""
	chan2 <- ""
	chan3 <- ""
	close(chan1)
	close(chan2)
	close(chan3)
	value1 := string(buffer1.Bytes())
	value2 := string(buffer2.Bytes())
	value3 := string(buffer3.Bytes())
	expected1 := "0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 "
	expected2 := "0 1 2 3 4 0 1 2 3 4 5 6 7 8 9 "
	expected3 := "0 1 2 3 4 0 1 2 3 4 "
	assert.Equal(t, expected1, value1, "Expected output for Job Group 1")
	assert.Equal(t, expected2, value2, "Expected output for Job Group 2")
	assert.Equal(t, expected3, value3, "Expected output for Job Group 3")
}
