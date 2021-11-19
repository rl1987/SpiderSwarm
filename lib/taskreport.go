package spsw

import (
	"github.com/google/uuid"
)

type TaskReport struct {
	UUID      string
	JobUUID   string
	TaskUUID  string
	TaskName  string
	Succeeded bool
	Error     error
}

func NewTaskReport(jobUUID string, taskUUID string, taskName string, succeeded bool, err error) *TaskReport {
	return &TaskReport{
		UUID:      uuid.New().String(),
		JobUUID:   jobUUID,
		TaskUUID:  taskUUID,
		TaskName:  taskName,
		Succeeded: succeeded,
		Error:     err,
	}
}
