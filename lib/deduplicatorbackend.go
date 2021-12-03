package spsw

type DeduplicatorBackend interface {
	IsScheduledTaskDuplicated(scheduledTask *ScheduledTask) bool
	NoteScheduledTask(scheduledTask *ScheduledTask) error
}

type AbstractDeduplicatorBackend struct {
	DeduplicatorBackend
}
