package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

type Task struct {
	Name         string
	UUID         string
	CreatedAt    time.Time
	WorkflowName string
	JobUUID      string

	Inputs    map[string]*DataPipe
	Outputs   map[string]*DataPipe
	Actions   []Action
	DataPipes []*DataPipe
}

func NewTask(name string, workflowName string, jobUUID string) *Task {
	return &Task{
		Name:         name,
		UUID:         uuid.New().String(),
		CreatedAt:    time.Now(),
		WorkflowName: workflowName,
		JobUUID:      jobUUID,

		Inputs:    map[string]*DataPipe{},
		Outputs:   map[string]*DataPipe{},
		Actions:   []Action{},
		DataPipes: []*DataPipe{},
	}
}

func (t *Task) addDataPipeFromTemplate(dataPipeTemplate *DataPipeTemplate, nameToAction map[string]Action) {
	var newDP *DataPipe
	newDP = nil

	if len(dataPipeTemplate.SourceActionName) > 0 && len(dataPipeTemplate.DestActionName) > 0 {
		fromAction := nameToAction[dataPipeTemplate.SourceActionName]
		toAction := nameToAction[dataPipeTemplate.DestActionName]

		if fromAction != nil && toAction != nil {
			newDP = NewDataPipeBetweenActions(fromAction, toAction)
			fromAction.AddOutput(dataPipeTemplate.SourceOutputName, newDP)
			toAction.AddInput(dataPipeTemplate.DestInputName, newDP)
		}
	} else if len(dataPipeTemplate.TaskInputName) > 0 {
		toAction := nameToAction[dataPipeTemplate.DestActionName]

		if toAction != nil {
			newDP = NewDataPipe()
			newDP.ToAction = toAction
			toAction.AddInput(dataPipeTemplate.DestInputName, newDP)
			t.Inputs[dataPipeTemplate.TaskInputName] = newDP
		}
	} else if len(dataPipeTemplate.TaskOutputName) > 0 {
		fromAction := nameToAction[dataPipeTemplate.SourceActionName]

		if fromAction != nil {
			newDP = NewDataPipe()
			newDP.FromAction = fromAction
			fromAction.AddOutput(dataPipeTemplate.SourceOutputName, newDP)
			t.Outputs[dataPipeTemplate.TaskOutputName] = newDP
		}
	}

	if newDP != nil {
		t.DataPipes = append(t.DataPipes, newDP)
	}
}

func NewTaskFromTemplate(taskTempl *TaskTemplate, workflow *Workflow, jobUUID string) *Task {
	task := NewTask(taskTempl.TaskName, workflow.Name, jobUUID)

	nameToAction := map[string]Action{}

	for _, actionTempl := range taskTempl.ActionTemplates {
		newAction := NewActionFromTemplate(&actionTempl, workflow, jobUUID)
		task.Actions = append(task.Actions, newAction)
		nameToAction[actionTempl.Name] = newAction
	}

	for _, dataPipeTemplate := range taskTempl.DataPipeTemplates {
		task.addDataPipeFromTemplate(&dataPipeTemplate, nameToAction)
	}

	return task
}

func NewTaskFromPromise(promise *TaskPromise, workflow *Workflow, jobUUID string) *Task {
	// TODO: implement
	return &Task{}
}

func (t *Task) AddInput(name string, action Action, actionInputName string, dataPipe *DataPipe) {
	t.Inputs[name] = dataPipe
	t.DataPipes = append(t.DataPipes, dataPipe)

	action.AddInput(actionInputName, dataPipe)

	dataPipe.ToAction = action
}

func (t *Task) AddOutput(name string, action Action, actionOutputName string, dataPipe *DataPipe) {
	t.Outputs[name] = dataPipe
	t.DataPipes = append(t.DataPipes, dataPipe)

	action.AddOutput(actionOutputName, dataPipe)

	dataPipe.FromAction = action
}

func (t *Task) AddAction(action Action) {
	t.Actions = append(t.Actions, action)
}

func (t *Task) AddDataPipeBetweenActions(fromAction Action, fromOutputName string, toAction Action, toInputName string) {
	// TODO: check if both actions are in Actions array and if Input/Output names
	// are allowed.

	dataPipe := NewDataPipeBetweenActions(fromAction, toAction)

	fromAction.AddOutput(fromOutputName, dataPipe)
	toAction.AddInput(toInputName, dataPipe)

	t.DataPipes = append(t.DataPipes, dataPipe)
}

// Based on: https://github.com/adonovan/gopl.io/blob/master/ch5/toposort/main.go
func (t *Task) sortActionsTopologically() []Action {
	order := []Action{}
	seen := make(map[string]bool)
	var visitAll func(items []Action)

	visitAll = func(actions []Action) {
		for _, action := range actions {
			if action != nil && !seen[action.GetUniqueID()] {
				seen[action.GetUniqueID()] = true
				precedingActions := action.GetPrecedingActions()
				visitAll(precedingActions)
				order = append(order, action)
			}
		}
	}

	lastActions := []Action{}

	for _, output := range t.Outputs {
		if output.FromAction != nil {
			lastActions = append(lastActions, output.FromAction)
		}
	}

	visitAll(lastActions)

	return order
}

func (t *Task) Run() error {
	order := t.sortActionsTopologically()
	fmt.Println("order")
	spew.Dump(order)

	for _, action := range order {
		fmt.Println("Running action:")
		spew.Dump(action)
		err := action.Run()
		if err != nil && !action.IsFailureAllowed() {
			return err
		}
	}

	for _, output := range t.Outputs {
		if len(output.Queue) >= 1 && output.Queue[0].Type == DataChunkTypeStrings {
			strings, ok := output.Remove().([]string)
			if ok {
				for _, s := range strings {
					output.Add(s)
				}
			}
		}
	}

	return nil
}
