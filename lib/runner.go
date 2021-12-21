package spsw

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Runner struct {
	BackendAddr string
}

func NewRunner(backendAddr string) *Runner {
	return &Runner{BackendAddr: backendAddr}
}

func (r *Runner) initLogging() {
	// Based on: https://github.com/sirupsen/logrus/issues/63#issuecomment-476486166
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			repopath := fmt.Sprintf("%s/src/github.com/spiderswarm/spiderswarm/", os.Getenv("GOPATH"))
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	logLevelStr := os.Getenv("SPSW_LOGLEVEL")

	if logLevelStr != "" {
		logLevel, err := log.ParseLevel(logLevelStr)

		if err != nil {
			panic(err)
		}

		log.SetLevel(logLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func (r *Runner) setupSpiderBus() *SpiderBus {
	spiderBusBackend := NewRedisSpiderBusBackend(r.BackendAddr, "")
	spiderBus := NewSpiderBus()
	spiderBus.Backend = spiderBusBackend

	return spiderBus
}

func (r *Runner) RunManager(workflow *Workflow) *Manager {
	r.initLogging()

	deduplicator := NewDeduplicator(r.BackendAddr)

	manager := NewManager(deduplicator)

	if workflow != nil {
		manager.StartScrapingJob(workflow)
	}

	spiderBus := r.setupSpiderBus()

	managerAdapter := NewSpiderBusAdapterForManager(spiderBus, manager)
	managerAdapter.Start()

	if workflow != nil {
		log.Info(fmt.Sprintf("Starting Manager %v", manager))
		go manager.Run()
	}

	return manager
}

func (r *Runner) RunExporter(outputDirPath string) *Exporter {
	r.initLogging()

	exporter := NewExporter()

	exporterBackend := NewCSVExporterBackend(outputDirPath)

	exporter.AddBackend(exporterBackend)

	spiderBus := r.setupSpiderBus()

	exporterAdapter := NewSpiderBusAdapterForExporter(spiderBus, exporter)
	exporterAdapter.Start()

	log.Info(fmt.Sprintf("Starting Exporter %v", exporter))
	go exporter.Run()

	return exporter
}

func (r *Runner) RunWorkers(n int) []*Worker {
	r.initLogging()
	var workers []*Worker

	workers = []*Worker{}

	for i := 0; i < n; i++ {
		worker := NewWorker()
		workers = append(workers, worker)
	}

	for _, worker := range workers {
		go func(worker *Worker) {
			spiderBus := r.setupSpiderBus()

			adapter := NewSpiderBusAdapterForWorker(spiderBus, worker)
			adapter.Start()
			log.Info(fmt.Sprintf("Starting Worker %v", worker))
			worker.Run()
		}(worker)
	}

	return workers
}

func (r *Runner) RunSingleNode(nWorkers int, outputDirPath string, workflow *Workflow) {
	r.RunWorkers(nWorkers)
	r.RunExporter(outputDirPath)
	manager := r.RunManager(nil)

	manager.StartScrapingJob(workflow)
	manager.Run()
}
