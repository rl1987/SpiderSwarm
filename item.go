package main

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	UUID         string
	WorkflowName string
	JobUUID      string
	TaskUUID     string
	CreatedAt    time.Time
	Name         string

	Fields map[string]interface{}
}

func NewItem(name string, workflowName string, jobUUID string, taskUUID string) *Item {
	return &Item{
		UUID:         uuid.New().String(),
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
		TaskUUID:     taskUUID,
		CreatedAt:    time.Now(),
		Fields:       map[string]interface{}{},
		Name:         name,
	}
}

func (i *Item) IsSplayable() bool {
	hasLists := false
	equalLen := true
	lastLen := -1

	for _, value := range i.Fields {
		rt := reflect.TypeOf(value) // XXX: this is bad for performance!

		if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
			hasLists = true

			if lastLen != -1 && lastLen != len(value.([]interface{})) {
				equalLen = false
				lastLen = len(value.([]interface{}))
			}
		}
	}

	return hasLists && equalLen
}

func (i *Item) SetField(name string, value interface{}) {
	i.Fields[name] = value
}
