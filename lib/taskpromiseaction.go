package spsw

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type TaskPromiseAction struct {
	AbstractAction
	UUID          string
	TaskName      string
	WorkflowName  string
	JobUUID       string
	RequireFields []string
}

const TaskPromiseActionInputRefrain = "TaskPromiseActionInputRefrain"

const TaskPromiseActionOutputPromise = "TaskPromiseActionOutputPromise"

func NewTaskPromiseAction(inputNames []string, taskName string, jobUUID string, requireFields []string) *TaskPromiseAction {
	if inputNames == nil {
		inputNames = []string{}
	}

	return &TaskPromiseAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  append(inputNames, TaskPromiseActionInputRefrain),
			AllowedOutputNames: []string{TaskPromiseActionOutputPromise},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
		TaskName:      taskName,
		JobUUID:       jobUUID,
		RequireFields: requireFields,
	}
}

func NewTaskPromiseActionFromTemplate(actionTempl *ActionTemplate) Action {
	var inputNames []string
	var taskName string
	var requireFields []string

	inputNames = actionTempl.ConstructorParams["inputNames"].StringsValue
	taskName = actionTempl.ConstructorParams["taskName"].StringValue
	requireFields = actionTempl.ConstructorParams["requireFields"].StringsValue

	action := NewTaskPromiseAction(inputNames, taskName, "", requireFields)

	action.Name = actionTempl.Name

	return action
}

func (tpa *TaskPromiseAction) String() string {
	return fmt.Sprintf("<TaskPromiseAction %s Name: %s, AllowedInputNames: %v, TaskName: %s, WorkflowName: %s, JobUUID: %s, RequireFields: %v>",
		tpa.UUID, tpa.Name, tpa.AllowedInputNames, tpa.TaskName, tpa.WorkflowName, tpa.JobUUID, tpa.RequireFields)
}

func (tpa *TaskPromiseAction) Run() error {
	inputDataChunksByInputName := map[string]*DataChunk{}

	if tpa.Inputs[TaskPromiseActionInputRefrain] != nil {
		refrain, ok := tpa.Inputs[TaskPromiseActionInputRefrain].Remove().(bool)
		if refrain && ok {
		  	log.Info(fmt.Sprintf("Refraining from making a TaskPromise for task %s in TaskPromise %s (%s)",
				tpa.TaskName, tpa.UUID, tpa.Name))
			return nil
		}
	}

	for name, input := range tpa.Inputs {
		if name == TaskPromiseActionInputRefrain {
			continue
		}

		if len(input.Queue) > 0 {
			x := input.Remove()
			newChunk, _ := NewDataChunk(x)
			inputDataChunksByInputName[name] = newChunk
		}
	}

	for _, rf := range tpa.RequireFields {
		if inputDataChunksByInputName[rf] == nil {
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
