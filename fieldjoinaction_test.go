package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldJoinActionRun(t *testing.T) {
	action := NewFieldJoinAction([]string{"Name", "Surname", "Phone", "Email"})

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

	expectedItem := map[string]string{
		"Name":    "John",
		"Surname": "Smith",
		"Phone":   "555-1212",
		"Email":   "john@smith.int",
	}

	item, ok := itemOut.Remove().(map[string]string)
	assert.True(t, ok)

	assert.Equal(t, expectedItem, item)
}
