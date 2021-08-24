package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorker(t *testing.T) {
	worker := NewWorker()

	assert.NotNil(t, worker)
	assert.NotNil(t, worker.ScheduledTasksIn)
	assert.NotNil(t, worker.DataChunksOut)
	assert.NotNil(t, worker.Done)
}
