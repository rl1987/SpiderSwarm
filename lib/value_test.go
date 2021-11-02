package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValueFromInt(t *testing.T) {
	expectValue := &Value{
		ValueType: ValueTypeInt,
		IntValue:  42,
	}

	value := NewValueFromInt(42)

	assert.Equal(t, expectValue, value)
}

func TestNewValueFromBool(t *testing.T) {
	expectValue := &Value{
		ValueType: ValueTypeBool,
		BoolValue: true,
	}

	value := NewValueFromBool(true)

	assert.Equal(t, expectValue, value)
}
