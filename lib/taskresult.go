package spsw

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

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
	OutputDataChunks  map[string][]*DataChunk
}

func NewTaskResult(jobUUID string, taskUUID string, scheduledTaskUUID string, succeeded bool, err error) *TaskResult {
	return &TaskResult{
		UUID:              uuid.New().String(),
		JobUUID:           jobUUID,
		TaskUUID:          taskUUID,
		ScheduledTaskUUID: scheduledTaskUUID,
		Succeeded:         succeeded,
		Error:             err,
		OutputDataChunks:  map[string][]*DataChunk{},
	}
}

func NewTaskResultFromJSON(raw []byte) *TaskResult {
	taskResult := &TaskResult{}

	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(taskResult)
	if err != nil {
		return nil
	}

	return taskResult
}

func (tr *TaskResult) EncodeToJSON() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(tr)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}
