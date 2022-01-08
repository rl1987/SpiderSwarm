package spsw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskResult(t *testing.T) {
	jobUUID := "A0"
	taskUUID := "B1"
	scheduledTaskUUID := "C2"
	succeeded := false
	err := errors.New("Test error")

	taskResult := NewTaskResult(jobUUID, taskUUID, scheduledTaskUUID, succeeded, err)

	assert.NotNil(t, taskResult)
	assert.Equal(t, jobUUID, taskResult.JobUUID)
	assert.Equal(t, taskUUID, taskResult.TaskUUID)
	assert.Equal(t, scheduledTaskUUID, taskResult.ScheduledTaskUUID)
	assert.Equal(t, succeeded, taskResult.Succeeded)
	assert.Equal(t, err, taskResult.Error)

	assert.NotNil(t, taskResult.OutputDataChunks)
	assert.Equal(t, 0, len(taskResult.OutputDataChunks))
}

func TestTaskResultAddOutputItem(t *testing.T) {
	jobUUID := "A0"
	taskUUID := "B1"
	scheduledTaskUUID := "C2"
	succeeded := true

	taskResult := NewTaskResult(jobUUID, taskUUID, scheduledTaskUUID, succeeded, nil)

	item := &Item{Name: "TestItem"}

	taskResult.AddOutputItem("testOut", item)

	assert.Equal(t, 1, len(taskResult.OutputDataChunks["testOut"]))

	chunk := taskResult.OutputDataChunks["testOut"][0]

	assert.Equal(t, DataChunkTypeItem, chunk.Type)
	assert.Equal(t, item, chunk.PayloadItem)
}

func TestTaskResultAddOutputTaskPromise(t *testing.T) {
	jobUUID := "A0"
	taskUUID := "B1"
	scheduledTaskUUID := "C2"
	succeeded := true

	taskResult := NewTaskResult(jobUUID, taskUUID, scheduledTaskUUID, succeeded, nil)

	promise1 := &TaskPromise{TaskName: "DoIt"}
	taskResult.AddOutputTaskPromise("promises", promise1)
	promise2 := &TaskPromise{TaskName: "DoItAgain"}
	taskResult.AddOutputTaskPromise("promises", promise2)

	assert.Equal(t, 2, len(taskResult.OutputDataChunks["promises"]))
	assert.Equal(t, "DoIt", taskResult.OutputDataChunks["promises"][0].PayloadPromise.TaskName)
	assert.Equal(t, "DoItAgain", taskResult.OutputDataChunks["promises"][1].PayloadPromise.TaskName)
}
