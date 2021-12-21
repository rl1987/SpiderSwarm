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

func TestNewSpiderBusAdapterForExporter(t *testing.T) {
	spiderBus := NewSpiderBus()
	exporter := NewExporter()

	adapter := NewSpiderBusAdapterForExporter(spiderBus, exporter)

	assert.NotNil(t, adapter)

	assert.Equal(t, spiderBus, adapter.Bus)
	assert.Equal(t, exporter.ItemsIn, adapter.ItemsOut)
	assert.Nil(t, adapter.ScheduledTasksIn)
	assert.Nil(t, adapter.ScheduledTasksOut)
	assert.Nil(t, adapter.TaskPromisesIn)
	assert.Nil(t, adapter.TaskPromisesOut)
	assert.Nil(t, adapter.TaskResultsIn)
	assert.Nil(t, adapter.TaskResultsOut)
	assert.Nil(t, adapter.ItemsIn)
}

func TestNewSpiderBusAdapterForManager(t *testing.T) {
	spiderBus := NewSpiderBus()
	manager := NewManager(nil)

	adapter := NewSpiderBusAdapterForManager(spiderBus, manager)

	assert.NotNil(t, adapter)

	assert.Equal(t, spiderBus, adapter.Bus)
	assert.Equal(t, manager.TaskPromisesIn, adapter.TaskPromisesOut)
	assert.Equal(t, manager.TaskResultsIn, adapter.TaskResultsOut)
	assert.Equal(t, manager.ScheduledTasksOut, adapter.ScheduledTasksIn)
	assert.Equal(t, manager.ItemsOut, adapter.ItemsIn)
	assert.Nil(t, adapter.TaskPromisesIn)
	assert.Nil(t, adapter.ScheduledTasksOut)
	assert.Nil(t, adapter.TaskResultsIn)
	assert.Nil(t, adapter.ItemsOut)
}
