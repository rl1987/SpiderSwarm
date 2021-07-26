package main

type ActionTemplate struct {
	Name              string
	StructName        string
	ConstructorParams map[string]interface{}
}

type DataPipeTemplate struct {
	SourceActionName string
	DestActionName   string
}

type TaskTemplate struct {
	TaskName          string
	ActionTemplates   []ActionTemplate
	DataPipeTemplates []DataPipeTemplate
}

type Workflow struct {
	Name          string
	Version       string
	UUID          string
	TaskTemplates []TaskTemplate
}
