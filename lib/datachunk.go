package spsw

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

const DataChunkTypeString = "DataChunkTypeString"
const DataChunkTypeMapStringToString = "DataChunkTypeMapStringToString"
const DataChunkTypeMapStringToStrings = "DataChunkTypeMapStringToStrings"
const DataChunkTypeBytes = "DataChunkTypeBytes"
const DataChunkTypeHTTPHeader = "DataChunkTypeHTTPHeader"
const DataChunkTypeInt = "DataChunkTypeInt"
const DataChunkTypeItem = "DataChunkTypeItem"
const DataChunkTypePromise = "DataChunkTypePromise"
const DataChunkTypeValue = "DataChunkTypeValue"

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

	if _, okItem := payload.(*Item); okItem {
		return NewDataChunkWithType(DataChunkTypeItem, payload), nil
	}

	if _, okPromise := payload.(*TaskPromise); okPromise {
		return NewDataChunkWithType(DataChunkTypePromise, payload), nil
	}

	if _, okStrings := payload.([]string); okStrings {
		return NewDataChunkWithType(DataChunkTypeValue, NewValueFromStrings(payload.([]string))), nil
	}

	if _, okValue := payload.(*Value); okValue {
		return NewDataChunkWithType(DataChunkTypeValue, payload), nil
	}

	return nil, errors.New("Unsupported payload type")
}
