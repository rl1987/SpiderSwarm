package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("spiderswarm")

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
