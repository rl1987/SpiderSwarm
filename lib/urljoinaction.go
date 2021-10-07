package spsw

import (
	"errors"
	"net/url"

	"github.com/google/uuid"
)

const URLJoinActionInputBaseURL = "URLJoinActionInputBaseURL"
const URLJoinActionInputRelativeURL = "URLJoinActionInputRelativeURL"
const URLJoinActionOutputAbsoluteURL = "URLJoinActionOutputAbsoluteURL"

type URLJoinAction struct {
	AbstractAction
	BaseURL string
}

func NewURLJoinAction(baseURL string) *URLJoinAction {
	return &URLJoinAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				URLJoinActionInputBaseURL,
				URLJoinActionInputRelativeURL,
			},
			AllowedOutputNames: []string{
				URLJoinActionOutputAbsoluteURL,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
		BaseURL: baseURL,
	}
}

func NewURLJoinActionFromTemplate(actionTempl *ActionTemplate) *URLJoinAction {
	baseURL := actionTempl.ConstructorParams["baseURL"].StringValue

	action := NewURLJoinAction(baseURL)

	action.Name = actionTempl.Name

	return action
}

func (uja *URLJoinAction) Run() error {
	var baseURL *url.URL
	var absoluteURL *url.URL
	var err error

	baseURL, err = url.Parse(uja.BaseURL)
	if err != nil {
		return err
	}

	if uja.Inputs[URLJoinActionInputRelativeURL] == nil {
		return errors.New("URLJoinActionInputRelativeURL not connected")
	}

	if uja.Inputs[URLJoinActionInputBaseURL] != nil {
		if baseURLStr, ok1 := uja.Inputs[URLJoinActionInputBaseURL].Remove().(string); ok1 {
			baseURL, err = url.Parse(baseURLStr)
			if err != nil {
				return err
			}
		}
	}

	if relativeURLStr, ok2 := uja.Inputs[URLJoinActionInputRelativeURL].Remove().(string); ok2 {
		absoluteURL, err = baseURL.Parse(relativeURLStr)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Cannot get relative URL")
	}

	for _, outDP := range uja.Outputs[URLJoinActionOutputAbsoluteURL] {
		outDP.Add(absoluteURL.String())
	}

	return nil
}
