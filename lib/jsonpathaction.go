package spsw

import (
	"github.com/google/uuid"
)

type JSONPathAction struct {
	AbstractAction
	JSONPath string
	Decode   bool
}

const JSONPathActionInputJSONStr = "JSONPathActionInputJSONStr"
const JSONPathActionInputJSONBytes = "JSONPathActionInputJSONBytes"
const JSONPathActionOutputStr = "JSONPathActionOutputStr"

func NewJSONPathAction(jsonPath string, decode bool, expectMany bool) *JSONPathAction {
	return &JSONPathAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: expectMany,
			AllowedInputNames: []string{
				JSONPathActionInputJSONStr,
				JSONPathActionInputJSONBytes,
			},
			AllowedOutputNames: []string{
				JSONPathActionOutputStr,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
		JSONPath: jsonPath,
		Decode:   decode,
	}
}

func NewJSONPathActionFromTemplate(actionTempl *ActionTemplate, workflowName string) Action {
	jsonPath := actionTempl.ConstructorParams["jsonPath"].StringValue
	decode := actionTempl.ConstructorParams["decode"].BoolValue
	expectMany := actionTempl.ConstructorParams["expectMany"].BoolValue

	action := NewJSONPathAction(jsonPath, decode, expectMany)

	action.Name = actionTempl.Name

	return action
}

func (jpa *JSONPathAction) Run() error {
	return nil
}
