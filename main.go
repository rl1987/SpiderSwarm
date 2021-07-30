package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("spiderswarm")

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
						SourceActionName: "HTTP1",
						SourceOutputName: HTTPActionOutputBody,
						DestActionName:   "TitleExtraction",
						DestInputName:    XPathActionInputHTMLBytes,
					},
					DataPipeTemplate{
						SourceActionName: "HTTP1",
						SourceOutputName: HTTPActionOutputBody,
						DestActionName:   "LinkExtraction",
						DestInputName:    XPathActionInputHTMLBytes,
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
						TaskOutputName:   "HNOutput",
					},
				},
			},
		},
	}

	spew.Dump(workflow)
}
