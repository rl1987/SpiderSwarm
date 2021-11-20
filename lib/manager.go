package spsw

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	UUID              string
	TaskPromisesIn    chan *TaskPromise
	TaskReportsIn     chan *TaskReport
	ScheduledTasksOut chan *ScheduledTask
	CurrentWorkflow   *Workflow // TODO: support multiple scraping jobs running concurrently
	JobUUID           string
	NPendingTasks     int
}

func NewManager() *Manager {
	return &Manager{
		UUID:              uuid.New().String(),
		TaskPromisesIn:    make(chan *TaskPromise),
		TaskReportsIn:     make(chan *TaskReport),
		ScheduledTasksOut: make(chan *ScheduledTask),
		NPendingTasks:     0,
	}
}

func (m *Manager) String() string {
	return fmt.Sprintf("<Manager %s>", m.UUID)
}

func (m *Manager) StartScrapingJob(w *Workflow) {
	m.JobUUID = uuid.New().String()
	m.CurrentWorkflow = w
}

func (m *Manager) createScheduledTaskFromPromise(promise *TaskPromise, jobUUID string) *ScheduledTask {
	taskTempl := m.CurrentWorkflow.FindTaskTemplate(promise.TaskName)
	if taskTempl == nil {
		return nil
	}

	scheduledTask := NewScheduledTask(promise, taskTempl,
		m.CurrentWorkflow.Name, m.CurrentWorkflow.Version, jobUUID)

	// TODO: log this

	return scheduledTask
}

func (m *Manager) Run() error {
	log.Info(fmt.Sprintf("Starting runloop for manager %s", m.UUID))

	for _, taskTempl := range m.CurrentWorkflow.TaskTemplates {
		if !taskTempl.Initial {
			continue
		}

		newPromise := NewTaskPromise(taskTempl.TaskName,
			m.CurrentWorkflow.Name, m.JobUUID, map[string]*DataChunk{})
		log.Info(fmt.Sprintf("Fulfilling promise %v", newPromise))

		scheduledTask := NewScheduledTask(newPromise, &taskTempl,
			m.CurrentWorkflow.Name, m.CurrentWorkflow.Version, m.JobUUID)

		log.Info(fmt.Sprintf("Created scheduled task %v", scheduledTask))

		m.ScheduledTasksOut <- scheduledTask
		m.NPendingTasks++
	}

	for m.NPendingTasks > 0 {
		select {
		case promise := <-m.TaskPromisesIn:
			if promise == nil {
				continue
			}

			for _, p := range promise.Splay() {
				newScheduledTask := m.createScheduledTaskFromPromise(p, m.JobUUID)
				if newScheduledTask == nil {
					continue
				}

				log.Info(fmt.Sprintf("Created scheduled task %v", newScheduledTask))
				m.ScheduledTasksOut <- newScheduledTask
				m.NPendingTasks++
			}
		case report := <-m.TaskReportsIn:
			if report == nil {
				continue
			}

			m.NPendingTasks--
		}
	}

	return nil
}
