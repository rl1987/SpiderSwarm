package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

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

func NewHTTPActionFromTemplate(actionTempl *ActionTemplate) *HTTPAction {
	var baseURL string
	var method string
	var canFail bool

	baseURL, _ = actionTempl.ConstructorParams["baseURL"].(string)
	method, _ = actionTempl.ConstructorParams["method"].(string)
	canFail, _ = actionTempl.ConstructorParams["canFail"].(bool)

	return NewHTTPAction(baseURL, method, canFail)
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
