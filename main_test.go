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

func ExampleTask2() {
	httpAction := NewHTTPAction("https://news.ycombinator.com/news", "GET", false)

	titleXpathAction := NewXPathAction("//a[@class='storylink']/text()", true)
	linkXpathAction := NewXPathAction("//a[@class='storylink']/@href", true)

	task := NewTask("HN", "", "")

	titlesOut := NewDataPipe()
	linksOut := NewDataPipe()

	task.AddAction(httpAction)
	task.AddAction(titleXpathAction)
	task.AddAction(linkXpathAction)

	task.AddOutput("titles", titleXpathAction, XPathActionOutputStr, titlesOut)
	task.AddOutput("links", linkXpathAction, XPathActionOutputStr, linksOut)

	task.AddDataPipeBetweenActions(httpAction, HTTPActionOutputBody, titleXpathAction, XPathActionInputHTMLBytes)
	task.AddDataPipeBetweenActions(httpAction, HTTPActionOutputBody, linkXpathAction, XPathActionInputHTMLBytes)

	err := task.Run()
	if err != nil {
		spew.Dump(err)
		return
	}

	var titles []string
	var links []string

	for {
		if title, ok := titlesOut.Remove().(string); ok {
			titles = append(titles, title)
		} else {
			break
		}
	}

	for {
		if link, ok := linksOut.Remove().(string); ok {
			links = append(links, link)
		} else {
			break
		}
	}

	spew.Dump(titles)
	spew.Dump(links)

}
