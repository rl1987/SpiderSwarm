package spsw

import (
	"fmt"
	"time"

	//"github.com/davecgh/go-spew/spew"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type ActionTemplate struct {
	Name              string
	StructName        string
	ConstructorParams map[string]interface{} // XXX: should this be defined more strictly?
}

type DataPipeTemplate struct {
	SourceActionName string
	SourceOutputName string
	DestActionName   string
	DestInputName    string
	TaskInputName    string
	TaskOutputName   string
}

type TaskTemplate struct {
	TaskName          string
	Initial           bool
	ActionTemplates   []ActionTemplate
	DataPipeTemplates []DataPipeTemplate
}

type Workflow struct {
	Name          string
	Version       string
	TaskTemplates []TaskTemplate
}

func (w *Workflow) FindTaskTemplate(taskName string) *TaskTemplate {
	var taskTempl *TaskTemplate
	taskTempl = nil

	for _, tt := range w.TaskTemplates {
		if tt.TaskName == taskName {
			taskTempl = &tt
			break
		}
	}

	return taskTempl
}

func (w *Workflow) createScheduledTaskFromPromise(promise *TaskPromise, jobUUID string) *ScheduledTask {
	taskTempl := w.FindTaskTemplate(promise.TaskName)
	if taskTempl == nil {
		return nil
	}

	scheduledTask := NewScheduledTask(promise, taskTempl, w.Name, w.Version, jobUUID)

	// TODO: log this

	return scheduledTask
}

func (w *Workflow) Run() ([]*Item, error) {
	spiderBusBackend := NewSQLiteSpiderBusBackend("")

	spiderBus := &SpiderBus{}
	spiderBus.Backend = spiderBusBackend

	jobUUID := uuid.New().String()
	startedAt := time.Now()

	log.Info(fmt.Sprintf("Job %s started from workflow %s:%s at %v", jobUUID, w.Name, w.Version,
		startedAt))

	var scheduledTask *ScheduledTask
	var scheduledTasks []*ScheduledTask
	var gotItem *Item
	var gotPromise *TaskPromise

	for _, taskTempl := range w.TaskTemplates {
		if !taskTempl.Initial {
			continue
		}

		newPromise := NewTaskPromise(taskTempl.TaskName, w.Name, jobUUID, map[string]*DataChunk{})
		log.Info(fmt.Sprintf("Enqueing promise %v", newPromise))

		scheduledTask = NewScheduledTask(newPromise, &taskTempl, w.Name, w.Version, jobUUID)
		scheduledTasks = append(scheduledTasks, scheduledTask)
	}

	scheduledTasksIn := make(chan *ScheduledTask)

	worker1 := NewWorker()
	worker2 := NewWorker()
	worker3 := NewWorker()
	worker4 := NewWorker()

	spiderBusAdapter1 := NewSpiderBusAdapterForWorker(spiderBus, worker1)
	spiderBusAdapter2 := NewSpiderBusAdapterForWorker(spiderBus, worker2)
	spiderBusAdapter3 := NewSpiderBusAdapterForWorker(spiderBus, worker3)
	spiderBusAdapter4 := NewSpiderBusAdapterForWorker(spiderBus, worker4)

	exporter := NewExporter()
	// TODO: make ExporterBackend API more abstract to enable plugin architecture.
	exporterBackend := NewCSVExporterBackend("/tmp")
	// FIXME: refrain from hardcoding field names; consider finding them from
	// Workflow.

	spiderBusAdapter5 := NewSpiderBusAdapterForExporter(spiderBus, exporter)

	err := exporterBackend.StartExporting(jobUUID, []string{"filer_id", "legal_name", "dba", "phone"})
	if err != nil {
		spew.Dump(err)
		return nil, err
	}

	exporter.AddBackend(exporterBackend)

	go exporter.Run()

	adapters := []*SpiderBusAdapter{spiderBusAdapter1, spiderBusAdapter2, spiderBusAdapter3,
		spiderBusAdapter4, spiderBusAdapter5}

	for _, adapter := range adapters {
		adapter.Start()
	}

	workers := []*Worker{worker1, worker2, worker3, worker4}

	for _, worker := range workers {
		worker.ScheduledTasksIn = scheduledTasksIn
		go worker.Run()
	}

	time.Sleep(1)

	for {
		if len(scheduledTasks) == 0 {
			break
		}

		scheduledTask, scheduledTasks = scheduledTasks[0], scheduledTasks[1:]

		scheduledTasksIn <- scheduledTask

		gotItem = nil
		gotPromise = nil

		select {
		case i := <-worker1.ItemsOut:
			gotItem = i
		case i := <-worker2.ItemsOut:
			gotItem = i
		case i := <-worker3.ItemsOut:
			gotItem = i
		case i := <-worker4.ItemsOut:
			gotItem = i

		case tp := <-worker1.TaskPromisesOut:
			gotPromise = tp
		case tp := <-worker2.TaskPromisesOut:
			gotPromise = tp
		case tp := <-worker3.TaskPromisesOut:
			gotPromise = tp
		case tp := <-worker4.TaskPromisesOut:
			gotPromise = tp
		}

		if gotItem != nil {
			for _, i := range gotItem.Splay() {
				log.Info(fmt.Sprintf("Got item %v", i))
				exporter.ItemsIn <- i
			}
		} else if gotPromise != nil {
			for _, p := range gotPromise.Splay() {
				newScheduledTask := w.createScheduledTaskFromPromise(p, jobUUID)
				if newScheduledTask == nil {
					continue
				}

				scheduledTasks = append(scheduledTasks, newScheduledTask)
			}
		}
	}

	worker1.Done <- nil
	worker2.Done <- nil
	worker3.Done <- nil
	worker4.Done <- nil

	exporterBackend.FinishExporting(jobUUID)

	return []*Item{}, nil
}
