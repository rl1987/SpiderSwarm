package spiderswarm

import (
	"github.com/google/uuid"
)

type ScheduledTask struct {
	UUID            string
	Promise         *TaskPromise
	Template        *TaskTemplate
	WorkflowName    string
	WorkflowVersion string
	JobUUID         string
}

func NewScheduledTask(promise *TaskPromise, template *TaskTemplate, workflowName string, workflowVersion string, jobUUID string) *ScheduledTask {
	return &ScheduledTask{
		UUID:            uuid.New().String(),
		Promise:         promise,
		Template:        template,
		WorkflowName:    workflowName,
		WorkflowVersion: workflowVersion,
		JobUUID:         jobUUID,
	}
}
