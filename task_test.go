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
