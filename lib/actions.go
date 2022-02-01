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

type InitFunc func(*ActionTemplate) Action

var ActionConstructorTable = map[string]InitFunc{
	"CSVParseAction":        NewCSVParseActionFromTemplate,
	"FormExtractionAction":  NewFormExtractionActionFromTemplate,
	"HTTPAction":            NewHTTPActionFromTemplate,
	"XPathAction":           NewXPathActionFromTemplate,
	"FieldJoinAction":       NewFieldJoinActionFromTemplate,
	"TaskPromiseAction":     NewTaskPromiseActionFromTemplate,
	"UTF8DecodeAction":      NewUTF8DecodeActionFromTemplate,
	"UTF8EncodeAction":      NewUTF8EncodeActionFromTemplate,
	"ConstAction":           NewConstActionFromTemplate,
	"URLJoinAction":         NewURLJoinActionFromTemplate,
	"StringMapUpdateAction": NewStringMapUpdateActionFromTemplate,
	"URLParseAction":        NewURLParseActionFromTemplate,
	"StringCutAction":       NewStringCutActionFromTemplate,
	"JSONPathAction":        NewJSONPathActionFromTemplate,
	"LengthThresholdAction": NewLengthThresholdActionFromTemplate,
}

var AllowedInputNameTable = map[string][]string{
	"CSVParseAction": []string{
		CSVParseActionInputCSVBytes,
		CSVParseActionInputCSVStr,
	},
	"FormExtractionAction": []string{
		FormExtractionActionInputHTMLBytes,
		FormExtractionActionInputHTMLStr,
	},
	"HTTPAction": []string{
		HTTPActionInputBaseURL,
		HTTPActionInputBody,
		HTTPActionInputCookies,
		HTTPActionInputFormData,
		HTTPActionInputHeaders,
		HTTPActionInputURLParams,
	},
	"XPathAction": []string{
		XPathActionInputHTMLBytes,
		XPathActionInputHTMLStr,
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
	"StringMapUpdateAction": []string{
		StringMapUpdateActionInputNew,
		StringMapUpdateActionInputOld,
		StringMapUpdateActionInputOverridenValue,
	},
	"URLParseAction": []string{
		URLParseActionInputURL,
	},
	"StringCutAction": []string{
		StringCutActionInputStr,
	},
	"JSONPathAction": []string{
		JSONPathActionInputJSONBytes,
		JSONPathActionInputJSONStr,
	},
	"LengthThresholdAction": []string{
		LengthThresholdActionInputSlice,
	},
}

var AllowedOutputNameTable = map[string][]string{
	"CSVParseAction": []string{
		CSVParseActionOutputMap,
	},
	"FormExtractionAction": []string{
		FormExtractionActionOutputFormData,
	},
	"HTTPAction": []string{
		HTTPActionOutputBody,
		HTTPActionOutputCookies,
		HTTPActionOutputHeaders,
		HTTPActionOutputResponseURL,
		HTTPActionOutputStatusCode,
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
	"StringMapUpdateAction": []string{
		StringMapUpdateActionOutputUpdated,
	},
	"StringCutAction": []string{
		StringCutActionOutputStr,
	},
	"JSONPathAction": []string{
		JSONPathActionOutputStr,
	},
	"LengthThresholdAction": []string{
		LengthThresholdActionOutputThresholdUnmet,
	},
}

func RegisterAction(structName string, initFunc InitFunc, allowedInputNames []string, allowedOutputNames []string) {
	ActionConstructorTable[structName] = initFunc
	AllowedInputNameTable[structName] = allowedInputNames
	AllowedOutputNameTable[structName] = allowedOutputNames
}

func NewActionFromTemplate(actionTempl *ActionTemplate) Action {
	initFunc := ActionConstructorTable[actionTempl.StructName]
	if initFunc != nil {
		return initFunc(actionTempl)
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
