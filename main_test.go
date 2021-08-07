package main

import (
	spsw "github.com/rl1987/spiderswarm/lib"

	"github.com/davecgh/go-spew/spew"
)

func ExampleHTTPAction() {
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

func ExampleTask() {
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

func ExampleWorkflow() {
	workflow := &spsw.Workflow{
		Name:    "testWorkflow",
		Version: "v0.0.0.0.1",
		TaskTemplates: []spsw.TaskTemplate{
			spsw.TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []spsw.ActionTemplate{
					spsw.ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://news.ycombinator.com/",
							"method":  "GET",
							"canFail": false,
						},
					},
					spsw.ActionTemplate{
						Name:              "UTF8Decode",
						StructName:        "UTF8DecodeAction",
						ConstructorParams: map[string]interface{}{},
					},
					spsw.ActionTemplate{
						Name:       "MakePromise",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"htmlStr1", "htmlStr2"},
							"taskName":   "ParseHTML",
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						SourceActionName: "HTTP1",
						SourceOutputName: spsw.HTTPActionOutputBody,
						DestActionName:   "UTF8Decode",
						DestInputName:    spsw.UTF8DecodeActionInputBytes,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: spsw.UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr1",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: spsw.UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr2",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "MakePromise",
						SourceOutputName: spsw.TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
			spsw.TaskTemplate{
				TaskName: "ParseHTML",
				Initial:  false,
				ActionTemplates: []spsw.ActionTemplate{
					spsw.ActionTemplate{
						Name:       "TitleExtraction",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//a[@class='storylink']/text()",
							"expectMany": true,
						},
					},
					spsw.ActionTemplate{
						Name:       "LinkExtraction",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//a[@class='storylink']/@href",
							"expectMany": true,
						},
					},
					spsw.ActionTemplate{
						Name:       "YieldItem",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"title", "link"},
							"itemName":   "HNItem",
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						TaskInputName:  "htmlStr1",
						DestActionName: "TitleExtraction",
						DestInputName:  spsw.XPathActionInputHTMLStr,
					},
					spsw.DataPipeTemplate{
						TaskInputName:  "htmlStr2",
						DestActionName: "LinkExtraction",
						DestInputName:  spsw.XPathActionInputHTMLStr,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "TitleExtraction",
						SourceOutputName: spsw.XPathActionOutputStr,
						DestActionName:   "YieldItem",
						DestInputName:    "title",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "LinkExtraction",
						SourceOutputName: spsw.XPathActionOutputStr,
						DestActionName:   "YieldItem",
						DestInputName:    "link",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "YieldItem",
						SourceOutputName: spsw.FieldJoinActionOutputItem,
						TaskOutputName:   "items",
					},
				},
			},
		},
	}

	spew.Dump(workflow)

	items, err := workflow.Run()
	if err != nil {
		spew.Dump(err)
	} else {
		spew.Dump(items)
	}

}
