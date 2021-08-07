package spiderswarm

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
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

func (tp *TaskPromise) splayOff() *TaskPromise {
	newPromise := &TaskPromise{
		UUID:                       uuid.New().String(),
		TaskName:                   tp.TaskName,
		WorkflowName:               tp.WorkflowName,
		JobUUID:                    tp.JobUUID,
		InputDataChunksByInputName: map[string]*DataChunk{},
		CreatedAt:                  time.Now(),
	}

	for name, chunk := range tp.InputDataChunksByInputName {
		if chunk.Type == DataChunkTypeStrings {
			var s string
			s, chunk.Payload = chunk.Payload.([]string)[0], chunk.Payload.([]string)[1:]
			newChunk, _ := NewDataChunk(s)
			newPromise.InputDataChunksByInputName[name] = newChunk
		} else {
			newPromise.InputDataChunksByInputName[name] = chunk
		}
	}

	return newPromise
}

func (tp *TaskPromise) Splay() []*TaskPromise {
	if !tp.IsSplayable() {
		return []*TaskPromise{tp}
	}

	promises := []*TaskPromise{}

	promiseCopy := &TaskPromise{}

	copier.Copy(promiseCopy, tp)

	for {
		if !tp.IsSplayable() {
			break
		}

		newPromise := promiseCopy.splayOff()

		promises = append(promises, newPromise)
	}

	return promises
}
