package spsw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFieldJoinActionFromTemplate(t *testing.T) {
	itemName := "testItem"

	actionTempl := &ActionTemplate{
		Name:       "testAction",
		StructName: "FieldJoinAction",
		ConstructorParams: map[string]Value{
			"inputNames": Value{
				ValueType:    ValueTypeStrings,
				StringsValue: []string{"title", "link"},
			},
			"itemName": Value{
				ValueType:   ValueTypeString,
				StringValue: itemName,
			},
		},
	}

	action, ok := NewFieldJoinActionFromTemplate(actionTempl).(*FieldJoinAction)
	assert.True(t, ok)

	assert.Equal(t, itemName, action.ItemName)
	assert.Equal(t, []string{"title", "link"}, action.AllowedInputNames)
}

func TestFieldJoinActionRun(t *testing.T) {
	jobUUID := "17C67CA0-35C6-488D-9C7B-F1AB4BAF5274"
	taskUUID := "D6887944-5ECA-44A4-87D5-C7E364E53271"
	itemName := "testItem"

	action := NewFieldJoinAction([]string{"Name", "Surname", "Phone", "Email"},
		jobUUID, taskUUID, itemName)

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
	mapOut := NewDataPipe()

	err = action.AddOutput(FieldJoinActionOutputItem, itemOut)
	assert.Nil(t, err)

	err = action.AddOutput(FieldJoinActionOutputMap, mapOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	expectedItemFields := map[string]*Value{
		"Name":    &Value{ValueType: ValueTypeString, StringValue: "John"},
		"Surname": &Value{ValueType: ValueTypeString, StringValue: "Smith"},
		"Phone":   &Value{ValueType: ValueTypeString, StringValue: "555-1212"},
		"Email":   &Value{ValueType: ValueTypeString, StringValue: "john@smith.int"},
	}

	item, ok := itemOut.Remove().(*Item)
	assert.True(t, ok)

	assert.Equal(t, jobUUID, item.JobUUID)
	assert.Equal(t, taskUUID, item.TaskUUID)
	assert.Equal(t, itemName, item.Name)

	assert.Equal(t, expectedItemFields, item.Fields)

	m, ok := mapOut.Remove().(map[string]string)
	assert.True(t, ok)

	for key, value := range expectedItemFields {
		assert.Equal(t, value.StringValue, m[key])
	}
}

func TestFieldJoinActionRunErr(t *testing.T) {
	jobUUID := "17C67CA0-35C6-488D-9C7B-F1AB4BAF5274"
	taskUUID := "D6887944-5ECA-44A4-87D5-C7E364E53271"
	itemName := "testItem"

	action := NewFieldJoinAction([]string{"Name", "Surname", "Phone", "Email"},
		jobUUID, taskUUID, itemName)

	err := action.Run()

	assert.NotNil(t, err)
	assert.Equal(t, errors.New("No output connected"), err)

	action.AddOutput(FieldJoinActionOutputItem, NewDataPipe())

	err = action.Run()

	assert.Equal(t, errors.New("No inputs connected"), err)
}
