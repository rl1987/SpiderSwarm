package spsw

import (
	"errors"

	"github.com/google/uuid"
)

type TaskPromiseAction struct {
	AbstractAction
	UUID          string
	TaskName      string
	WorkflowName  string
	JobUUID       string
	RequireFields []string
}

const TaskPromiseActionOutputPromise = "TaskPromiseActionOutputPromise"

func NewTaskPromiseAction(inputNames []string, taskName string, workflowName string, jobUUID string, requireFields []string) *TaskPromiseAction {
	return &TaskPromiseAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  inputNames,
			AllowedOutputNames: []string{TaskPromiseActionOutputPromise},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
		TaskName:      taskName,
		WorkflowName:  workflowName,
		JobUUID:       jobUUID,
		RequireFields: requireFields,
	}
}

func NewTaskPromiseActionFromTemplate(actionTempl *ActionTemplate, workflowName string) *TaskPromiseAction {
	var inputNames []string
	var taskName string
	var requireFields []string

	inputNames = actionTempl.ConstructorParams["inputNames"].StringsValue
	taskName = actionTempl.ConstructorParams["taskName"].StringValue
	requireFields = actionTempl.ConstructorParams["requireFields"].StringsValue

	action := NewTaskPromiseAction(inputNames, taskName, workflowName, "", requireFields)

	action.Name = actionTempl.Name

	return action
}

func (tpa *TaskPromiseAction) Run() error {
	inputDataChunksByInputName := map[string]*DataChunk{}

	for name, input := range tpa.Inputs {
		if len(input.Queue) > 0 {
			x := input.Remove()
			newChunk, _ := NewDataChunk(x)
			inputDataChunksByInputName[name] = newChunk
		}
	}

	for _, rf := range tpa.RequireFields {
		if inputDataChunksByInputName[rf] == nil ||
			(len(inputDataChunksByInputName[rf].PayloadValue.StringValue) == 0 && // XXX: this seems bit awkward
				len(inputDataChunksByInputName[rf].PayloadValue.StringsValue) == 0) {
			return nil
		}
	}

	if len(tpa.Inputs) == 0 {
		return errors.New("No inputs connected")
	}

	if tpa.Outputs[TaskPromiseActionOutputPromise] == nil || len(tpa.Outputs[TaskPromiseActionOutputPromise]) == 0 {
		return errors.New("No outputs connected")
	}

	promise := NewTaskPromise(tpa.TaskName, tpa.WorkflowName, tpa.JobUUID, inputDataChunksByInputName)

	for _, output := range tpa.Outputs[TaskPromiseActionOutputPromise] {
		output.Add(promise)
	}

	return nil
}
