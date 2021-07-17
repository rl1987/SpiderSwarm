package main

import (
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
	assert.Equal(t, len(httpAction.AbstractAction.AllowedInputNames), 3)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[0], HTTPActionInputURLParams)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[1], HTTPActionInputHeaders)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[2], HTTPActionInputCookies)
	assert.Equal(t, len(httpAction.AbstractAction.AllowedOutputNames), 3)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[0], HTTPActionOutputBody)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[1], HTTPActionOutputHeaders)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[2], HTTPActionOutputStatusCode)

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
	assert.Equal(t, httpAction.AbstractAction.Outputs[HTTPActionOutputBody], dp)
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
	assert.Equal(t, a4, actions[3])
}
