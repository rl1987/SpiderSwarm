package spsw

import (
	"errors"
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
	for _, tt := range w.TaskTemplates {
		for _, actionTempl := range tt.ActionTemplates {
			structName := actionTempl.StructName
			if ActionConstructorTable[structName] == nil {
				return fmt.Errorf("No entry found in ActionConstructorTable for struct name %s", structName)
			}
		}
	}

	return nil
}

// TODO: unit-testing
func (w *Workflow) validateActionConnectedness() error {
	// XXX: We're instantiating Task because we don't know upfront what allowed inputs/outputs will be for each action
	// Perhaps there'a better way. We could make global tables for allowed input/output names.
	for _, tt := range w.TaskTemplates {
		task := NewTaskFromTemplate(&tt, "", "")

		sortedActions := task.sortActionsTopologically()

		if len(task.Actions) != len(sortedActions) {
			return fmt.Errorf("Task %s seems to be not fully connected - count of actions not matching after topological sorting", task.Name)
		}
	}

	return nil
}

func (w *Workflow) validateDataPipeConnectedness() error {
	for _, tt := range w.TaskTemplates {
		task := NewTaskFromTemplate(&tt, "", "")

		for _, dp := range task.DataPipes {
			hasFromAction := (dp.FromAction != nil)
			hasToAction := (dp.ToAction != nil)

			if hasFromAction && hasToAction {
				continue
			}

			// We don't allow short-circuiting input and output.
			if !hasFromAction && !hasToAction {
				return errors.New("Found disconnected data pipe")
			}

			if !hasFromAction {
				isTaskInput := false

				for _, inputs := range task.Inputs {
					for _, inDP := range inputs {
						if dp == inDP {
							isTaskInput = true
							break
						}
					}

					if isTaskInput {
						break
					}
				}

				if !isTaskInput {
					return fmt.Errorf("DataPipe to action %s is disconnected", dp.ToAction.GetName())
				} else {
					continue
				}
			}

			if !hasToAction {
				isTaskOutput := false

				for _, outDP := range task.Outputs {
					if dp == outDP {
						isTaskOutput = true
						break
					}
				}

				if !isTaskOutput {
					return fmt.Errorf("DataPipe from action %s is disconnected", dp.FromAction.GetName())
				} else {
					continue
				}
			}
		}
	}

	return nil
}

func (w *Workflow) Validate() (bool, error) {
	var err error

	err = w.validateActionStructNames()
	if err != nil {
		return false, err
	}

	err = w.validateActionConnectedness()
	if err != nil {
		return false, err
	}

	err = w.validateDataPipeConnectedness()
	if err != nil {
		return false, err
	}

	return true, nil
}
