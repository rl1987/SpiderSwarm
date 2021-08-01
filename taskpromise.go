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

func (tp *TaskPromise) IsSplayable() bool {
	hasLists := false
	equalLen := true
	lastLen := -1

	for _, chunk := range tp.InputDataChunksByInputName {
		if chunk.Type == DataChunkTypeStrings {
			hasLists = true

			if lastLen != -1 && lastLen != len(chunk.Payload.([]string)) {
				equalLen = false
				break
			}

			lastLen = len(chunk.Payload.([]string))
		}
	}

	return hasLists && equalLen && lastLen != 0
}
