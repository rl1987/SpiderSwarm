package spsw

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

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

func NewTaskReportFromJSON(raw []byte) *TaskReport {
	taskReport := &TaskReport{}

	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(taskReport)
	if err != nil {
		return nil
	}

	return taskReport
}

func (tr *TaskReport) EncodeToJSON() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(tr)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}
