package spsw

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

const URLParseActionInputURL = "URLParseActionInputURL"

const URLParseActionOutputScheme = "URLParseActionOutputScheme"
const URLParseActionOutputHost = "URLParseActionOutputHost"
const URLParseActionOutputPath = "URLParseActionOutputPath"
const URLParseActionOutputParams = "URLParseActionOutputParams"

type URLParseAction struct {
	AbstractAction
}

func NewURLParseAction() *URLParseAction {
	return &URLParseAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				URLParseActionInputURL,
			},
			AllowedOutputNames: []string{
				URLParseActionOutputScheme,
				URLParseActionOutputHost,
				URLParseActionOutputPath,
				URLParseActionOutputParams,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
	}
}

func NewURLParseActionFromTemplate(actionTempl *ActionTemplate) Action {
	action := NewURLParseAction()

	action.Name = actionTempl.Name

	return action
}

func (upa *URLParseAction) String() string {
	return fmt.Sprintf("<URLParseAction %s>", upa.UUID)
}

func (upa *URLParseAction) fixParams(params url.Values) map[string][]string {
	newParams := map[string][]string{}

	for key, values := range params {
		newParams[key] = values
	}

	return newParams
}

func (upa *URLParseAction) Run() error {
	if upa.Inputs[URLParseActionInputURL] == nil {
		return errors.New("Input not connected")
	}

	urlStr, ok := upa.Inputs[URLParseActionInputURL].Remove().(string)
	if !ok {
		return nil // XXX: is this error condition?
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if upa.Outputs[URLParseActionOutputScheme] != nil {
		for _, outDP := range upa.Outputs[URLParseActionOutputScheme] {
			outDP.Add(parsed.Scheme)
		}
	}

	if upa.Outputs[URLParseActionOutputHost] != nil {
		for _, outDP := range upa.Outputs[URLParseActionOutputHost] {
			outDP.Add(parsed.Host)
		}
	}

	if upa.Outputs[URLParseActionOutputPath] != nil {
		for _, outDP := range upa.Outputs[URLParseActionOutputPath] {
			outDP.Add(parsed.Path)
		}
	}

	if upa.Outputs[URLParseActionOutputParams] != nil {
		params, _ := url.ParseQuery(parsed.RawQuery)
		for _, outDP := range upa.Outputs[URLParseActionOutputParams] {
			outDP.Add(upa.fixParams(params))
		}
	}

	return nil
}
