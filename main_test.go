package spiderswarm

import (
	"github.com/davecgh/go-spew/spew"
)

func ExampleHTTPAction() {
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

func ExampleTask() {
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

func ExampleWorkflow() {
	workflow := &Workflow{
		Name:    "testWorkflow",
		Version: "v0.0.0.0.1",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://news.ycombinator.com/",
							"method":  "GET",
							"canFail": false,
						},
					},
					ActionTemplate{
						Name:              "UTF8Decode",
						StructName:        "UTF8DecodeAction",
						ConstructorParams: map[string]interface{}{},
					},
					ActionTemplate{
						Name:       "MakePromise",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"htmlStr1", "htmlStr2"},
							"taskName":   "ParseHTML",
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						SourceActionName: "HTTP1",
						SourceOutputName: HTTPActionOutputBody,
						DestActionName:   "UTF8Decode",
						DestInputName:    UTF8DecodeActionInputBytes,
					},
					DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr1",
					},
					DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr2",
					},
					DataPipeTemplate{
						SourceActionName: "MakePromise",
						SourceOutputName: TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
			TaskTemplate{
				TaskName: "ParseHTML",
				Initial:  false,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "TitleExtraction",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//a[@class='storylink']/text()",
							"expectMany": true,
						},
					},
					ActionTemplate{
						Name:       "LinkExtraction",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//a[@class='storylink']/@href",
							"expectMany": true,
						},
					},
					ActionTemplate{
						Name:       "YieldItem",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"title", "link"},
							"itemName":   "HNItem",
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "htmlStr1",
						DestActionName: "TitleExtraction",
						DestInputName:  XPathActionInputHTMLStr,
					},
					DataPipeTemplate{
						TaskInputName:  "htmlStr2",
						DestActionName: "LinkExtraction",
						DestInputName:  XPathActionInputHTMLStr,
					},
					DataPipeTemplate{
						SourceActionName: "TitleExtraction",
						SourceOutputName: XPathActionOutputStr,
						DestActionName:   "YieldItem",
						DestInputName:    "title",
					},
					DataPipeTemplate{
						SourceActionName: "LinkExtraction",
						SourceOutputName: XPathActionOutputStr,
						DestActionName:   "YieldItem",
						DestInputName:    "link",
					},
					DataPipeTemplate{
						SourceActionName: "YieldItem",
						SourceOutputName: FieldJoinActionOutputItem,
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
