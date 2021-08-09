package spiderswarm

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

const HTTPActionInputBaseURL = "HTTPActionInputBaseURL"
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
				HTTPActionInputBaseURL,
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

	// TODO: unit-test this part
	if ha.Inputs[HTTPActionInputBaseURL] != nil {
		baseURLStr, ok := ha.Inputs[HTTPActionInputBaseURL].Remove().(string)
		if ok {
			baseURL, err := url.Parse(baseURLStr)
			if err == nil {
				request.URL = baseURL
			}
		}
	}

	if ha.Inputs[HTTPActionInputBody] != nil && ha.Method != http.MethodGet {
		bodyBytes, ok := ha.Inputs[HTTPActionInputBody].Remove().([]byte)
		if ok {
			body = bytes.NewBuffer(bodyBytes)
			request.Body = ioutil.NopCloser(body)
		}
	}

	if ha.Inputs[HTTPActionInputURLParams] != nil {
		x := ha.Inputs[HTTPActionInputURLParams].Remove()
		urlParamsOneToMany, ok1 := x.(map[string][]string)
		// TODO: unit-test this part
		if ok1 {
			for key, values := range urlParamsOneToMany {
				for _, value := range values {
					q.Add(key, value)
				}
			}
		} else {
			urlParamsOneToOne, ok2 := x.(map[string]string)
			if ok2 {
				for key, value := range urlParamsOneToOne {
					q.Add(key, value)
				}
			}
		}
	}

	if ha.Inputs[HTTPActionInputHeaders] != nil {
		request.Header = http.Header{}
		headers, ok := ha.Inputs[HTTPActionInputHeaders].Remove().(http.Header)

		if ok {
			for key, values := range headers {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}
		}
	}

	if ha.Inputs[HTTPActionInputCookies] != nil {
		cookies, ok := ha.Inputs[HTTPActionInputCookies].Remove().(map[string]string)

		if ok {
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
