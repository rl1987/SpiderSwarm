package spsw

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type SpiderBusAdapter struct {
	UUID              string
	Bus               *SpiderBus
	ScheduledTasksIn  chan *ScheduledTask
	ScheduledTasksOut chan *ScheduledTask
	TaskPromisesIn    chan *TaskPromise
	TaskPromisesOut   chan *TaskPromise
	TaskResultsIn     chan *TaskResult
	TaskResultsOut    chan *TaskResult
	ItemsIn           chan *Item
	ItemsOut          chan *Item
}

func NewSpiderBusAdapterForWorker(sb *SpiderBus, w *Worker) *SpiderBusAdapter {
	return &SpiderBusAdapter{
		UUID:              uuid.New().String(),
		Bus:               sb,
		ScheduledTasksOut: w.ScheduledTasksIn,
		TaskPromisesIn:    w.TaskPromisesOut,
		TaskResultsIn:     w.TaskResultsOut,
		ItemsIn:           w.ItemsOut,
	}
}

func NewSpiderBusAdapterForExporter(sb *SpiderBus, e *Exporter) *SpiderBusAdapter {
	return &SpiderBusAdapter{
		UUID:     uuid.New().String(),
		Bus:      sb,
		ItemsOut: e.ItemsIn,
	}
}

func NewSpiderBusAdapterForManager(sb *SpiderBus, m *Manager) *SpiderBusAdapter {
	return &SpiderBusAdapter{
		UUID:             uuid.New().String(),
		Bus:              sb,
		TaskPromisesOut:  m.TaskPromisesIn,
		TaskResultsOut:   m.TaskResultsIn,
		ScheduledTasksIn: m.ScheduledTasksOut,
	}
}

func (sba *SpiderBusAdapter) Start() {
	log.Info(fmt.Sprintf("SpiderBusAdapter %s starting run loops", sba.UUID))

	if sba.ScheduledTasksIn != nil {
		go func() {
			for scheduledTask := range sba.ScheduledTasksIn {
				if scheduledTask != nil {
					sba.Bus.Enqueue(scheduledTask)
				}
			}
		}()
	}

	if sba.ScheduledTasksOut != nil {
		go func() {
			for {
				scheduledTask, err := sba.Bus.Dequeue(SpiderBusEntryTypeScheduledTask)

				if scheduledTask == nil || err != nil {
					time.Sleep(1)
					continue
				}

				sba.ScheduledTasksOut <- scheduledTask.(*ScheduledTask)
			}
		}()
	}

	if sba.TaskPromisesIn != nil {
		go func() {
			for taskPromise := range sba.TaskPromisesIn {
				sba.Bus.Enqueue(taskPromise)
			}
		}()
	}

	if sba.TaskPromisesOut != nil {
		go func() {
			for {
				taskPromise, err := sba.Bus.Dequeue(SpiderBusEntryTypeTaskPromise)

				if taskPromise == nil || err != nil {
					time.Sleep(10 * time.Second)
					continue
				}

				sba.TaskPromisesOut <- taskPromise.(*TaskPromise)
			}
		}()
	}

	if sba.TaskResultsIn != nil {
		go func() {
			for taskResult := range sba.TaskResultsIn {
				sba.Bus.Enqueue(taskResult)
			}
		}()
	}

	if sba.TaskResultsOut != nil {
		go func() {
			for {
				taskResult, err := sba.Bus.Dequeue(SpiderBusEntryTypeTaskResult)

				if taskResult == nil || err != nil {
					time.Sleep(1)
					continue
				}

				sba.TaskResultsOut <- taskResult.(*TaskResult)
			}
		}()
	}

	if sba.ItemsIn != nil {
		go func() {
			for item := range sba.ItemsIn {
				sba.Bus.Enqueue(item)
			}
		}()
	}

	if sba.ItemsOut != nil {
		go func() {
			for {
				item, err := sba.Bus.Dequeue(SpiderBusEntryTypeItem)

				if item == nil || err != nil {
					time.Sleep(1)
					continue
				}

				sba.ItemsOut <- item.(*Item)
			}
		}()
	}
}
