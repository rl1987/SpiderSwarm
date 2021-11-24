package spsw

import (
//"github.com/google/uuid"
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
