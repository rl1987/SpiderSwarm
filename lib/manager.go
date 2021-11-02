package spsw

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	UUID              string
	TaskPromisesIn    chan *TaskPromise
	ScheduledTasksOut chan *ScheduledTask
	CurrentWorkflow   *Workflow // TODO: support multiple scraping jobs running concurrently
	JobUUID           string
}

func NewManager() *Manager {
	return &Manager{
		UUID:              uuid.New().String(),
		TaskPromisesIn:    make(chan *TaskPromise),
		ScheduledTasksOut: make(chan *ScheduledTask),
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

		spew.Dump(scheduledTask)

		m.ScheduledTasksOut <- scheduledTask
	}

	// TODO: workers should report back about success/failure of the task;
	//       managers should report back about status of the scraping job.
	for promise := range m.TaskPromisesIn {
		if promise == nil {
			continue
		}

		for _, p := range promise.Splay() {
			newScheduledTask := m.createScheduledTaskFromPromise(p, m.JobUUID)
			if newScheduledTask == nil {
				continue
			}

			spew.Dump(newScheduledTask)
			m.ScheduledTasksOut <- newScheduledTask
		}

	}

	return nil
}
