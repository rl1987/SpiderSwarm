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

func initLogging() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func printUsage() {
	fmt.Println("Read the code for now")
}

func setupSpiderBus(backendAddr string) *spsw.SpiderBus {
	// TODO: make RedisSpiderBusBackend functional. Remove backends that are based on relational
	// databases.
	spiderBusBackend := spsw.NewRedisSpiderBusBackend(backendAddr, "")
	//spiderBusBackend := spsw.NewSQLiteSpiderBusBackend(backendAddr)
	spiderBus := spsw.NewSpiderBus()
	spiderBus.Backend = spiderBusBackend

	return spiderBus
}

func runManager(workflow *spsw.Workflow, backendAddr string) *spsw.Manager {
	manager := spsw.NewManager()

	if workflow != nil {
		manager.StartScrapingJob(workflow)
	}

	spiderBus := setupSpiderBus(backendAddr)

	managerAdapter := spsw.NewSpiderBusAdapterForManager(spiderBus, manager)
	managerAdapter.Start()

	if workflow != nil {
		log.Info(fmt.Sprintf("Starting Manager %v", manager))
		go manager.Run()
	}

	return manager
}

func runExporter(outputDirPath string, backendAddr string) *spsw.Exporter {
	exporter := spsw.NewExporter()

	exporterBackend := spsw.NewCSVExporterBackend(outputDirPath)

	exporter.AddBackend(exporterBackend)

	spiderBus := setupSpiderBus(backendAddr)

	exporterAdapter := spsw.NewSpiderBusAdapterForExporter(spiderBus, exporter)
	exporterAdapter.Start()

	log.Info(fmt.Sprintf("Starting Exporter %v", exporter))
	go exporter.Run()

	return exporter
}

func runWorkers(n int, backendAddr string) []*spsw.Worker {
	var workers []*spsw.Worker

	workers = []*spsw.Worker{}

	for i := 0; i < n; i++ {
		worker := spsw.NewWorker()
		workers = append(workers, worker)
	}

	for _, worker := range workers {
		go func(worker *spsw.Worker) {
			spiderBus := setupSpiderBus(backendAddr)

			adapter := spsw.NewSpiderBusAdapterForWorker(spiderBus, worker)
			adapter.Start()
			log.Info(fmt.Sprintf("Starting Worker %v", worker))
			worker.Run()
		}(worker)
	}

	return workers
}

func getWorkflow() *spsw.Workflow {
	return &spsw.Workflow{
		Name:    "testWorkflow",
		Version: "v0.0.0.0.1",
		TaskTemplates: []spsw.TaskTemplate{
			spsw.TaskTemplate{
				TaskName: "ScrapeListPage",
				Initial:  true,
				ActionTemplates: []spsw.ActionTemplate{
					spsw.ActionTemplate{
						Name:       "GetListPage",
						StructName: "HTTPAction",
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "http://104.248.27.41:8000/",
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
						Name:       "ExtractNextPageURL",
						StructName: "XPathAction",
						ConstructorParams: map[string]spsw.Value{
							"xpath": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "//li[@class=\"next\"]/a/@href",
							},
							"expectMany": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: false,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "MakeNextPageURLAbsolute",
						StructName: "URLJoinAction",
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "http://104.248.27.41:8000/",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "ExtractBookLinks",
						StructName: "XPathAction",
						ConstructorParams: map[string]spsw.Value{
							"xpath": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "//div[@class=\"image_container\"]/a/@href",
							},
							"expectMany": spsw.Value{
								ValueType: spsw.ValueTypeBool,
								BoolValue: true,
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "MakeBookURLAbsolute",
						StructName: "URLJoinAction",
						ConstructorParams: map[string]spsw.Value{
							"baseURL": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "http://104.248.27.41:8000/",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "PromiseToScrapeBookPage",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"url"},
							},
							"taskName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "ScrapeBookPage",
							},
						},
					},
					spsw.ActionTemplate{
						Name:       "PromiseToScrapeBookList",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]spsw.Value{
							"inputNames": spsw.Value{
								ValueType:    spsw.ValueTypeStrings,
								StringsValue: []string{"url"},
							},
							"taskName": spsw.Value{
								ValueType:   spsw.ValueTypeString,
								StringValue: "ScrapeListPage",
							},
						},
					},
				},
				DataPipeTemplates: []spsw.DataPipeTemplate{
					spsw.DataPipeTemplate{
						TaskInputName:  "url",
						DestActionName: "GetListPage",
						DestInputName:  spsw.HTTPActionInputBaseURL,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "GetListPage",
						SourceOutputName: spsw.HTTPActionOutputBody,
						DestActionName:   "ExtractNextPageURL",
						DestInputName:    spsw.XPathActionInputHTMLBytes,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "GetListPage",
						SourceOutputName: spsw.HTTPActionOutputBody,
						DestActionName:   "ExtractBookLinks",
						DestInputName:    spsw.XPathActionInputHTMLBytes,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "ExtractNextPageURL",
						SourceOutputName: spsw.XPathActionOutputStr,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "ExtractBookLinks",
						SourceOutputName: spsw.XPathActionOutputStr,
						DestActionName:   "MakeBookURLAbsolute",
						DestInputName:    spsw.URLJoinActionInputRelativeURL,
					},
					spsw.DataPipeTemplate{
						SourceActionName: "MakeBookURLAbsolute",
						SourceOutputName: spsw.URLJoinActionOutputAbsoluteURL,
						DestActionName:   "PromiseToScrapeBookPage",
						DestInputName:    "url",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "MakeNextPageURLAbsolute",
						SourceOutputName: spsw.URLJoinActionOutputAbsoluteURL,
						DestActionName:   "PromiseToScrapeBookList",
						DestInputName:    "url",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "PromiseToScrapeBookPage",
						SourceOutputName: spsw.TaskPromiseActionOutputPromise,
						TaskOutputName:   "bookPagePromises",
					},
					spsw.DataPipeTemplate{
						SourceActionName: "PromiseToScrapeBookList",
						SourceOutputName: spsw.TaskPromiseActionOutputPromise,
						TaskOutputName:   "bookListPromises",
					},
				},
			},
		},
	}

}

func runTestWorkflow() {
	backendAddr := "/tmp/spiderbus.db"
	workflow := getWorkflow()

	workers := runWorkers(4, backendAddr)

	manager := runManager(workflow, backendAddr)
	exporter := runExporter("/tmp", backendAddr)

	spew.Dump(exporter)
	spew.Dump(workers)

	manager.StartScrapingJob(workflow)

	go manager.Run()

	manager.StartScrapingJob(workflow)

	// https://medium.com/@ashishstiwari/dont-simply-run-forever-loop-for-1594464040b1
	for {
		select {}
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

	workflow := getWorkflow()

	singleNodeCmd := flag.NewFlagSet("singlenode", flag.ExitOnError)
	singleNodeWorkers := singleNodeCmd.Int("workers", 1, "Number of worker goroutines")

	switch os.Args[1] {
	case "singlenode":
		singleNodeCmd.Parse(os.Args[2:])
		log.Info(fmt.Sprintf("Number of worker goroutines: %d", *singleNodeWorkers))
	case "worker":
		n, _ := strconv.Atoi(os.Args[2])
		backendAddr := os.Args[3]
		runWorkers(n, backendAddr)
		for {
			select {}
		}
	case "manager":
		backendAddr := os.Args[2]
		runManager(workflow, backendAddr)
		for {
			select {}
		}
	case "exporter":
		outputDir := os.Args[2]
		backendAddr := os.Args[3]
		runExporter(outputDir, backendAddr)
		for {
			select {}
		}
	case "client":
		// TODO: client for REST API
		fmt.Println("client part not implemented yet")
	default:
		runTestWorkflow()
	}
}
