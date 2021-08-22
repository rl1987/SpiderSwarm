package spiderswarm

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SpiderBus struct {
	UUID    string
	Backend SpiderBusBackend
}

func (sb *SpiderBus) NewSpiderBus() *SpiderBus {
	return &SpiderBus{
		UUID: uuid.New().String(),
	}
}

func (sb *SpiderBus) Enqueue(x interface{}) error {
	if scheduledTask, okTask := x.(*ScheduledTask); okTask {
		return sb.Backend.SendScheduledTask(scheduledTask)
	}

	if promise, okPromise := x.(*TaskPromise); okPromise {
		return sb.Backend.SendTaskPromise(promise)
	}

	if item, okItem := x.(*Item); okItem {
		return sb.Backend.SendItem(item)
	}

	return errors.New(fmt.Sprintf("SpiderBus.Enqueue: argument not recognised: %v", x))
}

const SpiderBusEntryTypeScheduledTask = "SpiderBusEntryTypeScheduledTask"
const SpiderBusEntryTypeTaskPromise = "SpiderBusEntryTypeTaskPromise"
const SpiderBusEntryTypeItem = "SpiderBusEntryTypeItem"

func (sb *SpiderBus) Dequeue(entryType string) (interface{}, error) {
	if entryType == SpiderBusEntryTypeScheduledTask {
		return sb.Backend.ReceiveScheduledTask(), nil
	}

	if entryType == SpiderBusEntryTypeTaskPromise {
		return sb.Backend.ReceiveTaskPromise(), nil
	}

	if entryType == SpiderBusEntryTypeItem {
		return sb.Backend.ReceiveItem(), nil
	}

	return nil, errors.New(fmt.Sprintf("SpiderBus.Dequeue: unrecognised entryType: %s", entryType))
}
