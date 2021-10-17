package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	spsw "github.com/spiderswarm/spiderswarm/lib"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func initLogging() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func printUsage() {
	fmt.Println("Read the code for now")
}

func runTestWorkflow() {
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
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "https://news.ycombinator.com/",
							},
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
						Name:              "UTF8Decode",
						StructName:        "UTF8DecodeAction",
						ConstructorParams: map[string]spsw.Value{},
					},
					spsw.ActionTemplate{
						Name:       "MakePromise",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"htmlStr1", "htmlStr2"},
							},
							"taskName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "ParseHTML",
							},
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
						ConstructorParams: map[string]spsw.Value{
							"xpath": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "//a[@class='storylink']/text()",
							},
							"expectMany": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "LinkExtraction",
						StructName: "XPathAction",
						ConstructorParams: map[string]spsw.Value{
							"xpath": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "//a[@class='storylink']/@href",
							},
							"expectMany": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "YieldItem",
						StructName: "FieldJoinAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"title", "link"},
							},
							"itemName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "HNItem",
							},
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

	spiderBusBackend := spsw.NewMySQLSpiderBusBackend("username:password@tcp(159.223.27.109:3306)/spsw")
	spiderBus := spsw.NewSpiderBus()
	spiderBus.Backend = spiderBusBackend

	manager := spsw.NewManager()

	manager.StartScrapingJob(workflow)

	exporter := spsw.NewExporter()
	// TODO: make ExporterBackend API more abstract to enable plugin architecture.
	exporterBackend := spsw.NewCSVExporterBackend("/tmp")

	// FIXME: refrain from hardcoding field names; consider finding them from
	// Workflow.
	err := exporterBackend.StartExporting(manager.JobUUID, []string{"link", "title"})
	if err != nil {
		spew.Dump(err)
		return
	}

	exporter.AddBackend(exporterBackend)

	managerAdapter := spsw.NewSpiderBusAdapterForManager(spiderBus, manager)
	managerAdapter.Start()

	exporterAdapter := spsw.NewSpiderBusAdapterForExporter(spiderBus, exporter)
	exporterAdapter.Start()

	go exporter.Run()
	go manager.Run()

	for i := 0; i < 4; i++ {
		go func() {
			worker := spsw.NewWorker()
			adapter := spsw.NewSpiderBusAdapterForWorker(spiderBus, worker)
			adapter.Start()
			worker.Run()
		}()
	}

	// HACK!
	// Since at this point we don't have a way to track the task execution state we
	// try to detect the end of scraping job by checking if all SQLite tables are empty.
	// This is unreliable as one or more Tasks might still be in progress.
	time.Sleep(100 * time.Second)

	for {
		time.Sleep(10 * time.Second)
		if spiderBusBackend.IsEmpty() {
			log.Info("It appears scraping job is done!")
			break
		}

	}
}

func main() {
	initLogging()
	log.Info("Starting spiderswarm instance...")

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	singleNodeCmd := flag.NewFlagSet("singlenode", flag.ExitOnError)
	singleNodeWorkers := singleNodeCmd.Int("workers", 1, "Number of worker goroutines")

	switch os.Args[1] {
	case "singlenode":
		singleNodeCmd.Parse(os.Args[2:])
		log.Info(fmt.Sprintf("Number of worker goroutines: %d", *singleNodeWorkers))
	case "client":
		// TODO: client for REST API
		fmt.Println("client part not implemented yet")
	default:
		runTestWorkflow()
	}
}
