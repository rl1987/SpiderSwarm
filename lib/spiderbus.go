package spsw

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SpiderBus struct {
	UUID    string
	Backend SpiderBusBackend
}

func NewSpiderBus() *SpiderBus {
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

	if report, okReport := x.(*TaskReport); okReport {
		return sb.Backend.SendTaskReport(report)
	}

	if result, okResult := x.(*TaskResult); okResult {
		return sb.Backend.SendTaskResult(result)
	}

	if item, okItem := x.(*Item); okItem {
		return sb.Backend.SendItem(item)
	}

	return errors.New(fmt.Sprintf("SpiderBus.Enqueue: argument not recognised: %v", x))
}

const SpiderBusEntryTypeScheduledTask = "SpiderBusEntryTypeScheduledTask"
const SpiderBusEntryTypeTaskPromise = "SpiderBusEntryTypeTaskPromise"
const SpiderBusEntryTypeTaskReport = "SpiderBusEntryTypeTaskReport"
const SpiderBusEntryTypeTaskResult = "SpiderBusEntryTypeTaskResult"
const SpiderBusEntryTypeItem = "SpiderBusEntryTypeItem"

func (sb *SpiderBus) Dequeue(entryType string) (interface{}, error) {
	if sb.Backend == nil {
		return nil, errors.New("SpiderBus has no backend assigned")
	}

	if entryType == SpiderBusEntryTypeScheduledTask {
		return sb.Backend.ReceiveScheduledTask(), nil
	}

	if entryType == SpiderBusEntryTypeTaskPromise {
		return sb.Backend.ReceiveTaskPromise(), nil
	}

	if entryType == SpiderBusEntryTypeTaskReport {
		return sb.Backend.ReceiveTaskReport(), nil
	}

	if entryType == SpiderBusEntryTypeItem {
		return sb.Backend.ReceiveItem(), nil
	}

	if entryType == SpiderBusEntryTypeTaskResult {
		return sb.Backend.ReceiveTaskResult(), nil
	}

	return nil, errors.New(fmt.Sprintf("SpiderBus.Dequeue: unrecognised entryType: %s", entryType))
}
