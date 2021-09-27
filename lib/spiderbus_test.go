package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestSpiderBusBackend struct {
	SpiderBusBackend
	ScheduledTasks []*ScheduledTask
	TaskPromises   []*TaskPromise
	Items          []*Item
}

func NewTestSpiderBusBackend() *TestSpiderBusBackend {
	return &TestSpiderBusBackend{
		ScheduledTasks: []*ScheduledTask{},
		TaskPromises:   []*TaskPromise{},
		Items:          []*Item{},
	}
}

func (tb *TestSpiderBusBackend) SendScheduledTask(scheduledTask *ScheduledTask) error {
	tb.ScheduledTasks = append(tb.ScheduledTasks, scheduledTask)
	return nil
}

func (tb *TestSpiderBusBackend) ReceiveScheduledTask() *ScheduledTask {
	if len(tb.ScheduledTasks) == 0 {
		return nil
	}

	var scheduledTask *ScheduledTask

	scheduledTask, tb.ScheduledTasks = tb.ScheduledTasks[0], tb.ScheduledTasks[1:]

	return scheduledTask
}

func (tb *TestSpiderBusBackend) SendTaskPromise(taskPromise *TaskPromise) error {
	tb.TaskPromises = append(tb.TaskPromises, taskPromise)
	return nil
}

func (tb *TestSpiderBusBackend) ReceiveTaskPromise() *TaskPromise {
	if len(tb.TaskPromises) == 0 {
		return nil
	}

	var taskPromise *TaskPromise

	taskPromise, tb.TaskPromises = tb.TaskPromises[0], tb.TaskPromises[1:]

	return taskPromise
}

func (tb *TestSpiderBusBackend) SendItem(item *Item) error {
	tb.Items = append(tb.Items, item)
	return nil
}

func (tb *TestSpiderBusBackend) ReceiveItem() *Item {
	if len(tb.Items) == 0 {
		return nil
	}

	var item *Item

	item, tb.Items = tb.Items[0], tb.Items[1:]

	return item
}

func TestNewSpiderBus(t *testing.T) {
	spiderBus := NewSpiderBus()
	assert.NotNil(t, spiderBus)
	assert.NotNil(t, spiderBus.UUID)
}

func TestSpiderBusEnqueueDequeueScheduledTask(t *testing.T) {
	testBackend := NewTestSpiderBusBackend()

	spiderBus := NewSpiderBus()
	spiderBus.Backend = testBackend

	gotScheduledTask, err := spiderBus.Dequeue(SpiderBusEntryTypeScheduledTask)
	assert.Nil(t, gotScheduledTask)
	assert.Nil(t, err)

	scheduledTask := &ScheduledTask{
		UUID:            "CAF7E6CB-C888-49B0-B78A-868C11C8BCAB",
		Promise:         TaskPromise{},
		Template:        TaskTemplate{},
		WorkflowName:    "WF0",
		WorkflowVersion: "v.2.5.0.1",
		JobUUID:         "5BD414FF-F827-4227-90A7-215001874BE6",
	}

	err = spiderBus.Enqueue(scheduledTask)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(testBackend.ScheduledTasks))

	gotScheduledTask, err = spiderBus.Dequeue(SpiderBusEntryTypeScheduledTask)
	assert.Nil(t, err)

	assert.Equal(t, scheduledTask, gotScheduledTask)
}

func TestSpiderBusEnqueueDequeueTaskPromise(t *testing.T) {
	testBackend := NewTestSpiderBusBackend()

	spiderBus := NewSpiderBus()
	spiderBus.Backend = testBackend

	gotTaskPromise, err := spiderBus.Dequeue(SpiderBusEntryTypeTaskPromise)
	assert.Nil(t, gotTaskPromise)
	assert.Nil(t, err)

	taskPromise := NewTaskPromise("Task0", "WF0", "2A6E6908-0547-4100-A543-4E99127D0C6D",
		nil)

	err = spiderBus.Enqueue(taskPromise)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(testBackend.TaskPromises))
	assert.Equal(t, taskPromise, testBackend.TaskPromises[0])

	gotTaskPromise, err = spiderBus.Dequeue(SpiderBusEntryTypeTaskPromise)
	assert.Nil(t, err)

	assert.Equal(t, taskPromise, gotTaskPromise)
}

func TestSpiderBusEnqueueDequeueItem(t *testing.T) {
	testBackend := NewTestSpiderBusBackend()

	spiderBus := NewSpiderBus()
	spiderBus.Backend = testBackend

	gotItem, err := spiderBus.Dequeue(SpiderBusEntryTypeItem)
	assert.Nil(t, gotItem)
	assert.Nil(t, err)

	item := NewItem("testItem", "WF0", "3087B17B-2E23-4ECD-A330-B03E63D622E6",
		"135AD4AB-1284-460B-9C0F-54C9FFFACE1A")

	err = spiderBus.Enqueue(item)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(testBackend.Items))
	assert.Equal(t, item, testBackend.Items[0])

	gotItem, err = spiderBus.Dequeue(SpiderBusEntryTypeItem)
	assert.Nil(t, err)

	assert.Equal(t, item, gotItem)
}
