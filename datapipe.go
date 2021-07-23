package main

import (
	"github.com/google/uuid"
)

type DataPipe struct {
	Done       bool
	Queue      []interface{}
	FromAction Action
	ToAction   Action
	UUID       string
}

func NewDataPipe() *DataPipe {
	return &DataPipe{false, []interface{}{}, nil, nil, uuid.New().String()}
}

func NewDataPipeBetweenActions(fromAction Action, toAction Action) *DataPipe {
	return &DataPipe{
		Done:       false,
		Queue:      []interface{}{},
		FromAction: fromAction,
		ToAction:   toAction,
	}
}

func (dp *DataPipe) Add(x interface{}) {
	dp.Queue = append(dp.Queue, x)
}

func (dp *DataPipe) Remove() interface{} {
	if len(dp.Queue) == 0 {
		return nil
	}

	lastIdx := len(dp.Queue) - 1
	x := dp.Queue[lastIdx]
	dp.Queue = dp.Queue[:lastIdx]

	return x
}
