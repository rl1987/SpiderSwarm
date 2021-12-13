package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testJSONStr = `
{
  "firstName": "Nancy",
  "lastName" : "Thompson",
  "age"      : 17,
  "address"  : {
    "streetAddress": "1428 Elm Street",
    "city"         : "Springwood, OH"
  },
  "phoneNumbers": [
    {
      "type"  : "iPhone",
      "number": "0123-4567-8888"
    },
    {
      "type"  : "home",
      "number": "0123-4567-8910"
    }
  ]
}
`

func TestNewJSONPathActionFromTemplate(t *testing.T) {
	jsonPath := "$.store.book[*].author"
	expectMany := true
	decode := true

	constructorParams := map[string]Value{
		"jsonPath": Value{
			ValueType:   ValueTypeString,
			StringValue: jsonPath,
		},
		"expectMany": Value{
			ValueType: ValueTypeBool,
			BoolValue: expectMany,
		},
		"decode": Value{
			ValueType: ValueTypeBool,
			BoolValue: decode,
		},
	}

	actionTempl := &ActionTemplate{
		StructName:        "JSONPathAction",
		ConstructorParams: constructorParams,
	}

	action, ok := NewJSONPathActionFromTemplate(actionTempl, "").(*JSONPathAction)
	assert.True(t, ok)

	assert.Equal(t, jsonPath, action.JSONPath)
	assert.Equal(t, expectMany, action.ExpectMany)
	assert.Equal(t, decode, action.Decode)
	assert.False(t, action.CanFail)
}

func TestJSONPathActionRunBasic(t *testing.T) {
	jsonStr := "{\"name\": \"John\", \"surname\": \"Smith\"}"

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(jsonStr)

	action := NewJSONPathAction("$.name", true, false)

	action.AddInput(JSONPathActionInputJSONStr, dataPipeIn)
	action.AddOutput(JSONPathActionOutputStr, dataPipeOut)

	err := action.Run()
	assert.Nil(t, err)

	resultStr, ok := dataPipeOut.Remove().(string)
	assert.True(t, ok)

	assert.Equal(t, "John", resultStr)
}

func TestJSONPathActionRunExpectMany(t *testing.T) {
	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add([]byte(testJSONStr))

	action := NewJSONPathAction("$..number", true, true)

	action.AddInput(JSONPathActionInputJSONBytes, dataPipeIn)
	action.AddOutput(JSONPathActionOutputStr, dataPipeOut)

	err := action.Run()
	assert.Nil(t, err)

	expectResults := []string{"0123-4567-8888", "0123-4567-8910"}

	results, ok := dataPipeOut.Remove().([]string)
	assert.True(t, ok)

	assert.Equal(t, expectResults, results)

	action.Decode = false
	dataPipeIn.Add([]byte(testJSONStr))

	err = action.Run()
	assert.Nil(t, err)

	result, ok2 := dataPipeOut.Remove().(string)
	assert.True(t, ok2)

	expectJSONStr := "[\"0123-4567-8888\",\"0123-4567-8910\"]"

	assert.Equal(t, expectJSONStr, result)
}
