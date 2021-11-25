package spsw

import (
	"github.com/google/uuid"
)

// TODO: use this instead of TaskReport
type TaskResult struct {
	UUID              string
	JobUUID           string
	TaskUUID          string
	ScheduledTaskUUID string
	Succeeded         bool
	Error             error
}

func NewTaskResult(jobUUID string, taskUUID string, succeeded bool, err error) *TaskResult {
	return &TaskResult{
		UUID:      uuid.New().String(),
		JobUUID:   jobUUID,
		TaskUUID:  taskUUID,
		Succeeded: succeeded,
		Error:     err,
	}
}
