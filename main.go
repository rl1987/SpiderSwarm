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
	spiderBusBackend := spsw.NewRedisSpiderBusBackend(backendAddr, "")
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
								StringsValue: []string{"'.ip_addr'"},
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

	fmt.Println(workflow)

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
