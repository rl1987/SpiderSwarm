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

	action := NewLengthThresholdActionFromTemplate(actionTempl, "").(*LengthThresholdAction)

	assert.NotNil(t, action)
	assert.Equal(t, threshold, action.Threshold)
	assert.Equal(t, actionTempl.Name, action.Name)
}
