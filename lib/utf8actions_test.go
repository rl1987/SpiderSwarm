package spsw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUTF8EncodeActionRun(t *testing.T) {
	str := "abc"

	dataPipeIn := NewDataPipe()
	dataPipeOut := NewDataPipe()

	dataPipeIn.Add(str)

	utf8EncodeAction := NewUTF8EncodeAction()

	utf8EncodeAction.AddInput(UTF8EncodeActionInputStr, dataPipeIn)
	utf8EncodeAction.AddOutput(UTF8EncodeActionOutputBytes, dataPipeOut)

	err := utf8EncodeAction.Run()
	assert.Nil(t, err)

	binData, ok := dataPipeOut.Remove().([]byte)
	assert.True(t, ok)

	assert.Equal(t, binData, []byte{0x61, 0x62, 0x63})
}

func TestUTF8DecodeActionMultipleOutputs(t *testing.T) {
	action := NewUTF8DecodeAction()

	input := NewDataPipe()
	output1 := NewDataPipe()
	output2 := NewDataPipe()

	err := action.AddInput(UTF8DecodeActionInputBytes, input)
	assert.Nil(t, err)

	err = action.AddOutput(UTF8DecodeActionOutputStr, output1)
	assert.Nil(t, err)

	err = action.AddOutput(UTF8DecodeActionOutputStr, output2)
	assert.Nil(t, err)

	b := []byte("123")

	input.Add(b)

	err = action.Run()
	assert.Nil(t, err)

	s1, ok1 := output1.Remove().(string)
	assert.True(t, ok1)
	assert.Equal(t, "123", s1)

	s2, ok2 := output2.Remove().(string)
	assert.True(t, ok2)
	assert.Equal(t, "123", s2)
}

func TestUTF8DecodeActionRunErrors(t *testing.T) {
	action := NewUTF8DecodeAction()

	err := action.Run()

	assert.Equal(t, errors.New("Input not connected"), err)

	action.AddInput(UTF8DecodeActionInputBytes, NewDataPipe())

	err = action.Run()

	assert.Equal(t, errors.New("Output not connected"), err)

	action.AddOutput(UTF8DecodeActionOutputStr, NewDataPipe())

	err = action.Run()

	assert.Equal(t, errors.New("Failed to get binary data"), err)
}
