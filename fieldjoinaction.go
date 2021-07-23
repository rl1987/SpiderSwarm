package main

import (
	"errors"

	"github.com/google/uuid"
)

const FieldJoinActionOutputItem = "FieldJoinActionOutputItem"

type FieldJoinAction struct {
	AbstractAction
}

func NewFieldJoinAction(inputNames []string) *FieldJoinAction {
	return &FieldJoinAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  inputNames,
			AllowedOutputNames: []string{FieldJoinActionOutputItem},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
	}
}

func (fja *FieldJoinAction) Run() error {
	if fja.Outputs[FieldJoinActionOutputItem] == nil {
		return errors.New("Output not connected")
	}

	if len(fja.Inputs) == 0 {
		return errors.New("No inputs connected")
	}

	// TODO: develop a proper data model for items
	item := map[string]string{}

	for key, inDP := range fja.Inputs {
		value, ok := inDP.Remove().(string)
		if ok {
			item[key] = value
		}
	}

	for _, outDP := range fja.Outputs[FieldJoinActionOutputItem] {
		outDP.Add(item)
	}

	return nil
}
