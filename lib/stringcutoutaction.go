package spiderswarm

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

const StringCutActionInputStr = "StringCutActionInputStr"
const StringCutActionOutputStr = "StringCutActionOutputStr"

type StringCutAction struct {
	AbstractAction
	From string
	To   string
}

func NewStringCutAction(from string, to string) *StringCutAction {
	return &StringCutAction{
		AbstractAction: AbstractAction{
			CanFail:            false,
			ExpectMany:         false,
			AllowedInputNames:  []string{StringCutActionInputStr},
			AllowedOutputNames: []string{StringCutActionOutputStr},
			Inputs:             map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				XPathActionOutputStr: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
		From: from,
		To:   to,
	}
}

func NewStringCutActionFromTemplate(actionTempl *ActionTemplate) *StringCutAction {
	var from string
	var to string

	from, _ = actionTempl.ConstructorParams["from"].(string)
	to, _ = actionTempl.ConstructorParams["to"].(string)

	return NewStringCutAction(from, to)
}

func (sca *StringCutAction) Run() error {
	if sca.Inputs[StringCutActionInputStr] == nil {
		return errors.New("Input not connected")
	}

	if sca.Outputs[StringCutActionOutputStr] == nil || len(sca.Outputs[StringCutActionOutputStr]) == 0 {
		return errors.New("No outputs connected")
	}

	var fromIdx int
	var toIdx int

	inputStr, ok := sca.Inputs[StringCutActionInputStr].Remove().(string)
	if ok {
		fromIdx = strings.Index(inputStr, sca.From)
		if fromIdx == -1 {
			return errors.New(".From not found")
		}

		fromIdx += len(sca.From)

		toIdx = strings.Index(inputStr[fromIdx:], sca.To)
		if toIdx == -1 {
			return errors.New(".To not found")
		}

		outputStr := inputStr[fromIdx:toIdx]

		for _, output := range sca.Outputs[StringCutActionOutputStr] {
			output.Add(outputStr)
		}
	} else {
		return errors.New("Cannot get input string")
	}

	return nil
}
