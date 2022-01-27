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

	action := NewStringMapUpdateAction("")

	action.AddInput(StringMapUpdateActionInputOld, oldCookiesIn)
	action.AddInput(StringMapUpdateActionInputNew, newCookiesIn)
	action.AddOutput(StringMapUpdateActionOutputUpdated, updatedCookiesOut)

	err := action.Run()

	assert.Nil(t, err)

	gotUpdatedCookies, ok := updatedCookiesOut.Remove().(map[string]string)

	assert.True(t, ok)
	assert.Equal(t, updatedCookies, gotUpdatedCookies)
}

func TestStringMapUpdateActionRunWithOverride(t *testing.T) {
	old := map[string]string{
		"page": "1",
		"session": "5514",
	}	

	key := "page"
	newValue := "2"

	expectUpdated := map[string]string{
		"page": "2",
		"session": "5514",
	}

	inDP := NewDataPipe()
	inDP2 := NewDataPipe()
	outDP := NewDataPipe()

	inDP.Add(newValue)
	inDP2.Add(old)

	action := NewStringMapUpdateAction(key)

	action.AddInput(StringMapUpdateActionInputOverridenValue, inDP)
	action.AddInput(StringMapUpdateActionInputOld, inDP2)
	action.AddOutput(StringMapUpdateActionOutputUpdated, outDP)

	err := action.Run()

	assert.Nil(t, err)

	gotUpdated, ok := outDP.Remove().(map[string]string)

	assert.True(t, ok)
	assert.Equal(t, expectUpdated, gotUpdated)
}

