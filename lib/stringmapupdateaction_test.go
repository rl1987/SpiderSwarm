package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMapUpdateActionRun(t *testing.T) {
	oldCookies := map[string]string{
		"session": "12211",
		"n":       "1",
	}

	newCookies := map[string]string{
		"n": "2",
		"x": "1.11",
	}

	updatedCookies := map[string]string{
		"session": "12211",
		"n":       "2",
		"x":       "1.11",
	}

	oldCookiesIn := NewDataPipe()
	newCookiesIn := NewDataPipe()
	updatedCookiesOut := NewDataPipe()

	oldCookiesIn.Add(oldCookies)
	newCookiesIn.Add(newCookies)

	action := NewStringMapUpdateAction()

	action.AddInput(StringMapUpdateActionInputOld, oldCookiesIn)
	action.AddInput(StringMapUpdateActionInputNew, newCookiesIn)
	action.AddOutput(StringMapUpdateActionOutputUpdated, updatedCookiesOut)

	err := action.Run()

	assert.Nil(t, err)

	gotUpdatedCookies, ok := updatedCookiesOut.Remove().(map[string]string)

	assert.True(t, ok)
	assert.Equal(t, updatedCookies, gotUpdatedCookies)
}
