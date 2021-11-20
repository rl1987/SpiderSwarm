package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	spsw "github.com/spiderswarm/spiderswarm/lib"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func printUsage() {
	fmt.Println("Read the code for now")
}

func getWorkflow() *spsw.Workflow {
	return &spsw.Workflow{
		Name:    "testWorkflow",
		Version: "v0.0.0.0.1",
		TaskTemplates: []spsw.TaskTemplate{
			spsw.TaskTemplate{
				TaskName: "TestTask",
				Initial:  true,
				ActionTemplates: []spsw.ActionTemplate{
					spsw.ActionTemplate{
						Name:       "GetJSON",
						StructName: "HTTPAction",
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "https://ifconfig.me/all.json",
							},
							"method": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "GET",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "ExtractIP",
						StructName: "JQAction",
						ConstructorParams: map[string]spsw.Value{
							"jqArgs": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{".ip_addr"},
							},
							"expectMany": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: false,
							},
							"canFail": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: false,
							},
							"decodeOutput": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "MakeItem",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"ip_addr"},
							},
							"itemName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "ip_addr",
							},
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						SourceActionName: "GetJSON",
						SourceOutputName: spsw.HTTPActionOutputBody,
						DestActionName:   "ExtractIP",
						DestInputName:    spsw.JQActionInputJQStdinBytes,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "ExtractIP",
						SourceOutputName: spsw.JQActionOutputJQStdoutStr,
						DestActionName:   "MakeItem",
						DestInputName:    "ip_addr",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "MakeItem",
						SourceOutputName: spsw.FieldJoinActionOutputItem,
						TaskOutputName:   "ip_addr_item",
					},
				},
			},
		},
	}
}

func runTestWorkflow() {
	backendAddr := "127.0.0.1:6379"
	workflow := getWorkflow()

	runner := spsw.NewRunner(backendAddr)

	runner.RunSingleNode(4, "/tmp", workflow)
}

func NewAbstractAction(actionTempl *spsw.ActionTemplate, workflowName string) spsw.Action {
	return &spsw.AbstractAction{}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	spsw.ActionConstructorTable["dummy"] = NewAbstractAction

	spew.Dump(spsw.ActionConstructorTable)

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
