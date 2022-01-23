package spsw

import (
  	"net/http"
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

func TestNewValue(t *testing.T) {
	table := []struct{
		Input interface{}
		ExpectedValue *Value
	}{
		{42, &Value{ValueType: ValueTypeInt, IntValue: 42}},
		{false, &Value{ValueType: ValueTypeBool, BoolValue: false}},
		{"s", &Value{ValueType: ValueTypeString, StringValue: "s"}},
		{[]string{"a", "b"}, &Value{ValueType: ValueTypeStrings, StringsValue: []string{"a", "b"}}},
		{map[string]string{"a":"b"}, &Value{ValueType: ValueTypeMapStringToString, MapStringToStringValue: map[string]string{"a": "b"}}},
		{map[string][]string{"x": []string{"1", "2"}}, &Value{ValueType: ValueTypeMapStringToStrings, MapStringToStringsValue: map[string][]string{"x": []string{"1", "2"}}}},
		{[]byte("\xde\xea\xbe\xef"), &Value{ValueType: ValueTypeBytes, BytesValue: []byte("\xde\xea\xbe\xef")}},
		{http.Header{}, &Value{ValueType: ValueTypeHTTPHeaders, HTTPHeadersValue: http.Header{}}},
		{4.2, nil},
	}

	for _, entry := range table {
		value := NewValue(entry.Input)
		assert.Equal(t, entry.ExpectedValue, value)
	}

}
