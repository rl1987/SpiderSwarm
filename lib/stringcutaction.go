package spsw

import (
	"errors"
	"fmt"
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
			Outputs:            map[string][]*DataPipe{},
			UUID:               uuid.New().String(),
		},
		From: from,
		To:   to,
	}
}

func NewStringCutActionFromTemplate(actionTempl *ActionTemplate) Action {
	var from string
	var to string

	from = actionTempl.ConstructorParams["from"].StringValue
	to = actionTempl.ConstructorParams["to"].StringValue

	action := NewStringCutAction(from, to)

	action.Name = actionTempl.Name

	return action
}

func (sca *StringCutAction) String() string {
	return fmt.Sprintf("<StringCutAction %s Name: %s, From: %s, To: %s>", sca.UUID, sca.Name, sca.From, sca.To)
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

	x := sca.Inputs[StringCutActionInputStr].Remove()

	if inputStr, ok := x.(string); ok {
		fromIdx = strings.Index(inputStr, sca.From)
		if fromIdx == -1 {
			return errors.New(".From not found")
		}

		fromIdx += len(sca.From)

		toIdx = strings.Index(inputStr[fromIdx:], sca.To)
		if toIdx == -1 {
			return errors.New(".To not found")
		}

		toIdx += fromIdx

		outputStr := inputStr[fromIdx:toIdx]

		for _, output := range sca.Outputs[StringCutActionOutputStr] {
			output.Add(outputStr)
		}
	} else if inputStrings, ok2 := x.([]string); ok2 {
		outputStrings := []string{}

		for _, inputStr := range inputStrings {
			fromIdx = strings.Index(inputStr, sca.From)
			if fromIdx == -1 {
				continue
			}
		  	fromIdx += len(sca.From)

		 	toIdx = strings.Index(inputStr[fromIdx:], sca.To)
			if toIdx == -1 {
				continue
			}
		  	toIdx += fromIdx

		  	outputStr := inputStr[fromIdx:toIdx]
			
			outputStrings = append(outputStrings, outputStr)
		}
		
		for _, outDP := range sca.Outputs[StringCutActionOutputStr] {
			outDP.Add(outputStrings)
		}
	} else {
		return errors.New("Cannot get input string")
	}

	return nil
}
