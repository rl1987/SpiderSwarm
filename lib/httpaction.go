package spsw

import (
	"bytes"
	"fmt"
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
const HTTPActionOutputResponseURL = "HTTPActionOutputResponseURL"

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
				HTTPActionOutputResponseURL,
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

	baseURL = actionTempl.ConstructorParams["baseURL"].StringValue
	method = actionTempl.ConstructorParams["method"].StringValue
	canFail = actionTempl.ConstructorParams["canFail"].BoolValue

	action := NewHTTPAction(baseURL, method, canFail)

	action.Name = actionTempl.Name

	return action
}

func (ha *HTTPAction) String() string {
	return fmt.Sprintf("<HTTPAction %s Name: %s CanFail: %v, BaseURL: %s, Method: %s>", ha.UUID, ha.Name, ha.CanFail, ha.BaseURL, ha.Method)
}

func (ha *HTTPAction) Run() error {
	var body *bytes.Buffer
	body = nil

	request, err := http.NewRequest(ha.Method, ha.BaseURL, nil)
	if err != nil {
		return err
	}

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

	q := request.URL.Query()

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

	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")

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

	// https://stackoverflow.com/questions/51845690/how-to-program-go-to-use-a-proxy-when-using-a-custom-transport
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Transport: transport,
	}

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

	if ha.Outputs[HTTPActionOutputResponseURL] != nil {
		responseURL := resp.Request.URL // XXX: what about redirects?

		for _, outDP := range ha.Outputs[HTTPActionOutputResponseURL] {
			outDP.Add(responseURL.String())
		}
	}

	return nil
}
