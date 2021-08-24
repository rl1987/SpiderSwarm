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

func NewActionFromTemplate(actionTempl *ActionTemplate, workflowName string, jobUUID string) Action {
	if actionTempl.StructName == "HTTPAction" {
		return NewHTTPActionFromTemplate(actionTempl)
	} else if actionTempl.StructName == "XPathAction" {
		return NewXPathActionFromTemplate(actionTempl)
	} else if actionTempl.StructName == "FieldJoinAction" {
		return NewFieldJoinActionFromTemplate(actionTempl, workflowName)
	} else if actionTempl.StructName == "TaskPromiseAction" {
		return NewTaskPromiseActionFromTemplate(actionTempl, workflowName)
	} else if actionTempl.StructName == "UTF8DecodeAction" {
		return NewUTF8DecodeActionFromTemplate(actionTempl)
	} else if actionTempl.StructName == "UTF8EncodeAction" {
		return NewUTF8EncodeActionFromTemplate(actionTempl)
	} else if actionTempl.StructName == "ConstAction" {
		return NewConstActionFromTemplate(actionTempl)
	} else if actionTempl.StructName == "URLJoinAction" {
		return NewURLJoinActionFromTemplate(actionTempl)
	} else if actionTempl.StructName == "StringCutAction" {
		return NewStringCutActionFromTemplate(actionTempl)
	} else if actionTempl.StructName == "HTTPCookieJoinAction" {
		return NewHTTPCookieJoinActionFromTemplate(actionTempl)
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
