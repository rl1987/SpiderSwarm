package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.TaskPromisesIn)
	assert.NotNil(t, manager.ScheduledTasksOut)
	assert.Nil(t, manager.CurrentWorkflow)
}

func TestManagerStartScrapingJob(t *testing.T) {
	manager := NewManager()

	workflow := &Workflow{}

	manager.StartScrapingJob(workflow)

	assert.Equal(t, workflow, manager.CurrentWorkflow)
}
