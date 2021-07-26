package main

import (
	"time"

	"github.com/google/uuid"
)

type TaskPromise struct {
	UUID                       string
	TaskName                   string
	WorkflowName               string
	JobUUID                    string
	InputDataChunksByInputName map[string]*DataChunk
	CreatedAt                  time.Time
}

func NewTaskPromise(taskName string, workflowName string, jobUUID string, inputDataChunksByInputName map[string]*DataChunk) *TaskPromise {
	return &TaskPromise{
		UUID:                       uuid.New().String(),
		TaskName:                   taskName,
		WorkflowName:               workflowName,
		JobUUID:                    jobUUID,
		InputDataChunksByInputName: inputDataChunksByInputName,
		CreatedAt:                  time.Now(),
	}
}
