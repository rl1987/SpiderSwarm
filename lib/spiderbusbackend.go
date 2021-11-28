package spsw

type SpiderBusBackend interface {
	IsScheduledTaskDuplicated(scheduledTask *ScheduledTask, jobUUID string) bool
	SendScheduledTask(scheduledTask *ScheduledTask) error
	ReceiveScheduledTask() *ScheduledTask
	IsTaskPromiseDuplicated(taskPromise *TaskPromise, jobUUID string) bool
	SendTaskPromise(taskPromise *TaskPromise) error
	ReceiveTaskPromise() *TaskPromise
	IsItemDuplicated(item *Item, jobUUID string) bool
	SendItem(item *Item) error
	ReceiveItem() *Item
	SendTaskResult(taskResult *TaskResult) error
	ReceiveTaskResult() *TaskResult
}

type AbstractSpiderBusBackend struct {
	SpiderBusBackend
}
