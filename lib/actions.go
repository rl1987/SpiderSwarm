package spsw

import (
	"errors"
)

// Action is a single stateless operation that is used as building block for Task.
type Action interface {
	Run() error
	AddInput(name string, dataPipe *DataPipe) error
	AddOutput(name string, dataPipe *DataPipe) error
	GetUniqueID() string
	GetName() string
	GetPrecedingActions() []Action
	IsFailureAllowed() bool
}

// AbstractAction an equivalent of abstract class for all structs that will conform to Action interface.
type AbstractAction struct {
	Action
	Name               string
	Inputs             map[string]*DataPipe
	Outputs            map[string][]*DataPipe
	CanFail            bool
	ExpectMany         bool
	AllowedInputNames  []string
	AllowedOutputNames []string
	UUID               string
}

type InitFunc func(*ActionTemplate, string) Action

var ActionConstructorTable = map[string]InitFunc{
	"HTTPAction":           NewHTTPActionFromTemplate,
	"XPathAction":          NewXPathActionFromTemplate,
	"FieldJoinAction":      NewFieldJoinActionFromTemplate,
	"TaskPromiseAction":    NewTaskPromiseActionFromTemplate,
	"UTF8DecodeAction":     NewUTF8DecodeActionFromTemplate,
	"UTF8EncodeAction":     NewUTF8EncodeActionFromTemplate,
	"ConstAction":          NewConstActionFromTemplate,
	"URLJoinAction":        NewURLJoinActionFromTemplate,
	"HTTPCookieJoinAction": NewHTTPCookieJoinActionFromTemplate,
	"URLParseAction":       NewURLParseActionFromTemplate,
	"StringCutAction":      NewStringCutActionFromTemplate,
	"JSONPathAction":       NewJSONPathActionFromTemplate,
}

var AllowedInputNameTable = map[string][]string{
	"HTTPAction": []string{
		HTTPActionInputBaseURL,
		HTTPActionInputURLParams,
		HTTPActionInputHeaders,
		HTTPActionInputCookies,
		HTTPActionInputBody,
	},
	"XPathAction": []string{
		XPathActionInputHTMLStr,
		XPathActionInputHTMLBytes,
	},
	"FieldJoinAction":   []string{},
	"TaskPromiseAction": []string{},
	"UTF8DecodeAction": []string{
		UTF8DecodeActionInputBytes,
	},
	"UTF8EncodeAction": []string{
		UTF8EncodeActionInputStr,
	},
	"ConstAction": []string{},
	"URLJoinAction": []string{
		URLJoinActionInputBaseURL,
		URLJoinActionInputRelativeURL,
	},
	"HTTPCookieJoinAction": []string{
		HTTPCookieJoinActionInputOldCookies,
		HTTPCookieJoinActionInputNewCookies,
	},
	"URLParseAction": []string{
		URLParseActionInputURL,
	},
	"StringCutAction": []string{
		StringCutActionInputStr,
	},
	"JSONPathAction": []string{
		JSONPathActionInputJSONStr,
		JSONPathActionInputJSONBytes,
	},
}

var AllowedOutputNameTable = map[string][]string{
	"HTTPAction": []string{
		HTTPActionOutputBody,
		HTTPActionOutputHeaders,
		HTTPActionOutputStatusCode,
		HTTPActionOutputCookies,
		HTTPActionOutputResponseURL,
	},
	"XPathAction": []string{
		XPathActionOutputStr,
	},
	"FieldJoinAction": []string{
		FieldJoinActionOutputItem,
		FieldJoinActionOutputMap,
	},
	"TaskPromiseAction": []string{
		TaskPromiseActionOutputPromise,
	},
	"UTF8DecodeAction": []string{
		UTF8DecodeActionOutputStr,
	},
	"UTF8EncodeAction": []string{
		UTF8EncodeActionOutputBytes,
	},
	"ConstAction": []string{
		ConstActionOutput,
	},
	"URLJoinAction": []string{
		URLJoinActionOutputAbsoluteURL,
	},
	"StringCutAction": []string{
		StringCutActionOutputStr,
	},
	"JSONPathAction": []string{
		JSONPathActionOutputStr,
	},
}

func NewActionFromTemplate(actionTempl *ActionTemplate, workflowName string, jobUUID string) Action {
	initFunc := ActionConstructorTable[actionTempl.StructName]
	if initFunc != nil {
		return initFunc(actionTempl, workflowName)
	}

	return nil
}

// AddInput adds input data pipe of given name to Inputs map iff name is in AllowedInputNames.
// Return error otherwise.
func (a *AbstractAction) AddInput(name string, dataPipe *DataPipe) error {
	for _, n := range a.AllowedInputNames {
		if n == name {
			a.Inputs[name] = dataPipe
			return nil
		}
	}

	return errors.New("input name not in AllowedInputNames")
}

func (a *AbstractAction) AddOutput(name string, dataPipe *DataPipe) error {
	for _, n := range a.AllowedOutputNames {
		if n == name {
			if _, ok := a.Outputs[name]; ok {
				a.Outputs[name] = append(a.Outputs[name], dataPipe)
			} else {
				a.Outputs[name] = []*DataPipe{dataPipe}
			}
			return nil
		}
	}

	return errors.New("input name not in AllowedOutputNames")
}

func (a *AbstractAction) GetUniqueID() string {
	return a.UUID
}

func (a *AbstractAction) GetName() string {
	return a.Name
}

func (a *AbstractAction) GetPrecedingActions() []Action {
	actions := []Action{}

	for _, dp := range a.Inputs {
		if dp.FromAction != nil {
			actions = append(actions, dp.FromAction)
		}
	}

	return actions

}

func (a *AbstractAction) Run() error {
	// To be implemented by concrete actions.
	return nil
}

func (a *AbstractAction) IsFailureAllowed() bool {
	return a.CanFail
}
