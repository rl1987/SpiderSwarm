package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataChunkStr(t *testing.T) {
	testStr := "test"

	strChunk, err := NewDataChunk(testStr)
	assert.Nil(t, err)
	assert.NotNil(t, strChunk)
	assert.Equal(t, testStr, strChunk.Payload)
	assert.Equal(t, DataChunkTypeString, strChunk.Type)
}

func TestNewDataChunkMapStringToStrings(t *testing.T) {
	params := map[string][]string{
		"a": []string{"1"},
		"b": []string{"2"},
	}

	chunk, err := NewDataChunk(params)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, params, chunk.Payload)
	assert.Equal(t, DataChunkTypeMapStringToStrings, chunk.Type)
}

func TestNewDataChunkHTTPHeader(t *testing.T) {
	headers := http.Header{
		"User-Agent": []string{"spiderswarm"},
		"Accept":     []string{"text/html"},
	}

	chunk, err := NewDataChunk(headers)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, headers, chunk.Payload)
	assert.Equal(t, DataChunkHTTPHeader, chunk.Type)
}
