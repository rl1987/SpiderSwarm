package spiderswarm

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

func NewTaskPromiseActionFromTemplate(actionTempl *ActionTemplate, workflow *Workflow) *TaskPromiseAction {
	var inputNames []string
	var taskName string

	inputNames, _ = actionTempl.ConstructorParams["inputNames"].([]string)
	taskName, _ = actionTempl.ConstructorParams["taskName"].(string)

	action := NewTaskPromiseAction(inputNames, taskName, workflow.Name, "")

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
