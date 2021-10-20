package spsw

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/google/uuid"
)

type ScheduledTask struct {
	UUID            string
	Promise         TaskPromise
	Template        TaskTemplate
	WorkflowName    string
	WorkflowVersion string
	JobUUID         string
}

func NewScheduledTask(promise *TaskPromise, template *TaskTemplate, workflowName string, workflowVersion string, jobUUID string) *ScheduledTask {
	return &ScheduledTask{
		UUID:            uuid.New().String(),
		Promise:         *promise,
		Template:        *template,
		WorkflowName:    workflowName,
		WorkflowVersion: workflowVersion,
		JobUUID:         jobUUID,
	}
}

func NewScheduledTaskFromJSON(raw []byte) *ScheduledTask {
	scheduledTask := &ScheduledTask{}

	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(scheduledTask)
	if err != nil {
		return nil
	}

	return scheduledTask
}

func (st *ScheduledTask) EncodeToJSON() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(st)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}
