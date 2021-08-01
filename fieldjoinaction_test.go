package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFieldJoinActionFromTemplate(t *testing.T) {
	workflowName := "testWorkflow"
	itemName := "testItem"

	actionTempl := &ActionTemplate{
		Name:       "testAction",
		StructName: "FieldJoinAction",
		ConstructorParams: map[string]interface{}{
			"inputNames": []string{"title", "link"},
			"itemName":   itemName,
		},
	}

	workflow := &Workflow{
		Name: workflowName,
	}

	action := NewFieldJoinActionFromTemplate(actionTempl, workflow)

	assert.Equal(t, workflowName, action.WorkflowName)
	assert.Equal(t, itemName, action.ItemName)
	assert.Equal(t, []string{"title", "link"}, action.AllowedInputNames)
}

func TestFieldJoinActionRun(t *testing.T) {
	workflowName := "testWorkflow"
	jobUUID := "17C67CA0-35C6-488D-9C7B-F1AB4BAF5274"
	taskUUID := "D6887944-5ECA-44A4-87D5-C7E364E53271"
	itemName := "testItem"

	action := NewFieldJoinAction([]string{"Name", "Surname", "Phone", "Email"},
		workflowName, jobUUID, taskUUID, itemName)

	nameIn := NewDataPipe()
	surnameIn := NewDataPipe()
	phoneIn := NewDataPipe()
	emailIn := NewDataPipe()

	nameIn.Add("John")
	surnameIn.Add("Smith")
	phoneIn.Add("555-1212")
	emailIn.Add("john@smith.int")

	err := action.AddInput("Name", nameIn)
	assert.Nil(t, err)
	err = action.AddInput("Surname", surnameIn)
	assert.Nil(t, err)
	err = action.AddInput("Phone", phoneIn)
	assert.Nil(t, err)
	err = action.AddInput("Email", emailIn)
	assert.Nil(t, err)

	itemOut := NewDataPipe()

	err = action.AddOutput(FieldJoinActionOutputItem, itemOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	expectedItemFields := map[string]interface{}{
		"Name":    "John",
		"Surname": "Smith",
		"Phone":   "555-1212",
		"Email":   "john@smith.int",
	}

	item, ok := itemOut.Remove().(*Item)
	assert.True(t, ok)

	assert.Equal(t, workflowName, item.WorkflowName)
	assert.Equal(t, jobUUID, item.JobUUID)
	assert.Equal(t, taskUUID, item.TaskUUID)
	assert.Equal(t, itemName, item.Name)

	assert.Equal(t, expectedItemFields, item.Fields)
}
