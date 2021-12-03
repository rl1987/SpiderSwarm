package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpiderBusAdapterForWorker(t *testing.T) {
	spiderBus := NewSpiderBus()
	worker := NewWorker()

	adapter := NewSpiderBusAdapterForWorker(spiderBus, worker)

	assert.NotNil(t, adapter)

	assert.Equal(t, spiderBus, adapter.Bus)
	assert.Equal(t, worker.ScheduledTasksIn, adapter.ScheduledTasksOut)
	assert.Equal(t, worker.TaskPromisesOut, adapter.TaskPromisesIn)
	assert.Nil(t, adapter.ScheduledTasksIn)
	assert.Nil(t, adapter.TaskPromisesOut)
	assert.Nil(t, adapter.ItemsOut)
}
