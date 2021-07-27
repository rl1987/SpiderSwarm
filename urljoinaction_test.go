package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLJoinActionRun(t *testing.T) {
	action := NewURLJoinAction("http://example.com/directory/")

	relativeURLStr := "../../..//search?q=golang"
	absoluteURLStr := "http://example.com/search?q=golang"

	inDP := NewDataPipe()
	outDP := NewDataPipe()

	err := action.AddInput(URLJoinActionInputRelativeURL, inDP)
	assert.Nil(t, err)

	err = action.AddOutput(URLJoinActionOutputAbsoluteURL, outDP)
	assert.Nil(t, err)

	inDP.Add(relativeURLStr)

	err = action.Run()

	gotAbsoluteURLStr, ok := outDP.Remove().(string)
	assert.True(t, ok)

	assert.Equal(t, absoluteURLStr, gotAbsoluteURLStr)
}
