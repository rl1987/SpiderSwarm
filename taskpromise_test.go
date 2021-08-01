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

func NewDataChunk_(payload interface{}) *DataChunk {
	chunk, _ := NewDataChunk(payload)
	return chunk
}

func TestTaskPromiseIsSplayable(t *testing.T) {
	promise1 := &TaskPromise{
		InputDataChunksByInputName: map[string]*DataChunk{
			"a": NewDataChunk_("a"),
			"b": NewDataChunk_([]string{"1", "2"}),
		},
	}

	assert.True(t, promise1.IsSplayable())

	promise2 := &TaskPromise{
		InputDataChunksByInputName: map[string]*DataChunk{
			"a": NewDataChunk_("a"),
			"b": NewDataChunk_([]string{"1", "2"}),
			"c": NewDataChunk_([]string{"x", "y"}),
		},
	}

	assert.True(t, promise2.IsSplayable())

	promise3 := &TaskPromise{
		InputDataChunksByInputName: map[string]*DataChunk{
			"a": NewDataChunk_("a"),
			"b": NewDataChunk_([]string{"1", "2"}),
			"c": NewDataChunk_([]string{"x", "y", "z"}),
		},
	}

	assert.False(t, promise3.IsSplayable())

	promise4 := &TaskPromise{
		InputDataChunksByInputName: map[string]*DataChunk{
			"a": NewDataChunk_("a"),
			"b": NewDataChunk_([]string{}),
		},
	}

	assert.False(t, promise4.IsSplayable())
}
