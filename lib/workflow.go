package spiderswarm

import (
	"fmt"
	"time"

	//"github.com/davecgh/go-spew/spew"
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
	jobUUID := uuid.New().String()
	startedAt := time.Now()

	var items []*Item

	log.Info(fmt.Sprintf("Job %s started from workflow %s:%s at %v", jobUUID, w.Name, w.Version,
		startedAt))

	var scheduledTask *ScheduledTask
	var scheduledTasks []*ScheduledTask
	var gotDataChunk *DataChunk

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

		select {
		case dc := <-worker1.DataChunksOut:
			gotDataChunk = dc
		case dc := <-worker2.DataChunksOut:
			gotDataChunk = dc
		case dc := <-worker3.DataChunksOut:
			gotDataChunk = dc
		case dc := <-worker4.DataChunksOut:
			gotDataChunk = dc
		}

		if gotDataChunk.Type == DataChunkTypeItem {
			item, _ := gotDataChunk.Payload.(*Item)

			for _, i := range item.Splay() {
				log.Info(fmt.Sprintf("Got item %v", i))
				items = append(items, i)
			}
		} else if gotDataChunk.Type == DataChunkTypePromise {
			promise, _ := gotDataChunk.Payload.(*TaskPromise)

			for _, p := range promise.Splay() {
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

	return items, nil
}
