package spsw

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const StringMapUpdateActionInputOld = "StringMapUpdateActionInputOld"
const StringMapUpdateActionInputNew = "StringMapUpdateActionInputNew"
const StringMapUpdateActionOutputUpdated = "StringMapUpdateActionOutputUpdated"

type StringMapUpdateAction struct {
	AbstractAction
}

func NewStringMapUpdateAction() *StringMapUpdateAction {
	return &StringMapUpdateAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				StringMapUpdateActionInputOld,
				StringMapUpdateActionInputNew,
			},
			AllowedOutputNames: []string{
				StringMapUpdateActionOutputUpdated,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
	}
}

func NewStringMapUpdateActionFromTemplate(actionTempl *ActionTemplate) Action {
	action := NewStringMapUpdateAction()

	action.Name = actionTempl.Name

	return action
}

func (smua *StringMapUpdateAction) String() string {
	return fmt.Sprintf("<StringMapUpdateAction %s Name: %s>", smua.UUID, smua.Name)
}

func (smua *StringMapUpdateAction) Run() error {
	if smua.Inputs[StringMapUpdateActionInputOld] == nil || smua.Inputs[StringMapUpdateActionInputNew] == nil {
		return errors.New("Both inputs must be connected")
	}

	if smua.Outputs[StringMapUpdateActionOutputUpdated] == nil || len(smua.Outputs[StringMapUpdateActionOutputUpdated]) == 0 {
		return errors.New("Output not connected")
	}

	updatedMap := map[string]string{}

	oldMap, ok1 := smua.Inputs[StringMapUpdateActionInputOld].Remove().(map[string]string)
	if !ok1 {
		return errors.New("Failed to get old cookies")
	}
	newMap, ok2 := smua.Inputs[StringMapUpdateActionInputNew].Remove().(map[string]string)
	if !ok2 {
		return errors.New("Failed to get new cookies")
	}

	for key, value := range oldMap {
		updatedMap[key] = value
	}

	for key, value := range newMap {
		updatedMap[key] = value
	}

	for _, output := range smua.Outputs[StringMapUpdateActionOutputUpdated] {
		output.Add(updatedMap)
	}

	return nil
}
