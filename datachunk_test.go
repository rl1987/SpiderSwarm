package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataChunkStr(t *testing.T) {
	testStr := "test"

	strChunk, err := NewDataChunk(testStr)
	assert.Nil(t, err)
	assert.NotNil(t, strChunk)
	assert.Equal(t, testStr, strChunk.Payload)

}

func TestNewDataChunkMapStringToStrings(t *testing.T) {
	headers := map[string][]string{
		"User-Agent": []string{"spiderswarm"},
		"Accept":     []string{"text/html"},
	}

	chunk, err := NewDataChunk(headers)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, headers, chunk.Payload)

}
