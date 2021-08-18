package spiderswarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScheduledTask(t *testing.T) {
	taskPromise := &TaskPromise{UUID: "D412D565-B2A8-4BE3-B3CB-B37008FDA099"}
	taskTemplate := &TaskTemplate{TaskName: "testTask"}

	workflowName := "testWorkflow"
	workflowVersion := "2.0"
	jobUUID := "5369DD61-E98E-465E-9619-4641D06728FB"

	scheduledTask := NewScheduledTask(taskPromise, taskTemplate, workflowName, workflowVersion, jobUUID)

	assert.NotNil(t, scheduledTask)
	assert.Equal(t, *taskPromise, scheduledTask.Promise)
	assert.Equal(t, *taskTemplate, scheduledTask.Template)
	assert.Equal(t, workflowName, scheduledTask.WorkflowName)
	assert.Equal(t, workflowVersion, scheduledTask.WorkflowVersion)
	assert.Equal(t, jobUUID, scheduledTask.JobUUID)
}
