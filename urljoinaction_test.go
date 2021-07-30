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

func TestURLJoinActionBaseWithRelative(t *testing.T) {
	action := NewURLJoinAction("http://example.com/")

	baseURLStr := "http://example.com/directory"
	relativeURLStr := "../../..//search?q=golang"
	absoluteURLStr := "http://example.com/search?q=golang"

	baseIn := NewDataPipe()
	relativeIn := NewDataPipe()
	absoluteOut := NewDataPipe()

	baseIn.Add(baseURLStr)
	relativeIn.Add(relativeURLStr)

	err := action.AddInput(URLJoinActionInputBaseURL, baseIn)
	assert.Nil(t, err)

	err = action.AddInput(URLJoinActionInputRelativeURL, relativeIn)
	assert.Nil(t, err)

	err = action.AddOutput(URLJoinActionOutputAbsoluteURL, absoluteOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	gotAbsoluteURLStr, ok := absoluteOut.Remove().(string)
	assert.True(t, ok)
	assert.Equal(t, absoluteURLStr, gotAbsoluteURLStr)
}
