package main

import (
	"github.com/davecgh/go-spew/spew"
)

func ExampleBasicTask() {
	httpAction := NewHTTPAction("https://cryptome.org", "GET", true)

	headers := map[string][]string{
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"},
	}

	headersIn := NewDataPipe()
	resultOut := NewDataPipe()

	httpAction.AddInput(HTTPActionInputHeaders, headersIn)
	xpathAction := NewXPathAction("//a/text()", true)

	task := NewTask("task1", "", "")
	task.AddAction(httpAction)
	task.AddAction(xpathAction)

	task.AddInput("headersIn", httpAction, HTTPActionInputHeaders, headersIn)
	task.AddOutput("resultOut", xpathAction, XPathActionOutputStr, resultOut)

	task.AddDataPipeBetweenActions(httpAction, HTTPActionOutputBody, xpathAction, XPathActionInputHTMLBytes)

	headersIn.Add(headers)

	spew.Dump(task)

	err := task.Run()
	if err != nil {
		spew.Dump(err)
	} else {
		spew.Dump(resultOut)
	}

}
