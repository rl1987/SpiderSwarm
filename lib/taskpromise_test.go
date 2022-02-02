package spsw

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

func TestTaskPromiseSplay(t *testing.T) {
	taskName := "HTTP1"
	workflowName := "testWorkflow"
	jobUUID := "97F34D30-7355-4C82-9480-A3B9CD086824"

	promise := &TaskPromise{
		TaskName:     taskName,
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
		InputDataChunksByInputName: map[string]*DataChunk{
			"param1": NewDataChunk_("aa"),
			"param2": NewDataChunk_(NewValueFromStrings([]string{"1", "2", "3"})),
		},
	}

	gotPromises := promise.Splay()

	for _, newPromise := range gotPromises {
		assert.Equal(t, promise.TaskName, newPromise.TaskName)
		assert.Equal(t, promise.WorkflowName, newPromise.WorkflowName)
		assert.Equal(t, promise.JobUUID, newPromise.JobUUID)

		assert.Equal(t, "aa", promise.InputDataChunksByInputName["param1"].PayloadValue.StringValue)
	}

	assert.Equal(t, "1", gotPromises[0].InputDataChunksByInputName["param2"].PayloadValue.StringValue)
	assert.Equal(t, "2", gotPromises[1].InputDataChunksByInputName["param2"].PayloadValue.StringValue)
	assert.Equal(t, "3", gotPromises[2].InputDataChunksByInputName["param2"].PayloadValue.StringValue)
}

func TestTaskPromiseSplay2(t *testing.T) {
	taskName := "HTTP2"
	workflowName := "testWorkflow"
	jobUUID := "97F34D30-7355-4C82-9480-A3B9CD086825"

	promise := &TaskPromise{
		TaskName:     taskName,
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
		InputDataChunksByInputName: map[string]*DataChunk{
			"param1": NewDataChunk_(NewValueFromString("aa")),
			"param2": NewDataChunk_(NewValueFromStrings([]string{"1", "2"})),
		},
	}

	gotPromises := promise.Splay()

	assert.Equal(t, 2, len(gotPromises))

	for _, newPromise := range gotPromises {
		assert.Equal(t, promise.TaskName, newPromise.TaskName)
		assert.Equal(t, promise.WorkflowName, newPromise.WorkflowName)
		assert.Equal(t, promise.JobUUID, newPromise.JobUUID)

		assert.Equal(t, "aa", promise.InputDataChunksByInputName["param1"].PayloadValue.StringValue)
	}

	assert.Equal(t, "1", gotPromises[0].InputDataChunksByInputName["param2"].PayloadValue.StringValue)
	assert.Equal(t, "2", gotPromises[1].InputDataChunksByInputName["param2"].PayloadValue.StringValue)
}

func TestTaskPromiseJSONE2E(t *testing.T) {
	promise := &TaskPromise{
		TaskName:     "HTTP1",
		WorkflowName: "TestFlow",
		JobUUID:      "",
		InputDataChunksByInputName: map[string]*DataChunk{
			"param1": NewDataChunk_(NewValueFromString("aa")),
			"param2": NewDataChunk_(NewValueFromStrings([]string{"1", "2"})),
		},
	}
	
	jsonBytes := promise.EncodeToJSON()

	gotPromise := NewTaskPromiseFromJSON(jsonBytes)

	assert.Equal(t, promise, gotPromise)
}

