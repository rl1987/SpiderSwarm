package spsw

const ValueTypeInt = "ValueTypeInt"
const ValueTypeBool = "ValueTypeBool"
const ValueTypeString = "ValueTypeString"
const ValueTypeStrings = "ValueTypeStrings"

type Value struct {
	ValueType    string
	BoolValue    bool
	IntValue     int
	StringValue  string
	StringsValue []string
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
