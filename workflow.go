package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
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

func (w *Workflow) Run() ([]*Item, error) {
	jobUUID := uuid.New().String()
	startedAt := time.Now()

	var items []*Item

	fmt.Printf("Job %s started from workflow %s:%s at %v\n", jobUUID, w.Name, w.Version,
		startedAt)

	var tasks []*Task

	for _, taskTempl := range w.TaskTemplates {
		if !taskTempl.Initial {
			continue
		}

		newTask := NewTaskFromTemplate(&taskTempl, w.Name, jobUUID) // TODO: implement this one

		tasks = append(tasks, newTask)
	}

	for {
		if len(tasks) == 0 {
			break
		}

		task, tasks := tasks[0], tasks[1:]

		err := task.Run()
		if err != nil {
			spew.Dump(task)
			spew.Dump(err)
		} else {
			for _, outDP := range task.Outputs {
				for {
					if len(outDP.Queue) == 0 {
						break
					}

					x := outDP.Remove()

					if item, okItem := x.(*Item); okItem {
						items = append(items, item)
					}

					if promise, okPromise := x.(*TaskPromise); okPromise {
						newTask := NewTaskFromPromise(promise, w.Name, jobUUID)
						tasks = append(tasks, newTask)
					}
				}
			}
		}
	}

	return items, nil
}
