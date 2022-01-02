package spsw

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type LengthThresholdAction struct {
	AbstractAction
	UUID      string
	Threshold int
}

const LengthThresholdActionInputSlice = "LengthThresholdActionInputSlice"

const LengthThresholdActionOutputThresholdUnmet = "LengthThresholdActionOutputThresholdUnmet"

func NewLengthThresholdAction(threshold int) *LengthThresholdAction {
	return &LengthThresholdAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  []string{LengthThresholdActionInputSlice},
			AllowedOutputNames: []string{LengthThresholdActionOutputThresholdUnmet},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
		Threshold: threshold,
	}
}

func NewLengthThresholdActionFromTemplate(actionTempl *ActionTemplate, workflowName string) Action {
	threshold := actionTempl.ConstructorParams["threshold"].IntValue

	action := NewLengthThresholdAction(threshold)

	action.Name = actionTempl.Name

	return action
}

func (lta *LengthThresholdAction) String() string {
	return fmt.Sprintf("<LengthThresholdAction %s Name: %s, Threshold: %d>", lta.UUID, lta.Name, lta.Threshold)
}

func (lta *LengthThresholdAction) Run() error {
	if lta.Inputs[LengthThresholdActionInputSlice] == nil {
		return errors.New("Input not connected")
	}

	if lta.Outputs[LengthThresholdActionOutputThresholdUnmet] == nil ||
		len(lta.Outputs[LengthThresholdActionOutputThresholdUnmet]) == 0 {
		return errors.New("Output not connected")
	}

	x := lta.Inputs[LengthThresholdActionInputSlice].Remove()

	val := reflect.ValueOf(x)

	var unmet bool

	if !val.IsValid() || val.IsNil() || val.Len() < lta.Threshold {
		unmet = true
	} else {
		unmet = false
	}

	for _, outDP := range lta.Outputs[LengthThresholdActionOutputThresholdUnmet] {
		outDP.Add(unmet)
	}

	return nil
}
