package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	resultStrings, ok := dataPipeOut.Remove().([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"1", "2", "3"}, resultStrings)

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
