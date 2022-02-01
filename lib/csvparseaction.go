package spsw

import (
	"github.com/google/uuid"
)

const CSVParseActionInputCSVBytes = "CSVParseActionInputCSVBytes"
const CSVParseActionInputCSVStr = "CSVParseActionInputCSVStr"
const CSVParseActionOutputMap = "CSVParseActionOutputMap"

type CSVParseAction struct {
	AbstractAction
}

func NewCSVParseAction() *CSVParseAction {
	return &CSVParseAction{
		AbstractAction: AbstractAction{
			AllowedInputNames: []string{CSVParseActionInputCSVBytes, CSVParseActionInputCSVStr},
			AllowedOutputNames: []string{CSVParseActionOutputMap},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			CanFail: false,
			UUID: uuid.New().String(),
		},
	}
}

func NewCSVParseActionFromTemplate(actionTempl *ActionTemplate) Action {
	action := NewCSVParseAction()

	action.Name = actionTempl.Name

	return action
}

