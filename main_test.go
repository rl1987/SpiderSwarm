package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPAction(t *testing.T) {
	baseURL := "https://httpbin.org/post"
	method := "POST"

	httpAction := NewHTTPAction(baseURL, method, false)

	assert.NotNil(t, httpAction)
	assert.False(t, httpAction.AbstractAction.ExpectMany)
	assert.Equal(t, httpAction.BaseURL, baseURL)
	assert.Equal(t, httpAction.Method, method)
	assert.Equal(t, len(httpAction.AbstractAction.AllowedInputNames), 4)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[0], HTTPActionInputURLParams)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[1], HTTPActionInputHeaders)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[2], HTTPActionInputCookies)
	assert.Equal(t, len(httpAction.AbstractAction.AllowedOutputNames), 4)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[0], HTTPActionOutputBody)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[1], HTTPActionOutputHeaders)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[2], HTTPActionOutputStatusCode)

}

func TestHTTPActionRunGET(t *testing.T) {
	testHeaders := http.Header{
		"User-Agent": []string{"spiderswarm"},
		"Accept":     []string{"text/plain"},
	}

	testParams := map[string][]string{
		"a": []string{"1"},
		"b": []string{"2"},
	}

	expectedBody := []byte("Test Payload")

	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			u := req.URL
			m, _ := url.ParseQuery(u.RawQuery)
			assert.Equal(t, 2, len(m))
			assert.Equal(t, "1", m["a"][0])
			assert.Equal(t, "2", m["b"][0])

			assert.Equal(t, "spiderswarm", req.Header["User-Agent"][0])
			assert.Equal(t, "text/plain", req.Header["Accept"][0])

			res.Header()["Server"] = []string{"TestServer"}
			res.WriteHeader(200)
			res.Write(expectedBody)
		}))

	defer testServer.Close()

	httpAction := NewHTTPAction(testServer.URL, http.MethodGet, false)

	headersIn := NewDataPipe()
	err := headersIn.Add(testHeaders)
	assert.Nil(t, err)

	err = httpAction.AddInput(HTTPActionInputHeaders, headersIn)
	assert.Nil(t, err)

	paramsIn := NewDataPipe()
	paramsIn.Add(testParams)

	err = httpAction.AddInput(HTTPActionInputURLParams, paramsIn)
	assert.Nil(t, err)

	headersOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputHeaders, headersOut)
	assert.Nil(t, err)

	bodyOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputBody, bodyOut)
	assert.Nil(t, err)

	statusOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputStatusCode, statusOut)
	assert.Nil(t, err)

	err = httpAction.Run()
	assert.Nil(t, err)

	gotBody, ok1 := bodyOut.Remove().([]byte)
	assert.True(t, ok1)
	assert.Equal(t, expectedBody, gotBody)

	gotHeaders, ok2 := headersOut.Remove().(http.Header)
	assert.True(t, ok2)
	assert.True(t, len(gotHeaders) > 1)
	assert.Equal(t, "TestServer", gotHeaders["Server"][0])

	gotStatus, ok3 := statusOut.Remove().(int)
	assert.True(t, ok3)
	assert.Equal(t, 200, gotStatus)

}

func TestHTTPActionRunHEAD(t *testing.T) {
	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodHead, req.Method)

			res.WriteHeader(200)
		}))

	defer testServer.Close()

	httpAction := NewHTTPAction(testServer.URL, http.MethodHead, false)

	err := httpAction.Run()
	assert.Nil(t, err)
}

func TestHTTPActionRunPOST(t *testing.T) {
	expectedBody := []byte("Test Payload")

	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodPost, req.Method)

			body, err := io.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, expectedBody, body)

			res.WriteHeader(201)
		}))

	defer testServer.Close()

	bodyIn := NewDataPipe()
	bodyIn.Add(expectedBody)

	httpAction := NewHTTPAction(testServer.URL, http.MethodPost, false)

	err := httpAction.AddInput(HTTPActionInputBody, bodyIn)
	assert.Nil(t, err)

	err = httpAction.Run()
	assert.Nil(t, err)
}

func TestHTTPActionRunCookies(t *testing.T) {
	cookieName := "SessionID"
	cookieValue1 := "DFB32DF6-ABB5-4877-9B6A-D7F8D9791880"
	cookieValue2 := "2FED9EDE-E3DD-40EA-8F53-A9FBD8959F49"

	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodGet, req.Method)

			foundExpectedCookie := false

			for _, cookie := range req.Cookies() {
				if cookie.Name == cookieName && cookie.Value == cookieValue1 {
					foundExpectedCookie = true
					break
				}
			}

			assert.True(t, foundExpectedCookie)

			newCookie := &http.Cookie{
				Name:  cookieName,
				Value: cookieValue2,
			}

			http.SetCookie(res, newCookie)

			res.WriteHeader(200)
		}))

	defer testServer.Close()

	cookieDict := map[string]string{cookieName: cookieValue1}

	cookiesIn := NewDataPipe()
	cookiesIn.Add(cookieDict)

	httpAction := NewHTTPAction(testServer.URL, http.MethodGet, false)
	err := httpAction.AddInput(HTTPActionInputCookies, cookiesIn)
	assert.Nil(t, err)

	cookiesOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputCookies, cookiesOut)
	assert.Nil(t, err)

	err = httpAction.Run()
	assert.Nil(t, err)

	var gotCookies map[string]string
	var ok bool

	gotCookies, ok = cookiesOut.Remove().(map[string]string)
	assert.True(t, ok)

	foundExpectedCookie := false

	for name, value := range gotCookies {
		if name == cookieName && value == cookieValue2 {
			foundExpectedCookie = true
		}
	}

	assert.True(t, foundExpectedCookie)
}

func TestAddInput(t *testing.T) {
	baseURL := "https://httpbin.org/post"
	method := "POST"

	httpAction := NewHTTPAction(baseURL, method, false)

	dp := NewDataPipe()

	err := httpAction.AddInput("bad_name", dp)
	assert.NotNil(t, err)

	err = httpAction.AddInput(HTTPActionInputURLParams, dp)
	assert.Nil(t, err)
	assert.Equal(t, httpAction.AbstractAction.Inputs[HTTPActionInputURLParams], dp)

}

func TestAddOutput(t *testing.T) {
	baseURL := "https://httpbin.org/post"
	method := "POST"

	httpAction := NewHTTPAction(baseURL, method, false)

	dp := NewDataPipe()

	err := httpAction.AddOutput("bad_name", dp)
	assert.NotNil(t, err)

	err = httpAction.AddOutput(HTTPActionOutputBody, dp)
	assert.Nil(t, err)
	assert.Equal(t, dp, httpAction.AbstractAction.Outputs[HTTPActionOutputBody][0])
}

func TestUTF8EncodeActionRun(t *testing.T) {
	str := "abc"

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(str)

	utf8EncodeAction := NewUTF8EncodeAction()

	utf8EncodeAction.AddInput(UTF8EncodeActionInputStr, dataPipeIn)
	utf8EncodeAction.AddOutput(UTF8EncodeActionOutputBytes, dataPipeOut)

	err := utf8EncodeAction.Run()
	assert.Nil(t, err)

	binData, ok := dataPipeOut.Remove().([]byte)
	assert.True(t, ok)

	assert.Equal(t, binData, []byte{0x61, 0x62, 0x63})
}

func TestXPathActionRunBasic(t *testing.T) {
	htmlStr := "<html><body><title>This is title!</title></body></html>"

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(htmlStr)

	xpathAction := NewXPathAction("//title/text()", false)

	xpathAction.AddInput(XPathActionInputHTMLStr, dataPipeIn)
	xpathAction.AddOutput(XPathActionOutputStr, dataPipeOut)

	err := xpathAction.Run()
	assert.Nil(t, err)

	resultStr, ok := dataPipeOut.Remove().(string)
	assert.True(t, ok)

	assert.Equal(t, "This is title!", resultStr)
}

func TestXPathActionRunMultipleResults(t *testing.T) {
	htmlStr := "<p>1</p><p>2</p><p>3</p>"

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(htmlStr)

	xpathAction := NewXPathAction("//p/text()", true)

	xpathAction.AddInput(XPathActionInputHTMLStr, dataPipeIn)
	xpathAction.AddOutput(XPathActionOutputStr, dataPipeOut)

	err := xpathAction.Run()
	assert.Nil(t, err)

	resultStr, ok := dataPipeOut.Remove().(string)
	assert.True(t, ok)
	assert.Equal(t, "3", resultStr)

	resultStr, ok = dataPipeOut.Remove().(string)
	assert.True(t, ok)
	assert.Equal(t, "2", resultStr)

	resultStr, ok = dataPipeOut.Remove().(string)
	assert.True(t, ok)
	assert.Equal(t, "1", resultStr)

	_, ok = dataPipeOut.Remove().(string)
	assert.False(t, ok)
}

func TestXPathActionBadInput(t *testing.T) {
	// https://datatracker.ietf.org/doc/html/rfc5735
	inputStr := "192.0.2.16"

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(inputStr)

	xpathAction := NewXPathAction("//a/@href", true)

	xpathAction.AddInput(XPathActionInputHTMLStr, dataPipeIn)
	xpathAction.AddOutput(XPathActionOutputStr, dataPipeOut)

	xpathAction.Run() // Must not crash.
}

func TestXPathActionBadXPath(t *testing.T) {
	inputStr := "<html><body><a href=\"/next-gen-product\">Next gen product</a></body></html>"

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(inputStr)

	// Missing bracket in XPath.
	xpathAction := NewXPathAction("//a[contains(@href, \"next\")", true)

	xpathAction.AddInput(XPathActionInputHTMLStr, dataPipeIn)
	xpathAction.AddOutput(XPathActionOutputStr, dataPipeOut)

	err := xpathAction.Run() // Must not crash.
	assert.NotNil(t, err)
}

func TestSortActionsTopologically(t *testing.T) {
	task := NewTask("testTask", "", "")

	a1 := NewHTTPAction("https://cryptome.org/", "GET", false)
	a2 := NewXPathAction("//title/text()", true)
	a3 := NewXPathAction("//body", true)
	a4 := NewXPathAction("//h1/text()", true)

	task.Actions = []Action{a1, a2, a3, a4}

	in1 := NewDataPipe()
	in2 := NewDataPipe()

	err := a1.AddInput(HTTPActionInputURLParams, in1)
	assert.Nil(t, err)
	err = a1.AddInput(HTTPActionInputHeaders, in2)
	assert.Nil(t, err)

	in1.ToAction = a1
	in2.ToAction = a2

	out1 := NewDataPipe()
	out2 := NewDataPipe()

	out1.FromAction = a2
	out2.FromAction = a4

	err = a2.AddOutput(XPathActionInputHTMLBytes, out1)
	assert.NotNil(t, err)
	err = a4.AddOutput(XPathActionInputHTMLBytes, out2)
	assert.NotNil(t, err)

	task.Inputs["in1"] = in1
	task.DataPipes = append(task.DataPipes, in1)
	task.Inputs["in2"] = in2
	task.DataPipes = append(task.DataPipes, in2)

	task.Outputs["out1"] = out1
	task.DataPipes = append(task.DataPipes, out1)
	task.Outputs["out2"] = out2
	task.DataPipes = append(task.DataPipes, out2)

	dpA1ToA2 := NewDataPipeBetweenActions(a1, a2)
	err = a1.AddOutput(HTTPActionOutputBody, dpA1ToA2)
	assert.Nil(t, err)
	err = a2.AddInput(XPathActionInputHTMLBytes, dpA1ToA2)
	assert.Nil(t, err)
	task.DataPipes = append(task.DataPipes, dpA1ToA2)

	dpA1ToA3 := NewDataPipeBetweenActions(a1, a3)
	// HACK: this is invalid
	// TODO: make Action support multiple outputs for the same name
	a1.AddOutput(HTTPActionOutputHeaders, dpA1ToA3)
	a3.AddInput(XPathActionInputHTMLBytes, dpA1ToA3)
	task.DataPipes = append(task.DataPipes, dpA1ToA3)

	dpA3ToA4 := NewDataPipeBetweenActions(a3, a4)
	err = a3.AddOutput(XPathActionOutputStr, dpA3ToA4)
	assert.Nil(t, err)
	err = a4.AddInput(XPathActionInputHTMLStr, dpA3ToA4)
	assert.Nil(t, err)
	task.DataPipes = append(task.DataPipes, dpA3ToA4)

	actions := task.sortActionsTopologically()

	assert.NotNil(t, actions)
	assert.Equal(t, 4, len(actions))
	assert.Equal(t, a1, actions[0])
	assert.True(t, actions[3] == a4 || actions[3] == a2)
}

func TestUTF8DecodeActionMultipleOutputs(t *testing.T) {
	action := NewUTF8DecodeAction()

	input := NewDataPipe()
	output1 := NewDataPipe()
	output2 := NewDataPipe()

	err := action.AddInput(UTF8DecodeActionInputBytes, input)
	assert.Nil(t, err)

	err = action.AddOutput(UTF8DecodeActionOutputStr, output1)
	assert.Nil(t, err)

	err = action.AddOutput(UTF8DecodeActionOutputStr, output2)
	assert.Nil(t, err)

	b := []byte("123")

	input.Add(b)

	err = action.Run()
	assert.Nil(t, err)

	s1, ok1 := output1.Remove().(string)
	assert.True(t, ok1)
	assert.Equal(t, "123", s1)

	s2, ok2 := output2.Remove().(string)
	assert.True(t, ok2)
	assert.Equal(t, "123", s2)
}

func TestFieldJoinActionRun(t *testing.T) {
	action := NewFieldJoinAction([]string{"Name", "Surname", "Phone", "Email"})

	nameIn := NewDataPipe()
	surnameIn := NewDataPipe()
	phoneIn := NewDataPipe()
	emailIn := NewDataPipe()

	nameIn.Add("John")
	surnameIn.Add("Smith")
	phoneIn.Add("555-1212")
	emailIn.Add("john@smith.int")

	err := action.AddInput("Name", nameIn)
	assert.Nil(t, err)
	err = action.AddInput("Surname", surnameIn)
	assert.Nil(t, err)
	err = action.AddInput("Phone", phoneIn)
	assert.Nil(t, err)
	err = action.AddInput("Email", emailIn)
	assert.Nil(t, err)

	itemOut := NewDataPipe()

	err = action.AddOutput(FieldJoinActionOutputItem, itemOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	expectedItem := map[string]string{
		"Name":    "John",
		"Surname": "Smith",
		"Phone":   "555-1212",
		"Email":   "john@smith.int",
	}

	item, ok := itemOut.Remove().(map[string]string)
	assert.True(t, ok)

	assert.Equal(t, expectedItem, item)
}
