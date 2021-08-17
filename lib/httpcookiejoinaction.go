package spiderswarm

import (
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
	// TODO: implement
	return nil
}
