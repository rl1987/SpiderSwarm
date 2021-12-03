package spsw

import (
	"github.com/google/uuid"
)

type Deduplicator struct {
	UUID    string
	Backend DeduplicatorBackend
}

func NewDeduplicator(backendAddr string) *Deduplicator {
	return &Deduplicator{
		UUID:    uuid.New().String(),
		Backend: NewRedisDeduplicatorBackend(backendAddr, ""),
	}
}

func (d *Deduplicator) IsScheduledTaskDuplicated(scheduledTask *ScheduledTask) bool {
	return d.Backend.IsScheduledTaskDuplicated(scheduledTask)
}

func (d *Deduplicator) NoteScheduledTask(scheduledTask *ScheduledTask) error {
	return d.Backend.NoteScheduledTask(scheduledTask)
}
