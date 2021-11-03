package spsw

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func NewTaskPromiseFromJSON(raw []byte) *TaskPromise {
	taskPromise := &TaskPromise{}

	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(taskPromise)
	if err != nil {
		return nil
	}

	return taskPromise
}

func (tp *TaskPromise) Hash() []byte {
	h := sha256.New()

	h.Write([]byte(tp.TaskName))
	h.Write([]byte(tp.WorkflowName))

	for key, dc := range tp.InputDataChunksByInputName {
		h.Write([]byte(key))
		h.Write(dc.Hash())
	}

	return h.Sum(nil)
}

func (tp *TaskPromise) String() string {
	return fmt.Sprintf("<TaskPromise %s TaskName: %s, WorkflowName: %s, JobUUID: %s, InputDataChunksByInputName: %v, CreatedAt: %v",
		tp.UUID, tp.TaskName, tp.WorkflowName, tp.JobUUID, tp.InputDataChunksByInputName, tp.CreatedAt)
}

func (tp *TaskPromise) IsSplayable() bool {
	hasLists := false
	equalLen := true
	lastLen := -1

	for _, chunk := range tp.InputDataChunksByInputName {
		if chunk.Type == DataChunkTypeValue {
			chunkValue := chunk.PayloadValue

			if chunkValue.ValueType == ValueTypeStrings {
				hasLists = true

				if lastLen != -1 && lastLen != len(chunkValue.StringsValue) {
					equalLen = false
					break
				}

				lastLen = len(chunkValue.StringsValue)
			}
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
		if chunk.Type == DataChunkTypeValue {
			chunkValue := chunk.PayloadValue
			if chunkValue.ValueType == ValueTypeStrings {
				var s string
				s, chunkValue.StringsValue = chunkValue.StringsValue[0], chunkValue.StringsValue[1:]
				newChunk, _ := NewDataChunk(s)
				newPromise.InputDataChunksByInputName[name] = newChunk
			} else {
				newPromise.InputDataChunksByInputName[name] = chunk
			}
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

func (tp *TaskPromise) EncodeToJSON() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(tp)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}
