package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskPromise(t *testing.T) {
	taskName := "testTask"
	workflowName := "testWorkflow"
	jobUUID := "AB5A4B8F-2815-4870-A928-25A0A3E965A1"

	dataChunk, _ := NewDataChunk(map[string][]string{
		"q": []string{"Free Julian Assange"},
	})

	inputDataChunksByInputName := map[string]*DataChunk{
		HTTPActionInputURLParams: dataChunk,
	}

	taskPromise := NewTaskPromise(taskName, workflowName, jobUUID, inputDataChunksByInputName)

	assert.NotNil(t, taskPromise)
	assert.Equal(t, taskName, taskPromise.TaskName)
	assert.Equal(t, workflowName, taskPromise.WorkflowName)
	assert.Equal(t, jobUUID, taskPromise.JobUUID)
	assert.Equal(t, inputDataChunksByInputName, taskPromise.InputDataChunksByInputName)
}
