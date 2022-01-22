package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"strconv"

	spsw "github.com/spiderswarm/spiderswarm/lib"

	log "github.com/sirupsen/logrus"
)

func printUsage() {
	fmt.Println("Read the code for now")
}

func getWorkflow() *spsw.Workflow {
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
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "https://apps.fcc.gov/cgb/form499/499a.cfm",
							},
							"method": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "GET",
							},
							"canFail": spsw.Value{
								ValueType: spsw.ValueTypeString,
								BoolValue: false,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "XPath_states",
						StructName: "XPathAction",
						ConstructorParams: map[string]spsw.Value{
							"xpath": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "//select[@name=\"state\"]/option[not(@selected)]/@value",
							},
							"expectMany": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "TaskPromise_ScrapeList",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"state", "cookies"},
							},
							"taskName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "ScrapeCompanyList",
							},
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
						SourceActionName: "HTTP_Form",
						SourceOutputName: spsw.HTTPActionOutputCookies,
						DestActionName:   "TaskPromise_ScrapeList",
						DestInputName:    "cookies",
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
						ConstructorParams: map[string]spsw.Value{
							"c": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "Any Type",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "Const_R1",
						StructName: "ConstAction",
						ConstructorParams: map[string]spsw.Value{
							"c": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "and",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "Const_XML",
						StructName: "ConstAction",
						ConstructorParams: map[string]spsw.Value{
							"c": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "FALSE",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "JoinParams",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"comm_type", "R1", "state", "XML"},
							},
							"itemName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "params",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "HTTP_List",
						StructName: "HTTPAction",
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "https://apps.fcc.gov/cgb/form499/499results.cfm",
							},
							"canFail": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: false,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "XPath_Companies",
						StructName: "XPathAction",
						ConstructorParams: map[string]spsw.Value{
							"xpath": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "//table[@border=\"1\"]//a/@href",
							},
							"expectMany": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					spsw.ActionTemplate{
						Name:              "JoinCookies",
						StructName:        "StringMapUpdateAction",
						ConstructorParams: map[string]spsw.Value{},
					},
					spsw.ActionTemplate{
						Name:       "TaskPromise_ScrapeCompanyPage",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"relativeURL", "cookies"},
							},
							"taskName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "ScrapeCompanyPage",
							},
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						TaskInputName:  "state",
						DestActionName: "JoinParams",
						DestInputName:  "state",
					},
					spsw.DataPipeTemplate{
						TaskInputName:  "cookies",
						DestActionName: "HTTP_List",
						DestInputName:  spsw.HTTPActionInputCookies,
					},
					spsw.DataPipeTemplate{
						TaskInputName:  "cookies",
						DestActionName: "JoinCookies",
						DestInputName:  spsw.StringMapUpdateActionInputOld,
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
						DestInputName:    "comm_type",
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
						SourceActionName: "HTTP_List",
						SourceOutputName: spsw.HTTPActionOutputCookies,
						DestActionName:   "JoinCookies",
						DestInputName:    spsw.StringMapUpdateActionInputNew,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "JoinCookies",
						SourceOutputName: spsw.StringMapUpdateActionoOutputUpdated,
						DestActionName:   "TaskPromise_ScrapeCompanyPage",
						DestInputName:    "cookies",
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
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "https://apps.fcc.gov/cgb/form499/",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "HTTP_Company",
						StructName: "HTTPAction",
						ConstructorParams: map[string]spsw.Value{
							"method": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "GET",
							},
							"canFail": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: false,
							},
						},
					},
					spsw.ActionTemplate{
						Name:              "BodyBytesToStr",
						StructName:        "UTF8DecodeAction",
						ConstructorParams: map[string]spsw.Value{},
					},
					spsw.ActionTemplate{
						Name:       "GetFilerID",
						StructName: "StringCutAction",
						ConstructorParams: map[string]spsw.Value{
							"from": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "499 Filer ID Number:                <b>",
							},
							"to": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "</b>",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "GetLegalName",
						StructName: "StringCutAction",
						ConstructorParams: map[string]spsw.Value{
							"from": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "Legal Name of Reporting Entity:     <b>",
							},
							"to": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "</b>",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "GetDBA",
						StructName: "StringCutAction",
						ConstructorParams: map[string]spsw.Value{
							"from": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "Doing Business As:                  <b>",
							},
							"to": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "</b>",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "GetPhone",
						StructName: "StringCutAction",
						ConstructorParams: map[string]spsw.Value{
							"from": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "Customer Inquiries Telephone:       <b>",
							},
							"to": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "</b>",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "MakeItem",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"filer_id", "legal_name", "dba", "phone"},
							},
							"itemName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "company",
							},
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						TaskInputName:  "relativeURL",
						DestActionName: "URLJoin",
						DestInputName:  spsw.URLJoinActionInputRelativeURL,
					},
					spsw.DataPipeTemplate{
						TaskInputName:  "cookies",
						DestActionName: "HTTP_Company",
						DestInputName:  spsw.HTTPActionInputCookies,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "URLJoin",
						SourceOutputName: spsw.URLJoinActionOutputAbsoluteURL,
						DestActionName:   "HTTP_Company",
						DestInputName:    spsw.HTTPActionInputBaseURL,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "HTTP_Company",
						SourceOutputName: spsw.HTTPActionOutputBody,
						DestActionName:   "BodyBytesToStr",
						DestInputName:    spsw.UTF8DecodeActionInputBytes,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "BodyBytesToStr",
						SourceOutputName: spsw.UTF8DecodeActionOutputStr,
						DestActionName:   "GetFilerID",
						DestInputName:    spsw.StringCutActionInputStr,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "BodyBytesToStr",
						SourceOutputName: spsw.UTF8DecodeActionOutputStr,
						DestActionName:   "GetLegalName",
						DestInputName:    spsw.StringCutActionInputStr,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "BodyBytesToStr",
						SourceOutputName: spsw.UTF8DecodeActionOutputStr,
						DestActionName:   "GetDBA",
						DestInputName:    spsw.StringCutActionInputStr,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "BodyBytesToStr",
						SourceOutputName: spsw.UTF8DecodeActionOutputStr,
						DestActionName:   "GetPhone",
						DestInputName:    spsw.StringCutActionInputStr,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "GetFilerID",
						SourceOutputName: spsw.StringCutActionOutputStr,
						DestActionName:   "MakeItem",
						DestInputName:    "filer_id",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "GetLegalName",
						SourceOutputName: spsw.StringCutActionOutputStr,
						DestActionName:   "MakeItem",
						DestInputName:    "legal_name",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "GetDBA",
						SourceOutputName: spsw.StringCutActionOutputStr,
						DestActionName:   "MakeItem",
						DestInputName:    "dba",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "GetPhone",
						SourceOutputName: spsw.StringCutActionOutputStr,
						DestActionName:   "MakeItem",
						DestInputName:    "phone",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "MakeItem",
						SourceOutputName: spsw.FieldJoinActionOutputItem,
						TaskOutputName:   "items",
					},
				},
			},
		},
	}

	return workflow
}
func runTestWorkflow() {
	backendAddr := "127.0.0.1:6379"
	workflow := getWorkflow()

	runner := spsw.NewRunner(backendAddr)

	runner.RunSingleNode(4, "./", workflow)
}

func NewAbstractAction(actionTempl *spsw.ActionTemplate, workflowName string) spsw.Action {
	return &spsw.AbstractAction{}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	workflow := getWorkflow()

	singleNodeCmd := flag.NewFlagSet("singlenode", flag.ExitOnError)
	singleNodeWorkers := singleNodeCmd.Int("workers", 1, "Number of worker goroutines")

	runner := &spsw.Runner{}

	switch os.Args[1] {
	case "singlenode":
		singleNodeCmd.Parse(os.Args[2:])
		log.Info(fmt.Sprintf("Number of worker goroutines: %d", *singleNodeWorkers))
		log.Error("Not implemented")
	case "worker":
		n, _ := strconv.Atoi(os.Args[2])
		backendAddr := os.Args[3]
		runner.BackendAddr = backendAddr
		runner.RunWorkers(n)
		for {
			select {}
		}
	case "manager":
		backendAddr := os.Args[2]
		runner.BackendAddr = backendAddr
		runner.RunManager(workflow)
		for {
			select {}
		}
	case "exporter":
		outputDir := os.Args[2]
		backendAddr := os.Args[3]
		runner.BackendAddr = backendAddr
		runner.RunExporter(outputDir)
		for {
			select {}
		}
	case "client":
		// TODO: client for REST API
		log.Error("client part not implemented yet")
	default:
		runTestWorkflow()
	}
}
