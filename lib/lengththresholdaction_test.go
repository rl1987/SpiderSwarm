package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLengthThresholdActionFromTemplate(t *testing.T) {
	threshold := 10

	actionTempl := &ActionTemplate{
		Name:       "TresholdAction",
		StructName: "LengthThresholdAction",
		ConstructorParams: map[string]Value{
			"threshold": Value{
				ValueType: ValueTypeInt,
				IntValue:  threshold,
			},
		},
	}

	action := NewLengthThresholdActionFromTemplate(actionTempl).(*LengthThresholdAction)

	assert.NotNil(t, action)
	assert.Equal(t, threshold, action.Threshold)
	assert.Equal(t, actionTempl.Name, action.Name)
}

func TestLengthThresholdActionRun(t *testing.T) {
	threshold := 2

	action := NewLengthThresholdAction(threshold)

	dpIn := NewDataPipe()
	action.AddInput(LengthThresholdActionInputSlice, dpIn)

	dpOut := NewDataPipe()
	action.AddOutput(LengthThresholdActionOutputThresholdUnmet, dpOut)

	dpIn.Add([]string{"a"})

	err := action.Run()
	assert.Nil(t, err)

	result, ok := dpOut.Remove().(bool)
	assert.True(t, ok)
	assert.True(t, result)

	dpIn.Add([]string{"a", "b"})

	action.Run()

	result, ok = dpOut.Remove().(bool)
	assert.True(t, ok)
	assert.False(t, result)
}
