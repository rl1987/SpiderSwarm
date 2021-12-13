package spsw

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
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
	if jpa.Inputs[JSONPathActionInputJSONStr] == nil && jpa.Inputs[JSONPathActionInputJSONBytes] == nil {
		return errors.New("Input not connected")
	}

	if jpa.Outputs[JSONPathActionOutputStr] == nil {
		return errors.New("Output not connected")
	}

	var jsonStr string

	if jpa.Inputs[JSONPathActionInputJSONStr] != nil {
		jsonStr = jpa.Inputs[JSONPathActionInputJSONStr].Remove().(string)
	} else {
		jsonBytes := jpa.Inputs[JSONPathActionInputJSONBytes].Remove().([]byte)
		jsonStr = string(jsonBytes)
	}

	parsed, err := oj.ParseString(jsonStr)
	if err != nil {
		return err
	}

	query, err := jp.ParseString(jpa.JSONPath)
	if err != nil {
		return err
	}

	var result interface{}

	result = query.Get(parsed)

	if !jpa.Decode {
		jsonStr2 := oj.JSON(result)

		for _, outDP := range jpa.Outputs[JSONPathActionOutputStr] {
			outDP.Add(jsonStr2)
		}

		return nil
	}

	if jpa.ExpectMany {
		if resultStrings, okStrings := result.([]string); okStrings {
			for _, outDP := range jpa.Outputs[JSONPathActionOutputStr] {
				outDP.Add(resultStrings)
			}
		} else if resultIntfs, okIntfs := result.([]interface{}); okIntfs {
			resultStrings := []string{}

			for _, x := range resultIntfs {
				resultStrings = append(resultStrings, fmt.Sprintf("%v", x))
			}
		}
	} else {
		if resultStr, okStr := result.(string); okStr {
			for _, outDP := range jpa.Outputs[JSONPathActionOutputStr] {
				outDP.Add(resultStr)
			}
		} else {
			resultStr := fmt.Sprintf("%v", result)

			for _, outDP := range jpa.Outputs[JSONPathActionOutputStr] {
				outDP.Add(resultStr)
			}
		}
	}

	return nil
}
