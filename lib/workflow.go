package spsw

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type ActionTemplate struct {
	Name              string
	StructName        string
	ConstructorParams map[string]Value
}

func (at ActionTemplate) String() string {
	return fmt.Sprintf("<ActionTemplate Name: %s, StructName: %s, ConstructorParams: %v>",
		at.Name, at.StructName, at.ConstructorParams)
}

type DataPipeTemplate struct {
	SourceActionName string
	SourceOutputName string
	DestActionName   string
	DestInputName    string
	TaskInputName    string
	TaskOutputName   string
}

func (dpt *DataPipeTemplate) String() string {
	return fmt.Sprintf("<DataPipeTemplate SourceActionName: %s, SourceOutputName: %s, DestActionName: %s, DestInputName: %s, TaskInputName: %s, TaskOutputName: %s>",
		dpt.SourceActionName, dpt.SourceOutputName, dpt.DestActionName, dpt.DestInputName, dpt.TaskInputName,
		dpt.TaskOutputName)
}

type TaskTemplate struct {
	TaskName          string
	Initial           bool
	ActionTemplates   []ActionTemplate
	DataPipeTemplates []DataPipeTemplate
}

func (tt TaskTemplate) String() string {
	return fmt.Sprintf("<TaskTemplate TaskName: %s, Initial: %v, ActionTemplates: %s, DataPipeTemplates: %v>",
		tt.TaskName, tt.Initial, &tt.ActionTemplates, tt.DataPipeTemplates)
}

type Workflow struct {
	Name          string
	Version       string
	TaskTemplates []TaskTemplate
}

func (w *Workflow) String() string {
	return fmt.Sprintf("<Workflow Name: %s, Version: %s, TaskTemplates: %v>", w.Name, w.Version, &w.TaskTemplates)
}

func NewWorkflowFromYAML(yamlStr string) *Workflow {
	yamlBytes := []byte(yamlStr)

	workflow := &Workflow{}

	err := yaml.Unmarshal(yamlBytes, workflow)

	if err != nil {
		panic(err)
	}

	return workflow
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

func (w *Workflow) ToYAML() string {
	yamlBytes, err := yaml.Marshal(w)

	if err != nil {
		panic(err)
	}

	return string(yamlBytes)
}
