package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringCutActionFromTemplate(t *testing.T) {
	constructorParams := map[string]interface{}{
		"from": "here",
		"to":   "eternity",
	}

	actionTempl := &ActionTemplate{
		StructName:        "StringCutAction",
		ConstructorParams: constructorParams,
	}

	action := NewStringCutActionFromTemplate(actionTempl)

	assert.NotNil(t, action)
	assert.Equal(t, "here", action.From)
	assert.Equal(t, "eternity", action.To)
}

func TestStringCutActionRun(t *testing.T) {
	inputStr := "... latitude: '12.222';"
	expectedOutStr := "12.222"

	action := NewStringCutAction("latitude: '", "'")

	strIn := NewDataPipe()
	strOut := NewDataPipe()

	strIn.Add(inputStr)

	err := action.AddInput(StringCutActionInputStr, strIn)
	assert.Nil(t, err)

	err = action.AddOutput(StringCutActionOutputStr, strOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	gotStr, ok := action.Outputs[StringCutActionOutputStr][0].Remove().(string)
	assert.True(t, ok)
	assert.Equal(t, expectedOutStr, gotStr)
}
