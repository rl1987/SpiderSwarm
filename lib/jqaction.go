package spsw

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
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

func (jqa *JQAction) outputBytes(outBytes []byte) {
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
}

func (jqa *JQAction) String() string {
	return fmt.Sprintf("<JQAction %s Name: %s, CanFail: %v, ExpectMany: %v, DecodeOutput: %v, JQArgs: %v>",
		jqa.UUID, jqa.Name, jqa.CanFail, jqa.ExpectMany, jqa.DecodeOutput, jqa.JQArgs)
}

func (jqa *JQAction) Run() error {
	if jqa.Inputs[JQActionInputJQStdinStr] == nil && jqa.Inputs[JQActionInputJQStdinBytes] == nil {
		return errors.New("No inputs connected")
	}

	if jqa.Outputs[JQActionOutputJQStdoutStr] == nil && jqa.Outputs[JQActionOutputJQStdoutBytes] == nil {
		return errors.New("No outputs connected")
	}

	var inBuf bytes.Buffer

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

	cmd := exec.Command("jq", jqa.JQArgs...)
	cmd.Stdin = &inBuf

	outBytes, err := cmd.Output()
	if err != nil {
		return err
	}

	if jqa.DecodeOutput {
		if jqa.ExpectMany {
			var strings []string

			err := json.Unmarshal(outBytes, &strings)
			if err != nil {
				spew.Dump(outBytes)
				return err
			}

			for _, s := range strings {
				jqa.outputBytes([]byte(s))
			}
		} else {
			var s string

			err := json.Unmarshal(outBytes, &s)
			if err != nil {
				spew.Dump(outBytes)
				return err
			}

			jqa.outputBytes([]byte(s))
		}
	} else {
		jqa.outputBytes(outBytes)
	}

	return nil
}
