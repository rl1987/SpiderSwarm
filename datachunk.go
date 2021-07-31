package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

const DataChunkTypeString = "DataChunkTypeString"
const DataChunkTypeStrings = "DataChunkTypeStrings"
const DataChunkTypeMapStringToString = "DataChunkTypeMapStringToString"
const DataChunkTypeMapStringToStrings = "DataChunkTypeMapStringToStrings"
const DataChunkTypeBytes = "DataChunkTypeBytes"
const DataChunkTypeHTTPHeader = "DataChunkTypeHTTPHeader"
const DataChunkTypeInt = "DataChunkTypeInt"
const DataChunkTypeItem = "DataChunkTypeItem"
const DataChunkTypePromise = "DataChunkTypePromise"

type DataChunk struct {
	Type    string
	Payload interface{}
	UUID    string
}

func NewDataChunkWithType(t string, payload interface{}) *DataChunk {
	return &DataChunk{
		Type:    t,
		Payload: payload,
		UUID:    uuid.New().String(),
	}
}

func NewDataChunk(payload interface{}) (*DataChunk, error) {
	if _, okStr := payload.(string); okStr {
		return NewDataChunkWithType(DataChunkTypeString, payload), nil
	}

	if _, okStrings := payload.([]string); okStrings {
		return NewDataChunkWithType(DataChunkTypeStrings, payload), nil
	}

	if _, okMapString := payload.(map[string]string); okMapString {
		return NewDataChunkWithType(DataChunkTypeMapStringToString, payload), nil
	}

	if _, okMapStrings := payload.(map[string][]string); okMapStrings {
		return NewDataChunkWithType(DataChunkTypeMapStringToStrings, payload), nil
	}

	if _, okBytes := payload.([]byte); okBytes {
		return NewDataChunkWithType(DataChunkTypeBytes, payload), nil
	}

	if _, okHeader := payload.(http.Header); okHeader {
		return NewDataChunkWithType(DataChunkTypeHTTPHeader, payload), nil
	}

	if _, okInt := payload.(int); okInt {
		return NewDataChunkWithType(DataChunkTypeInt, payload), nil
	}

	if _, okPromise := payload.(*TaskPromise); okPromise {
		return NewDataChunkWithType(DataChunkTypePromise, payload), nil
	}

	return nil, errors.New("Unsupported payload type")
}
