package main

import (
	"github.com/google/uuid"
)

type TaskPromiseAction struct {
	AbstractAction
	UUID         string
	TaskName     string
	WorkflowName string
	JobUUID      string
}

const TaskPromiseActionOutputPromise = "TaskPromiseActionOutputPromise"

func NewTaskPromiseAction(inputNames []string, taskName string, workflowName string, jobUUID string) *TaskPromiseAction {
	return &TaskPromiseAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  inputNames,
			AllowedOutputNames: []string{TaskPromiseActionOutputPromise},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
	}
}
