package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func initLogging() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	initLogging()
	log.Info("Starting spiderswarm instance...")

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
