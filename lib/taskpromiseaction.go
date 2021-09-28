package spsw

import (
	"errors"

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
		TaskName:     taskName,
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
	}
}

func NewTaskPromiseActionFromTemplate(actionTempl *ActionTemplate, workflowName string) *TaskPromiseAction {
	var inputNames []string
	var taskName string

	inputNames, _ = actionTempl.ConstructorParams["inputNames"].([]string)
	taskName, _ = actionTempl.ConstructorParams["taskName"].(string)

	if inputNames == nil {
		// HACK to work around the issues of inputNames sometimes being of type
		// []interface{}
		inputNamesIntf, okHack := actionTempl.ConstructorParams["inputNames"].([]interface{})
		if okHack && len(inputNamesIntf) > 0 {
			if _, okStr := inputNamesIntf[0].(string); okStr {
				inputNames = []string{}
				for i, _ := range inputNamesIntf {
					s := inputNamesIntf[i].(string)
					inputNames = append(inputNames, s)
				}
			}
		}
	}

	if inputNames == nil {
		panic("Fatal error in NewTaskPromiseActionFromTemplate: inputNames is nil")
	}

	action := NewTaskPromiseAction(inputNames, taskName, workflowName, "")

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
