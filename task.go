package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

func (t *Task) Run() error {

	return errors.New("Not implemented")
}

type Task struct {
	Name         string
	UUID         string
	CreatedAt    time.Time
	WorkflowName string
	WorkflowUUID string

	Inputs    map[string]*DataPipe
	Outputs   map[string]*DataPipe
	Actions   []Action
	DataPipes []*DataPipe
}

func NewTask(name string, workflowName string, workflowUUID string) *Task {
	return &Task{
		Name:         name,
		UUID:         uuid.New().String(),
		CreatedAt:    time.Now(),
		WorkflowName: workflowName,
		WorkflowUUID: workflowUUID,

		Inputs:    map[string]*DataPipe{},
		Outputs:   map[string]*DataPipe{},
		Actions:   []Action{},
		DataPipes: []*DataPipe{},
	}
}

func (t *Task) indexActions() map[string]*Action {
	var index map[string]*Action

	for _, a := range t.Actions {
		actionUUID := a.GetUniqueID()
		index[actionUUID] = &a
	}

	return index
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
