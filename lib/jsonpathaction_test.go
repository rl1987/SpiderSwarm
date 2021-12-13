package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
