package spiderswarm

import (
	"errors"

	"github.com/google/uuid"
)

const HTTPCookieJoinActionInputOldCookies = "HTTPCookieJoinActionInputOldCookies"
const HTTPCookieJoinActionInputNewCookies = "HTTPCookieJoinActionInputNewCookies"
const HTTPCookieJoinActionOutputUpdatedCookies = "HTTPCookieJoinActionOutputUpdatedCookies"

type HTTPCookieJoinAction struct {
	AbstractAction
}

func NewHTTPCookieJoinAction() *HTTPCookieJoinAction {
	return &HTTPCookieJoinAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				HTTPCookieJoinActionInputOldCookies,
				HTTPCookieJoinActionInputNewCookies,
			},
			AllowedOutputNames: []string{
				HTTPCookieJoinActionOutputUpdatedCookies,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
	}
}

func (hcja *HTTPCookieJoinAction) Run() error {
	if hcja.Inputs[HTTPCookieJoinActionInputOldCookies] == nil || hcja.Inputs[HTTPCookieJoinActionInputNewCookies] == nil {
		return errors.New("Both inputs must be connected")
	}

	if hcja.Outputs[HTTPCookieJoinActionOutputUpdatedCookies] == nil || len(hcja.Outputs[HTTPCookieJoinActionOutputUpdatedCookies]) == 0 {
		return errors.New("Output not connected")
	}

	updatedCookies := map[string]string{}

	oldCookies, ok1 := hcja.Inputs[HTTPCookieJoinActionInputOldCookies].Remove().(map[string]string)
	if !ok1 {
		return errors.New("Failed to get old cookies")
	}

	newCookies, ok2 := hcja.Inputs[HTTPCookieJoinActionInputNewCookies].Remove().(map[string]string)
	if !ok2 {
		return errors.New("Failed to get new cookies")
	}

	for key, value := range oldCookies {
		updatedCookies[key] = value
	}

	for key, value := range newCookies {
		updatedCookies[key] = value
	}

	for _, output := range hcja.Outputs[HTTPCookieJoinActionOutputUpdatedCookies] {
		output.Add(updatedCookies)
	}

	return nil
}
