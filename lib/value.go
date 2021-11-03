package spsw

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

const ValueTypeInt = "ValueTypeInt"
const ValueTypeBool = "ValueTypeBool"
const ValueTypeString = "ValueTypeString"
const ValueTypeStrings = "ValueTypeStrings"
const ValueTypeMapStringToString = "ValueTypeStringToString"
const ValueTypeMapStringToStrings = "ValueTypeStringToStrings"
const ValueTypeBytes = "ValueTypeBytes"
const ValueTypeHTTPHeaders = "ValueTypeHTTPHeaders"

type Value struct {
	ValueType               string
	BoolValue               bool                `yaml:"BoolValue,omitempty"`
	IntValue                int                 `yaml:"IntValue,omitempty"`
	StringValue             string              `yaml:"StringValue,omitempty"`
	StringsValue            []string            `yaml:"StringsValue,omitempty"`
	MapStringToStringValue  map[string]string   `yaml:"MapStringToStringValue,omitempty"`
	MapStringToStringsValue map[string][]string `yaml:"MapStringToStringsValue,omitempty"`
	BytesValue              []byte              `yaml:"BytesValue,omitempty"`
	HTTPHeadersValue        http.Header         `yaml:"HTTPHeadersValue,omitempty"`
}

func NewValueFromInt(i int) *Value {
	return &Value{
		ValueType: ValueTypeInt,
		IntValue:  i,
	}
}

func NewValueFromBool(b bool) *Value {
	return &Value{
		ValueType: ValueTypeBool,
		BoolValue: b,
	}
}

func NewValueFromString(s string) *Value {
	return &Value{
		ValueType:   ValueTypeString,
		StringValue: s,
	}
}

func NewValueFromStrings(s []string) *Value {
	return &Value{
		ValueType:    ValueTypeStrings,
		StringsValue: s,
	}
}

func NewValueFromMapStringToString(m map[string]string) *Value {
	return &Value{
		ValueType:              ValueTypeMapStringToString,
		MapStringToStringValue: m,
	}
}

func NewValueFromMapStringToStrings(m map[string][]string) *Value {
	return &Value{
		ValueType:               ValueTypeMapStringToStrings,
		MapStringToStringsValue: m,
	}
}

func NewValueFromBytes(b []byte) *Value {
	return &Value{
		ValueType:  ValueTypeBytes,
		BytesValue: b,
	}
}

func NewValueFromHTTPHeaders(h http.Header) *Value {
	return &Value{
		ValueType:        ValueTypeHTTPHeaders,
		HTTPHeadersValue: h,
	}
}

func (value *Value) Hash() []byte {
	h := sha256.New()

	h.Write([]byte(value.ValueType))
	h.Write([]byte(fmt.Sprintf("%v", value.GetUnderlyingValue())))

	return h.Sum(nil)
}

func (value *Value) String() string {
	return fmt.Sprintf("<Value %v>", value.GetUnderlyingValue())
}

func (value *Value) GetUnderlyingValue() interface{} {
	if value.ValueType == ValueTypeInt {
		return value.IntValue
	} else if value.ValueType == ValueTypeBool {
		return value.BoolValue
	} else if value.ValueType == ValueTypeString {
		return value.StringValue
	} else if value.ValueType == ValueTypeStrings {
		return value.StringsValue
	} else if value.ValueType == ValueTypeMapStringToString {
		return value.MapStringToStringValue
	} else if value.ValueType == ValueTypeMapStringToStrings {
		return value.MapStringToStringsValue
	} else if value.ValueType == ValueTypeBytes {
		return value.BytesValue
	} else if value.ValueType == ValueTypeHTTPHeaders {
		return value.HTTPHeadersValue
	}

	return nil
}
