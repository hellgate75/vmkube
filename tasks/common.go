package tasks

import (
	"os/exec"
	"vmkube/model"
	"vmkube/procedures"
	"vmkube/state"
	"vmkube/term"
)

type TaskSignal struct {
	TaskId      string
	Complete    bool
	Errors      bool
	Interrupted bool
	Done        int
	Jobs        int
}

func (Signal *TaskSignal) IsDone() bool {
	return Signal.Complete && (Signal.Done == Signal.Jobs || Signal.Errors || Signal.Interrupted)
}

const MachineReadOperationTimeout = 900

type MachineOperationsJob struct {
	Name             string
	State            bool
	OutChan          chan *MachineOperationsJob
	OwnState         term.KeyValueElement
	Project          model.Project
	Infra            model.Infrastructure
	InstanceId       string
	Activity         ActivityCouple
	ActivityGroup    ActivityGroup
	MachineMessage   procedures.MachineMessage
	Index            int
	PartOf           int
	SendStartMessage bool
	Command          string
	Machine          string
	commandPipe      chan procedures.MachineMessage
	commandChannel   chan *exec.Cmd
	control          procedures.MachineControlStructure
}

type Runnable interface {
	Start(exitChannel chan bool)
	Stop()
	Status() bool
	IsInterrupted() bool
	IsError() bool
	Response() interface{}
	WaitFor()
}

type SchedulerState struct {
	Active bool
	Paused bool
	Pool   []SchedulerTask
}

type JobProcess interface {
	Run()
	IsRunning() bool
	IsAsync() bool
	HasErrors() bool
	Abort()
	Init(Sequence int, Global int)
	WaitFor()
	GetRunnable() Runnable
}

type Job struct {
	Id       string
	Name     string
	Runnable Runnable
	Async    bool
	Sequence int
	Of       int
	State    bool
}

type TaskProcess interface {
	Run()
	IsRunning() bool
	Abort()
	Init(context *state.StateContext)
	Wait()
	Deactivate()
}

type SchedulerTask struct {
	Id        string
	Jobs      []JobProcess
	Active    bool
	Working   bool
	LastIndex int
	State     *state.StateContext
}
