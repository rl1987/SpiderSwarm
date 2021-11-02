package spsw

import (
	"fmt"

	"github.com/google/uuid"
)

type DataPipe struct {
	Done       bool
	Queue      []*DataChunk
	FromAction Action
	ToAction   Action
	UUID       string
}

func NewDataPipe() *DataPipe {
	return &DataPipe{false, []*DataChunk{}, nil, nil, uuid.New().String()}
}

func NewDataPipeBetweenActions(fromAction Action, toAction Action) *DataPipe {
	return &DataPipe{
		Done:       false,
		Queue:      []*DataChunk{},
		FromAction: fromAction,
		ToAction:   toAction,
		UUID:       uuid.New().String(),
	}
}

func (dp *DataPipe) String() string {
	fromActionName := "<null>"
	if dp.FromAction != nil {
		fromActionName = dp.FromAction.GetName()
	}

	toActionName := "<null>"
	if dp.ToAction != nil {
		toActionName = dp.ToAction.GetName()
	}

	return fmt.Sprintf("<DataPipe %v FromAction.Name: %s, ToAction.Name: %s>", dp.UUID,
		fromActionName, toActionName)
}

func (dp *DataPipe) Add(x interface{}) error {
	if chunk, err := NewDataChunk(x); err == nil {
		dp.Queue = append(dp.Queue, chunk)
	} else {
		return err
	}

	return nil
}

func (dp *DataPipe) AddItem(item *Item) error {
	chunk := NewDataChunkWithType(DataChunkTypeItem, item)

	dp.Queue = append(dp.Queue, chunk)

	return nil
}

func (dp *DataPipe) Remove() interface{} {
	if len(dp.Queue) == 0 {
		return nil
	}

	lastIdx := len(dp.Queue) - 1
	lastChunk := dp.Queue[lastIdx]
	dp.Queue = dp.Queue[:lastIdx]

	if lastChunk.Type == DataChunkTypeItem {
		return lastChunk.PayloadItem
	} else if lastChunk.Type == DataChunkTypePromise {
		return lastChunk.PayloadPromise
	} else if lastChunk.Type == DataChunkTypeValue {
		value := lastChunk.PayloadValue
		return value.GetUnderlyingValue()
	}

	return nil
}
