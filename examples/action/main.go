package main

import (
	spsw "github.com/rl1987/spiderswarm/lib"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	httpAction := spsw.NewHTTPAction("https://cryptome.org", "GET", true)

	headers := map[string][]string{
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"},
	}

	headersIn := spsw.NewDataPipe()
	resultOut := spsw.NewDataPipe()

	httpAction.AddInput(spsw.HTTPActionInputHeaders, headersIn)
	xpathAction := spsw.NewXPathAction("//a/text()", true)

	task := spsw.NewTask("task1", "", "")
	task.AddAction(httpAction)
	task.AddAction(xpathAction)

	task.AddInput("headersIn", httpAction, spsw.HTTPActionInputHeaders, headersIn)
	task.AddOutput("resultOut", xpathAction, spsw.XPathActionOutputStr, resultOut)

	task.AddDataPipeBetweenActions(httpAction, spsw.HTTPActionOutputBody, xpathAction, spsw.XPathActionInputHTMLBytes)

	headersIn.Add(headers)

	spew.Dump(task)

	err := task.Run()
	if err != nil {
		spew.Dump(err)
	} else {
		spew.Dump(resultOut)
	}
}
