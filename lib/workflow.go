package spsw

import (
	//"errors"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type ActionTemplate struct {
	Name              string           `yaml:"Name"`
	StructName        string           `yaml:"StructName"`
	ConstructorParams map[string]Value `yaml:"ConstructorParams"`
}

func (at ActionTemplate) String() string {
	return fmt.Sprintf("<ActionTemplate Name: %s, StructName: %s, ConstructorParams: %v>",
		at.Name, at.StructName, at.ConstructorParams)
}

type DataPipeTemplate struct {
	SourceActionName string `yaml:"SourceActionName"`
	SourceOutputName string `yaml:"SourceOutputName"`
	DestActionName   string `yaml:"DestActionName"`
	DestInputName    string `yaml:"DestInputName"`
	TaskInputName    string `yaml:"TaskInputName"`
	TaskOutputName   string `yaml:"TaskOutputName"`
}

func (dpt *DataPipeTemplate) String() string {
	return fmt.Sprintf("<DataPipeTemplate SourceActionName: %s, SourceOutputName: %s, DestActionName: %s, DestInputName: %s, TaskInputName: %s, TaskOutputName: %s>",
		dpt.SourceActionName, dpt.SourceOutputName, dpt.DestActionName, dpt.DestInputName, dpt.TaskInputName,
		dpt.TaskOutputName)
}

type TaskTemplate struct {
	TaskName          string             `yaml:"TaskName"`
	Initial           bool               `yaml:"Initial"`
	ActionTemplates   []ActionTemplate   `yaml:"ActionTemplates"`
	DataPipeTemplates []DataPipeTemplate `yaml:"DataPipeTemplates"`
}

func (tt TaskTemplate) String() string {
	return fmt.Sprintf("<TaskTemplate TaskName: %s, Initial: %v, ActionTemplates: %s, DataPipeTemplates: %v>",
		tt.TaskName, tt.Initial, &tt.ActionTemplates, tt.DataPipeTemplates)
}

type Workflow struct {
	Name          string         `yaml:"Name"`
	Version       string         `yaml:"Version"`
	TaskTemplates []TaskTemplate `yaml:"TaskTemplates"`
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

func (w *Workflow) validateActionStructNames() error {
	return nil
}

func (w *Workflow) Validate() (bool, error) {
	var err error

	err = w.validateActionStructNames()
	if err != nil {
		return false, err
	}

	return true, nil
}
