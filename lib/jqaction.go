package spsw

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"

	"github.com/google/uuid"
)

const JQActionInputJQStdinStr = "JQActionInputJQStdinStr"
const JQActionInputJQStdinBytes = "JQActionInputJQStdinBytes"
const JQActionOutputJQStdoutStr = "JQActionOutputJQStdoutStr"
const JQActionOutputJQStdoutBytes = "JQActionOutputJQStdoutBytes"

type JQAction struct {
	AbstractAction
	JQArgs       []string
	DecodeOutput bool
}

func NewJQAction(jqArgs []string, decodeOutput bool, canFail bool, expectMany bool) *JQAction {
	return &JQAction{
		AbstractAction: AbstractAction{
			CanFail:    canFail,
			ExpectMany: expectMany,
			AllowedInputNames: []string{
				JQActionInputJQStdinStr,
				JQActionInputJQStdinBytes,
			},
			AllowedOutputNames: []string{
				JQActionOutputJQStdoutStr,
				JQActionOutputJQStdoutBytes,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
		JQArgs:       jqArgs,
		DecodeOutput: decodeOutput,
	}
}

func NewJQActionFromTemplate(actionTempl *ActionTemplate) *JQAction {
	var jqArgs []string
	var decodeOutput bool
	var canFail bool
	var expectMany bool

	jqArgs = actionTempl.ConstructorParams["jqArgs"].StringsValue
	decodeOutput = actionTempl.ConstructorParams["decodeOutput"].BoolValue
	canFail = actionTempl.ConstructorParams["canFail"].BoolValue
	expectMany = actionTempl.ConstructorParams["expectMany"].BoolValue

	action := NewJQAction(jqArgs, decodeOutput, canFail, expectMany)

	action.Name = actionTempl.Name

	return action
}

func (jqa *JQAction) Run() error {
	if jqa.Inputs[JQActionInputJQStdinStr] == nil && jqa.Inputs[JQActionInputJQStdinBytes] == nil {
		return errors.New("No inputs connected")
	}

	if jqa.Outputs[JQActionOutputJQStdoutStr] == nil && jqa.Outputs[JQActionOutputJQStdoutBytes] == nil {
		return errors.New("No outputs connected")
	}

	var inBuf bytes.Buffer
	var outBuf bytes.Buffer

	if jqa.Inputs[JQActionInputJQStdinBytes] != nil {
		inBytes, ok := jqa.Inputs[JQActionInputJQStdinBytes].Remove().([]byte)
		if !ok {
			return errors.New("Failed to read bytes from input")
		}

		inBuf.Write(inBytes)
	} else if jqa.Inputs[JQActionInputJQStdinStr] != nil {
		inStr, ok := jqa.Inputs[JQActionInputJQStdinStr].Remove().(string)
		if !ok {
			return errors.New("Failed to read JSON string from input")
		}

		inBuf.WriteString(inStr)
	}

	cmd := exec.Cmd{
		Path:   "jq",
		Args:   jqa.JQArgs,
		Stdin:  &inBuf,
		Stdout: &outBuf,
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	outBytes, _ := ioutil.ReadAll(&outBuf)

	if jqa.Outputs[JQActionOutputJQStdoutStr] != nil {
		for _, output := range jqa.Outputs[JQActionOutputJQStdoutStr] {
			output.Add(outBytes)
		}
	}

	if jqa.Outputs[JQActionOutputJQStdoutBytes] != nil {
		outStr := string(outBytes)

		for _, output := range jqa.Outputs[JQActionOutputJQStdoutBytes] {
			output.Add(outStr)
		}
	}

	return nil
}
