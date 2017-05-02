package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"vmkube/state"
)

func TestNewStateContext(t *testing.T) {
	context := state.NewStateContext()
	assert.Equal(t, false, context.HasValue("1"), "Non registered Id has no value")
	assert.Equal(t, false, context.State("1"), "Non registered Id has false state")
}

func TestStateContextUse(t *testing.T) {
	context := state.NewStateContext()
	context.Collect("1") <- state.StateReferenceData{
		Id:     "1",
		Status: true,
	}
	time.Sleep(time.Second * 1)
	context.Collect("2") <- state.StateReferenceData{
		Id:     "2",
		Status: false,
	}
	time.Sleep(time.Second * 1)
	assert.Equal(t, true, context.HasValue("1"), "Registered Id has value")
	assert.Equal(t, true, context.State("1"), "Registered Id has right state")
	assert.Equal(t, true, context.HasValue("2"), "Registered Id has value")
	assert.Equal(t, false, context.State("2"), "Registered Id has right state")
}

func TestNewGenericContext(t *testing.T) {
	context := state.NewGenericContext()
	assert.Equal(t, false, context.HasKey("1"), "Non registered Id has no value")
	assert.Nil(t, context.Value("1"), "Non registered Id has No Reference Value")
}

func TestGenericContextUse(t *testing.T) {
	context := state.NewGenericContext()
	context.Collect("1") <- state.ReferenceEntry{
		Key:   "1",
		Value: "A",
	}
	time.Sleep(time.Second * 1)
	context.Collect("2") <- state.ReferenceEntry{
		Key:   "2",
		Value: "B",
	}
	time.Sleep(time.Second * 1)
	assert.Equal(t, true, context.HasKey("1"), "Registered Id key exists")
	assert.Equal(t, "A", *context.Value("1"), "Registered Id has right value")
	assert.Equal(t, true, context.HasKey("2"), "Registered Id key exists")
	assert.Equal(t, "B", *context.Value("2"), "Registered Id has right value")
}
