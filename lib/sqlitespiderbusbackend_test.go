package spiderswarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLiteSpiderBusBackendScheduledTaskE2E(t *testing.T) {
	taskPromise := &TaskPromise{UUID: "D412D565-B2A8-4BE3-B3CB-B37008FDA099"}
	taskTemplate := &TaskTemplate{TaskName: "testTask"}

	workflowName := "testWorkflow"
	workflowVersion := "2.0"
	jobUUID := "5369DD61-E98E-465E-9619-4641D06728FB"

	scheduledTask := NewScheduledTask(taskPromise, taskTemplate, workflowName, workflowVersion, jobUUID)

	assert.NotNil(t, scheduledTask)

	backend := NewSQLiteSpiderBusBackend("")

	assert.NotNil(t, backend)

	gotScheduledTask := backend.ReceiveScheduledTask()
	assert.Nil(t, gotScheduledTask)

	err := backend.SendScheduledTask(scheduledTask)
	assert.Nil(t, err)

	gotScheduledTask = backend.ReceiveScheduledTask()
	assert.Equal(t, scheduledTask, gotScheduledTask)

	gotScheduledTask = backend.ReceiveScheduledTask()
	assert.Nil(t, gotScheduledTask)

}
