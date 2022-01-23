package spsw

import (
	"github.com/google/uuid"
)

const FormExtractionActionInputHTMLStr = "FormExtractionActionInputHTMLStr"
const FormExtractionActionInputHTMLBytes = "FormExtractionActionInputHTMLBytes"

const FormExtractionActionOutputFormData = "FormExtractionActionOutputFormData"

type FormExtractionAction struct {
	AbstractAction
	FormID string
}

func NewFormExtractionAction(formID string) *FormExtractionAction {
	return &FormExtractionAction{
		AbstractAction: AbstractAction{
			CanFail: false,
			ExpectMany: false,
			AllowedInputNames: []string{
				FormExtractionActionInputHTMLStr,
				FormExtractionActionInputHTMLBytes,
			},
			AllowedOutputNames: []string{
				FormExtractionActionOutputFormData,
			},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				FormExtractionActionOutputFormData: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
		FormID: formID,
	}
}

func NewFormExtractionActionFromTemplate(actionTempl *ActionTemplate, workflowName string) Action {
	var formID string

	formID = actionTempl.ConstructorParams["formID"].StringValue

	action := NewFormExtractionAction(formID)

	action.Name = actionTempl.Name

	return action
}



