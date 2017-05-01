package state

import (
	"sync"
	"time"
)

type StateContextData map[string]bool

type StateReferenceData struct {
	Id			string
	Status	bool
}

type StateContext interface {
	Collect(RequestId string) 	chan StateReferenceData
	State(RequestId string)			bool
	HasValue(RequestId string)	bool
}

type stateCtxData struct {
	mutex     sync.RWMutex
	Responses StateContextData
}

func (data *stateCtxData) Collect(RequestId string) chan StateReferenceData {
	var Channel chan StateReferenceData = make(chan StateReferenceData, 1)
	go func(data *stateCtxData, Channel chan StateReferenceData, RequestId string) {
		select {
			case Response, ok := <- Channel:
				if ok {
					data.mutex.Lock()
					data.Responses[Response.Id] = Response.Status
					data.mutex.Unlock()
				}
			case <-time.After(time.Second * RESPONSE_TIMEOUT):
				data.mutex.Lock()
				data.Responses[RequestId] = false
				data.mutex.Unlock()
		}
		close(Channel)
	}(data, Channel, RequestId)
	return  Channel
}

func (data *stateCtxData) State(RequestId string) bool {
	defer  data.mutex.RUnlock()
	data.mutex.RLock()
	val, ok := data.Responses[RequestId]
	if ok {
		return val
	}
	return false
}

func (data *stateCtxData) HasValue(RequestId string) bool {
	defer  data.mutex.RUnlock()
	data.mutex.RLock()
	_,ok := data.Responses[RequestId]
	return  ok
}

func NewStateContext() StateContext {
	return  StateContext(&stateCtxData{
		Responses: make(StateContextData),
	})
}