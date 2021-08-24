package spsw

import (
	"errors"

	"github.com/google/uuid"
)

const UTF8DecodeActionInputBytes = "UTF8DecodeActionInputBytes"
const UTF8DecodeActionOutputStr = "UTF8DecodeActionOutputStr"

type UTF8DecodeAction struct {
	AbstractAction
}

func NewUTF8DecodeAction() *UTF8DecodeAction {
	return &UTF8DecodeAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				UTF8DecodeActionInputBytes,
			},
			AllowedOutputNames: []string{
				UTF8DecodeActionOutputStr,
			},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				UTF8DecodeActionOutputStr: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
	}
}

func NewUTF8DecodeActionFromTemplate(actionTempl *ActionTemplate) *UTF8DecodeAction {
	action := NewUTF8DecodeAction()
	action.Name = actionTempl.Name
	return action
}

func (ua *UTF8DecodeAction) Run() error {
	if ua.Inputs[UTF8DecodeActionInputBytes] == nil {
		return errors.New("Input not connected")
	}

	if ua.Outputs[UTF8DecodeActionOutputStr] == nil {
		return errors.New("Output not connected")
	}

	binData, ok := ua.Inputs[UTF8DecodeActionInputBytes].Remove().([]byte)
	if !ok {
		return errors.New("Failed to get binary data")
	}

	str := string(binData)

	for _, outDP := range ua.Outputs[UTF8DecodeActionOutputStr] {
		outDP.Add(str)
	}

	return nil
}

type UTF8EncodeAction struct {
	AbstractAction
}

const UTF8EncodeActionInputStr = "UTF8EncodeActionInputStr"
const UTF8EncodeActionOutputBytes = "UTF8EncodeActionOutputBytes"

func NewUTF8EncodeAction() *UTF8EncodeAction {
	return &UTF8EncodeAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				UTF8EncodeActionInputStr,
			},
			AllowedOutputNames: []string{
				UTF8EncodeActionOutputBytes,
			},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				UTF8EncodeActionOutputBytes: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
	}
}

func NewUTF8EncodeActionFromTemplate(actionTempl *ActionTemplate) *UTF8EncodeAction {
	action := NewUTF8EncodeAction()
	action.Name = actionTempl.Name
	return action
}

func (ua *UTF8EncodeAction) Run() error {
	if ua.Inputs[UTF8EncodeActionInputStr] == nil {
		return errors.New("Input not connected")
	}

	if ua.Outputs[UTF8EncodeActionOutputBytes] == nil {
		return errors.New("Output not connected")
	}

	str, ok := ua.Inputs[UTF8EncodeActionInputStr].Remove().(string)
	if !ok {
		return errors.New("Failed to get string")
	}

	binData := []byte(str)

	for _, outDP := range ua.Outputs[UTF8EncodeActionOutputBytes] {
		outDP.Add(binData)
	}

	return nil
}
