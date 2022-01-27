package spsw

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActionFromTemplate(t *testing.T) {
	actionTempl1 := &ActionTemplate{Name: "HTTPAction", StructName: "HTTPAction"}
	action1, ok1 := NewActionFromTemplate(actionTempl1).(*HTTPAction)

	assert.True(t, ok1)
	assert.NotNil(t, action1)
	assert.Equal(t, actionTempl1.Name, action1.Name)

	actionTempl2 := &ActionTemplate{
		Name:       "XPathAction",
		StructName: "XPathAction",
		ConstructorParams: map[string]Value{
			"xpath": Value{
				ValueType:   ValueTypeString,
				StringValue: "//title",
			},
			"expectMany": Value{
				ValueType: ValueTypeBool,
				BoolValue: false,
			},
		},
	}
	action2, ok2 := NewActionFromTemplate(actionTempl2).(*XPathAction)

	assert.True(t, ok2)
	assert.NotNil(t, action2)
	assert.Equal(t, actionTempl2.Name, action2.Name)

	actionTempl3 := &ActionTemplate{Name: "FieldJoinAction", StructName: "FieldJoinAction"}
	action3, ok3 := NewActionFromTemplate(actionTempl3).(*FieldJoinAction)

	assert.True(t, ok3)
	assert.NotNil(t, action3)
	assert.Equal(t, actionTempl3.Name, action3.Name)

	actionTempl4 := &ActionTemplate{
		Name:       "TaskPromiseAction",
		StructName: "TaskPromiseAction",
		ConstructorParams: map[string]Value{
			"inputNames": Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"page", "query"},
			},
		},
	}
	action4, ok4 := NewActionFromTemplate(actionTempl4).(*TaskPromiseAction)

	assert.True(t, ok4)
	assert.NotNil(t, action4)
	assert.Equal(t, actionTempl4.Name, action4.Name)
	assert.Equal(t, []string{"page", "query", TaskPromiseActionInputRefrain}, action4.AllowedInputNames)

	actionTempl5 := &ActionTemplate{Name: "UTF8DecodeAction", StructName: "UTF8DecodeAction"}
	action5, ok5 := NewActionFromTemplate(actionTempl5).(*UTF8DecodeAction)

	assert.True(t, ok5)
	assert.NotNil(t, action5)
	assert.Equal(t, actionTempl5.Name, action5.Name)

	actionTempl6 := &ActionTemplate{Name: "UTF8EncodeAction", StructName: "UTF8EncodeAction"}
	action6, ok6 := NewActionFromTemplate(actionTempl6).(*UTF8EncodeAction)

	assert.True(t, ok6)
	assert.NotNil(t, action6)
	assert.Equal(t, actionTempl6.Name, action6.Name)

	actionTempl7 := &ActionTemplate{Name: "ConstAction", StructName: "ConstAction"}
	action7, ok7 := NewActionFromTemplate(actionTempl7).(*ConstAction)

	assert.True(t, ok7)
	assert.NotNil(t, action7)
	assert.Equal(t, actionTempl7.Name, action7.Name)

	actionTempl8 := &ActionTemplate{Name: "URLJoinAction", StructName: "URLJoinAction"}
	action8, ok8 := NewActionFromTemplate(actionTempl8).(*URLJoinAction)

	assert.True(t, ok8)
	assert.NotNil(t, action8)
	assert.Equal(t, actionTempl8.Name, action8.Name)

	actionTempl9 := &ActionTemplate{Name: "StringCutAction", StructName: "StringCutAction"}
	action9, ok9 := NewActionFromTemplate(actionTempl9).(*StringCutAction)

	assert.True(t, ok9)
	assert.NotNil(t, action9)
	assert.Equal(t, actionTempl9.Name, action9.Name)

	actionTempl10 := &ActionTemplate{Name: "StringMapUpdateAction", StructName: "StringMapUpdateAction"}
	action10, ok10 := NewActionFromTemplate(actionTempl10).(*StringMapUpdateAction)

	assert.True(t, ok10)
	assert.NotNil(t, action10)
	assert.Equal(t, actionTempl10.Name, action10.Name)
}

func TestInputTableSortedness(t *testing.T) {
	for _, names := range AllowedInputNameTable {
		assert.True(t, sort.StringsAreSorted(names))
	}

	for _, names := range AllowedOutputNameTable {
		assert.True(t, sort.StringsAreSorted(names))
	}
}

func NewTestAction(*ActionTemplate) Action {
	return &AbstractAction{}
}

func TestRegisterAction(t *testing.T) {
	structName := "TestAction"
	allowedInputNames := []string{"in1", "in2"}
	allowedOutputNames := []string{"out1", "out2"}

	RegisterAction(structName, NewTestAction, allowedInputNames, allowedOutputNames)

	assert.NotNil(t, ActionConstructorTable[structName])
	assert.Equal(t, allowedInputNames, AllowedInputNameTable[structName])
	assert.Equal(t, allowedOutputNames, AllowedOutputNameTable[structName])
}
