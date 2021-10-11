package spsw

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

const DataChunkTypeItem = "DataChunkTypeItem"
const DataChunkTypePromise = "DataChunkTypePromise"
const DataChunkTypeValue = "DataChunkTypeValue"

type DataChunk struct {
	Type           string
	PayloadItem    *Item
	PayloadPromise *TaskPromise
	PayloadValue   *Value
	UUID           string
}

func NewDataChunkWithType(t string, payload interface{}) *DataChunk {
	dc := &DataChunk{
		Type: t,
		UUID: uuid.New().String(),
	}

	if t == DataChunkTypeItem {
		dc.PayloadItem = payload.(*Item)
	} else if t == DataChunkTypePromise {
		dc.PayloadPromise = payload.(*TaskPromise)
	} else if t == DataChunkTypeValue {
		dc.PayloadValue = payload.(*Value)
	}

	return dc
}

func NewDataChunk(payload interface{}) (*DataChunk, error) {
	if _, okStr := payload.(string); okStr {
		return NewDataChunkWithType(DataChunkTypeValue, NewValueFromString(payload.(string))), nil
	}

	if _, okMapString := payload.(map[string]string); okMapString {
		return NewDataChunkWithType(DataChunkTypeValue, NewValueFromMapStringToString(payload.(map[string]string))), nil
	}

	if _, okMapStrings := payload.(map[string][]string); okMapStrings {
		return NewDataChunkWithType(DataChunkTypeValue, NewValueFromMapStringToStrings(payload.(map[string][]string))), nil
	}

	if _, okBytes := payload.([]byte); okBytes {
		return NewDataChunkWithType(DataChunkTypeValue, NewValueFromBytes(payload.([]byte))), nil
	}

	if _, okHeader := payload.(http.Header); okHeader {
		return NewDataChunkWithType(DataChunkTypeValue, NewValueFromHTTPHeaders(payload.(http.Header))), nil
	}

	if _, okInt := payload.(int); okInt {
		return NewDataChunkWithType(DataChunkTypeValue, NewValueFromInt(payload.(int))), nil
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
