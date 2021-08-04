package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskPromiseActionFromTemplate(t *testing.T) {
	taskName := "HTTP2"

	actionTempl := &ActionTemplate{
		Name:       "testAction",
		StructName: "TaskPromiseAction",
		ConstructorParams: map[string]interface{}{
			"inputNames": []string{"page", "session"},
			"taskName":   taskName,
		},
	}

	workflow := &Workflow{
		Name: "testWorkflow",
	}

	action := NewTaskPromiseActionFromTemplate(actionTempl, workflow)

	assert.NotNil(t, action)
	assert.Equal(t, taskName, action.TaskName)
	assert.Equal(t, actionTempl.ConstructorParams["inputNames"], action.AllowedInputNames)
	assert.Equal(t, workflow.Name, action.WorkflowName)
	assert.Equal(t, []string{TaskPromiseActionOutputPromise}, action.AllowedOutputNames)
}
