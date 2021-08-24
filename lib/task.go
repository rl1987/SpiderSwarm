package spsw

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Task struct {
	Name         string
	UUID         string
	CreatedAt    time.Time
	WorkflowName string
	JobUUID      string

	Inputs    map[string][]*DataPipe
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

		Inputs:    map[string][]*DataPipe{},
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

			if t.Inputs[dataPipeTemplate.TaskInputName] == nil {
				t.Inputs[dataPipeTemplate.TaskInputName] = []*DataPipe{newDP}
			} else {
				t.Inputs[dataPipeTemplate.TaskInputName] = append(t.Inputs[dataPipeTemplate.TaskInputName], newDP)
			}
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

func NewTaskFromTemplate(taskTempl *TaskTemplate, workflowName string, jobUUID string) *Task {
	task := NewTask(taskTempl.TaskName, workflowName, jobUUID)

	nameToAction := map[string]Action{}

	for _, actionTempl := range taskTempl.ActionTemplates {
		newAction := NewActionFromTemplate(&actionTempl, workflowName, jobUUID)
		task.Actions = append(task.Actions, newAction)
		nameToAction[actionTempl.Name] = newAction
	}

	for _, dataPipeTemplate := range taskTempl.DataPipeTemplates {
		task.addDataPipeFromTemplate(&dataPipeTemplate, nameToAction)
	}

	return task
}

func (t *Task) populateTaskInputsFromPromise(promise *TaskPromise) {
	for inputName, chunk := range promise.InputDataChunksByInputName {
		inputs := t.Inputs[inputName]
		if inputs == nil {
			continue
		}

		for _, inDP := range inputs {
			inDP.Queue = append(inDP.Queue, chunk)
		}
	}
}

func NewTaskFromPromise(promise *TaskPromise, workflow *Workflow) *Task {
	taskTempl := workflow.FindTaskTemplate(promise.TaskName)

	if taskTempl == nil {
		return nil
	}

	task := NewTaskFromTemplate(taskTempl, workflow.Name, promise.JobUUID)

	task.JobUUID = promise.JobUUID

	return task
}

func NewTaskFromScheduledTask(scheduledTask *ScheduledTask) *Task {
	task := NewTaskFromTemplate(&scheduledTask.Template, scheduledTask.WorkflowName, scheduledTask.JobUUID)

	task.populateTaskInputsFromPromise(&scheduledTask.Promise)

	return task
}

func (t *Task) AddInput(name string, action Action, actionInputName string, dataPipe *DataPipe) {
	if t.Inputs[name] == nil {
		t.Inputs[name] = []*DataPipe{dataPipe}
	} else {
		t.Inputs[name] = append(t.Inputs[name], dataPipe)
	}

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

	for _, action := range order {
		log.Info(fmt.Sprintf("Running action: %v", action))
		err := action.Run()
		if err != nil && !action.IsFailureAllowed() {
			log.Error(fmt.Sprintf("Action failed with error: %v", err))
			return err
		}
	}

	for _, outDP := range t.Outputs {
		for _, chunk := range outDP.Queue {
			if item, okItem := chunk.Payload.(*Item); okItem {
				item.JobUUID = t.JobUUID
				item.TaskUUID = t.UUID
				item.WorkflowName = t.WorkflowName
			}

			if promise, okPromise := chunk.Payload.(*TaskPromise); okPromise {
				promise.JobUUID = t.JobUUID
				promise.WorkflowName = t.WorkflowName
			}
		}
	}

	return nil
}
