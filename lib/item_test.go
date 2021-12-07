package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemIsSplayable(t *testing.T) {
	item1 := &Item{
		Fields: map[string]*Value{
			"field1": &Value{
				ValueType:   ValueTypeString,
				StringValue: "a",
			},
			"field2": &Value{
				ValueType:   ValueTypeString,
				StringValue: "b",
			},
		},
	}

	assert.False(t, item1.IsSplayable())

	item2 := &Item{
		Fields: map[string]*Value{
			"field1": &Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"a", "b", "c"},
			},
			"field2": &Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"1", "2", "3"},
			},
		},
	}

	assert.True(t, item2.IsSplayable())

	item3 := &Item{
		Fields: map[string]*Value{
			"field1": &Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"x", "y", "z"},
			},
			"field2": &Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"0", "1"},
			},
		},
	}

	assert.False(t, item3.IsSplayable())

	item4 := &Item{
		Fields: map[string]*Value{
			"field1": &Value{
				ValueType:   ValueTypeString,
				StringValue: "a",
			},
			"field2": &Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{},
			},
		},
	}

	assert.False(t, item4.IsSplayable())
}

func TestItemSplay(t *testing.T) {
	simpleItem := &Item{
		Fields: map[string]*Value{
			"field1": &Value{ValueType: ValueTypeString, StringValue: "a"},
			"field2": &Value{ValueType: ValueTypeString, StringValue: "1"},
		},
	}

	items := simpleItem.Splay()

	assert.Equal(t, []*Item{simpleItem}, items)

	compositeItem := &Item{
		Fields: map[string]*Value{
			"field1": &Value{ValueType: ValueTypeStrings, StringsValue: []string{"a", "b", "c"}},
			"field2": &Value{ValueType: ValueTypeStrings, StringsValue: []string{"1", "2", "3"}},
			"field3": &Value{ValueType: ValueTypeString, StringValue: "C"},
		},
	}

	expectItems := []*Item{
		&Item{
			Fields: map[string]*Value{
				"field1": &Value{ValueType: ValueTypeString, StringValue: "a"},
				"field2": &Value{ValueType: ValueTypeString, StringValue: "1"},
				"field3": &Value{ValueType: ValueTypeString, StringValue: "C"},
			},
		},
		&Item{
			Fields: map[string]*Value{
				"field1": &Value{ValueType: ValueTypeString, StringValue: "b"},
				"field2": &Value{ValueType: ValueTypeString, StringValue: "2"},
				"field3": &Value{ValueType: ValueTypeString, StringValue: "C"},
			},
		},
		&Item{
			Fields: map[string]*Value{
				"field1": &Value{ValueType: ValueTypeString, StringValue: "c"},
				"field2": &Value{ValueType: ValueTypeString, StringValue: "3"},
				"field3": &Value{ValueType: ValueTypeString, StringValue: "C"},
			},
		},
	}

	items = compositeItem.Splay()

	assert.Equal(t, expectItems[0].Fields, items[0].Fields)
	assert.Equal(t, expectItems[1].Fields, items[1].Fields)
	assert.Equal(t, expectItems[2].Fields, items[2].Fields)
}

func TestItemSetField(t *testing.T) {
	item := &Item{
		Fields: map[string]*Value{},
	}

	item.SetField("testStr", "testStr")
	item.SetField("testStrings", []string{"1", "2"})
	item.SetField("testInt", 42)
	item.SetField("testBool", false)

	assert.Equal(t, 4, len(item.Fields))

	expectedFields := map[string]*Value{
		"testStr": &Value{
			ValueType:   ValueTypeString,
			StringValue: "testStr",
		},
		"testStrings": &Value{
			ValueType:    ValueTypeStrings,
			StringsValue: []string{"1", "2"},
		},
		"testInt": &Value{
			ValueType: ValueTypeInt,
			IntValue:  42,
		},
		"testBool": &Value{
			ValueType: ValueTypeBool,
			BoolValue: false,
		},
	}

	assert.Equal(t, expectedFields, item.Fields)

}
