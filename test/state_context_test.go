package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"vmkube/state"
	"time"
)

func TestNewStateContext(t *testing.T) {
	context := state.NewStateContext()
	assert.Equal(t, false, context.HasValue("1"), "Non registered Id has no value")
	assert.Equal(t, false, context.State("1"), "Non registered Id has false state")
}

func TestStateContextUse(t *testing.T) {
	context := state.NewStateContext()
	context.Collect("1") <- state.StateReferenceData{
		Id: "1",
		Status: true,
	}
	time.Sleep(time.Second*1)
	context.Collect("2") <- state.StateReferenceData{
		Id: "2",
		Status: false,
	}
	time.Sleep(time.Second*1)
	assert.Equal(t, true, context.HasValue("1"), "Registered Id has value")
	assert.Equal(t, true, context.State("1"), "Registered Id has right state")
	assert.Equal(t, true, context.HasValue("2"), "Registered Id has value")
	assert.Equal(t, false, context.State("2"), "Registered Id has right state")
}
