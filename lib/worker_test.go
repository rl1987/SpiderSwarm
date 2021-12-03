package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorker(t *testing.T) {
	worker := NewWorker()

	assert.NotNil(t, worker)
	assert.NotNil(t, worker.ScheduledTasksIn)
	assert.NotNil(t, worker.TaskPromisesOut)
	assert.NotNil(t, worker.Done)
}

func TestWorkerRunHeedDone(t *testing.T) {
	worker := NewWorker()

	finished := false

	go func() {
		worker.Run()

		finished = true
	}()

	worker.Done <- true

	assert.True(t, finished)
}

func TestWorkerExecuteTaskNoError(t *testing.T) {
	t.Skip()

	testTask := &Task{}

	worker := NewWorker()

	gotErr := worker.executeTask(testTask)

	assert.Nil(t, gotErr)
}
