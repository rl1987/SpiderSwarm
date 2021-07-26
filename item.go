package main

import (
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
	}
}

func (i *Item) SetField(name string, value interface{}) {
	i.Fields[name] = value
}
