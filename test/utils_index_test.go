package test

import (
	"vmkube/utils"
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
	"fmt"
	"log"
)

var (
	index1, index2, index3 utils.Index
)


func TestNewIndex(t *testing.T) {
	index1.New(10)
	assert.NotEmpty(t, index1.Value, "Index Instance Must be correct")
	assert.Equal(t, 10, len(index1.Value), "Length of Index Content Must be correct")
	expected := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 10)
	assert.Equal(t, expected, index1.Value, "Index Instance Must have a correct value")
}

func TestIndexNext(t *testing.T) {
	index1.New(10)
	index2 = index1.Next()
	index3 = index2.Next()
	
	assert.NotEmpty(t, index2.Value, "Next Index from Zero Instance Must be correct")
	assert.NotEmpty(t, index3.Value, "Next of Next Index from Zero Instance Must be correct")
	assert.Equal(t, 10, len(index1.Value), "Length of Original Index Content Must be correct")
	assert.Equal(t, 10, len(index2.Value), "Length of Next Index from Zero Content Must be correct")
	assert.Equal(t, 10, len(index3.Value), "Length of Next of Next Index from Zero Content Must be correct")
	expected := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 10)
	assert.Equal(t, expected, index1.Value, "Index Instance have a correct value (0)")
	expected2 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected2 += fmt.Sprintf("%c", (utils.INDEX_START + 1) )
	assert.Equal(t, expected2, index2.Value, "Next Index from Zero Instance Must have a correct value (1)")
	expected3 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected3 += fmt.Sprintf("%c", (utils.INDEX_START + 2) )
	assert.Equal(t, expected3, index3.Value, "Next of Next Index from Zero Instance Must have a correct value (2)")
}

func TestIndexPrevious(t *testing.T) {
	index1.New(10)
	index2 = index1.Next()
	index2 = index2.Next()
	expected0 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected0 += fmt.Sprintf("%c", (utils.INDEX_START + 2) )
	assert.Equal(t, expected0, index2.Value, "Test Index Instance (Next of Next of Zero) Must be 2")
	index2 = index2.Previous()
	index3 = index2.Previous()
	
	assert.NotEmpty(t, index2.Value, "Previous Index from Two Instance Must be correct")
	assert.NotEmpty(t, index3.Value, "Previous of Previous Index from Two Instance Must be correct")
	assert.Equal(t, 10, len(index1.Value), "Length of Original Index Content Must be correct")
	assert.Equal(t, 10, len(index2.Value), "Length of Previous Index from Two Content Must be correct")
	assert.Equal(t, 10, len(index3.Value), "Length of Previous of Previous Index from Two Content Must be correct")
	expected := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 10)
	assert.Equal(t, expected, index1.Value, "Index Instance have a correct value")
	expected2 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected2 += fmt.Sprintf("%c", (utils.INDEX_START + 1) )
	assert.Equal(t, expected2, index2.Value, "Previous Index from Two Instance Must have a correct value (1)")
	assert.Equal(t, index1.Value, index3.Value, "Previous of Previous Index from Two Instance Must have a correct value (0)")
}


func TestIndexFromInt(t *testing.T) {
	index1.FromInt(1, 10)
	index2.FromInt(10, 10)
	index3.FromInt(int(utils.INDEX_END)-int(utils.INDEX_START)+1, 10)
	
	assert.NotEmpty(t, index1.Value, "From Int 1 Index Instance Must be correct")
	assert.NotEmpty(t, index2.Value, "From Int 10 Index Instance Must be correct")
	assert.NotEmpty(t, index3.Value, "From Int One over the INDEX_END Index Instance Must be correct")
	assert.Equal(t, 10, len(index1.Value), "Length of From Int 1 Index Content Must be correct")
	assert.Equal(t, 10, len(index2.Value), "Length of From Int 10 Index Content Must be correct")
	assert.Equal(t, 10, len(index3.Value), "Length of From Int One over the INDEX_END Index Content Must be correct")
	expected := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	assert.Equal(t, expected, index1.Value, "From Int 1 Index Instance have a correct value (x01)")
	expected2 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected2 += fmt.Sprintf("%c", (utils.INDEX_START + byte(10)) )
	log.Println(fmt.Sprintf("%s", expected2))
	log.Println(fmt.Sprintf("%s", index2.Value))
	assert.Equal(t, expected2, index2.Value, "From Int 10 Index Instance Must have a correct value (x10)")
	expected3 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 8)
	expected3 += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	expected3 += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	assert.Equal(t, expected3, index3.Value, "From Int One over the INDEX_END Index Instance Must have a correct value (1x01)")
}

func TestIndexSum(t *testing.T) {
	index1.FromInt(1, 10)
	index2.FromInt(10, 10)
	index2.Sum(index1)
	index3.FromInt(int(utils.INDEX_END)-int(utils.INDEX_START), 10)
	index3.Sum(index1)
	assert.NotEmpty(t, index1.Value, "From Int 1 Index Instance Must be correct")
	assert.NotEmpty(t, index2.Value, "Add Int 10 and One Index Instance Must be correct")
	assert.NotEmpty(t, index3.Value, "Subtract In 10 from Eleven over the INDEX_END Index Instance Must be correct")
	assert.Equal(t, 10, len(index1.Value), "Length of From Int 1 Index Content Must be correct")
	assert.Equal(t, 10, len(index2.Value), "Length of Add Int 10 and One Index Content Must be correct")
	assert.Equal(t, 10, len(index3.Value), "Length of Subtract In 10 from Eleven over the INDEX_END Index Content Must be correct")
	expected := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	assert.Equal(t, expected, index1.Value, "From Int 1 Index Instance have a correct value (x01)")
	expected2 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected2 += fmt.Sprintf("%c", (utils.INDEX_START + byte(11)) )
	log.Println(fmt.Sprintf("%s", expected2))
	log.Println(fmt.Sprintf("%s", index2.Value))
	assert.Equal(t, expected2, index2.Value, "Add Int 10 and One Index Instance Must have a correct value (x10)")
	expected3 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 8)
	expected3 += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	expected3 += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	assert.Equal(t, expected3, index3.Value, "Subtract In 10 from Eleven over the INDEX_END Index Instance Must have a correct value (1x01)")
}

func TestIndexSubtract(t *testing.T) {
	index1.FromInt(10, 10)
	index2.FromInt(21, 10)
	index2.Subtract(index1)
	index3.FromInt(int(utils.INDEX_END)-int(utils.INDEX_START) + 11, 10)
	index3.Subtract(index1)
	assert.NotEmpty(t, index1.Value, "From Int 1 Index Instance Must be correct")
	assert.NotEmpty(t, index2.Value, "Subtract Int 10 From TwentyOne Index Instance Must be correct")
	assert.NotEmpty(t, index3.Value, "Subtract In 10 from Eleven over the INDEX_END Index Instance Must be correct")
	assert.Equal(t, 10, len(index1.Value), "Length of From Int 1 Index Content Must be correct")
	assert.Equal(t, 10, len(index2.Value), "Length of Subtract Int 10 From TwentyOne Index Content Must be correct")
	assert.Equal(t, 10, len(index3.Value), "Length of Subtract In 10 from Eleven over the INDEX_END Index Content Must be correct")
	expected := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected += fmt.Sprintf("%c", (utils.INDEX_START + byte(10)) )
	assert.Equal(t, expected, index1.Value, "From Int 1 Index Instance have a correct value (x01)")
	expected2 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 9)
	expected2 += fmt.Sprintf("%c", (utils.INDEX_START + byte(11)) )
	log.Println(fmt.Sprintf("%s", expected2))
	log.Println(fmt.Sprintf("%s", index2.Value))
	assert.Equal(t, expected2, index2.Value, "Subtract Int 10 From TwentyOne Index Instance Must have a correct value (x10)")
	expected3 := strings.Repeat(fmt.Sprintf("%c", utils.INDEX_START), 8)
	expected3 += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	expected3 += fmt.Sprintf("%c", (utils.INDEX_START + byte(1)) )
	assert.Equal(t, expected3, index3.Value, "Subtract In 10 from Eleven over the INDEX_END Index Instance Must have a correct value (1x01)")
}

func TestIndexCompare(t *testing.T) {
	index1.FromInt(10, 10)
	index2.FromInt(11, 10)
	index3.FromInt(9, 10)
	
	assert.NotEmpty(t, index1.Value, "From Int 10 Index Instance Must be correct")
	assert.NotEmpty(t, index2.Value, "From Int 11 Index Instance Must be correct")
	assert.NotEmpty(t, index3.Value, "From Int  9 Index Instance Must be correct")
	assert.Equal(t, 10, len(index1.Value), "Length of From Int 10 Index Content Must be correct")
	assert.Equal(t, 10, len(index2.Value), "Length of From Int 11 Index Content Must be correct")
	assert.Equal(t, 10, len(index3.Value), "Length of From Int  9 Index Content Must be correct")
	assert.Equal(t, 0, index1.Compare(index1), "From Int 10  Index Copare to 10 Index Instances have result 0 - equals to  the compared index ...")
	assert.Equal(t, -1, index1.Compare(index2), "From Int 10  Index Copare to 10 Index Instances have result -1 - less than the compared index ..")
	assert.Equal(t, 1, index1.Compare(index3), "From Int 10  Index Copare to 10 Index Instances have result 1 -  more than  the compared index ...")
}
