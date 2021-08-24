package spsw

type SpiderBusBackend interface {
	SendScheduledTask(scheduledTask *ScheduledTask) error
	ReceiveScheduledTask() *ScheduledTask
	SendTaskPromise(taskPromise *TaskPromise) error
	ReceiveTaskPromise() *TaskPromise
	SendItem(item *Item) error
	ReceiveItem() *Item
}

type AbstractSpiderBusBackend struct {
	SpiderBusBackend
}
