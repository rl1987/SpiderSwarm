package spsw

import (
	//"net/url"

	"github.com/google/uuid"
)

type URLParseAction struct {
	AbstractAction
}

func NewURLParseAction() *URLParseAction {
	return &URLParseAction{
		AbstractAction: AbstractAction{
			CanFail:            false,
			ExpectMany:         false,
			AllowedInputNames:  []string{},
			AllowedOutputNames: []string{},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			UUID:               uuid.New().String(),
		},
	}
}

func NewURLParseActionFromTemplate(actionTempl *ActionTemplate) *URLParseAction {
	action := NewURLParseAction()

	action.Name = actionTempl.Name

	return action
}
