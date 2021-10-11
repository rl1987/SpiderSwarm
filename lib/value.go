package spsw

import "net/http"

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
	BoolValue               bool
	IntValue                int
	StringValue             string
	StringsValue            []string
	MapStringToStringValue  map[string]string
	MapStringToStringsValue map[string][]string
	BytesValue              []byte
	HTTPHeadersValue        http.Header
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
