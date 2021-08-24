package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstActionRun(t *testing.T) {
	c := "A flower raised in a greenhouse is still beautiful, even though it knows no adversity. But a flower growing in the field that has braved wind, rain, cold, and heat possesses something more than just beauty."

	action := NewConstAction(c)

	dataOut := NewDataPipe()

	err := action.AddOutput(ConstActionOutput, dataOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	gotC, ok := dataOut.Remove().(string)
	assert.True(t, ok)
	assert.Equal(t, c, gotC)
}
