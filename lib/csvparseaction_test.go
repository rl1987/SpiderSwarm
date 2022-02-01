package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVParseActionRun(t *testing.T) {
	csvStr := "id,name,age\r\n1,John,25\r\n2,Jane,22\r\n"

	expectMap := map[string][]string{
		"id": []string{"1", "2"},
		"name": []string{"John", "Jane"},
		"age": []string{"25", "22"},
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

