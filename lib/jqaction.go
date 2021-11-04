package spsw

import (
	//"encoding/json"
	//"io"
	//"os/exec"

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

func (jq *JQAction) Run() error {
	//var stdinReader io.Reader
	//var stoutWriter io.Writer

	//_ := exec.Cmd{
	//	Path: "jq",
	//}

	// TODO: implement this

	return nil
}
