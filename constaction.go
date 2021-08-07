package spiderswarm

import (
	"errors"

	"github.com/google/uuid"
)

const ConstActionOutput = "ConstActionOutput"

type ConstAction struct {
	AbstractAction
	C interface{}
}

func NewConstAction(c interface{}) *ConstAction {
	return &ConstAction{
		AbstractAction: AbstractAction{
			CanFail:            false,
			ExpectMany:         false,
			AllowedInputNames:  []string{},
			AllowedOutputNames: []string{ConstActionOutput},
			Inputs:             map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				XPathActionOutputStr: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
		C: c,
	}
}

func NewConstActionFromTemplate(actionTempl *ActionTemplate) *ConstAction {
	c, _ := actionTempl.ConstructorParams["c"]
	return NewConstAction(c)
}

func (ca *ConstAction) Run() error {
	if ca.Outputs[ConstActionOutput] == nil {
		return errors.New("Output not connected")
	}

	for _, output := range ca.Outputs[ConstActionOutput] {
		output.Add(ca.C)
	}

	return nil
}
