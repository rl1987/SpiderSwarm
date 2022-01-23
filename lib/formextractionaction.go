package spsw

import (
	"errors"
	"fmt"
	"golang.org/x/net/html" // XXX
	"strings"

	"github.com/antchfx/htmlquery"
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
			CanFail:    false,
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

func (fea *FormExtractionAction) Run() error {
	if fea.Inputs[FormExtractionActionInputHTMLStr] == nil && fea.Inputs[FormExtractionActionInputHTMLBytes] == nil {
		return errors.New("Input not connected")
	}

	if fea.Outputs[FormExtractionActionOutputFormData] == nil || len(fea.Outputs[FormExtractionActionOutputFormData]) == 0 {
		return errors.New("Output not connected")
	}

	var htmlStr string

	if fea.Inputs[FormExtractionActionInputHTMLStr] != nil {
		htmlStr, _ = fea.Inputs[FormExtractionActionInputHTMLStr].Remove().(string)
	} else {
		htmlBytes, _ := fea.Inputs[FormExtractionActionInputHTMLStr].Remove().([]byte)
		htmlStr = string(htmlBytes)
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return err
	}

	// TODO: support multi-valued inputs (e.g. checkboxes) later.
	formData := map[string]string{}

	var inputNodes []*html.Node
	// TODO: make form ID optional - not all forms have id attribute.
	xpath := fmt.Sprintf("//form[@id=\"%s\"]//input", fea.FormID)

	inputNodes, err = htmlquery.QueryAll(doc, xpath)
	if err != nil {
		return err
	}

	for _, input := range inputNodes {
		if input == nil {
			continue
		}

		name := ""
		value := ""

		for _, attrib := range input.Attr {
			if attrib.Key == "name" {
				name = attrib.Val
			}

			if attrib.Key == "value" {
				value = attrib.Val
			}
		}

		if name != "" {
			formData[name] = value
		}
	}

	for _, outDP := range fea.Outputs[FormExtractionActionOutputFormData] {
		outDP.Add(formData)
	}

	return nil
}
