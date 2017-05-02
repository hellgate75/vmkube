package state

import (
	"sync"
	"time"
)

const RESPONSE_TIMEOUT = 900

type GenericContextData map[string]interface{}

type ReferenceEntry struct {
	Key   string
	Value interface{}
}

type GenericContext interface {
	Collect(Key string) chan ReferenceEntry
	Value(Key string) *interface{}
	HasKey(Key string) bool
}

type genericCtxData struct {
	mutex   sync.RWMutex
	Entries GenericContextData
}

func (data *genericCtxData) Collect(RequestId string) chan ReferenceEntry {
	var Channel chan ReferenceEntry = make(chan ReferenceEntry, 1)
	go func(data *genericCtxData, Channel chan ReferenceEntry, RequestId string) {
		select {
		case Response, ok := <-Channel:
			if ok {
				data.mutex.Lock()
				data.Entries[Response.Key] = Response.Value
				data.mutex.Unlock()
			}
		case <-time.After(time.Second * RESPONSE_TIMEOUT):
			break
		}
		close(Channel)
	}(data, Channel, RequestId)
	return Channel
}

func (data *genericCtxData) Value(RequestId string) *interface{} {
	defer data.mutex.RUnlock()
	data.mutex.RLock()
	val, ok := data.Entries[RequestId]
	if ok {
		return &val
	}
	return nil
}

func (data *genericCtxData) HasKey(RequestId string) bool {
	defer data.mutex.RUnlock()
	data.mutex.RLock()
	_, ok := data.Entries[RequestId]
	return ok
}

func NewGenericContext() GenericContext {
	return GenericContext(&genericCtxData{
		Entries: make(GenericContextData),
	})
}
