package spsw

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataChunkStr(t *testing.T) {
	testStr := "test"

	strChunk, err := NewDataChunk(testStr)
	assert.Nil(t, err)
	assert.NotNil(t, strChunk)
	assert.Equal(t, testStr, strChunk.PayloadValue.StringValue)
	assert.Equal(t, DataChunkTypeValue, strChunk.Type)
}

func TestNewDataChunkMapStringToStrings(t *testing.T) {
	params := map[string][]string{
		"a": []string{"1"},
		"b": []string{"2"},
	}

	chunk, err := NewDataChunk(params)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, params, chunk.PayloadValue.MapStringToStringsValue)
	assert.Equal(t, DataChunkTypeValue, chunk.Type)
}

func TestNewDataChunkHTTPHeader(t *testing.T) {
	headers := http.Header{
		"User-Agent": []string{"spiderswarm"},
		"Accept":     []string{"text/html"},
	}

	chunk, err := NewDataChunk(headers)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, headers, chunk.PayloadValue.HTTPHeadersValue)
	assert.Equal(t, DataChunkTypeValue, chunk.Type)
}

func TestNewDataChunkInt(t *testing.T) {
	i := 123

	chunk, err := NewDataChunk(i)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, i, chunk.PayloadValue.IntValue)
	assert.Equal(t, DataChunkTypeValue, chunk.Type)
}

func TestNewDataChunkItem(t *testing.T) {
	item := &Item{}

	chunk, err := NewDataChunk(item)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, item, chunk.PayloadItem)
	assert.Equal(t, DataChunkTypeItem, chunk.Type)
}

func TestNewDataChunkTaskPromise(t *testing.T) {
	promise := &TaskPromise{}

	chunk, err := NewDataChunk(promise)
	assert.Nil(t, err)
	assert.NotNil(t, chunk)
	assert.Equal(t, promise, chunk.PayloadPromise)
	assert.Equal(t, DataChunkTypePromise, chunk.Type)
}

func TestNewDataChunkFail(t *testing.T) {
	err := errors.New("Unsupported payload type")

	chunk, gotErr := NewDataChunk(err)
	assert.Nil(t, chunk)
	assert.NotNil(t, err)
	assert.Equal(t, err, gotErr)
}
