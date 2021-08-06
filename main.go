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

	// https://apps.fcc.gov/cgb/form499/499a.cfm
	// https://apps.fcc.gov/cgb/form499/499results.cfm?comm_type=Any+Type&state=alaska&R1=and&XML=FALSE

	workflow := &Workflow{
		Name: "FCC_telecom",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "ScrapeStates",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP_Form",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://apps.fcc.gov/cgb/form499/499a.cfm",
							"method":  "GET",
							"canFail": false,
						},
					},
					ActionTemplate{
						Name:       "XPath_states",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//select[@name=\"state\"]/option[not(@selected)]/@value",
							"expectMany": true,
						},
					},
					ActionTemplate{
						Name:       "TaskPromise_ScrapeList",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"states"},
							"taskName":   "ScrapeCompanyList",
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						SourceActionName: "HTTP_Form",
						SourceOutputName: HTTPActionOutputBody,
						DestActionName:   "XPath_states",
						DestInputName:    XPathActionInputHTMLBytes,
					},
					DataPipeTemplate{
						SourceActionName: "XPath_states",
						SourceOutputName: XPathActionOutputStr,
						DestActionName:   "TaskPromise_ScrapeList",
						DestInputName:    "states",
					},
					DataPipeTemplate{
						SourceActionName: "TaskPromise_ScrapeList",
						SourceOutputName: TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
			TaskTemplate{
				TaskName: "ScrapeCompanyList",
				Initial:  false,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "Const_commType",
						StructName: "ConstAction",
						ConstructorParams: map[string]interface{}{
							"c": "Any Type",
						},
					},
					ActionTemplate{
						Name:       "Const_R1",
						StructName: "ConstAction",
						ConstructorParams: map[string]interface{}{
							"c": "and",
						},
					},
					ActionTemplate{
						Name:       "Const_XML",
						StructName: "ConstAction",
						ConstructorParams: map[string]interface{}{
							"c": "FALSE",
						},
					},
					ActionTemplate{
						Name:       "JoinParams",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"commType", "R1", "state", "XML"},
							"itemName":   "params",
						},
					},
					ActionTemplate{
						Name:       "HTTP_List",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://apps.fcc.gov/cgb/form499/499results.cfm",
							"canFail": false,
						},
					},
					ActionTemplate{
						Name:       "XPath_Companies",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//table[@border=\"1\"]//a/@href",
							"expectMany": true,
						},
					},
					ActionTemplate{
						Name:       "TaskPromise_ScrapeCompanyPage",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"relativeURL"},
							"taskName":   "ScrapeCompanyPage",
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "relativeURL",
						DestActionName: "JoinParams",
						DestInputName:  "state",
					},
					DataPipeTemplate{
						SourceActionName: "Const_R1",
						SourceOutputName: ConstActionOutput,
						DestActionName:   "JoinParams",
						DestInputName:    "R1",
					},
					DataPipeTemplate{
						SourceActionName: "Const_XML",
						SourceOutputName: ConstActionOutput,
						DestActionName:   "JoinParams",
						DestInputName:    "XML",
					},
					DataPipeTemplate{
						SourceActionName: "Const_commType",
						SourceOutputName: ConstActionOutput,
						DestActionName:   "JoinParams",
						DestInputName:    "commType",
					},
					DataPipeTemplate{
						SourceActionName: "JoinParams",
						SourceOutputName: FieldJoinActionOutputMap,
						DestActionName:   "HTTP_List",
						DestInputName:    HTTPActionInputURLParams,
					},
					DataPipeTemplate{
						SourceActionName: "HTTP_Action",
						SourceOutputName: HTTPActionOutputBody,
						DestActionName:   "XPath_Companies",
						DestInputName:    XPathActionInputHTMLBytes,
					},
					DataPipeTemplate{
						SourceActionName: "XPath_Companies",
						SourceOutputName: XPathActionOutputStr,
						DestActionName:   "TaskPromise_ScrapeCompanyPage",
						DestInputName:    "relativeURL",
					},
					DataPipeTemplate{
						SourceActionName: "TaskPromise_ScrapeCompanyPage",
						SourceOutputName: TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
			TaskTemplate{
				TaskName: "ScrapeCompanyPage",
				Initial:  false,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "URLJoin",
						StructName: "URLJoinAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://apps.fcc.gov/cgb/form499/",
						},
					},
					ActionTemplate{
						Name:       "HTTP_Company",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"method":  "GET",
							"canFail": false,
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{},
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
