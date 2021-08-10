package spiderswarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActionFromTemplate(t *testing.T) {
	workflow := &Workflow{}
	jobUUID := "B3FB46CC-E3AC-41FB-AC88-A15B23C2F299"

	actionTempl1 := &ActionTemplate{Name: "HTTPAction", StructName: "HTTPAction"}
	action1, ok1 := NewActionFromTemplate(actionTempl1, workflow, jobUUID).(*HTTPAction)

	assert.True(t, ok1)
	assert.NotNil(t, action1)
	assert.Equal(t, actionTempl1.Name, action1.Name)

	actionTempl2 := &ActionTemplate{
		Name:       "XPathAction",
		StructName: "XPathAction",
		ConstructorParams: map[string]interface{}{
			"xpath":      "//title",
			"expectMany": false,
		},
	}
	action2, ok2 := NewActionFromTemplate(actionTempl2, workflow, jobUUID).(*XPathAction)

	assert.True(t, ok2)
	assert.NotNil(t, action2)
	assert.Equal(t, actionTempl2.Name, action2.Name)

	actionTempl3 := &ActionTemplate{Name: "FieldJoinAction", StructName: "FieldJoinAction"}
	action3, ok3 := NewActionFromTemplate(actionTempl3, workflow, jobUUID).(*FieldJoinAction)

	assert.True(t, ok3)
	assert.NotNil(t, action3)
	assert.Equal(t, actionTempl3.Name, action3.Name)

	actionTempl4 := &ActionTemplate{Name: "TaskPromiseAction", StructName: "TaskPromiseAction"}
	action4, ok4 := NewActionFromTemplate(actionTempl4, workflow, jobUUID).(*TaskPromiseAction)

	assert.True(t, ok4)
	assert.NotNil(t, action4)
	assert.Equal(t, actionTempl4.Name, action4.Name)

	actionTempl5 := &ActionTemplate{Name: "UTF8DecodeAction", StructName: "UTF8DecodeAction"}
	action5, ok5 := NewActionFromTemplate(actionTempl5, workflow, jobUUID).(*UTF8DecodeAction)

	assert.True(t, ok5)
	assert.NotNil(t, action5)
	assert.Equal(t, actionTempl5.Name, action5.Name)

	actionTempl6 := &ActionTemplate{Name: "UTF8EncodeAction", StructName: "UTF8EncodeAction"}
	action6, ok6 := NewActionFromTemplate(actionTempl6, workflow, jobUUID).(*UTF8EncodeAction)

	assert.True(t, ok6)
	assert.NotNil(t, action6)
	assert.Equal(t, actionTempl6.Name, action6.Name)

	actionTempl7 := &ActionTemplate{Name: "ConstAction", StructName: "ConstAction"}
	action7, ok7 := NewActionFromTemplate(actionTempl7, workflow, jobUUID).(*ConstAction)

	assert.True(t, ok7)
	assert.NotNil(t, action7)
	assert.Equal(t, actionTempl7.Name, action7.Name)

	actionTempl8 := &ActionTemplate{Name: "URLJoinAction", StructName: "URLJoinAction"}
	action8, ok8 := NewActionFromTemplate(actionTempl8, workflow, jobUUID).(*URLJoinAction)

	assert.True(t, ok8)
	assert.NotNil(t, action8)
	assert.Equal(t, actionTempl8.Name, action8.Name)

	actionTempl9 := &ActionTemplate{Name: "StringCutAction", StructName: "StringCutAction"}
	action9, ok9 := NewActionFromTemplate(actionTempl9, workflow, jobUUID).(*StringCutAction)

	assert.True(t, ok9)
	assert.NotNil(t, action9)
	assert.Equal(t, actionTempl9.Name, action9.Name)

}
