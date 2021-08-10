package spiderswarm

import (
	"bytes"
	"errors"
	"golang.org/x/net/html" // XXX
	"io"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/google/uuid"
)

const XPathActionInputHTMLStr = "XPathActionInputHTMLStr"
const XPathActionInputHTMLBytes = "XPathActionInputHTMLBytes"
const XPathActionOutputStr = "XPathActionOutputStr"

type XPathAction struct {
	AbstractAction
	XPath string
}

func NewXPathAction(xpath string, expectMany bool) *XPathAction {
	return &XPathAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: expectMany,
			AllowedInputNames: []string{
				XPathActionInputHTMLStr,
				XPathActionInputHTMLBytes,
			},
			AllowedOutputNames: []string{
				XPathActionOutputStr,
			},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				XPathActionOutputStr: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
		XPath: xpath,
	}
}

func NewXPathActionFromTemplate(actionTempl *ActionTemplate) *XPathAction {
	var xpath string
	var expectMany bool

	xpath, _ = actionTempl.ConstructorParams["xpath"].(string)
	expectMany = actionTempl.ConstructorParams["expectMany"].(bool)

	action := NewXPathAction(xpath, expectMany)

	action.Name = actionTempl.Name

	return action
}

// https://stackoverflow.com/a/38855264
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	if n != nil {
		html.Render(w, n)
	}
	return buf.String()
}

func extractAttribute(result string, attribName string) string {
	// HACK to clean up attribute string
	result = strings.Replace(result, "<"+attribName+">", "", -1)
	result = strings.Replace(result, "</"+attribName+">", "", -1)

	return result
}

func (xa *XPathAction) Run() error {
	if xa.Inputs[XPathActionInputHTMLStr] == nil && xa.Inputs[XPathActionInputHTMLBytes] == nil {
		return errors.New("Input not connected")
	}

	if xa.Outputs[XPathActionOutputStr] == nil {
		return errors.New("Output not connected")
	}

	var htmlStr string

	if xa.Inputs[XPathActionInputHTMLStr] != nil {
		htmlStr, _ = xa.Inputs[XPathActionInputHTMLStr].Remove().(string)
	} else if xa.Inputs[XPathActionInputHTMLBytes] != nil {
		htmlBytes, ok := xa.Inputs[XPathActionInputHTMLBytes].Remove().([]byte)
		if ok {
			htmlStr = string(htmlBytes)
		}
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return err
	}

	var attribName string
	var extractAttrib bool

	// HACK to clean up attribute string
	splitXPath := strings.Split(xa.XPath, "/")
	if len(splitXPath) > 0 {
		if len(splitXPath[len(splitXPath)-1]) > 0 {
			if splitXPath[len(splitXPath)-1][0] == '@' {
				attribName = splitXPath[len(splitXPath)-1][1:]
				extractAttrib = true
			}
		}
	}

	if !xa.ExpectMany {
		var n *html.Node
		n, err = htmlquery.Query(doc, xa.XPath)
		if err != nil {
			return err
		}

		result := renderNode(n)
		if extractAttrib {
			result = extractAttribute(result, attribName)
		}

		for _, outDP := range xa.Outputs[XPathActionOutputStr] {
			outDP.Add(result)
		}
	} else {
		var nodes []*html.Node
		nodes, err = htmlquery.QueryAll(doc, xa.XPath)
		if err != nil {
			return err
		}

		var results []string

		for _, n := range nodes {
			if n == nil {
				continue
			}

			result := renderNode(n)
			if extractAttrib {
				result = extractAttribute(result, attribName)
			}

			results = append(results, result)
		}

		for _, outDP := range xa.Outputs[XPathActionOutputStr] {
			outDP.Add(results)
		}
	}

	return nil
}
