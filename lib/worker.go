package spsw

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	UUID             string
	ScheduledTasksIn chan *ScheduledTask
	TaskPromisesOut  chan *TaskPromise
	TaskReportsOut   chan *TaskReport
	ItemsOut         chan *Item
	Done             chan interface{}
}

func NewWorker() *Worker {
	return &Worker{
		UUID:             uuid.New().String(),
		ScheduledTasksIn: make(chan *ScheduledTask),
		TaskPromisesOut:  make(chan *TaskPromise),
		TaskReportsOut:   make(chan *TaskReport),
		ItemsOut:         make(chan *Item),
		Done:             make(chan interface{}),
	}
}

func (w *Worker) String() string {
	return fmt.Sprintf("<Worker %s>", w.UUID)
}

func (w *Worker) executeTask(task *Task) error {
	err := task.Run()
	if err != nil {
		log.Error(fmt.Sprintf("Task %v failed with error: %v", task, err))

		taskReport := NewTaskReport(task.JobUUID, task.UUID, task.Name, false, err)
		w.TaskReportsOut <- taskReport

		return err
	}

	for _, outDP := range task.Outputs {
		if len(outDP.Queue) == 0 {
			continue
		}

		x := outDP.Remove()

		if item, okItem := x.(*Item); okItem {
			w.ItemsOut <- item
		}

		if promise, okPromise := x.(*TaskPromise); okPromise {
			w.TaskPromisesOut <- promise
		}
	}

	taskReport := NewTaskReport(task.JobUUID, task.UUID, task.Name, true, nil)
	w.TaskReportsOut <- taskReport

	return nil
}

func (w *Worker) Run() error {
	log.Info(fmt.Sprintf("Starting runloop for worker %s", w.UUID))

	for {
		select {
		case scheduledTask := <-w.ScheduledTasksIn:
			if scheduledTask == nil {
				continue
			}

			log.Info(fmt.Printf("Worker %s got scheduled task %v", w.UUID, scheduledTask))

			task := NewTaskFromScheduledTask(scheduledTask)
			log.Info(fmt.Sprintf("Worker %s running task %v", w.UUID, task))

			w.executeTask(task)
		case <-w.Done:
			return nil
		}
	}

}
