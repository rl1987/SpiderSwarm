package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("SpiderSwarm")
	httpAction := NewHTTPAction("https://ifconfig.me", "GET", true)

	headers := map[string][]string{
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"},
	}

	headersIn := NewDataPipe()

	headersIn.Add(headers)

	httpAction.AddInput(HTTPActionInputHeaders, headersIn)

	bodyOut := NewDataPipe()
	httpAction.AddOutput(HTTPActionOutputBody, bodyOut)

	headersOut := NewDataPipe()
	httpAction.AddOutput(HTTPActionOutputHeaders, headersOut)

	statusCodeOut := NewDataPipe()
	httpAction.AddOutput(HTTPActionOutputStatusCode, statusCodeOut)

	err := httpAction.Run()
	if err != nil {
		fmt.Println(err)
	}

	spew.Dump(bodyOut)

	xpathAction := NewXPathAction("//title/text()", false)

	_ = xpathAction.AddInput(XPathActionInputHTMLBytes, bodyOut)

	resultOut := NewDataPipe()

	_ = xpathAction.AddOutput(XPathActionOutputStr, resultOut)

	err = xpathAction.Run()
	if err != nil {
		fmt.Println(err)
	}

	spew.Dump(xpathAction)
}
