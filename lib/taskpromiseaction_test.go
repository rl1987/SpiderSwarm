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

	workflow := &Workflow{
		Name: "testWorkflow",
	}

	action := NewTaskPromiseActionFromTemplate(actionTempl, workflow.Name)

	assert.NotNil(t, action)
	assert.Equal(t, taskName, action.TaskName)
	assert.Equal(t, actionTempl.ConstructorParams["inputNames"], action.AllowedInputNames)
	assert.Equal(t, workflow.Name, action.WorkflowName)
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
				ValueType:   ValueTypeString,
				StringValue: taskName,
			},
		},
	}

	workflow := &Workflow{
		Name: "testWorkflow",
	}

	action := NewTaskPromiseActionFromTemplate(actionTempl, workflow.Name)

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
	assert.Equal(t, workflow.Name, promise.WorkflowName)
	assert.Equal(t, page, promise.InputDataChunksByInputName["page"].Payload)
	assert.Equal(t, session, promise.InputDataChunksByInputName["session"].Payload)
}
