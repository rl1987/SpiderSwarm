package main

import (
	spsw "github.com/rl1987/spiderswarm/lib"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	httpAction := spsw.NewHTTPAction("https://news.ycombinator.com/news", "GET", false)

	titleXpathAction := spsw.NewXPathAction("//a[@class='storylink']/text()", true)
	linkXpathAction := spsw.NewXPathAction("//a[@class='storylink']/@href", true)

	task := spsw.NewTask("HN", "", "")

	titlesOut := spsw.NewDataPipe()
	linksOut := spsw.NewDataPipe()

	task.AddAction(httpAction)
	task.AddAction(titleXpathAction)
	task.AddAction(linkXpathAction)

	task.AddOutput("titles", titleXpathAction, spsw.XPathActionOutputStr, titlesOut)
	task.AddOutput("links", linkXpathAction, spsw.XPathActionOutputStr, linksOut)

	task.AddDataPipeBetweenActions(httpAction, spsw.HTTPActionOutputBody, titleXpathAction, spsw.XPathActionInputHTMLBytes)
	task.AddDataPipeBetweenActions(httpAction, spsw.HTTPActionOutputBody, linkXpathAction, spsw.XPathActionInputHTMLBytes)

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
