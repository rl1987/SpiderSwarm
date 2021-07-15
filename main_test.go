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
