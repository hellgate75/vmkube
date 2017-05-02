package test

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"vmkube/utils"
)

//
//os.Exit(0)

type Test struct {
	Name    string
	Surname string
	Age     int
}

var (
	myStruct []interface{} = []interface{}{
		Test{
			Name:    "Fabiana",
			Surname: "XXXXXXXXXXX",
			Age:     42,
		},
		Test{
			Name:    "Fabrizio",
			Surname: "XXXXXXXXXXX",
			Age:     42,
		},
		Test{
			Name:    "Francesco",
			Surname: "XXXXXXXXXXX",
			Age:     38,
		},
	}
)

func TestReduce(t *testing.T) {
	interfaceX, err := utils.ReduceStruct("Name", myStruct)
	assert.Equal(t, nil, err, "Expected no errors from reduce")
	assert.Equal(t, reflect.Slice, reflect.TypeOf(interfaceX).Kind(), "Expected Array of Elements from reduce")
	assert.Equal(t, 3, len(interfaceX), "Expected Array of three Elements from reduce")
}

func TestReduceToStringsSlice(t *testing.T) {
	interfaceX, err := utils.ReduceStruct("Name", myStruct)
	myArray := utils.ReducedToStringsSlice(interfaceX)
	assert.Equal(t, nil, err, "Expected no errors from reduce")
	assert.Equal(t, 3, len(myArray), "Expected Conveted Array of three Elements from reduce")
	assert.Equal(t, "Fabiana", myArray[0], "Expected first element in the array is the same of the Name in the structure")
	assert.Equal(t, "Fabrizio", myArray[1], "Expected second element in the array is the same of the Name in the structure")
	assert.Equal(t, "Francesco", myArray[2], "Expected third element in the array is the same of the Name in the structure")
}

func TestReduceToIntsSlice(t *testing.T) {
	interfaceX, err := utils.ReduceStruct("Age", myStruct)
	myArray := utils.ReducedToIntsSlice(interfaceX)
	assert.Equal(t, nil, err, "Expected no errors from reduce")
	assert.Equal(t, 3, len(myArray), "Expected Conveted Array of three Elements from reduce")
	assert.Equal(t, 42, myArray[0], "Expected first element in the array is the same of the Age in the structure")
	assert.Equal(t, 42, myArray[1], "Expected second element in the array is the same of the Age in the structure")
	assert.Equal(t, 38, myArray[2], "Expected third element in the array is the same of the Age in the structure")
}
