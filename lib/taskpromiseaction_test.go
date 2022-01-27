package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskPromiseActionFromTemplate(t *testing.T) {
	taskName := "HTTP2"

	actionTempl := &ActionTemplate{
		Name:       "testAction",
		StructName: "TaskPromiseAction",
		ConstructorParams: map[string]Value{
			"inputNames": Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"page", "session"},
			},
			"taskName": Value{
				ValueType:   ValueTypeString,
				StringValue: taskName,
			},
		},
	}

	action, ok := NewTaskPromiseActionFromTemplate(actionTempl).(*TaskPromiseAction)
	assert.True(t, ok)

	expectInputNames := []string{"page", "session", TaskPromiseActionInputRefrain}

	assert.NotNil(t, action)
	assert.Equal(t, taskName, action.TaskName)
	assert.Equal(t, expectInputNames, action.AllowedInputNames)
	assert.Equal(t, []string{TaskPromiseActionOutputPromise}, action.AllowedOutputNames)
}

func TestTaskPromiseActionRun(t *testing.T) {
	taskName := "HTTP2"

	actionTempl := &ActionTemplate{
		Name:       "testAction",
		StructName: "TaskPromiseAction",
		ConstructorParams: map[string]Value{
			"inputNames": Value{
				ValueType:    ValueTypeString,
				StringsValue: []string{"page", "session"},
			},
			"taskName": Value{
				ValueType:   ValueTypeStrings,
				StringValue: taskName,
			},
		},
	}

	action := NewTaskPromiseActionFromTemplate(actionTempl)

	pageIn := NewDataPipe()
	sessionIn := NewDataPipe()
	promiseOut := NewDataPipe()

	err := action.AddInput("page", pageIn)
	assert.Nil(t, err)

	err = action.AddInput("session", sessionIn)
	assert.Nil(t, err)

	err = action.AddOutput(TaskPromiseActionOutputPromise, promiseOut)
	assert.Nil(t, err)

	page := "2"
	session := "session_111"

	pageIn.Add(page)
	sessionIn.Add(session)

	err = action.Run()
	assert.Nil(t, err)

	promise, ok := promiseOut.Remove().(*TaskPromise)
	assert.True(t, ok)
	assert.Equal(t, taskName, promise.TaskName)
	assert.Equal(t, page, promise.InputDataChunksByInputName["page"].PayloadValue.StringValue)
	assert.Equal(t, session, promise.InputDataChunksByInputName["session"].PayloadValue.StringValue)
}

func TestTaskPromiseActionRunRequireFields(t *testing.T) {
	taskName := "HTTP3"

	actionTempl := &ActionTemplate{
		Name:       "testAction",
		StructName: "TaskPromiseAction",
		ConstructorParams: map[string]Value{
			"inputNames": Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"page", "session"},
			},
			"taskName": Value{
				ValueType:   ValueTypeString,
				StringValue: taskName,
			},
			"requireFields": Value{
				ValueType: ValueTypeStrings,
				StringsValue: []string{"page", "session"},
			},
		},
	}

	action := NewTaskPromiseActionFromTemplate(actionTempl)

	pageIn := NewDataPipe()
	sessionIn := NewDataPipe()
	promiseOut := NewDataPipe()

	err := action.AddInput("page", pageIn)
	assert.Nil(t, err)

	err = action.AddInput("session", sessionIn)
	assert.Nil(t, err)

	err = action.AddOutput(TaskPromiseActionOutputPromise, promiseOut)
	assert.Nil(t, err)

	page := "2"
	session := "session_111"

	pageIn.Add(page)

	err = action.Run()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(promiseOut.Queue))

	pageIn.Add(page)
	sessionIn.Add(session)

	err = action.Run()
	assert.Nil(t, err)

	promise, ok := promiseOut.Remove().(*TaskPromise)
	assert.True(t, ok)
	assert.Equal(t, taskName, promise.TaskName)
}

