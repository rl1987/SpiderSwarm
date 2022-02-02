package spsw

import (
  	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVParseAction(t *testing.T) {
	action := NewCSVParseAction()

	assert.NotNil(t, action)

	assert.Equal(t, []string{CSVParseActionInputCSVBytes, CSVParseActionInputCSVStr},
		action.AllowedInputNames)
	assert.Equal(t, []string{CSVParseActionOutputMap}, action.AllowedOutputNames)
}

func TestNewCSVParseActionFromTemplate(t *testing.T) {
	actionTempl := &ActionTemplate{
		Name: "TestAction",
		StructName: "CSVParseAction",
		ConstructorParams: map[string]Value{},
	}

	action := NewCSVParseActionFromTemplate(actionTempl)

	assert.NotNil(t, action)
}

func TestCSVParseActionRun(t *testing.T) {
	csvStr := "id,name,age\r\n1,John,25\r\n2,Jane,22\r\n"

	expectMap := map[string][]string{
		"id":   []string{"1", "2"},
		"name": []string{"John", "Jane"},
		"age":  []string{"25", "22"},
	}

	inDP := NewDataPipe()
	outDP := NewDataPipe()

	action := NewCSVParseAction()

	inDP.Add(csvStr)

	action.AddInput(CSVParseActionInputCSVStr, inDP)
	action.AddOutput(CSVParseActionOutputMap, outDP)

	err := action.Run()

	assert.Nil(t, err)
	gotMap, ok := outDP.Remove().(map[string][]string)

	assert.True(t, ok)
	assert.Equal(t, expectMap, gotMap)
}

func TestCSVParseActionRunError(t *testing.T) {
	t.Skip()

	badCSVStr := "\xaa"

	inDP := NewDataPipe()
	outDP := NewDataPipe()

	action := NewCSVParseAction()

	inDP.Add([]byte(badCSVStr))

	action.AddInput(CSVParseActionInputCSVBytes, inDP)
	action.AddOutput(CSVParseActionOutputMap, outDP)

	err := action.Run()

	fmt.Println(err)

	assert.NotNil(t, err)
}

