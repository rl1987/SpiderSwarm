package main

import (
	"bytes"
	"errors"
	"golang.org/x/net/html" // XXX
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/google/uuid"
)

type Action interface {
	Run() error
	AddInput(name string, dataPipe *DataPipe) error
	AddOutput(name string, dataPipe *DataPipe) error
	GetUniqueID() string
	GetPrecedingActions() []Action
}

type AbstractAction struct {
	Action
	Inputs             map[string]*DataPipe
	Outputs            map[string][]*DataPipe
	CanFail            bool
	ExpectMany         bool
	AllowedInputNames  []string
	AllowedOutputNames []string
	UUID               string
}

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

const HTTPActionInputURLParams = "HTTPActionInputURLParams"
const HTTPActionInputHeaders = "HTTPActionInputHeaders"
const HTTPActionInputCookies = "HTTPActionInputCookies"
const HTTPActionInputBody = "HTTPActionInputBody"

const HTTPActionOutputBody = "HTTPActionOutputBody"
const HTTPActionOutputHeaders = "HTTPActionOutputHeaders"
const HTTPActionOutputStatusCode = "HTTPActionOutputStatusCode"
const HTTPActionOutputCookies = "HTTPActionOutputCookies"

type HTTPAction struct {
	AbstractAction
	BaseURL string
	Method  string
}

func NewHTTPAction(baseURL string, method string, canFail bool) *HTTPAction {
	return &HTTPAction{
		AbstractAction: AbstractAction{
			CanFail:    canFail,
			ExpectMany: false,
			AllowedInputNames: []string{
				HTTPActionInputURLParams,
				HTTPActionInputHeaders,
				HTTPActionInputCookies,
				HTTPActionInputBody,
			},
			AllowedOutputNames: []string{
				HTTPActionOutputBody,
				HTTPActionOutputHeaders,
				HTTPActionOutputStatusCode,
				HTTPActionOutputCookies,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
		BaseURL: baseURL,
		Method:  method,
	}
}

func (ha *HTTPAction) Run() error {
	var body *bytes.Buffer
	body = nil

	request, err := http.NewRequest(ha.Method, ha.BaseURL, nil)
	if err != nil {
		return err
	}

	q := request.URL.Query()

	if ha.Inputs[HTTPActionInputBody] != nil && ha.Method != http.MethodGet {
		bodyBytes, ok := ha.Inputs[HTTPActionInputBody].Remove().([]byte)
		if ok {
			body = bytes.NewBuffer(bodyBytes)
			request.Body = ioutil.NopCloser(body)
		}
	}

	if ha.Inputs[HTTPActionInputURLParams] != nil {
		for {
			urlParams, ok := ha.Inputs[HTTPActionInputURLParams].Remove().(map[string][]string)
			if !ok {
				break
			}

			for key, values := range urlParams {
				for _, value := range values {
					q.Add(key, value)
				}
			}
		}
	}

	if ha.Inputs[HTTPActionInputHeaders] != nil {
		request.Header = http.Header{}
		for {
			headers, ok := ha.Inputs[HTTPActionInputHeaders].Remove().(http.Header)

			if !ok {
				break
			}

			for key, values := range headers {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}

		}
	}

	if ha.Inputs[HTTPActionInputCookies] != nil {
		for {
			cookies, ok := ha.Inputs[HTTPActionInputCookies].Remove().(map[string]string)

			if !ok {
				break
			}

			for key, value := range cookies {
				c := &http.Cookie{Name: key, Value: value}
				request.AddCookie(c)
			}

		}
	}

	request.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if ha.Outputs[HTTPActionOutputBody] != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			for _, outDP := range ha.Outputs[HTTPActionOutputBody] {
				outDP.Add(body)
			}
		}
	}

	if ha.Outputs[HTTPActionOutputHeaders] != nil {
		headers := resp.Header

		for _, outDP := range ha.Outputs[HTTPActionOutputHeaders] {
			outDP.Add(headers)
		}
	}

	if ha.Outputs[HTTPActionOutputStatusCode] != nil {
		statusCode := resp.StatusCode

		for _, outDP := range ha.Outputs[HTTPActionOutputStatusCode] {
			outDP.Add(statusCode)
		}
	}

	if ha.Outputs[HTTPActionOutputCookies] != nil {
		// XXX: maybe a slice of http.Cookie structs should go into the output?
		cookieDict := map[string]string{}

		for _, cookie := range resp.Cookies() {
			cookieDict[cookie.Name] = cookie.Value
		}

		for _, outDP := range ha.Outputs[HTTPActionOutputCookies] {
			outDP.Add(cookieDict)
		}
	}

	return nil
}

type UTF8DecodeAction struct {
	AbstractAction
}

const UTF8DecodeActionInputBytes = "UTF8DecodeActionInputBytes"
const UTF8DecodeActionOutputStr = "UTF8DecodeActionOutputStr"

func NewUTF8DecodeAction() *UTF8DecodeAction {
	return &UTF8DecodeAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				UTF8DecodeActionInputBytes,
			},
			AllowedOutputNames: []string{
				UTF8DecodeActionOutputStr,
			},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				UTF8DecodeActionOutputStr: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
	}
}

func (ua *UTF8DecodeAction) Run() error {
	if ua.Inputs[UTF8DecodeActionInputBytes] == nil {
		return errors.New("Input not connected")
	}

	if ua.Outputs[UTF8DecodeActionOutputStr] == nil {
		return errors.New("Output not connected")
	}

	binData, ok := ua.Inputs[UTF8DecodeActionInputBytes].Remove().([]byte)
	if !ok {
		return errors.New("Failed to get binary data")
	}

	str := string(binData)

	for _, outDP := range ua.Outputs[UTF8DecodeActionOutputStr] {
		outDP.Add(str)
	}

	return nil
}

type UTF8EncodeAction struct {
	AbstractAction
}

const UTF8EncodeActionInputStr = "UTF8EncodeActionInputStr"
const UTF8EncodeActionOutputBytes = "UTF8EncodeActionOutputBytes"

func NewUTF8EncodeAction() *UTF8EncodeAction {
	return &UTF8EncodeAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				UTF8EncodeActionInputStr,
			},
			AllowedOutputNames: []string{
				UTF8EncodeActionOutputBytes,
			},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				UTF8EncodeActionOutputBytes: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
	}
}

func (ua *UTF8EncodeAction) Run() error {
	if ua.Inputs[UTF8EncodeActionInputStr] == nil {
		return errors.New("Input not connected")
	}

	if ua.Outputs[UTF8EncodeActionOutputBytes] == nil {
		return errors.New("Output not connected")
	}

	str, ok := ua.Inputs[UTF8EncodeActionInputStr].Remove().(string)
	if !ok {
		return errors.New("Failed to get string")
	}

	binData := []byte(str)

	for _, outDP := range ua.Outputs[UTF8EncodeActionOutputBytes] {
		outDP.Add(binData)
	}

	return nil
}

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

// https://stackoverflow.com/a/38855264
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	if n != nil {
		html.Render(w, n)
	}
	return buf.String()
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

	if !xa.ExpectMany {
		var n *html.Node
		n, err = htmlquery.Query(doc, xa.XPath)
		if err != nil {
			return err
		}

		result := renderNode(n)

		for _, outDP := range xa.Outputs[XPathActionOutputStr] {
			outDP.Add(result)
		}
	} else {
		var nodes []*html.Node
		nodes, err = htmlquery.QueryAll(doc, xa.XPath)
		if err != nil {
			return err
		}

		for _, n := range nodes {
			if n == nil {
				continue
			}

			result := renderNode(n)
			for _, outDP := range xa.Outputs[XPathActionOutputStr] {
				outDP.Add(result)
			}
		}
	}

	return nil
}

const FieldJoinActionOutputItem = "FieldJoinActionOutputItem"

type FieldJoinAction struct {
	AbstractAction
}

func NewFieldJoinAction(inputNames []string) *FieldJoinAction {
	return &FieldJoinAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  inputNames,
			AllowedOutputNames: []string{FieldJoinActionOutputItem},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
	}
}

func (fja *FieldJoinAction) Run() error {
	if fja.Outputs[FieldJoinActionOutputItem] == nil {
		return errors.New("Output not connected")
	}

	if len(fja.Inputs) == 0 {
		return errors.New("No inputs connected")
	}

	// TODO: develop a proper data model for items
	item := map[string]string{}

	for key, inDP := range fja.Inputs {
		value, ok := inDP.Remove().(string)
		if ok {
			item[key] = value
		}
	}

	for _, outDP := range fja.Outputs[FieldJoinActionOutputItem] {
		outDP.Add(item)
	}

	return nil
}

type NullAction struct {
	AbstractAction
}

func NewNullAction() *NullAction {
	return &NullAction{
		AbstractAction: AbstractAction{
			UUID: uuid.New().String(),
		},
	}
}
