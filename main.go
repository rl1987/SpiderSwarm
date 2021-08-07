package main

import (
	"os"

	spsw "github.com/rl1987/spiderswarm/lib"

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

	workflow := &spsw.Workflow{
		Name: "FCC_telecom",
		TaskTemplates: []spsw.TaskTemplate{
			spsw.TaskTemplate{
				TaskName: "ScrapeStates",
				Initial:  true,
				ActionTemplates: []spsw.ActionTemplate{
					spsw.ActionTemplate{
						Name:       "HTTP_Form",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://apps.fcc.gov/cgb/form499/499a.cfm",
							"method":  "GET",
							"canFail": false,
						},
					},
					spsw.ActionTemplate{
						Name:       "XPath_states",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//select[@name=\"state\"]/option[not(@selected)]/@value",
							"expectMany": true,
						},
					},
					spsw.ActionTemplate{
						Name:       "TaskPromise_ScrapeList",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"state"},
							"taskName":   "ScrapeCompanyList",
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						SourceActionName: "HTTP_Form",
						SourceOutputName: spsw.HTTPActionOutputBody,
						DestActionName:   "XPath_states",
						DestInputName:    spsw.XPathActionInputHTMLBytes,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "XPath_states",
						SourceOutputName: spsw.XPathActionOutputStr,
						DestActionName:   "TaskPromise_ScrapeList",
						DestInputName:    "state",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "TaskPromise_ScrapeList",
						SourceOutputName: spsw.TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
			spsw.TaskTemplate{
				TaskName: "ScrapeCompanyList",
				Initial:  false,
				ActionTemplates: []spsw.ActionTemplate{
					spsw.ActionTemplate{
						Name:       "Const_commType",
						StructName: "ConstAction",
						ConstructorParams: map[string]interface{}{
							"c": "Any Type",
						},
					},
					spsw.ActionTemplate{
						Name:       "Const_R1",
						StructName: "ConstAction",
						ConstructorParams: map[string]interface{}{
							"c": "and",
						},
					},
					spsw.ActionTemplate{
						Name:       "Const_XML",
						StructName: "ConstAction",
						ConstructorParams: map[string]interface{}{
							"c": "FALSE",
						},
					},
					spsw.ActionTemplate{
						Name:       "JoinParams",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"commType", "R1", "state", "XML"},
							"itemName":   "params",
						},
					},
					spsw.ActionTemplate{
						Name:       "HTTP_List",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://apps.fcc.gov/cgb/form499/499results.cfm",
							"canFail": false,
						},
					},
					spsw.ActionTemplate{
						Name:       "XPath_Companies",
						StructName: "XPathAction",
						ConstructorParams: map[string]interface{}{
							"xpath":      "//table[@border=\"1\"]//a/@href",
							"expectMany": true,
						},
					},
					spsw.ActionTemplate{
						Name:       "TaskPromise_ScrapeCompanyPage",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"relativeURL"},
							"taskName":   "ScrapeCompanyPage",
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						TaskInputName:  "relativeURL",
						DestActionName: "JoinParams",
						DestInputName:  "state",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "Const_R1",
						SourceOutputName: spsw.ConstActionOutput,
						DestActionName:   "JoinParams",
						DestInputName:    "R1",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "Const_XML",
						SourceOutputName: spsw.ConstActionOutput,
						DestActionName:   "JoinParams",
						DestInputName:    "XML",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "Const_commType",
						SourceOutputName: spsw.ConstActionOutput,
						DestActionName:   "JoinParams",
						DestInputName:    "commType",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "JoinParams",
						SourceOutputName: spsw.FieldJoinActionOutputMap,
						DestActionName:   "HTTP_List",
						DestInputName:    spsw.HTTPActionInputURLParams,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "HTTP_List",
						SourceOutputName: spsw.HTTPActionOutputBody,
						DestActionName:   "XPath_Companies",
						DestInputName:    spsw.XPathActionInputHTMLBytes,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "XPath_Companies",
						SourceOutputName: spsw.XPathActionOutputStr,
						DestActionName:   "TaskPromise_ScrapeCompanyPage",
						DestInputName:    "relativeURL",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "TaskPromise_ScrapeCompanyPage",
						SourceOutputName: spsw.TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
			spsw.TaskTemplate{
				TaskName: "ScrapeCompanyPage",
				Initial:  false,
				ActionTemplates: []spsw.ActionTemplate{
					spsw.ActionTemplate{
						Name:       "URLJoin",
						StructName: "URLJoinAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://apps.fcc.gov/cgb/form499/",
						},
					},
					spsw.ActionTemplate{
						Name:       "HTTP_Company",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"method":  "GET",
							"canFail": false,
						},
					},
					spsw.ActionTemplate{
						Name:              "BodyBytesToStr",
						StructName:        "UTF8DecodeAction",
						ConstructorParams: map[string]interface{}{},
					},
					spsw.ActionTemplate{
						Name:       "GetFilerID",
						StructName: "StringCutAction",
						ConstructorParams: map[string]interface{}{
							"from": "499 Filer ID Number:                <b>",
							"to":   "</b>",
						},
					},
					spsw.ActionTemplate{
						Name:       "GetLegalName",
						StructName: "StringCutAction",
						ConstructorParams: map[string]interface{}{
							"from": "Legal Name of Reporting Entity:     <b>",
							"to":   "</b>",
						},
					},
					spsw.ActionTemplate{
						Name:       "GetDBA",
						StructName: "StringCutAction",
						ConstructorParams: map[string]interface{}{
							"from": "Doing Business As:                  <b>",
							"to":   "</b>",
						},
					},
					spsw.ActionTemplate{
						Name:       "GetPhone",
						StructName: "StringCutAction",
						ConstructorParams: map[string]interface{}{
							"from": "Customer Inquiries Telephone:       <b>",
							"to":   "</b>",
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{},
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
