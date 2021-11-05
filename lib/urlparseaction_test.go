package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLParseActionRun(t *testing.T) {
	action := NewURLParseAction()

	urlStr := "http://localhost/dev?a=1&b=2"

	expectParams := map[string][]string{
		"a": []string{"1"},
		"b": []string{"2"},
	}

	urlStrIn := NewDataPipe()
	urlStrIn.Add(urlStr)

	err := action.AddInput(URLParseActionInputURL, urlStrIn)
	assert.Nil(t, err)

	schemeOut := NewDataPipe()
	err = action.AddOutput(URLParseActionOutputScheme, schemeOut)
	assert.Nil(t, err)

	hostOut := NewDataPipe()
	err = action.AddOutput(URLParseActionOutputHost, hostOut)
	assert.Nil(t, err)

	pathOut := NewDataPipe()
	err = action.AddOutput(URLParseActionOutputPath, pathOut)
	assert.Nil(t, err)

	paramsOut := NewDataPipe()
	err = action.AddOutput(URLParseActionOutputParams, paramsOut)
	assert.Nil(t, err)

	err = action.Run()
	assert.Nil(t, err)

	gotScheme, ok1 := schemeOut.Remove().(string)
	assert.True(t, ok1)
	assert.Equal(t, "http", gotScheme)

	gotHost, ok2 := hostOut.Remove().(string)
	assert.True(t, ok2)
	assert.Equal(t, "localhost", gotHost)

	gotPath, ok3 := pathOut.Remove().(string)
	assert.True(t, ok3)
	assert.Equal(t, "/dev", gotPath)

	gotParams, ok4 := paramsOut.Remove().(map[string][]string)
	assert.True(t, ok4)
	assert.Equal(t, expectParams, gotParams)
}
