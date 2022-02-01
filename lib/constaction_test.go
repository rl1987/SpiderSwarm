package spsw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstActionRun(t *testing.T) {
	c := "A flower raised in a greenhouse is still beautiful, even though it knows no adversity. But a flower growing in the field that has braved wind, rain, cold, and heat possesses something more than just beauty."

	action := NewConstAction(&Value{ValueType: ValueTypeString, StringValue: c})

	dataOut := NewDataPipe()

	err := action.AddOutput(ConstActionOutput, dataOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	gotC, ok := dataOut.Remove().(string)
	assert.True(t, ok)
	assert.Equal(t, c, gotC)
}

func TestConstActionRunErr(t *testing.T) {
	c := "..."

	action := NewConstAction(&Value{ValueType: ValueTypeString, StringValue: c})

	err := action.Run()
	assert.NotNil(t, err) // fails because output is not connected.
	assert.Equal(t, errors.New("Output not connected"), err)
}
