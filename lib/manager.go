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
	TaskResultsIn     chan *TaskResult
	ScheduledTasksOut chan *ScheduledTask
	CurrentWorkflow   *Workflow // TODO: support multiple scraping jobs running concurrently
	JobUUID           string
	NPendingTasks     int
	NFinishedTasks    int
	NFailedTasks      int
	NScheduledTasks   int
}

func NewManager() *Manager {
	return &Manager{
		UUID:              uuid.New().String(),
		TaskPromisesIn:    make(chan *TaskPromise),
		TaskResultsIn:     make(chan *TaskResult),
		ScheduledTasksOut: make(chan *ScheduledTask),
		NPendingTasks:     0,
		NFinishedTasks:    0,
		NFailedTasks:      0,
		NScheduledTasks:   0,
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

	return scheduledTask
}

func (m *Manager) logPendingTasks() {
	log.Info(fmt.Sprintf("Manager %s tasks: %d pending, %d finished out of %d scheduled", m.UUID, m.NPendingTasks,
		m.NFinishedTasks, m.NScheduledTasks))
}

func (m *Manager) handleTaskPromise(promise *TaskPromise) {
	if promise == nil {
		return
	}

	for _, p := range promise.Splay() {
		newScheduledTask := m.createScheduledTaskFromPromise(p, m.JobUUID)
		if newScheduledTask == nil {
			continue
		}

		log.Info(fmt.Sprintf("Created scheduled task %v", newScheduledTask))
		m.ScheduledTasksOut <- newScheduledTask
		m.NPendingTasks++
		m.NScheduledTasks++
		m.logPendingTasks()
	}
}

func (m *Manager) handleItem(item *Item) {
	// TODO: push item to SpiderBus
}

func (m *Manager) processTaskResult(taskResult *TaskResult) {
	if taskResult == nil {
		return
	}

	m.NPendingTasks--

	if !taskResult.Succeeded {
		spew.Dump(taskResult.Error)
		return
	}

	m.NFinishedTasks++

	for _, chunks := range taskResult.OutputDataChunks {
		for _, chunk := range chunks {
			if chunk.Type == DataChunkTypePromise {
				m.handleTaskPromise(chunk.PayloadPromise)
			} else if chunk.Type == DataChunkTypeItem {
				m.handleItem(chunk.PayloadItem)
			}
		}
	}
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
		m.NScheduledTasks++

		m.logPendingTasks()
	}

	for taskResult := range m.TaskResultsIn {
		spew.Dump(taskResult)

		m.processTaskResult(taskResult)

		if m.NPendingTasks == 0 {
			break
		}

		m.logPendingTasks()
	}

	return nil
}
