package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortActionsTopologically(t *testing.T) {
	task := NewTask("testTask", "", "")

	a1 := NewHTTPAction("https://cryptome.org/", "GET", false)
	a2 := NewXPathAction("//title/text()", true)
	a3 := NewXPathAction("//body", true)
	a4 := NewXPathAction("//h1/text()", true)

	task.Actions = []Action{a1, a2, a3, a4}

	in1 := NewDataPipe()
	in2 := NewDataPipe()

	err := a1.AddInput(HTTPActionInputURLParams, in1)
	assert.Nil(t, err)
	err = a1.AddInput(HTTPActionInputHeaders, in2)
	assert.Nil(t, err)

	in1.ToAction = a1
	in2.ToAction = a2

	out1 := NewDataPipe()
	out2 := NewDataPipe()

	out1.FromAction = a2
	out2.FromAction = a4

	err = a2.AddOutput(XPathActionInputHTMLBytes, out1)
	assert.NotNil(t, err)
	err = a4.AddOutput(XPathActionInputHTMLBytes, out2)
	assert.NotNil(t, err)

	task.Inputs["in1"] = in1
	task.DataPipes = append(task.DataPipes, in1)
	task.Inputs["in2"] = in2
	task.DataPipes = append(task.DataPipes, in2)

	task.Outputs["out1"] = out1
	task.DataPipes = append(task.DataPipes, out1)
	task.Outputs["out2"] = out2
	task.DataPipes = append(task.DataPipes, out2)

	dpA1ToA2 := NewDataPipeBetweenActions(a1, a2)
	err = a1.AddOutput(HTTPActionOutputBody, dpA1ToA2)
	assert.Nil(t, err)
	err = a2.AddInput(XPathActionInputHTMLBytes, dpA1ToA2)
	assert.Nil(t, err)
	task.DataPipes = append(task.DataPipes, dpA1ToA2)

	dpA1ToA3 := NewDataPipeBetweenActions(a1, a3)
	// HACK: this is invalid
	// TODO: make Action support multiple outputs for the same name
	a1.AddOutput(HTTPActionOutputHeaders, dpA1ToA3)
	a3.AddInput(XPathActionInputHTMLBytes, dpA1ToA3)
	task.DataPipes = append(task.DataPipes, dpA1ToA3)

	dpA3ToA4 := NewDataPipeBetweenActions(a3, a4)
	err = a3.AddOutput(XPathActionOutputStr, dpA3ToA4)
	assert.Nil(t, err)
	err = a4.AddInput(XPathActionInputHTMLStr, dpA3ToA4)
	assert.Nil(t, err)
	task.DataPipes = append(task.DataPipes, dpA3ToA4)

	actions := task.sortActionsTopologically()

	assert.NotNil(t, actions)
	assert.Equal(t, 4, len(actions))
	assert.Equal(t, a1, actions[0])
	assert.True(t, actions[3] == a4 || actions[3] == a2)
}

func TestTaskAddInput(t *testing.T) {
	task := NewTask("testTask", "", "")

	dataPipe := NewDataPipe()

	httpAction := NewHTTPAction("https://news.ycombinator.com/news", "GET", false)

	task.AddAction(httpAction)

	task.AddInput("headersIn", httpAction, HTTPActionInputHeaders, dataPipe)

	assert.Equal(t, dataPipe, task.Inputs["headersIn"])
	assert.Equal(t, dataPipe, httpAction.Inputs[HTTPActionInputHeaders])
	assert.Equal(t, dataPipe, task.DataPipes[0])
	assert.Equal(t, httpAction, dataPipe.ToAction)
}

func TestTaskAddOutput(t *testing.T) {
	task := NewTask("testTask", "", "")

	dataPipe := NewDataPipe()

	httpAction := NewHTTPAction("https://news.ycombinator.com/news", "GET", false)

	task.AddAction(httpAction)

	task.AddOutput("html", httpAction, HTTPActionOutputBody, dataPipe)

	assert.Equal(t, dataPipe, task.Outputs["html"])
	assert.Equal(t, dataPipe, httpAction.Outputs[HTTPActionOutputBody][0])
	assert.Equal(t, dataPipe, task.DataPipes[0])
	assert.Equal(t, httpAction, dataPipe.FromAction)
}

func TestTaskAddAction(t *testing.T) {
	task := NewTask("testTask", "", "")

	assert.Equal(t, 0, len(task.Actions))

	httpAction := NewHTTPAction("https://news.ycombinator.com/news", "GET", false)

	task.AddAction(httpAction)

	assert.Equal(t, 1, len(task.Actions))
	assert.Equal(t, httpAction, task.Actions[0])
}

func TestNewTaskFromTemplate(t *testing.T) {
	workflow := &Workflow{
		Name: "testWorkflow",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP1",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://news.ycombinator.com/",
							"method":  "GET",
							"canFail": false,
						},
					},
					ActionTemplate{
						Name:              "UTF8Decode",
						StructName:        "UTF8DecodeAction",
						ConstructorParams: map[string]interface{}{},
					},
					ActionTemplate{
						Name:       "MakePromise",
						StructName: "TaskPromiseAction",
						ConstructorParams: map[string]interface{}{
							"inputNames": []string{"htmlStr1", "htmlStr2"},
							"taskName":   "ParseHTML",
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "cookies",
						DestActionName: "HTTP1",
						DestInputName:  HTTPActionInputCookies,
					},
					DataPipeTemplate{
						SourceActionName: "HTTP1",
						SourceOutputName: HTTPActionOutputBody,
						DestActionName:   "UTF8Decode",
						DestInputName:    UTF8DecodeActionInputBytes,
					},
					DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr1",
					},
					DataPipeTemplate{
						SourceActionName: "UTF8Decode",
						SourceOutputName: UTF8DecodeActionOutputStr,
						DestActionName:   "MakePromise",
						DestInputName:    "htmlStr2",
					},
					DataPipeTemplate{
						SourceActionName: "MakePromise",
						SourceOutputName: TaskPromiseActionOutputPromise,
						TaskOutputName:   "promise",
					},
				},
			},
		},
	}
	jobUUID := "44ECE4B0-A1C9-4DE2-A456-7862F2A5B6CA"

	task := NewTaskFromTemplate(&workflow.TaskTemplates[0], workflow, jobUUID)

	assert.NotNil(t, task)
	assert.Equal(t, workflow.Name, task.WorkflowName)
	assert.Equal(t, jobUUID, task.JobUUID)
	assert.Equal(t, len(workflow.TaskTemplates[0].ActionTemplates), len(task.Actions))
	assert.Equal(t, len(workflow.TaskTemplates[0].DataPipeTemplates), len(task.DataPipes))

	httpAction, ok1 := task.Actions[0].(*HTTPAction)
	assert.True(t, ok1)
	assert.NotNil(t, httpAction)

	assert.Equal(t, "https://news.ycombinator.com/", httpAction.BaseURL)
	assert.Equal(t, "GET", httpAction.Method)
	assert.Equal(t, false, httpAction.CanFail)

	utf8DecodeAction, ok2 := task.Actions[1].(*UTF8DecodeAction)
	assert.True(t, ok2)
	assert.NotNil(t, utf8DecodeAction)

	promiseAction, ok3 := task.Actions[2].(*TaskPromiseAction)
	assert.True(t, ok3)
	assert.NotNil(t, promiseAction)

	assert.Equal(t, []string{"htmlStr1", "htmlStr2"}, promiseAction.AllowedInputNames)
	assert.Equal(t, "ParseHTML", promiseAction.TaskName)

	assert.Equal(t, 1, len(task.Inputs))
	assert.Equal(t, 1, len(task.Outputs))

	assert.Equal(t, 1, len(httpAction.Inputs))
	assert.Equal(t, 1, len(httpAction.Outputs))

	dataPipe0 := httpAction.Inputs[HTTPActionInputCookies]

	assert.Equal(t, httpAction, dataPipe0.ToAction)
	assert.Nil(t, dataPipe0.FromAction)
	assert.Equal(t, task.Inputs["cookies"], dataPipe0)

	dataPipe1 := httpAction.Outputs[HTTPActionOutputBody][0]

	assert.Equal(t, httpAction, dataPipe1.FromAction)
	assert.Equal(t, utf8DecodeAction, dataPipe1.ToAction)

	assert.Equal(t, 1, len(utf8DecodeAction.Inputs))
	assert.Equal(t, 1, len(utf8DecodeAction.Outputs))

	dataPipe2 := utf8DecodeAction.Outputs[UTF8DecodeActionOutputStr][0]

	assert.Equal(t, utf8DecodeAction, dataPipe2.FromAction)
	assert.Equal(t, promiseAction, dataPipe2.ToAction)

	dataPipe3 := utf8DecodeAction.Outputs[UTF8DecodeActionOutputStr][1]

	assert.Equal(t, utf8DecodeAction, dataPipe3.FromAction)
	assert.Equal(t, promiseAction, dataPipe3.ToAction)

	dataPipe4 := promiseAction.Outputs[TaskPromiseActionOutputPromise][0]

	assert.Equal(t, promiseAction, dataPipe4.FromAction)
	assert.Nil(t, dataPipe4.ToAction)

	assert.Equal(t, dataPipe4, task.Outputs["promise"])
}

func TestAddDataPipeBetweenActions(t *testing.T) {
	task := &Task{}

	httpAction := NewHTTPAction("https://www.example.org", "HEAD", false)
	xpathAction := NewXPathAction("//title/text()", false)

	task.AddAction(httpAction)
	task.AddAction(xpathAction)

	task.AddDataPipeBetweenActions(httpAction, HTTPActionOutputBody,
		xpathAction, XPathActionInputHTMLBytes)

	assert.Equal(t, 1, len(task.DataPipes))

	assert.Equal(t, httpAction, task.DataPipes[0].FromAction)
	assert.Equal(t, xpathAction, task.DataPipes[0].ToAction)
}

func TestNewTaskFromPromise(t *testing.T) {
	workflow := &Workflow{
		Name: "testWorkflow",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "GetHTML",
				Initial:  true,
				ActionTemplates: []ActionTemplate{
					ActionTemplate{
						Name:       "HTTP",
						StructName: "HTTPAction",
						ConstructorParams: map[string]interface{}{
							"baseURL": "https://news.ycombinator.com/",
							"method":  "GET",
							"canFail": false,
						},
					},
				},
				DataPipeTemplates: []DataPipeTemplate{
					DataPipeTemplate{
						TaskInputName:  "cookies",
						DestActionName: "HTTP",
						DestInputName:  HTTPActionInputCookies,
					},
					DataPipeTemplate{
						TaskOutputName:   "body",
						SourceActionName: "HTTP",
						SourceOutputName: HTTPActionOutputBody,
					},
				},
			},
		},
	}

	jobUUID := "45FF108C-CB87-40C4-A759-31577CC9567A"

	chunk := &DataChunk{
		Type: DataChunkTypeMapStringToString,
		Payload: map[string]string{
			"session": "S1234",
		},
	}

	promise := &TaskPromise{
		TaskName:     "GetHTML",
		WorkflowName: workflow.Name,
		JobUUID:      jobUUID,
		InputDataChunksByInputName: map[string]*DataChunk{
			"cookies": chunk,
		},
	}

	task := NewTaskFromPromise(promise, workflow)

	assert.NotNil(t, task)

	assert.Equal(t, 1, len(task.Actions))
	assert.Equal(t, 2, len(task.DataPipes))

	httpAction, ok := task.Actions[0].(*HTTPAction)

	assert.True(t, ok)
	assert.NotNil(t, httpAction)

	assert.Equal(t, chunk, task.Inputs["cookies"].Queue[0])
}
