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
	SendTaskReport(taskReport *TaskReport) error
	ReceiveTaskReport() *TaskReport
}

type AbstractSpiderBusBackend struct {
	SpiderBusBackend
}
