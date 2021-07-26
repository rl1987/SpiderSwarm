package main

import (
	"errors"

	"github.com/google/uuid"
)

const FieldJoinActionOutputItem = "FieldJoinActionOutputItem"

type FieldJoinAction struct {
	AbstractAction
	WorkflowName string
	JobUUID      string
	TaskUUID     string
	ItemName     string
}

func NewFieldJoinAction(inputNames []string, workflowName string, jobUUID string, taskUUID string, itemName string) *FieldJoinAction {
	return &FieldJoinAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  inputNames,
			AllowedOutputNames: []string{FieldJoinActionOutputItem},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
		TaskUUID:     taskUUID,
		ItemName:     itemName,
	}
}

func (fja *FieldJoinAction) Run() error {
	if fja.Outputs[FieldJoinActionOutputItem] == nil {
		return errors.New("Output not connected")
	}

	if len(fja.Inputs) == 0 {
		return errors.New("No inputs connected")
	}

	item := NewItem(fja.ItemName, fja.WorkflowName, fja.JobUUID, fja.TaskUUID)

	for key, inDP := range fja.Inputs {
		if len(inDP.Queue) > 0 {
			value := inDP.Remove()
			item.SetField(key, value)
		}
	}

	for _, outDP := range fja.Outputs[FieldJoinActionOutputItem] {
		outDP.AddItem(item)
	}

	return nil
}
