package spsw

import (
	"errors"
	"fmt"
	"sort"

	yaml "gopkg.in/yaml.v3"
)

type ActionTemplate struct {
	Name              string           `yaml:"Name"`
	StructName        string           `yaml:"StructName"`
	ConstructorParams map[string]Value `yaml:"ConstructorParams"`
}

func NewActionTemplate(name string, structName string, constructorParams map[string]interface{}) *ActionTemplate {
	constructorParamsV := map[string]Value{}

	if constructorParams != nil {
		for key, valIntf := range constructorParams {
			constructorParamsV[key] = *NewValue(valIntf)
		}
	}

	return &ActionTemplate{
		Name:              name,
		StructName:        structName,
		ConstructorParams: constructorParamsV,
	}
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

func NewTaskTemplate(taskName string, initial bool) *TaskTemplate {
	return &TaskTemplate{
		TaskName:          taskName,
		Initial:           initial,
		ActionTemplates:   []ActionTemplate{},
		DataPipeTemplates: []DataPipeTemplate{},
	}
}

func (tt TaskTemplate) String() string {
	return fmt.Sprintf("<TaskTemplate TaskName: %s, Initial: %v, ActionTemplates: %s, DataPipeTemplates: %v>",
		tt.TaskName, tt.Initial, &tt.ActionTemplates, tt.DataPipeTemplates)
}

func (tt *TaskTemplate) AddActionTemplate(actionTempl *ActionTemplate) error {
	for _, at := range tt.ActionTemplates {
		if at.Name == actionTempl.Name {
			return errors.New("Action template already exists with given name")
		}
	}

	tt.ActionTemplates = append(tt.ActionTemplates, *actionTempl)

	return nil
}

func (tt *TaskTemplate) RemoveActionTemplate(name string) error {
	idx := -1

	for i, at := range tt.ActionTemplates {
		if at.Name == name {
			idx = i
			break
		}
	}

	if idx == -1 {
		return errors.New("Action template not found with a given name")
	}

	// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
	tt.ActionTemplates = append(tt.ActionTemplates[:idx], tt.ActionTemplates[idx+1:]...)

	return nil
}

func (tt *TaskTemplate) AddDataPipeTemplate(dpt *DataPipeTemplate) error {
	// TODO: prevalidate this

	tt.DataPipeTemplates = append(tt.DataPipeTemplates, *dpt)
	return nil
}

func (tt *TaskTemplate) ConnectActionTemplates(sourceActionName string, sourceOutputName string, destActionName string, destInputName string) error {
	dpt := &DataPipeTemplate{
		SourceActionName: sourceActionName,
		SourceOutputName: sourceOutputName,
		DestActionName:   destActionName,
		DestInputName:    destInputName,
	}

	return tt.AddDataPipeTemplate(dpt)
}

func (tt *TaskTemplate) ConnectInputToActionTemplate(taskInputName string, destActionName string, destInputName string) error {
	dpt := &DataPipeTemplate{
		TaskInputName:  taskInputName,
		DestActionName: destActionName,
		DestInputName:  destInputName,
	}

	return tt.AddDataPipeTemplate(dpt)
}

func (tt *TaskTemplate) ConnectOutputToActionTemplate(sourceActionName string, sourceOutputName string, taskOutputName string) error {
	dpt := &DataPipeTemplate{
		SourceActionName: sourceActionName,
		SourceOutputName: sourceOutputName,
		TaskOutputName:   taskOutputName,
	}

	return tt.AddDataPipeTemplate(dpt)
}

func (tt *TaskTemplate) DisconnectActionTemplates(sourceActionName string, sourceOutputName string, destActionName string, destInputName string) error {
	idx := -1

	for i, dpt := range tt.DataPipeTemplates {
		if dpt.SourceActionName == sourceActionName && dpt.SourceOutputName == sourceOutputName &&
			dpt.DestActionName == destActionName && dpt.DestInputName == destInputName {
			idx = i
			break
		}
	}

	if idx == -1 {
		return errors.New("Not found")
	}

	tt.DataPipeTemplates = append(tt.DataPipeTemplates[:idx], tt.DataPipeTemplates[idx+1:]...)

	return nil
}

func (tt *TaskTemplate) DisconnectInput(taskInputName string) error {
	idx := -1

	for i, dpt := range tt.DataPipeTemplates {
		if dpt.TaskInputName == taskInputName {
			idx = i
			break
		}
	}

	if idx == -1 {
		return errors.New("Not found")
	}

	tt.DataPipeTemplates = append(tt.DataPipeTemplates[:idx], tt.DataPipeTemplates[idx+1:]...)

	return nil
}

// TODO: remove sourceOutputName and taskOutputName arguments as they are redundant.
func (tt *TaskTemplate) DisconnectOutput(sourceActionName string, sourceOutputName string, taskOutputName string) error {
	idx := -1

	for i, dpt := range tt.DataPipeTemplates {
		if dpt.SourceActionName == sourceActionName && dpt.SourceOutputName == sourceOutputName && dpt.TaskOutputName == taskOutputName {
			idx = i
			break
		}
	}

	if idx == -1 {
		return errors.New("Not found")
	}

	tt.DataPipeTemplates = append(tt.DataPipeTemplates[:idx], tt.DataPipeTemplates[idx+1:]...)

	return nil
}

type Workflow struct {
	Name          string         `yaml:"Name"`
	Version       string         `yaml:"Version"`
	TaskTemplates []TaskTemplate `yaml:"TaskTemplates"`
}

func NewWorkflow(name string, version string) *Workflow {
	return &Workflow{
		Name:          name,
		Version:       version,
		TaskTemplates: []TaskTemplate{},
	}
}

func (w *Workflow) AddTaskTemplate(taskTempl *TaskTemplate) {
	w.TaskTemplates = append(w.TaskTemplates, *taskTempl)
}

func (w *Workflow) SetInitial(taskName string) {
	idx := -1

	for i, tt := range w.TaskTemplates {
		if tt.TaskName == taskName {
			idx = i
		}
	}

	if idx == -1 {
		return
	}

	for i, _ := range w.TaskTemplates {
		w.TaskTemplates[i].Initial = false
	}

	w.TaskTemplates[idx].Initial = true
}

func (w *Workflow) RemoveTaskTemplate(taskName string) error {
	idx := -1

	for i, tt := range w.TaskTemplates {
		if tt.TaskName == taskName {
			idx = i
			break
		}
	}

	if idx == -1 {
		return errors.New("Not found")
	}

	w.TaskTemplates = append(w.TaskTemplates[:idx], w.TaskTemplates[idx+1:]...)

	return nil
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

// Based on: https://stackoverflow.com/a/33323321
func stringIsInSlice(needle string, haystack []string) bool {
	sort.Strings(haystack) // XXX: should we sort here or rely on haystack to be pre-sorted?
	i := sort.Search(len(haystack), func(i int) bool { return haystack[i] >= needle })
	if i < len(haystack) && haystack[i] == needle {
		return true
	}
	return false
}

func (w *Workflow) validateInputOutputNames() error {
	actionNameToStructName := map[string]string{}

	for _, tt := range w.TaskTemplates {
		for _, at := range tt.ActionTemplates {
			actionNameToStructName[at.Name] = at.StructName
		}
	}

	for _, tt := range w.TaskTemplates {
		for _, dpt := range tt.DataPipeTemplates {
			if dpt.SourceActionName != "" {
				structName := actionNameToStructName[dpt.SourceActionName]

				if !stringIsInSlice(dpt.SourceOutputName, AllowedOutputNameTable[structName]) {
					return fmt.Errorf("Output name %s is not allowed for %s", dpt.SourceOutputName,
						structName)
				}
			}

			if dpt.DestActionName != "" {
				structName := actionNameToStructName[dpt.DestActionName]

				if structName == "FieldJoinAction" || structName == "TaskPromiseAction" {
					continue
				}

				if !stringIsInSlice(dpt.DestInputName, AllowedInputNameTable[structName]) {
					return fmt.Errorf("Input name %s is not allowed for %s", dpt.DestInputName,
						structName)
				}
			}
		}
	}

	return nil
}

func (w *Workflow) GetInitialTaskTemplate() *TaskTemplate {
	var initialTaskTempl *TaskTemplate
	initialTaskTempl = nil

	for _, tt := range w.TaskTemplates {
		if tt.Initial {
			initialTaskTempl = &tt
			break
		}
	}

	return initialTaskTempl
}

func (w *Workflow) GetInitialTaskTemplateName() string {
	initialTaskTempl := w.GetInitialTaskTemplate()
	if initialTaskTempl != nil {
		return initialTaskTempl.TaskName
	}

	return ""
}

func (w *Workflow) Validate() (bool, error) {
	var err error

	err = w.validateInputOutputNames()
	if err != nil {
		return false, err
	}

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
